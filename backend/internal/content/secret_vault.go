package content

import (
	"encoding/json"
)

// ScrubSecretVaults walks a TipTap body and replaces the content of every
// secretVault node whose `author_id` is not the given current user. The
// node itself is preserved (so we can merge it back on save) but its
// content is swapped for a sealed placeholder so neither plain text nor
// wikilink slugs leak. Pass userID = 0 for "no current user" (anonymous /
// share-token / markdown export) which scrubs every vault.
func ScrubSecretVaults(body []byte, userID int64) []byte {
	if len(body) == 0 {
		return body
	}
	var doc map[string]any
	if err := json.Unmarshal(body, &doc); err != nil {
		return body
	}
	changed := false
	walkScrubVaults(doc, userID, &changed)
	if !changed {
		return body
	}
	out, err := json.Marshal(doc)
	if err != nil {
		return body
	}
	return out
}

func walkScrubVaults(node any, userID int64, changed *bool) {
	switch v := node.(type) {
	case map[string]any:
		if v["type"] == "secretVault" {
			attrs, _ := v["attrs"].(map[string]any)
			if attrs == nil {
				attrs = map[string]any{}
				v["attrs"] = attrs
			}
			authorID := readInt64(attrs["author_id"])
			if userID == 0 || authorID == 0 || authorID != userID {
				attrs["redacted"] = true
				v["content"] = []any{
					map[string]any{
						"type": "paragraph",
						"content": []any{
							map[string]any{
								"type": "text",
								"text": "▮▮▮ sealed ▮▮▮",
							},
						},
					},
				}
				*changed = true
			}
			return // don't descend
		}
		for _, child := range v {
			walkScrubVaults(child, userID, changed)
		}
	case []any:
		for _, child := range v {
			walkScrubVaults(child, userID, changed)
		}
	}
}

// MergeSecretVaults reconciles a new body that may have been edited by a
// non-author admin. The rules:
//
//   - For each secretVault in NEW body keyed by vault_id:
//       * if redacted=true, use the OLD body's content (caller couldn't
//         actually read it, so they couldn't have meaningfully edited it).
//       * else if author_id == userID, accept the new content.
//       * else (someone else's vault, somehow not redacted) — also restore
//         from OLD, defense in depth.
//   - Any vault present in OLD but missing in NEW:
//       * if current user IS its author: dropped intentionally — leave out.
//       * else: re-attach at end of NEW body so a non-author cannot delete
//         another author's vault by deleting its sealed placeholder.
//
// Vaults without a vault_id are passed through verbatim (legacy data).
func MergeSecretVaults(oldBody, newBody []byte, userID int64) []byte {
	if len(newBody) == 0 {
		return newBody
	}
	oldVaults := collectVaultsByID(oldBody)
	if len(oldVaults) == 0 {
		// Nothing to merge / restore — but still need to drop redacted flags
		// from new body so the saved JSON is clean.
		return cleanRedactedFlags(newBody)
	}
	var newDoc map[string]any
	if err := json.Unmarshal(newBody, &newDoc); err != nil {
		return newBody
	}

	seenIDs := map[string]bool{}
	walkMergeVaults(newDoc, oldVaults, userID, seenIDs)

	// Re-attach any OLD vault that the non-author dropped (the redacted
	// placeholder was deleted from the new body).
	for id, oldNode := range oldVaults {
		if seenIDs[id] {
			continue
		}
		attrs, _ := oldNode["attrs"].(map[string]any)
		if attrs == nil {
			continue
		}
		authorID := readInt64(attrs["author_id"])
		if authorID == userID {
			continue // user dropped their own vault — honor it
		}
		// append to doc.content
		if content, ok := newDoc["content"].([]any); ok {
			newDoc["content"] = append(content, oldNode)
		}
	}

	out, err := json.Marshal(newDoc)
	if err != nil {
		return newBody
	}
	return out
}

func collectVaultsByID(body []byte) map[string]map[string]any {
	out := map[string]map[string]any{}
	if len(body) == 0 {
		return out
	}
	var doc map[string]any
	if err := json.Unmarshal(body, &doc); err != nil {
		return out
	}
	var walk func(any)
	walk = func(n any) {
		switch v := n.(type) {
		case map[string]any:
			if v["type"] == "secretVault" {
				if attrs, ok := v["attrs"].(map[string]any); ok {
					if id, ok := attrs["vault_id"].(string); ok && id != "" {
						out[id] = v
					}
				}
				return
			}
			for _, child := range v {
				walk(child)
			}
		case []any:
			for _, child := range v {
				walk(child)
			}
		}
	}
	walk(doc)
	return out
}

func walkMergeVaults(node any, oldVaults map[string]map[string]any, userID int64, seenIDs map[string]bool) {
	switch v := node.(type) {
	case map[string]any:
		if v["type"] == "secretVault" {
			attrs, _ := v["attrs"].(map[string]any)
			if attrs == nil {
				return
			}
			id, _ := attrs["vault_id"].(string)
			if id == "" {
				delete(attrs, "redacted")
				return
			}
			seenIDs[id] = true
			redacted, _ := attrs["redacted"].(bool)
			authorID := readInt64(attrs["author_id"])
			delete(attrs, "redacted") // never persist the flag

			useOld := redacted || authorID == 0 || authorID != userID
			if !useOld {
				return // keep incoming content
			}
			if old, ok := oldVaults[id]; ok {
				if oldContent, ok := old["content"]; ok {
					v["content"] = oldContent
				}
				if oldAttrs, ok := old["attrs"].(map[string]any); ok {
					// preserve author_id from DB so it can't be spoofed
					if a := readInt64(oldAttrs["author_id"]); a != 0 {
						attrs["author_id"] = a
					}
				}
			}
			return
		}
		for _, child := range v {
			walkMergeVaults(child, oldVaults, userID, seenIDs)
		}
	case []any:
		for _, child := range v {
			walkMergeVaults(child, oldVaults, userID, seenIDs)
		}
	}
}

// cleanRedactedFlags strips the transient `redacted` attribute from every
// secretVault attrs map. Used when merging is a no-op because there's no
// existing body, but we still don't want the runtime flag in the DB.
func cleanRedactedFlags(body []byte) []byte {
	if len(body) == 0 {
		return body
	}
	var doc map[string]any
	if err := json.Unmarshal(body, &doc); err != nil {
		return body
	}
	changed := false
	var walk func(any)
	walk = func(n any) {
		switch v := n.(type) {
		case map[string]any:
			if v["type"] == "secretVault" {
				if attrs, ok := v["attrs"].(map[string]any); ok {
					if _, hadRedacted := attrs["redacted"]; hadRedacted {
						delete(attrs, "redacted")
						changed = true
					}
				}
				return
			}
			for _, child := range v {
				walk(child)
			}
		case []any:
			for _, child := range v {
				walk(child)
			}
		}
	}
	walk(doc)
	if !changed {
		return body
	}
	out, err := json.Marshal(doc)
	if err != nil {
		return body
	}
	return out
}

func readInt64(v any) int64 {
	switch n := v.(type) {
	case float64:
		return int64(n)
	case int64:
		return n
	case int:
		return int64(n)
	case json.Number:
		i, _ := n.Int64()
		return i
	}
	return 0
}
