package content

import (
	"encoding/json"
	"sort"
	"strings"
	"unicode"
)

type TitleSlug struct {
	Title string
	Slug  string
}

func AutoLink(body []byte, candidates []TitleSlug, selfSlug string) ([]byte, error) {
	if len(body) == 0 || len(candidates) == 0 {
		return body, nil
	}
	var root Node
	if err := json.Unmarshal(body, &root); err != nil {
		return body, err
	}

	entries := make([]TitleSlug, 0, len(candidates))
	for _, c := range candidates {
		if c.Slug == selfSlug || c.Title == "" {
			continue
		}
		entries = append(entries, c)
	}
	sort.Slice(entries, func(i, j int) bool {
		return len([]rune(entries[i].Title)) > len([]rune(entries[j].Title))
	})

	rewriteNode(&root, entries)
	out, err := json.Marshal(root)
	if err != nil {
		return body, err
	}
	return out, nil
}

func rewriteNode(n *Node, entries []TitleSlug) {
	if n == nil {
		return
	}
	// Skip nodes whose content is verbatim — wikilinks, code blocks, inline
	// code, and embedded graph widgets must not get their text rewritten.
	switch n.Type {
	case "wikilink", "codeBlock", "relationshipGraph":
		return
	}
	if len(n.Content) == 0 {
		return
	}
	next := make([]Node, 0, len(n.Content))
	for i := range n.Content {
		child := n.Content[i]
		if child.Type == "text" && len(child.Text) > 0 && !hasAnyMark(&child) {
			split := splitTextByTitles(child.Text, entries)
			next = append(next, split...)
			continue
		}
		rewriteNode(&child, entries)
		next = append(next, child)
	}
	n.Content = next
}

func hasAnyMark(n *Node) bool {
	if n == nil || n.Attrs == nil {
		return false
	}
	if marks, ok := n.Attrs["marks"]; ok {
		if arr, ok := marks.([]any); ok && len(arr) > 0 {
			return true
		}
	}
	return false
}

func splitTextByTitles(text string, entries []TitleSlug) []Node {
	if text == "" {
		return []Node{{Type: "text", Text: text}}
	}
	lower := strings.ToLower(text)
	bestStart := -1
	bestEnd := -1
	bestSlug := ""
	for _, e := range entries {
		needle := strings.ToLower(e.Title)
		if needle == "" {
			continue
		}
		idx := strings.Index(lower, needle)
		if idx < 0 {
			continue
		}
		end := idx + len(needle)
		if !isWordBoundary(text, idx, end) {
			continue
		}
		if bestStart < 0 || idx < bestStart || (idx == bestStart && (end-idx) > (bestEnd-bestStart)) {
			bestStart = idx
			bestEnd = end
			bestSlug = e.Slug
		}
	}
	if bestStart < 0 {
		return []Node{{Type: "text", Text: text}}
	}
	out := make([]Node, 0, 3)
	if bestStart > 0 {
		out = append(out, Node{Type: "text", Text: text[:bestStart]})
	}
	label := text[bestStart:bestEnd]
	out = append(out, Node{
		Type: "wikilink",
		Attrs: map[string]any{
			"slug":  bestSlug,
			"label": label,
		},
	})
	if bestEnd < len(text) {
		out = append(out, splitTextByTitles(text[bestEnd:], entries)...)
	}
	return out
}

func isWordBoundary(text string, start, end int) bool {
	if start > 0 {
		r, _ := lastRuneAt(text, start)
		if isWordRune(r) {
			return false
		}
	}
	if end < len(text) {
		r := []rune(text[end:])
		if len(r) > 0 && isWordRune(r[0]) {
			return false
		}
	}
	return true
}

func lastRuneAt(s string, before int) (rune, int) {
	if before <= 0 {
		return 0, 0
	}
	prefix := s[:before]
	r := []rune(prefix)
	if len(r) == 0 {
		return 0, 0
	}
	return r[len(r)-1], len(r) - 1
}

func isWordRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}
