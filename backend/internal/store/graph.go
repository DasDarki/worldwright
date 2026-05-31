package store

import (
	"context"
	"strings"
)

// RelationshipGraph captures every node and edge needed to render a
// relationship widget over a closed set of entity IDs. Edges are only emitted
// when BOTH endpoints are present in the visible node set, so visibility
// scrubbing happens automatically: any relationship that would point at a
// secret entity simply doesn't get returned to a viewer who can't see it.
type RelationshipGraph struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}

type GraphNode struct {
	ID           int64  `json:"id"`
	EntityTypeID int64  `json:"entity_type_id"`
	Title        string `json:"title"`
	Slug         string `json:"slug"`
	Summary      string `json:"summary,omitempty"`
}

type GraphEdge struct {
	From        int64  `json:"from"`
	To          int64  `json:"to"`
	TypeKey     string `json:"type_key"`
	TypeLabel   string `json:"type_label"`
	Inverse     string `json:"inverse_label,omitempty"`
	IsSymmetric bool   `json:"is_symmetric"`
	Category    string `json:"category"`
	Description string `json:"description,omitempty"`
}

func (s *Store) RelationshipGraphForIDs(ctx context.Context, ids []int64, visibility []string, lang string) (*RelationshipGraph, error) {
	out := &RelationshipGraph{Nodes: []GraphNode{}, Edges: []GraphEdge{}}
	if len(ids) == 0 {
		return out, nil
	}

	labelEn := lang != "de"

	idPlaceholders := strings.TrimRight(strings.Repeat("?,", len(ids)), ",")
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}

	visClause := ""
	if len(visibility) > 0 {
		vp := strings.TrimRight(strings.Repeat("?,", len(visibility)), ",")
		visClause = " AND visibility IN (" + vp + ")"
		for _, v := range visibility {
			args = append(args, v)
		}
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT id, entity_type_id, title, slug, COALESCE(summary, '')
		 FROM entities WHERE id IN (`+idPlaceholders+`)`+visClause+`
		 ORDER BY title`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	visible := map[int64]bool{}
	for rows.Next() {
		var n GraphNode
		if err := rows.Scan(&n.ID, &n.EntityTypeID, &n.Title, &n.Slug, &n.Summary); err != nil {
			return nil, err
		}
		out.Nodes = append(out.Nodes, n)
		visible[n.ID] = true
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(out.Nodes) < 2 {
		return out, nil
	}

	// Pull every relationship whose endpoints both sit in the visible set.
	visibleIDs := make([]any, 0, len(visible))
	for id := range visible {
		visibleIDs = append(visibleIDs, id)
	}
	visiblePlaceholders := strings.TrimRight(strings.Repeat("?,", len(visibleIDs)), ",")
	edgeArgs := append([]any{}, visibleIDs...)
	edgeArgs = append(edgeArgs, visibleIDs...)

	edgeRows, err := s.db.QueryContext(ctx,
		`SELECT r.from_entity_id, r.to_entity_id,
		        rt.key, rt.label_en, rt.label_de, rt.inverse_label_en, rt.inverse_label_de,
		        rt.is_symmetric, rt.category,
		        COALESCE(r.description, '')
		 FROM relationships r
		 JOIN relationship_types rt ON rt.id = r.relationship_type_id
		 WHERE r.from_entity_id IN (`+visiblePlaceholders+`)
		   AND r.to_entity_id   IN (`+visiblePlaceholders+`)`,
		edgeArgs...)
	if err != nil {
		return nil, err
	}
	defer edgeRows.Close()

	for edgeRows.Next() {
		var e GraphEdge
		var labelEnStr, labelDeStr, invEn, invDe string
		var isSym int
		if err := edgeRows.Scan(&e.From, &e.To, &e.TypeKey,
			&labelEnStr, &labelDeStr, &invEn, &invDe,
			&isSym, &e.Category, &e.Description); err != nil {
			return nil, err
		}
		e.IsSymmetric = isSym != 0
		if labelEn {
			e.TypeLabel = labelEnStr
			e.Inverse = invEn
		} else {
			e.TypeLabel = labelDeStr
			e.Inverse = invDe
		}
		out.Edges = append(out.Edges, e)
	}
	return out, edgeRows.Err()
}
