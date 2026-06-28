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

// ---------- Single entity --------------------------------------------------

// exportEntityMarkdown renders a single entity as a markdown file the
// caller can download. Visibility is enforced the same way as the view
// endpoint, so anonymous callers only get public entries. Secret vaults
// the caller authored are kept; everyone else's are sealed.
func exportEntityMarkdown(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		slug := c.Params("slug")
		e, err := st.EntityBySlug(c.UserContext(), slug, auth.VisibilityFor(c))
		if err != nil {
			return err
		}
		e.Body = content.ScrubSecretVaults(e.Body, viewerID(c))
		md := renderEntityMarkdown(e)
		c.Set(fiber.HeaderContentType, "text/markdown; charset=utf-8")
		c.Set(fiber.HeaderContentDisposition,
			fmt.Sprintf(`attachment; filename="%s.md"`, sanitizeFilename(e.Slug)))
		return c.SendString(md)
	}
}

// ---------- Wiki-wide exports ---------------------------------------------

// wikiExporter consolidates the four wiki-wide export modes (single .md,
// single .pdf, .md.zip, .pdf.zip) on top of the same content pipeline:
// collect visible entities → render markdown per entity → optionally
// concat / zip / convert to PDF.

func exportWikiMarkdown(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		md, _, err := buildWikiMarkdown(c, st)
		if err != nil {
			return err
		}
		c.Set(fiber.HeaderContentType, "text/markdown; charset=utf-8")
		c.Set(fiber.HeaderContentDisposition,
			fmt.Sprintf(`attachment; filename="worldwright-%s.md"`,
				time.Now().UTC().Format("20060102-150405")))
		return c.SendString(md)
	}
}

func exportWikiMarkdownZip(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		entries, types, err := collectWikiEntries(c, st)
		if err != nil {
			return err
		}
		typeFolder := folderForType(types)

		var buf bytes.Buffer
		zw := zip.NewWriter(&buf)

		index := strings.Builder{}
		index.WriteString("# Worldwright export\n\n")
		index.WriteString(fmt.Sprintf("Exported %s · %d entries\n\n",
			time.Now().UTC().Format(time.RFC3339), len(entries)))

		for _, e := range entries {
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

		idx, err := zw.Create("INDEX.md")
		if err != nil {
			return err
		}
		if _, err := idx.Write([]byte(index.String())); err != nil {
			return err
		}
		if err := zw.Close(); err != nil {
			return err
		}

		c.Set(fiber.HeaderContentType, "application/zip")
		c.Set(fiber.HeaderContentDisposition,
			fmt.Sprintf(`attachment; filename="worldwright-%s-md.zip"`,
				time.Now().UTC().Format("20060102-150405")))
		return c.Send(buf.Bytes())
	}
}

func exportWikiPDF(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		md, _, err := buildWikiMarkdown(c, st)
		if err != nil {
			return err
		}
		pdfBytes, err := content.MarkdownToPDF(md)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "pdf render failed: "+err.Error())
		}
		c.Set(fiber.HeaderContentType, "application/pdf")
		c.Set(fiber.HeaderContentDisposition,
			fmt.Sprintf(`attachment; filename="worldwright-%s.pdf"`,
				time.Now().UTC().Format("20060102-150405")))
		return c.Send(pdfBytes)
	}
}

func exportWikiPDFZip(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		entries, types, err := collectWikiEntries(c, st)
		if err != nil {
			return err
		}
		typeFolder := folderForType(types)

		var buf bytes.Buffer
		zw := zip.NewWriter(&buf)

		index := strings.Builder{}
		index.WriteString("Worldwright export — PDF archive\n")
		index.WriteString(fmt.Sprintf("Exported %s · %d entries\n\n",
			time.Now().UTC().Format(time.RFC3339), len(entries)))

		for _, e := range entries {
			folder := typeFolder[e.EntityTypeID]
			if folder == "" {
				folder = "_misc"
			}
			path := folder + "/" + sanitizeFilename(e.Slug) + ".pdf"

			pdfBytes, perr := content.MarkdownToPDF(renderEntityMarkdown(e))
			if perr != nil {
				// Skip the bad one but keep going — don't tank a whole export
				// because one entry couldn't render.
				index.WriteString(fmt.Sprintf("- [SKIPPED] %s (%s): %s\n", e.Title, e.Slug, perr.Error()))
				continue
			}
			f, err := zw.Create(path)
			if err != nil {
				return err
			}
			if _, err := f.Write(pdfBytes); err != nil {
				return err
			}
			index.WriteString(fmt.Sprintf("- %s → %s\n", e.Title, path))
		}

		idx, err := zw.Create("INDEX.txt")
		if err != nil {
			return err
		}
		if _, err := idx.Write([]byte(index.String())); err != nil {
			return err
		}
		if err := zw.Close(); err != nil {
			return err
		}

		c.Set(fiber.HeaderContentType, "application/zip")
		c.Set(fiber.HeaderContentDisposition,
			fmt.Sprintf(`attachment; filename="worldwright-%s-pdf.zip"`,
				time.Now().UTC().Format("20060102-150405")))
		return c.Send(buf.Bytes())
	}
}

// ---------- Shared helpers -------------------------------------------------

// collectWikiEntries pulls every entity the caller can see, loads it with
// its body / tags / fields, and applies secret-vault scrubbing keyed to the
// current user so they get to see their own sealed content but no-one
// else's.
func collectWikiEntries(c *fiber.Ctx, st *store.Store) ([]*store.Entity, []store.EntityType, error) {
	visibility := auth.VisibilityFor(c)
	summaries, err := st.ListEntities(c.UserContext(), store.EntityFilter{Visibility: visibility})
	if err != nil {
		return nil, nil, err
	}
	types, err := st.ListEntityTypes(c.UserContext())
	if err != nil {
		return nil, nil, err
	}
	uid := viewerID(c)
	out := make([]*store.Entity, 0, len(summaries))
	for _, s := range summaries {
		e, err := st.EntityByID(c.UserContext(), s.ID, visibility)
		if err != nil {
			continue
		}
		e.Body = content.ScrubSecretVaults(e.Body, uid)
		out = append(out, e)
	}
	return out, types, nil
}

func buildWikiMarkdown(c *fiber.Ctx, st *store.Store) (string, []*store.Entity, error) {
	entries, _, err := collectWikiEntries(c, st)
	if err != nil {
		return "", nil, err
	}
	var sb strings.Builder
	sb.WriteString("# Worldwright\n\n")
	sb.WriteString(fmt.Sprintf("*Exported %s · %d entries*\n\n",
		time.Now().UTC().Format(time.RFC3339), len(entries)))
	for i, e := range entries {
		if i > 0 {
			// real page break for the PDF pipeline + visual divider in MD
			sb.WriteString("\n\n---\n\n")
		}
		sb.WriteString(renderEntityMarkdown(e))
		sb.WriteString("\n")
	}
	return sb.String(), entries, nil
}

func folderForType(types []store.EntityType) map[int64]string {
	out := map[int64]string{}
	for _, t := range types {
		out[t.ID] = sanitizeFilename(t.Key)
	}
	return out
}

func viewerID(c *fiber.Ctx) int64 {
	if u := auth.UserFrom(c); u != nil {
		return u.ID
	}
	return 0
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
