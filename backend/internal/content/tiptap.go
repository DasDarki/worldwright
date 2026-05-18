package content

import (
	"encoding/json"
	"strings"
)

type Node struct {
	Type    string         `json:"type"`
	Text    string         `json:"text,omitempty"`
	Attrs   map[string]any `json:"attrs,omitempty"`
	Content []Node         `json:"content,omitempty"`
}

func ParseBody(body []byte) (string, []string, error) {
	if len(body) == 0 {
		return "", nil, nil
	}
	var root Node
	if err := json.Unmarshal(body, &root); err != nil {
		return "", nil, err
	}
	var sb strings.Builder
	var slugs []string
	walk(&root, &sb, &slugs)
	return strings.TrimSpace(sb.String()), dedupe(slugs), nil
}

func walk(n *Node, sb *strings.Builder, slugs *[]string) {
	if n == nil {
		return
	}
	switch n.Type {
	case "wikilink":
		if slug, ok := n.Attrs["slug"].(string); ok && slug != "" {
			*slugs = append(*slugs, slug)
			if label, ok := n.Attrs["label"].(string); ok && label != "" {
				sb.WriteString(label)
			} else {
				sb.WriteString(slug)
			}
			sb.WriteByte(' ')
		}
		return
	case "text":
		sb.WriteString(n.Text)
		sb.WriteByte(' ')
	}
	for i := range n.Content {
		walk(&n.Content[i], sb, slugs)
	}
	if n.Type == "paragraph" || n.Type == "heading" || n.Type == "blockquote" {
		sb.WriteByte('\n')
	}
}

func dedupe(in []string) []string {
	seen := make(map[string]bool, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		if seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	return out
}
