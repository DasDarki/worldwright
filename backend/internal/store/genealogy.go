package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func (s *Store) GenealogyForEntity(ctx context.Context, entityID int64, depth int, visibility []string) (*Genealogy, error) {
	if depth < 1 {
		depth = 1
	}
	if depth > 6 {
		depth = 6
	}

	visited := map[int64]bool{entityID: true}
	edges := []GenealogyEdge{}
	edgeKeys := map[string]bool{}

	addEdge := func(from, to int64, t string) {
		key := keyOf(from, to, t)
		if edgeKeys[key] {
			return
		}
		edgeKeys[key] = true
		edges = append(edges, GenealogyEdge{From: from, To: to, Type: t})
	}

	frontierUp := []int64{entityID}
	for d := 0; d < depth && len(frontierUp) > 0; d++ {
		next := []int64{}
		for _, id := range frontierUp {
			parents, err := s.parentsOf(ctx, id)
			if err != nil {
				return nil, err
			}
			for _, pid := range parents {
				addEdge(pid, id, "parent_of")
				if !visited[pid] {
					visited[pid] = true
					next = append(next, pid)
				}
			}
		}
		frontierUp = next
	}

	frontierDown := []int64{entityID}
	for d := 0; d < depth && len(frontierDown) > 0; d++ {
		next := []int64{}
		for _, id := range frontierDown {
			children, err := s.childrenOf(ctx, id)
			if err != nil {
				return nil, err
			}
			for _, cid := range children {
				addEdge(id, cid, "parent_of")
				if !visited[cid] {
					visited[cid] = true
					next = append(next, cid)
				}
			}
		}
		frontierDown = next
	}

	for id := range visited {
		spouses, err := s.symmetricKin(ctx, id, "spouse_of")
		if err != nil {
			return nil, err
		}
		for _, sid := range spouses {
			from, to := id, sid
			if from > to {
				from, to = to, from
			}
			addEdge(from, to, "spouse_of")
			if !visited[sid] {
				visited[sid] = true
			}
		}
		siblings, err := s.symmetricKin(ctx, id, "sibling_of")
		if err != nil {
			return nil, err
		}
		for _, sid := range siblings {
			from, to := id, sid
			if from > to {
				from, to = to, from
			}
			addEdge(from, to, "sibling_of")
			if !visited[sid] {
				visited[sid] = true
			}
		}
	}

	ids := make([]int64, 0, len(visited))
	for id := range visited {
		ids = append(ids, id)
	}

	nodes, err := s.entitiesByIDs(ctx, ids, visibility)
	if err != nil {
		return nil, err
	}

	allowed := map[int64]bool{}
	for _, n := range nodes {
		allowed[n.ID] = true
	}
	keptEdges := edges[:0]
	for _, e := range edges {
		if allowed[e.From] && allowed[e.To] {
			keptEdges = append(keptEdges, e)
		}
	}

	return &Genealogy{
		Focal: entityID,
		Nodes: nodes,
		Edges: keptEdges,
	}, nil
}

func (s *Store) parentsOf(ctx context.Context, entityID int64) ([]int64, error) {
	return s.kinByDirection(ctx, entityID, "parent_of", false)
}

func (s *Store) childrenOf(ctx context.Context, entityID int64) ([]int64, error) {
	return s.kinByDirection(ctx, entityID, "parent_of", true)
}

func (s *Store) kinByDirection(ctx context.Context, entityID int64, typeKey string, outgoing bool) ([]int64, error) {
	col := "from_entity_id"
	target := "to_entity_id"
	if outgoing {
		col, target = target, col
	}
	rows, err := s.db.QueryContext(ctx,
		`SELECT r.`+col+` FROM relationships r
		 JOIN relationship_types rt ON rt.id = r.relationship_type_id
		 WHERE rt.key = ? AND r.`+target+` = ?`,
		typeKey, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]int64, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
}

func (s *Store) symmetricKin(ctx context.Context, entityID int64, typeKey string) ([]int64, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT CASE WHEN r.from_entity_id = ? THEN r.to_entity_id ELSE r.from_entity_id END
		 FROM relationships r
		 JOIN relationship_types rt ON rt.id = r.relationship_type_id
		 WHERE rt.key = ? AND (r.from_entity_id = ? OR r.to_entity_id = ?)`,
		entityID, typeKey, entityID, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]int64, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
}

func (s *Store) entitiesByIDs(ctx context.Context, ids []int64, visibility []string) ([]GenealogyNode, error) {
	if len(ids) == 0 {
		return []GenealogyNode{}, nil
	}
	placeholders := strings.TrimRight(strings.Repeat("?,", len(ids)), ",")
	args := make([]any, 0, len(ids)+len(visibility))
	for _, id := range ids {
		args = append(args, id)
	}
	clause := ""
	if len(visibility) > 0 {
		vp := strings.TrimRight(strings.Repeat("?,", len(visibility)), ",")
		clause = " AND visibility IN (" + vp + ")"
		for _, v := range visibility {
			args = append(args, v)
		}
	}
	q := `SELECT id, entity_type_id, title, slug, COALESCE(summary, ''), visibility
	      FROM entities WHERE id IN (` + placeholders + `)` + clause
	rows, err := s.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]GenealogyNode, 0, len(ids))
	for rows.Next() {
		var n GenealogyNode
		if err := rows.Scan(&n.ID, &n.EntityTypeID, &n.Title, &n.Slug, &n.Summary, &n.Visibility); err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			return nil, err
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

func keyOf(from, to int64, t string) string {
	return fmt.Sprintf("%d-%d|%s", from, to, t)
}
