package store

import (
	"encoding/json"
	"time"
)

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	HasPassword  bool      `json:"has_password"`
	Role         string    `json:"role"`
	Locale       string    `json:"locale"`
	DisplayName  string    `json:"display_name,omitempty"`
	AvatarURL    string    `json:"avatar_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type NewUser struct {
	Email        string
	PasswordHash string
	Role         string
	Locale       string
	DisplayName  string
	AvatarURL    string
}

type OAuthIdentity struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Provider  string    `json:"provider"`
	Subject   string    `json:"subject"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type FieldDefinition struct {
	ID           int64           `json:"id"`
	EntityTypeID int64           `json:"entity_type_id"`
	Key          string          `json:"key"`
	LabelEn      string          `json:"label_en"`
	LabelDe      string          `json:"label_de"`
	DataType     string          `json:"data_type"`
	Config       json.RawMessage `json:"config,omitempty"`
	SortOrder    int             `json:"sort_order"`
	IsRequired   bool            `json:"is_required"`
}

type EntityType struct {
	ID       int64             `json:"id"`
	Key      string            `json:"key"`
	NameEn   string            `json:"name_en"`
	NameDe   string            `json:"name_de"`
	Icon     string            `json:"icon,omitempty"`
	Color    string            `json:"color,omitempty"`
	IsSystem bool              `json:"is_system"`
	Fields   []FieldDefinition `json:"fields,omitempty"`
}

type Entity struct {
	ID               int64             `json:"id"`
	EntityTypeID     int64             `json:"entity_type_id"`
	EntityType       *EntityType       `json:"entity_type,omitempty"`
	Title            string            `json:"title"`
	Slug             string            `json:"slug"`
	Summary          string            `json:"summary,omitempty"`
	Body             json.RawMessage   `json:"body"`
	BodyText         string            `json:"-"`
	ParentID         *int64            `json:"parent_id,omitempty"`
	MaterializedPath string            `json:"materialized_path"`
	Visibility       string            `json:"visibility"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	Tags             []string          `json:"tags"`
	FieldValues      map[string]string `json:"field_values"`
}

type EntitySummary struct {
	ID           int64  `json:"id"`
	EntityTypeID int64  `json:"entity_type_id"`
	Title        string `json:"title"`
	Slug         string `json:"slug"`
	Summary      string `json:"summary,omitempty"`
	ParentID     *int64 `json:"parent_id,omitempty"`
	Visibility   string `json:"visibility"`
}

type NewEntity struct {
	EntityTypeID int64
	Title        string
	Slug         string
	Summary      string
	Body         json.RawMessage
	BodyText     string
	ParentID     *int64
	Visibility   string
	Tags         []string
	FieldValues  map[string]string
	Wikilinks    []string
}

type TagWithCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type GenealogyNode struct {
	ID           int64  `json:"id"`
	EntityTypeID int64  `json:"entity_type_id"`
	Title        string `json:"title"`
	Slug         string `json:"slug"`
	Summary      string `json:"summary,omitempty"`
	Visibility   string `json:"visibility"`
}

type GenealogyEdge struct {
	From int64  `json:"from"`
	To   int64  `json:"to"`
	Type string `json:"type"`
}

type Genealogy struct {
	Focal int64           `json:"focal"`
	Nodes []GenealogyNode `json:"nodes"`
	Edges []GenealogyEdge `json:"edges"`
}

type Backlink struct {
	SourceEntityID int64  `json:"source_entity_id"`
	Title          string `json:"title"`
	Slug           string `json:"slug"`
	EntityTypeID   int64  `json:"entity_type_id"`
	Summary        string `json:"summary,omitempty"`
}

type Asset struct {
	ID        int64     `json:"id"`
	Filename  string    `json:"filename"`
	Path      string    `json:"-"`
	Mime      string    `json:"mime"`
	Size      int64     `json:"size"`
	Width     *int      `json:"width,omitempty"`
	Height    *int      `json:"height,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type NewAsset struct {
	Filename string
	Path     string
	Mime     string
	Size     int64
	Width    *int
	Height   *int
}

type MapAsset struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	AssetID        int64     `json:"asset_id"`
	Width          int       `json:"width"`
	Height         int       `json:"height"`
	ParentEntityID *int64    `json:"parent_entity_id,omitempty"`
	Pins           []MapPin  `json:"pins"`
}

type MapPin struct {
	ID               int64   `json:"id"`
	MapID            int64   `json:"map_id"`
	X                float64 `json:"x"`
	Y                float64 `json:"y"`
	Label            string  `json:"label,omitempty"`
	Icon             string  `json:"icon,omitempty"`
	TargetEntityID   *int64  `json:"target_entity_id,omitempty"`
	TargetEntitySlug string  `json:"target_entity_slug,omitempty"`
	TargetMapID      *int64  `json:"target_map_id,omitempty"`
	Visibility       string  `json:"visibility"`
}

type NewMapPin struct {
	X              float64
	Y              float64
	Label          string
	Icon           string
	TargetEntityID *int64
	TargetMapID    *int64
	Visibility     string
}

type Calendar struct {
	ID                int64             `json:"id"`
	Name              string            `json:"name"`
	EpochName         string            `json:"epoch_name,omitempty"`
	CurrentYear       int               `json:"current_year"`
	CurrentMonthIndex int               `json:"current_month_index"`
	CurrentDay        int               `json:"current_day"`
	LeapRule          string            `json:"leap_rule,omitempty"`
	Months            []CalendarMonth   `json:"months"`
	Weekdays          []CalendarWeekday `json:"weekdays"`
	Eras              []CalendarEra     `json:"eras"`
	Moons             []CalendarMoon    `json:"moons"`
}

type CalendarMonth struct {
	ID         int64  `json:"id"`
	CalendarID int64  `json:"calendar_id"`
	SortOrder  int    `json:"sort_order"`
	Name       string `json:"name"`
	Days       int    `json:"days"`
}

type CalendarWeekday struct {
	ID         int64  `json:"id"`
	CalendarID int64  `json:"calendar_id"`
	SortOrder  int    `json:"sort_order"`
	Name       string `json:"name"`
}

type CalendarEra struct {
	ID         int64  `json:"id"`
	CalendarID int64  `json:"calendar_id"`
	Name       string `json:"name"`
	StartYear  int    `json:"start_year"`
	Suffix     string `json:"suffix,omitempty"`
}

type CalendarMoon struct {
	ID         int64               `json:"id"`
	CalendarID int64               `json:"calendar_id"`
	Name       string              `json:"name"`
	CycleDays  float64             `json:"cycle_days"`
	OffsetDays float64             `json:"offset_days"`
	Phases     []CalendarMoonPhase `json:"phases"`
}

type CalendarMoonPhase struct {
	ID            int64   `json:"id"`
	MoonID        int64   `json:"moon_id"`
	SortOrder     int     `json:"sort_order"`
	Name          string  `json:"name"`
	CyclePosition float64 `json:"cycle_position"`
	Randomness    float64 `json:"randomness"`
	Icon          string  `json:"icon,omitempty"`
}

type Event struct {
	ID            int64            `json:"id"`
	Title         string           `json:"title"`
	Body          json.RawMessage  `json:"body"`
	BodyText      string           `json:"-"`
	CalendarID    int64            `json:"calendar_id"`
	EraID         *int64           `json:"era_id,omitempty"`
	Year          int              `json:"year"`
	MonthIndex    int              `json:"month_index"`
	Day           int              `json:"day"`
	EndYear       *int             `json:"end_year,omitempty"`
	EndMonthIndex *int             `json:"end_month_index,omitempty"`
	EndDay        *int             `json:"end_day,omitempty"`
	Importance    int              `json:"importance"`
	Visibility    string           `json:"visibility"`
	Participants  []EventEntity    `json:"participants"`
}

type EventEntity struct {
	EventID  int64          `json:"event_id"`
	EntityID int64          `json:"entity_id"`
	Role     string         `json:"role"`
	Entity   *EntitySummary `json:"entity,omitempty"`
}

type Timeline struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	EventIDs    []int64 `json:"event_ids"`
	Events      []Event `json:"events,omitempty"`
}

type NewTimeline struct {
	Name        string
	Description string
	EventIDs    []int64
}

type RelationshipType struct {
	ID             int64  `json:"id"`
	Key            string `json:"key"`
	LabelEn        string `json:"label_en"`
	LabelDe        string `json:"label_de"`
	InverseLabelEn string `json:"inverse_label_en,omitempty"`
	InverseLabelDe string `json:"inverse_label_de,omitempty"`
	IsSymmetric    bool   `json:"is_symmetric"`
	Category       string `json:"category"`
}

type Relationship struct {
	ID                 int64  `json:"id"`
	FromEntityID       int64  `json:"from_entity_id"`
	ToEntityID         int64  `json:"to_entity_id"`
	RelationshipTypeID int64  `json:"relationship_type_id"`
	Description        string `json:"description,omitempty"`
}

type RelationshipEdge struct {
	ID          int64            `json:"id"`
	Type        RelationshipType `json:"type"`
	Direction   string           `json:"direction"`
	Other       EntitySummary    `json:"other"`
	Description string           `json:"description,omitempty"`
}

type SearchHit struct {
	ID           int64   `json:"id"`
	Slug         string  `json:"slug"`
	Title        string  `json:"title"`
	Summary      string  `json:"summary,omitempty"`
	EntityTypeID int64   `json:"entity_type_id"`
	Snippet      string  `json:"snippet,omitempty"`
	Rank         float64 `json:"rank"`
}
