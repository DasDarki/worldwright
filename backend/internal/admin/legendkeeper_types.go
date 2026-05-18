package admin

import "encoding/json"

// LegendKeeper export schema (subset). Only fields the importer needs are
// modelled; unknown fields are silently ignored.

type lkExport struct {
	Version    int          `json:"version"`
	ExportID   string       `json:"exportId"`
	ExportedAt string       `json:"exportedAt"`
	Resources  []lkResource `json:"resources"`
}

type lkResource struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	ParentID   string       `json:"parentId"`
	Pos        string       `json:"pos"`
	Tags       []string     `json:"tags"`
	Aliases    []string     `json:"aliases"`
	IconGlyph  string       `json:"iconGlyph"`
	IconColor  string       `json:"iconColor"`
	Properties []lkProperty `json:"properties"`
	Documents  []lkDocument `json:"documents"`
}

type lkProperty struct {
	ID    string          `json:"id"`
	Title string          `json:"title"`
	Type  string          `json:"type"`
	Data  json.RawMessage `json:"data"`
}

type lkPropertyImageData struct {
	URL    string  `json:"url"`
	Height float64 `json:"height"`
	Scale  float64 `json:"scale"`
}

type lkPropertyTextData struct {
	Fragment json.RawMessage `json:"fragment"`
}

type lkDocument struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Type    string          `json:"type"`  // page | map | blank | time | board
	Pos     string          `json:"pos"`
	Content json.RawMessage `json:"content"`
	Map     *lkMapMeta      `json:"map,omitempty"`
}

type lkMapMeta struct {
	LocatorID string  `json:"locatorId"`
	MapID     string  `json:"mapId"`
	MinX      float64 `json:"min_x"`
	MaxX      float64 `json:"max_x"`
	MinY      float64 `json:"min_y"`
	MaxY      float64 `json:"max_y"`
	MaxZoom   int     `json:"max_zoom"`
}

type lkMapContent struct {
	Pins []lkMapPin `json:"pins"`
}

type lkMapPin struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Pos             []float64 `json:"pos"`
	Rank            string    `json:"rank"`
	IconGlyph       string    `json:"iconGlyph"`
	IconColor       string    `json:"iconColor"`
	LabelVisibility string    `json:"labelVisibility"`
	IsSynced        bool      `json:"isSynced"`
	ResourceID      string    `json:"resourceId"`
}

// Generic TipTap-ish node used during translation.
type lkNode struct {
	Type    string          `json:"type"`
	Text    string          `json:"text,omitempty"`
	Marks   []lkMark        `json:"marks,omitempty"`
	Attrs   json.RawMessage `json:"attrs,omitempty"`
	Content []lkNode        `json:"content,omitempty"`
}

type lkMark struct {
	Type  string          `json:"type"`
	Attrs json.RawMessage `json:"attrs,omitempty"`
}
