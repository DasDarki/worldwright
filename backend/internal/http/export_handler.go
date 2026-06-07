package http

import (
	"archive/zip"
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"worldwright/backend/internal/auth"
	"worldwright/backend/internal/content"
	"worldwright/backend/internal/store"
)

// exportEntityMarkdown renders a single entity as a markdown file the
// caller can download. Visibility is enforced the same way as the
// view endpoint, so anonymous callers only get public entries.
func exportEntityMarkdown(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		slug := c.Params("slug")
		e, err := st.EntityBySlug(c.UserContext(), slug, auth.VisibilityFor(c))
		if err != nil {
			return err
		}
		md := renderEntityMarkdown(e)
		c.Set(fiber.HeaderContentType, "text/markdown; charset=utf-8")
		c.Set(fiber.HeaderContentDisposition,
			fmt.Sprintf(`attachment; filename="%s.md"`, sanitizeFilename(e.Slug)))
		return c.SendString(md)
	}
}

// exportWikiZip streams a ZIP archive of every visible entity rendered as
// markdown. Files are grouped by entity-type folder so the tree is
// browsable in any text editor.
func exportWikiZip(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		visibility := auth.VisibilityFor(c)

		summaries, err := st.ListEntities(c.UserContext(), store.EntityFilter{Visibility: visibility})
		if err != nil {
			return err
		}
		types, err := st.ListEntityTypes(c.UserContext())
		if err != nil {
			return err
		}
		typeFolder := map[int64]string{}
		for _, t := range types {
			typeFolder[t.ID] = sanitizeFilename(t.Key)
		}

		var buf bytes.Buffer
		zw := zip.NewWriter(&buf)

		index := strings.Builder{}
		index.WriteString("# Worldwright export\n\n")
		index.WriteString(fmt.Sprintf("Exported %s · %d entries\n\n", time.Now().UTC().Format(time.RFC3339), len(summaries)))

		for _, s := range summaries {
			e, err := st.EntityByID(c.UserContext(), s.ID, visibility)
			if err != nil {
				continue
			}
			folder := typeFolder[e.EntityTypeID]
			if folder == "" {
				folder = "_misc"
			}
			path := folder + "/" + sanitizeFilename(e.Slug) + ".md"

			f, err := zw.Create(path)
			if err != nil {
				return err
			}
			if _, err := f.Write([]byte(renderEntityMarkdown(e))); err != nil {
				return err
			}
			index.WriteString(fmt.Sprintf("- [%s](%s)\n", e.Title, path))
		}

		idxF, err := zw.Create("INDEX.md")
		if err != nil {
			return err
		}
		if _, err := idxF.Write([]byte(index.String())); err != nil {
			return err
		}

		if err := zw.Close(); err != nil {
			return err
		}

		c.Set(fiber.HeaderContentType, "application/zip")
		c.Set(fiber.HeaderContentDisposition,
			fmt.Sprintf(`attachment; filename="worldwright-export-%s.zip"`,
				time.Now().UTC().Format("20060102-150405")))
		return c.Send(buf.Bytes())
	}
}

func renderEntityMarkdown(e *store.Entity) string {
	var sb strings.Builder
	sb.WriteString("# ")
	sb.WriteString(e.Title)
	sb.WriteString("\n\n")
	if e.Summary != "" {
		sb.WriteString("> ")
		sb.WriteString(e.Summary)
		sb.WriteString("\n\n")
	}
	if len(e.Tags) > 0 {
		sb.WriteString("**Tags:** ")
		for i, t := range e.Tags {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString("`")
			sb.WriteString(t)
			sb.WriteString("`")
		}
		sb.WriteString("\n\n")
	}
	if len(e.FieldValues) > 0 {
		sb.WriteString("| Field | Value |\n| --- | --- |\n")
		for k, v := range e.FieldValues {
			if v == "" {
				continue
			}
			sb.WriteString("| ")
			sb.WriteString(strings.ReplaceAll(k, "|", `\|`))
			sb.WriteString(" | ")
			sb.WriteString(strings.ReplaceAll(v, "|", `\|`))
			sb.WriteString(" |\n")
		}
		sb.WriteString("\n")
	}
	sb.WriteString(content.ToMarkdown(e.Body))
	return sb.String()
}

func sanitizeFilename(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "untitled"
	}
	r := strings.NewReplacer(
		"/", "-",
		"\\", "-",
		":", "-",
		"*", "-",
		"?", "-",
		"\"", "-",
		"<", "-",
		">", "-",
		"|", "-",
	)
	return r.Replace(s)
}
