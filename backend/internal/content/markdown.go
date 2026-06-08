package content

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ToMarkdown converts a TipTap body (the dialect Worldwright uses) into a
// GitHub-flavoured markdown string. Unknown nodes degrade gracefully to
// their text content so nothing is silently lost.
func ToMarkdown(body []byte) string {
	if len(body) == 0 {
		return ""
	}
	var root Node
	if err := json.Unmarshal(body, &root); err != nil {
		return ""
	}
	w := &mdWriter{}
	renderBlock(&root, w, 0)
	return strings.TrimSpace(w.b.String()) + "\n"
}

type mdWriter struct {
	b strings.Builder
}

func (w *mdWriter) Write(s string)                 { w.b.WriteString(s) }
func (w *mdWriter) Writef(f string, a ...any)      { fmt.Fprintf(&w.b, f, a...) }
func (w *mdWriter) Newline()                       { w.b.WriteByte('\n') }
func (w *mdWriter) BlockBreak()                    { w.b.WriteString("\n\n") }

func renderBlock(n *Node, w *mdWriter, depth int) {
	if n == nil {
		return
	}
	switch n.Type {
	case "doc":
		for i := range n.Content {
			renderBlock(&n.Content[i], w, depth)
			if i < len(n.Content)-1 {
				w.BlockBreak()
			}
		}
	case "paragraph":
		renderInline(n, w)
	case "heading":
		level := 1
		if v, ok := n.Attrs["level"].(float64); ok {
			level = int(v)
		} else if v, ok := n.Attrs["level"].(int); ok {
			level = v
		}
		if level < 1 {
			level = 1
		}
		if level > 6 {
			level = 6
		}
		w.Write(strings.Repeat("#", level) + " ")
		renderInline(n, w)
	case "blockquote":
		var inner mdWriter
		for i := range n.Content {
			renderBlock(&n.Content[i], &inner, depth)
			if i < len(n.Content)-1 {
				inner.BlockBreak()
			}
		}
		for _, line := range strings.Split(strings.TrimRight(inner.b.String(), "\n"), "\n") {
			w.Write("> " + line + "\n")
		}
	case "bulletList":
		for i := range n.Content {
			renderListItem(&n.Content[i], w, depth, "-", 0)
			if i < len(n.Content)-1 {
				w.Newline()
			}
		}
	case "orderedList":
		for i := range n.Content {
			renderListItem(&n.Content[i], w, depth, "", i+1)
			if i < len(n.Content)-1 {
				w.Newline()
			}
		}
	case "horizontalRule":
		w.Write("---")
	case "codeBlock":
		lang := ""
		if v, ok := n.Attrs["language"].(string); ok {
			lang = v
		}
		w.Write("```" + lang + "\n")
		w.Write(plainText(n))
		w.Write("\n```")
	case "callout":
		variant := "info"
		if v, ok := n.Attrs["variant"].(string); ok && v != "" {
			variant = v
		}
		header := map[string]string{
			"info": "ℹ Info",
			"warn": "⚠ Warning",
			"note": "✎ Note",
			"lore": "✦ Lore",
		}[variant]
		if header == "" {
			header = "Note"
		}
		var inner mdWriter
		for i := range n.Content {
			renderBlock(&n.Content[i], &inner, depth)
			if i < len(n.Content)-1 {
				inner.BlockBreak()
			}
		}
		w.Write("> **" + header + "**\n")
		for _, line := range strings.Split(strings.TrimRight(inner.b.String(), "\n"), "\n") {
			w.Write("> " + line + "\n")
		}
	case "spoiler":
		var inner mdWriter
		for i := range n.Content {
			renderBlock(&n.Content[i], &inner, depth)
			if i < len(n.Content)-1 {
				inner.BlockBreak()
			}
		}
		w.Write("<details><summary>Spoiler</summary>\n\n")
		w.Write(inner.b.String())
		w.Write("\n</details>")
	case "table":
		renderTable(n, w)
	case "wwimage":
		alt, _ := n.Attrs["alt"].(string)
		src, _ := n.Attrs["src"].(string)
		w.Writef("![%s](%s)", escapeMd(alt), src)
	case "relationshipGraph":
		count := 0
		if v, ok := n.Attrs["entityIds"].([]any); ok {
			count = len(v)
		}
		w.Writef("> *[relationship graph: %d entries]*", count)
	case "secretVault":
		w.Write("> 🔒 *[sealed section — visible only to the author]*")
	default:
		// Unknown block: just render inline text inside.
		if len(n.Content) > 0 {
			renderInline(n, w)
		}
	}
}

func renderListItem(item *Node, w *mdWriter, depth int, bullet string, order int) {
	prefix := ""
	if bullet != "" {
		prefix = bullet + " "
	} else {
		prefix = fmt.Sprintf("%d. ", order)
	}
	indent := strings.Repeat("  ", depth)
	var inner mdWriter
	for i := range item.Content {
		renderBlock(&item.Content[i], &inner, depth+1)
		if i < len(item.Content)-1 {
			inner.BlockBreak()
		}
	}
	lines := strings.Split(strings.TrimRight(inner.b.String(), "\n"), "\n")
	for i, line := range lines {
		if i == 0 {
			w.Write(indent + prefix + line + "\n")
		} else {
			pad := strings.Repeat(" ", len(prefix))
			w.Write(indent + pad + line + "\n")
		}
	}
}

func renderTable(n *Node, w *mdWriter) {
	rows := n.Content
	if len(rows) == 0 {
		return
	}
	// determine column count from first row
	headers := []string{}
	dataRows := [][]string{}
	for ri, row := range rows {
		cells := []string{}
		for _, cell := range row.Content {
			cells = append(cells, strings.TrimSpace(plainText(&cell)))
		}
		if ri == 0 {
			// if the first row's cells are all tableHeader, use as header
			allHeader := true
			for _, c := range row.Content {
				if c.Type != "tableHeader" {
					allHeader = false
					break
				}
			}
			if allHeader {
				headers = cells
				continue
			}
		}
		dataRows = append(dataRows, cells)
	}
	if len(headers) == 0 && len(dataRows) > 0 {
		// fabricate a blank header row so the markdown table is valid
		headers = make([]string, len(dataRows[0]))
		for i := range headers {
			headers[i] = fmt.Sprintf("col %d", i+1)
		}
	}
	if len(headers) == 0 {
		return
	}
	w.Write("| " + strings.Join(escapeCells(headers), " | ") + " |\n")
	w.Write("|" + strings.Repeat(" --- |", len(headers)) + "\n")
	for _, r := range dataRows {
		// pad short rows
		if len(r) < len(headers) {
			r = append(r, make([]string, len(headers)-len(r))...)
		}
		if len(r) > len(headers) {
			r = r[:len(headers)]
		}
		w.Write("| " + strings.Join(escapeCells(r), " | ") + " |\n")
	}
}

func escapeCells(cells []string) []string {
	out := make([]string, len(cells))
	for i, c := range cells {
		c = strings.ReplaceAll(c, "|", "\\|")
		c = strings.ReplaceAll(c, "\n", " ")
		out[i] = c
	}
	return out
}

func renderInline(n *Node, w *mdWriter) {
	for i := range n.Content {
		child := &n.Content[i]
		switch child.Type {
		case "text":
			marks := extractMarks(child)
			text := escapeMd(child.Text)
			text = wrapMarks(text, marks)
			w.Write(text)
		case "hardBreak":
			w.Write("  \n")
		case "wikilink":
			slug, _ := child.Attrs["slug"].(string)
			label, _ := child.Attrs["label"].(string)
			if label == "" {
				label = slug
			}
			w.Writef("[%s](./%s)", escapeMd(label), slug)
		case "wwimage":
			alt, _ := child.Attrs["alt"].(string)
			src, _ := child.Attrs["src"].(string)
			w.Writef("![%s](%s)", escapeMd(alt), src)
		default:
			renderInline(child, w)
		}
	}
}

func extractMarks(n *Node) []string {
	if n.Attrs == nil {
		return nil
	}
	raw, ok := n.Attrs["marks"]
	if !ok {
		return nil
	}
	arr, ok := raw.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(arr))
	for _, m := range arr {
		if obj, ok := m.(map[string]any); ok {
			if t, ok := obj["type"].(string); ok {
				out = append(out, t)
			}
		}
	}
	return out
}

func wrapMarks(text string, marks []string) string {
	for _, m := range marks {
		switch m {
		case "bold", "strong":
			text = "**" + text + "**"
		case "italic", "em":
			text = "*" + text + "*"
		case "strike", "strikethrough":
			text = "~~" + text + "~~"
		case "code":
			text = "`" + text + "`"
		}
	}
	return text
}

func escapeMd(s string) string {
	r := strings.NewReplacer(
		"\\", "\\\\",
		"`", "\\`",
		"*", "\\*",
		"_", "\\_",
		"[", "\\[",
		"]", "\\]",
		"<", "\\<",
		">", "\\>",
	)
	return r.Replace(s)
}

func plainText(n *Node) string {
	var sb strings.Builder
	var walk func(*Node)
	walk = func(x *Node) {
		if x == nil {
			return
		}
		if x.Type == "text" {
			sb.WriteString(x.Text)
		}
		if x.Type == "hardBreak" {
			sb.WriteByte('\n')
		}
		for i := range x.Content {
			walk(&x.Content[i])
		}
	}
	walk(n)
	return sb.String()
}
