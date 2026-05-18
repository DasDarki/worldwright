package content

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
)

// ScrubInvisibleWikilinks walks a TipTap body and replaces any wikilink whose
// target slug is not present in `visibleSlugs` with a plain text node carrying
// the link's label (or the slug itself if no label was set). Returns the new
// body bytes; if nothing changed, the input is returned unmodified.
func ScrubInvisibleWikilinks(body []byte, visibleSlugs map[string]bool) []byte {
	if len(body) == 0 {
		return body
	}
	var doc map[string]any
	if err := json.Unmarshal(body, &doc); err != nil {
		return body
	}
	changed := false
	walkScrub(doc, visibleSlugs, &changed)
	if !changed {
		return body
	}
	out, err := json.Marshal(doc)
	if err != nil {
		return body
	}
	return out
}

func walkScrub(node any, visible map[string]bool, changed *bool) {
	switch v := node.(type) {
	case map[string]any:
		if v["type"] == "wikilink" {
			if attrs, ok := v["attrs"].(map[string]any); ok {
				slug, _ := attrs["slug"].(string)
				if slug != "" && !visible[slug] {
					label, _ := attrs["label"].(string)
					if label == "" {
						label = slug
					}
					v["type"] = "text"
					v["text"] = label
					delete(v, "attrs")
					*changed = true
					return
				}
			}
		}
		for _, child := range v {
			walkScrub(child, visible, changed)
		}
	case []any:
		for _, child := range v {
			walkScrub(child, visible, changed)
		}
	}
}

// VisibleSlugSet returns the set of entity slugs visible to the caller, given
// the visibility filter list. Used to build the input for ScrubInvisibleWikilinks.
func VisibleSlugSet(ctx context.Context, db *sql.DB, visibility []string) (map[string]bool, error) {
	if len(visibility) == 0 {
		return map[string]bool{}, nil
	}
	placeholders := strings.TrimRight(strings.Repeat("?,", len(visibility)), ",")
	args := make([]any, 0, len(visibility))
	for _, v := range visibility {
		args = append(args, v)
	}
	rows, err := db.QueryContext(ctx,
		`SELECT slug FROM entities WHERE visibility IN (`+placeholders+`)`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	set := map[string]bool{}
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		set[s] = true
	}
	return set, rows.Err()
}
