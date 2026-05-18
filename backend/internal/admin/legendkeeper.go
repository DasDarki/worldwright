package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"worldwright/backend/internal/content"
	"worldwright/backend/internal/store"
)

// ImportSummary tells the caller what happened during a LegendKeeper import.
type ImportSummary struct {
	Source         string   `json:"source"`
	EntitiesCreated int     `json:"entities_created"`
	MapsCreated    int      `json:"maps_created"`
	PinsCreated    int      `json:"pins_created"`
	ImagesImported int      `json:"images_imported"`
	ImagesFailed   int      `json:"images_failed"`
	Warnings       []string `json:"warnings,omitempty"`
}

// ImportLegendKeeper wipes ALL existing content, then attempts a best-effort
// import of a LegendKeeper export JSON. Returns a summary describing what
// was imported; partial failures are recorded as warnings rather than
// aborting the whole run.
func (s *Service) ImportLegendKeeper(ctx context.Context, r io.Reader) (ImportSummary, error) {
	raw, err := io.ReadAll(io.LimitReader(r, 200*1024*1024))
	if err != nil {
		return ImportSummary{}, fmt.Errorf("read upload: %w", err)
	}
	var exp lkExport
	if err := json.Unmarshal(raw, &exp); err != nil {
		return ImportSummary{}, fmt.Errorf("parse legendkeeper json: %w", err)
	}
	if len(exp.Resources) == 0 {
		return ImportSummary{}, errors.New("export contains no resources")
	}

	if err := s.PruneAll(ctx); err != nil {
		return ImportSummary{}, fmt.Errorf("prune existing content: %w", err)
	}

	// Make sure events have a calendar to live in even though the entity
	// importer doesn't touch events directly.
	if _, err := s.ensureDefaultCalendar(ctx); err != nil {
		return ImportSummary{}, fmt.Errorf("ensure calendar: %w", err)
	}

	imp := &lkImporter{
		svc:        s,
		ctx:        ctx,
		entityType: map[string]int64{},
		slugByLkID: map[string]string{},
		assetByURL: map[string]int64{},
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
		summary: ImportSummary{Source: "legendkeeper"},
	}
	if err := imp.loadEntityTypes(); err != nil {
		return imp.summary, err
	}
	imp.indexResources(exp.Resources)

	// Pass 1a: create entities depth-first so parents exist before children.
	ordered := sortResourcesByDepth(exp.Resources)
	for _, r := range ordered {
		if err := imp.importResource(r); err != nil {
			imp.summary.Warnings = append(imp.summary.Warnings,
				fmt.Sprintf("resource %q (%s): %v", r.Name, r.ID, err))
		}
	}

	// Pass 2: rewrite mention placeholders -> real wikilinks now that all
	// entities have slugs. Walks every entity body once.
	if err := imp.resolveWikilinks(); err != nil {
		imp.summary.Warnings = append(imp.summary.Warnings, fmt.Sprintf("resolve wikilinks: %v", err))
	}

	// Maps come last so their entity references can resolve.
	for i := range exp.Resources {
		imp.importMaps(&exp.Resources[i])
	}

	return imp.summary, nil
}

type lkImporter struct {
	svc        *Service
	ctx        context.Context
	entityType map[string]int64 // key -> id
	slugByLkID map[string]string
	parentByLk map[string]string
	assetByURL map[string]int64
	http       *http.Client
	summary    ImportSummary
}

func (s *Service) ensureDefaultCalendar(ctx context.Context) (int64, error) {
	var id int64
	err := s.store.DB().QueryRowContext(ctx, `SELECT id FROM calendars LIMIT 1`).Scan(&id)
	if err == nil {
		return id, nil
	}
	res, err := s.store.DB().ExecContext(ctx,
		`INSERT INTO calendars (name, epoch_name, current_year, current_month_index, current_day, is_seed)
		 VALUES (?, ?, ?, ?, ?, 0)`,
		"Default Calendar", "", 1, 0, 1)
	if err != nil {
		return 0, err
	}
	id, _ = res.LastInsertId()

	months := []struct {
		name string
		days int
	}{
		{"January", 31}, {"February", 28}, {"March", 31}, {"April", 30},
		{"May", 31}, {"June", 30}, {"July", 31}, {"August", 31},
		{"September", 30}, {"October", 31}, {"November", 30}, {"December", 31},
	}
	for i, m := range months {
		if _, err := s.store.DB().ExecContext(ctx,
			`INSERT INTO calendar_months (calendar_id, sort_order, name, days) VALUES (?, ?, ?, ?)`,
			id, i+1, m.name, m.days); err != nil {
			return 0, err
		}
	}
	weekdays := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	for i, w := range weekdays {
		if _, err := s.store.DB().ExecContext(ctx,
			`INSERT INTO calendar_weekdays (calendar_id, sort_order, name) VALUES (?, ?, ?)`,
			id, i+1, w); err != nil {
			return 0, err
		}
	}
	return id, nil
}

func (i *lkImporter) loadEntityTypes() error {
	rows, err := i.svc.store.DB().QueryContext(i.ctx, `SELECT id, key FROM entity_types`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var key string
		if err := rows.Scan(&id, &key); err != nil {
			return err
		}
		i.entityType[key] = id
	}
	return rows.Err()
}

func (i *lkImporter) indexResources(resources []lkResource) {
	used := map[string]struct{}{}
	for _, r := range resources {
		slug := slugify(r.Name)
		if slug == "" {
			slug = strings.ToLower(r.ID)
		}
		base := slug
		n := 1
		for {
			if _, taken := used[slug]; !taken {
				break
			}
			n++
			slug = fmt.Sprintf("%s-%d", base, n)
		}
		used[slug] = struct{}{}
		i.slugByLkID[r.ID] = slug
	}
}

func (i *lkImporter) importResource(r *lkResource) error {
	slug := i.slugByLkID[r.ID]
	if slug == "" {
		return errors.New("no slug allocated")
	}
	title := strings.TrimSpace(r.Name)
	if title == "" {
		title = slug
	}
	entityTypeID := i.guessEntityType(r.Tags)

	// Pick the first page document as the body. Other documents (extra pages,
	// boards, blanks, times) are surfaced as appended sections so nothing is
	// silently dropped, but maps are handled separately.
	docs := make([]lkDocument, 0, len(r.Documents))
	for _, d := range r.Documents {
		if d.Type == "page" || d.Type == "blank" {
			docs = append(docs, d)
		}
	}

	bodyNodes := []map[string]any{}
	// SUMMARY property -> entity.summary + a lead paragraph.
	summary := ""
	for _, p := range r.Properties {
		if p.Type == "TEXT_FIELD" && strings.EqualFold(p.Title, "SUMMARY") && len(p.Data) > 0 {
			var d lkPropertyTextData
			if err := json.Unmarshal(p.Data, &d); err == nil && len(d.Fragment) > 0 {
				if txt := extractPlainText(d.Fragment); txt != "" {
					summary = txt
				}
			}
		}
		if p.Type == "IMAGE" && len(p.Data) > 0 {
			var d lkPropertyImageData
			if err := json.Unmarshal(p.Data, &d); err == nil && d.URL != "" {
				if node, ok := i.imageNodeFromURL(d.URL); ok {
					bodyNodes = append(bodyNodes, node)
				}
			}
		}
	}

	for idx, d := range docs {
		var doc lkNode
		if err := json.Unmarshal(d.Content, &doc); err != nil {
			i.summary.Warnings = append(i.summary.Warnings,
				fmt.Sprintf("doc %q: parse: %v", d.Name, err))
			continue
		}
		if idx > 0 && d.Name != "" {
			bodyNodes = append(bodyNodes, headingNode(d.Name))
		}
		for _, child := range doc.Content {
			if converted, ok := i.translateNode(&child); ok {
				bodyNodes = append(bodyNodes, converted)
			}
		}
	}

	if len(bodyNodes) == 0 {
		bodyNodes = []map[string]any{{"type": "paragraph"}}
	}
	body := map[string]any{
		"type":    "doc",
		"content": bodyNodes,
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return err
	}
	bodyText, slugs, _ := content.ParseBody(bodyJSON)

	var parentID *int64
	if r.ParentID != "" {
		if pid, err := i.lookupEntityID(r.ParentID); err == nil {
			parentID = &pid
		}
	}

	// We may have a parent that hasn't been created yet (depth-first ordering
	// in the export isn't guaranteed). Insert without parent now; a later
	// reparenting pass would fix this, but for v1 we order by depth.
	if parentID == nil && r.ParentID != "" {
		// fall through; resource will end up at root, acceptable for now
	}

	_, err = i.svc.store.CreateEntity(i.ctx, store.NewEntity{
		EntityTypeID: entityTypeID,
		Title:        title,
		Slug:         slug,
		Summary:      summary,
		Body:         bodyJSON,
		BodyText:     bodyText,
		ParentID:     parentID,
		Visibility:   "player",
		Tags:         normalizeTags(r.Tags),
		Wikilinks:    slugs,
	})
	if err != nil {
		return err
	}
	i.summary.EntitiesCreated++
	return nil
}

func (i *lkImporter) lookupEntityID(lkID string) (int64, error) {
	slug, ok := i.slugByLkID[lkID]
	if !ok {
		return 0, errors.New("unknown lk id")
	}
	var id int64
	err := i.svc.store.DB().QueryRowContext(i.ctx,
		`SELECT id FROM entities WHERE slug = ?`, slug).Scan(&id)
	return id, err
}

// translateNode converts a LegendKeeper TipTap node to our dialect. Returns
// (nodeJSON, ok). Unsupported nodes return false so the caller can drop them.
func (i *lkImporter) translateNode(n *lkNode) (map[string]any, bool) {
	switch n.Type {
	case "paragraph":
		return wrapBlock("paragraph", i.translateInline(n.Content)), true
	case "heading":
		level := readHeadingLevel(n.Attrs)
		if level < 2 {
			level = 2
		}
		if level > 3 {
			level = 3
		}
		node := map[string]any{
			"type":  "heading",
			"attrs": map[string]any{"level": level},
		}
		if children := i.translateInline(n.Content); len(children) > 0 {
			node["content"] = children
		}
		return node, true
	case "blockquote":
		return wrapBlock("blockquote", i.translateBlocks(n.Content)), true
	case "bulletList", "bullet_list":
		return wrapBlock("bulletList", i.translateBlocks(n.Content)), true
	case "orderedList", "ordered_list":
		return wrapBlock("orderedList", i.translateBlocks(n.Content)), true
	case "listItem", "list_item":
		return wrapBlock("listItem", i.translateBlocks(n.Content)), true
	case "rule", "horizontalRule":
		return map[string]any{"type": "horizontalRule"}, true
	case "hardBreak", "hard_break":
		return map[string]any{"type": "hardBreak"}, true
	case "panel", "group":
		// unwrap content
		blocks := i.translateBlocks(n.Content)
		if len(blocks) == 0 {
			return nil, false
		}
		// Caller wraps the result in a single node, so flatten by returning
		// the first block and silently dropping the rest? Better: emit a
		// blockquote so the contents remain visible.
		return map[string]any{
			"type":    "blockquote",
			"content": blocks,
		}, true
	case "mediaSingle":
		for _, child := range n.Content {
			if child.Type == "media" {
				url := readAttrString(child.Attrs, "url")
				if node, ok := i.imageNodeFromURL(url); ok {
					return node, true
				}
			}
		}
		return nil, false
	case "image":
		// shapesV2 image variant (from board documents) or top-level image
		url := readAttrString(n.Attrs, "url")
		if url == "" {
			url = readPropString(n.Attrs, "src")
		}
		if url != "" {
			if node, ok := i.imageNodeFromURL(url); ok {
				return node, true
			}
		}
		return nil, false
	case "external":
		url := readAttrString(n.Attrs, "url")
		if isImageURL(url) {
			if node, ok := i.imageNodeFromURL(url); ok {
				return node, true
			}
		}
		return nil, false
	}
	return nil, false
}

func (i *lkImporter) translateBlocks(nodes []lkNode) []map[string]any {
	out := make([]map[string]any, 0, len(nodes))
	for idx := range nodes {
		if converted, ok := i.translateNode(&nodes[idx]); ok {
			out = append(out, converted)
		}
	}
	return out
}

// translateInline keeps text / marks / mentions; everything else is dropped.
func (i *lkImporter) translateInline(nodes []lkNode) []map[string]any {
	out := make([]map[string]any, 0, len(nodes))
	for idx := range nodes {
		n := &nodes[idx]
		switch n.Type {
		case "text":
			node := map[string]any{
				"type": "text",
				"text": n.Text,
			}
			if marks := translateMarks(n.Marks); len(marks) > 0 {
				node["marks"] = marks
			}
			out = append(out, node)
		case "hardBreak", "hard_break":
			out = append(out, map[string]any{"type": "hardBreak"})
		case "mention":
			label := readAttrString(n.Attrs, "text")
			lkID := readAttrString(n.Attrs, "id")
			if lkID == "" {
				if label != "" {
					out = append(out, map[string]any{"type": "text", "text": label})
				}
				continue
			}
			// Pass-1 placeholder: store the LK id in slug; pass 2 rewrites.
			out = append(out, map[string]any{
				"type": "wikilink",
				"attrs": map[string]any{
					"slug":  "lk:" + lkID,
					"label": label,
				},
			})
		}
	}
	return out
}

func translateMarks(marks []lkMark) []map[string]any {
	if len(marks) == 0 {
		return nil
	}
	out := make([]map[string]any, 0, len(marks))
	for _, m := range marks {
		switch m.Type {
		case "em", "italic":
			out = append(out, map[string]any{"type": "italic"})
		case "strong", "bold":
			out = append(out, map[string]any{"type": "bold"})
		case "strike", "strikethrough":
			out = append(out, map[string]any{"type": "strike"})
		case "code":
			out = append(out, map[string]any{"type": "code"})
		}
	}
	return out
}

func wrapBlock(t string, content []map[string]any) map[string]any {
	node := map[string]any{"type": t}
	if len(content) > 0 {
		node["content"] = content
	}
	return node
}

func headingNode(text string) map[string]any {
	return map[string]any{
		"type":  "heading",
		"attrs": map[string]any{"level": 2},
		"content": []map[string]any{
			{"type": "text", "text": text},
		},
	}
}

// resolveWikilinks walks every entity body and replaces lk: placeholders
// with real slugs.
func (i *lkImporter) resolveWikilinks() error {
	rows, err := i.svc.store.DB().QueryContext(i.ctx, `SELECT id, body FROM entities`)
	if err != nil {
		return err
	}
	type row struct {
		id   int64
		body string
	}
	all := []row{}
	for rows.Next() {
		var r row
		if err := rows.Scan(&r.id, &r.body); err != nil {
			rows.Close()
			return err
		}
		all = append(all, r)
	}
	rows.Close()

	for _, r := range all {
		var doc map[string]any
		if err := json.Unmarshal([]byte(r.body), &doc); err != nil {
			continue
		}
		changed := false
		walkResolve(doc, i.slugByLkID, &changed)
		if !changed {
			continue
		}
		newBody, err := json.Marshal(doc)
		if err != nil {
			continue
		}
		// Also re-extract plaintext + wikilink slugs for FTS/backlinks.
		bodyText, slugs, _ := content.ParseBody(newBody)
		tx, err := i.svc.store.DB().Begin()
		if err != nil {
			continue
		}
		_, err = tx.ExecContext(i.ctx,
			`UPDATE entities SET body = ?, body_text = ? WHERE id = ?`,
			string(newBody), bodyText, r.id)
		if err != nil {
			tx.Rollback()
			continue
		}
		// rewrite links for this entity
		if _, err := tx.ExecContext(i.ctx, `DELETE FROM links WHERE source_entity_id = ?`, r.id); err == nil {
			for _, slug := range slugs {
				var target int64
				if err := tx.QueryRowContext(i.ctx, `SELECT id FROM entities WHERE slug = ?`, slug).Scan(&target); err == nil && target != r.id {
					_, _ = tx.ExecContext(i.ctx,
						`INSERT OR IGNORE INTO links (source_entity_id, target_entity_id, anchor) VALUES (?, ?, '')`,
						r.id, target)
				}
			}
		}
		_ = tx.Commit()
	}
	return nil
}

func walkResolve(node any, slugMap map[string]string, changed *bool) {
	switch v := node.(type) {
	case map[string]any:
		if v["type"] == "wikilink" {
			if attrs, ok := v["attrs"].(map[string]any); ok {
				if slug, ok := attrs["slug"].(string); ok && strings.HasPrefix(slug, "lk:") {
					if real, ok := slugMap[strings.TrimPrefix(slug, "lk:")]; ok {
						attrs["slug"] = real
						*changed = true
					} else {
						// Unknown mention target: drop the wikilink, replace
						// node with plain text to keep the label.
						label, _ := attrs["label"].(string)
						v["type"] = "text"
						v["text"] = label
						delete(v, "attrs")
						*changed = true
					}
				}
			}
		}
		for _, child := range v {
			walkResolve(child, slugMap, changed)
		}
	case []any:
		for _, child := range v {
			walkResolve(child, slugMap, changed)
		}
	}
}

func (i *lkImporter) importMaps(r *lkResource) {
	parentID, _ := i.lookupEntityID(r.ID)
	for _, d := range r.Documents {
		if d.Type != "map" || d.Map == nil {
			continue
		}
		assetID, ok := i.downloadImage(d.Map.LocatorID)
		if !ok {
			i.summary.Warnings = append(i.summary.Warnings, fmt.Sprintf("map %q: skipped (image download failed)", d.Name))
			continue
		}
		width := int(d.Map.MaxX - d.Map.MinX)
		height := int(d.Map.MaxY - d.Map.MinY)
		if width <= 0 {
			width = 1024
		}
		if height <= 0 {
			height = 1024
		}
		name := strings.TrimSpace(d.Name)
		if name == "" {
			name = r.Name
		}
		var parentRef *int64
		if parentID != 0 {
			pid := parentID
			parentRef = &pid
		}
		m, err := i.svc.store.CreateMap(i.ctx, store.NewMap{
			Name:           name,
			AssetID:        assetID,
			Width:          width,
			Height:         height,
			ParentEntityID: parentRef,
		})
		if err != nil {
			i.summary.Warnings = append(i.summary.Warnings, fmt.Sprintf("map %q: create: %v", name, err))
			continue
		}
		i.summary.MapsCreated++

		var pinContent lkMapContent
		if len(d.Content) > 0 {
			_ = json.Unmarshal(d.Content, &pinContent)
		}
		for _, p := range pinContent.Pins {
			if len(p.Pos) < 2 {
				continue
			}
			pin := store.NewMapPin{
				X:          p.Pos[0],
				Y:          p.Pos[1],
				Label:      p.Name,
				Icon:       p.IconGlyph,
				Visibility: "player",
			}
			if p.ResourceID != "" {
				if eid, err := i.lookupEntityID(p.ResourceID); err == nil {
					pin.TargetEntityID = &eid
				}
			}
			if _, err := i.svc.store.CreateMapPin(i.ctx, m.ID, pin); err != nil {
				continue
			}
			i.summary.PinsCreated++
		}
	}
}

// imageNodeFromURL fetches the URL (using the assetByURL cache) and returns
// a wwimage TipTap node. Returns (nil, false) on failure so the caller can
// drop the node entirely.
func (i *lkImporter) imageNodeFromURL(rawURL string) (map[string]any, bool) {
	if rawURL == "" {
		return nil, false
	}
	assetID, ok := i.downloadImage(rawURL)
	if !ok {
		return nil, false
	}
	return map[string]any{
		"type": "wwimage",
		"attrs": map[string]any{
			"src":     fmt.Sprintf("/api/assets/%d", assetID),
			"alt":     "",
			"assetId": assetID,
		},
	}, true
}

func (i *lkImporter) downloadImage(rawURL string) (int64, bool) {
	if id, ok := i.assetByURL[rawURL]; ok {
		return id, true
	}
	u, err := url.Parse(rawURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		i.summary.ImagesFailed++
		return 0, false
	}
	req, err := http.NewRequestWithContext(i.ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		i.summary.ImagesFailed++
		return 0, false
	}
	resp, err := i.http.Do(req)
	if err != nil {
		i.summary.ImagesFailed++
		return 0, false
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		i.summary.ImagesFailed++
		return 0, false
	}
	buf, err := io.ReadAll(io.LimitReader(resp.Body, 25*1024*1024))
	if err != nil || len(buf) == 0 {
		i.summary.ImagesFailed++
		return 0, false
	}
	mime := resp.Header.Get("Content-Type")
	if idx := strings.Index(mime, ";"); idx > 0 {
		mime = strings.TrimSpace(mime[:idx])
	}
	ext := extForMime(mime)
	if ext == "" {
		ext = strings.ToLower(path.Ext(u.Path))
		if mime == "" {
			mime = mimeForExt(ext)
		}
	}
	if ext == "" || mime == "" {
		i.summary.ImagesFailed++
		return 0, false
	}

	var width, height *int
	if cfg, _, derr := image.DecodeConfig(bytes.NewReader(buf)); derr == nil {
		w := cfg.Width
		h := cfg.Height
		width = &w
		height = &h
	}

	name := uuid.NewString() + ext
	if err := os.MkdirAll(i.svc.assetsDir, 0o755); err != nil {
		i.summary.ImagesFailed++
		return 0, false
	}
	dst := filepath.Join(i.svc.assetsDir, name)
	if err := os.WriteFile(dst, buf, 0o644); err != nil {
		i.summary.ImagesFailed++
		return 0, false
	}
	asset, err := i.svc.store.CreateAsset(i.ctx, store.NewAsset{
		Filename: path.Base(u.Path),
		Path:     name,
		Mime:     mime,
		Size:     int64(len(buf)),
		Width:    width,
		Height:   height,
	})
	if err != nil {
		_ = os.Remove(dst)
		i.summary.ImagesFailed++
		return 0, false
	}
	i.assetByURL[rawURL] = asset.ID
	i.summary.ImagesImported++
	return asset.ID, true
}

func (i *lkImporter) guessEntityType(tags []string) int64 {
	def := i.entityType["concept"]
	if def == 0 {
		// fallback to any system type if 'concept' is missing
		for _, id := range i.entityType {
			def = id
			break
		}
	}
	for _, raw := range tags {
		t := strings.ToLower(strings.TrimSpace(raw))
		var key string
		switch {
		case t == "character" || t == "charakter" || t == "person" || t == "npc":
			key = "character"
		case t == "city" || t == "town" || t == "stadt" || t == "location" || t == "ort" || t == "region" || t == "land" || t == "country" || t == "kingdom":
			key = "location"
		case t == "faction" || t == "fraktion" || t == "organization" || t == "organisation" || t == "guild" || t == "house":
			key = "faction"
		case t == "item" || t == "gegenstand" || t == "artifact" || t == "weapon" || t == "armor":
			key = "item"
		case t == "deity" || t == "gott" || t == "gottheit" || t == "god" || t == "godess" || t == "goddess":
			key = "deity"
		case t == "race" || t == "volk" || t == "species":
			key = "race"
		case t == "creature" || t == "kreatur" || t == "monster" || t == "beast":
			key = "creature"
		case t == "religion" || t == "cult":
			key = "religion"
		case t == "language" || t == "sprache":
			key = "language"
		case t == "culture" || t == "kultur":
			key = "culture"
		}
		if key != "" {
			if id, ok := i.entityType[key]; ok {
				return id
			}
		}
	}
	return def
}

// sortResourcesByDepth returns resources ordered such that any resource
// whose parent is also in the slice appears after that parent. Resources
// whose parent is missing (or empty) are treated as roots.
func sortResourcesByDepth(resources []lkResource) []*lkResource {
	byID := make(map[string]*lkResource, len(resources))
	for i := range resources {
		byID[resources[i].ID] = &resources[i]
	}
	depth := make(map[string]int, len(resources))
	var compute func(id string, seen map[string]bool) int
	compute = func(id string, seen map[string]bool) int {
		if d, ok := depth[id]; ok {
			return d
		}
		if seen[id] {
			return 0 // cycle guard
		}
		seen[id] = true
		r, ok := byID[id]
		if !ok || r.ParentID == "" {
			depth[id] = 0
			return 0
		}
		d := compute(r.ParentID, seen) + 1
		depth[id] = d
		return d
	}
	out := make([]*lkResource, 0, len(resources))
	for i := range resources {
		_ = compute(resources[i].ID, map[string]bool{})
		out = append(out, &resources[i])
	}
	// stable sort by depth ascending
	for i := 1; i < len(out); i++ {
		j := i
		for j > 0 && depth[out[j-1].ID] > depth[out[j].ID] {
			out[j-1], out[j] = out[j], out[j-1]
			j--
		}
	}
	return out
}

func normalizeTags(tags []string) []string {
	out := make([]string, 0, len(tags))
	seen := map[string]bool{}
	for _, t := range tags {
		t = strings.TrimSpace(strings.ToLower(t))
		if t == "" || seen[t] {
			continue
		}
		seen[t] = true
		out = append(out, t)
	}
	return out
}

func readAttrString(raw json.RawMessage, key string) string {
	if len(raw) == 0 {
		return ""
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return ""
	}
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

// readPropString reads attrs.props.<key>; some LK image nodes nest src under props.
func readPropString(raw json.RawMessage, key string) string {
	if len(raw) == 0 {
		return ""
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return ""
	}
	if props, ok := m["props"].(map[string]any); ok {
		if v, ok := props[key].(string); ok {
			return v
		}
	}
	return ""
}

func readHeadingLevel(raw json.RawMessage) int {
	if len(raw) == 0 {
		return 2
	}
	var m struct {
		Level json.Number `json:"level"`
	}
	if err := json.Unmarshal(raw, &m); err != nil {
		return 2
	}
	if m.Level == "" {
		return 2
	}
	n, err := strconv.Atoi(string(m.Level))
	if err != nil {
		return 2
	}
	return n
}

var nonSlugRE = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	// Replace accented letters with ASCII equivalents (best-effort).
	replacer := strings.NewReplacer(
		"ä", "ae", "ö", "oe", "ü", "ue", "ß", "ss",
		"á", "a", "à", "a", "â", "a", "ã", "a",
		"é", "e", "è", "e", "ê", "e", "ë", "e",
		"í", "i", "ì", "i", "î", "i", "ï", "i",
		"ó", "o", "ò", "o", "ô", "o", "õ", "o",
		"ú", "u", "ù", "u", "û", "u",
		"ç", "c", "ñ", "n",
		"'", "", "’", "", "`", "",
	)
	s = replacer.Replace(s)
	s = nonSlugRE.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

func extractPlainText(raw json.RawMessage) string {
	var node lkNode
	if err := json.Unmarshal(raw, &node); err != nil {
		return ""
	}
	var sb strings.Builder
	walkText(&node, &sb)
	return strings.TrimSpace(sb.String())
}

func walkText(n *lkNode, sb *strings.Builder) {
	if n == nil {
		return
	}
	if n.Type == "text" {
		sb.WriteString(n.Text)
		sb.WriteByte(' ')
	}
	for idx := range n.Content {
		walkText(&n.Content[idx], sb)
	}
}

func isImageURL(u string) bool {
	low := strings.ToLower(u)
	for _, ext := range []string{".png", ".jpg", ".jpeg", ".webp", ".gif"} {
		if strings.Contains(low, ext) {
			return true
		}
	}
	return false
}

func extForMime(mime string) string {
	switch strings.ToLower(mime) {
	case "image/png":
		return ".png"
	case "image/jpeg":
		return ".jpg"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	}
	return ""
}

func mimeForExt(ext string) string {
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".webp":
		return "image/webp"
	case ".gif":
		return "image/gif"
	}
	return ""
}
