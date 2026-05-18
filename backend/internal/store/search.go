package store

import (
	"context"
	"sort"
	"strings"
)

func (s *Store) Search(ctx context.Context, query string, visibility []string, limit int) ([]SearchHit, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []SearchHit{}, nil
	}
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	clause := ""
	args := []any{ftsQuery(query)}
	if len(visibility) > 0 {
		placeholders := strings.TrimRight(strings.Repeat("?,", len(visibility)), ",")
		clause = " AND e.visibility IN (" + placeholders + ")"
		for _, v := range visibility {
			args = append(args, v)
		}
	}
	args = append(args, limit)
	rows, err := s.db.QueryContext(ctx,
		`SELECT e.id, e.slug, e.title, COALESCE(e.summary, ''), e.entity_type_id,
		   snippet(entities_fts, 2, '<mark>', '</mark>', '…', 12), bm25(entities_fts)
		 FROM entities_fts JOIN entities e ON e.id = entities_fts.rowid
		 WHERE entities_fts MATCH ?`+clause+`
		 ORDER BY bm25(entities_fts) LIMIT ?`,
		args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]SearchHit, 0)
	seen := map[int64]bool{}
	for rows.Next() {
		var hit SearchHit
		if err := rows.Scan(&hit.ID, &hit.Slug, &hit.Title, &hit.Summary, &hit.EntityTypeID, &hit.Snippet, &hit.Rank); err != nil {
			return nil, err
		}
		seen[hit.ID] = true
		out = append(out, hit)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(out) < limit/2 || len(out) == 0 {
		need := limit - len(out)
		fuzzy, ferr := s.fuzzySearchTitles(ctx, query, visibility, seen, need)
		if ferr == nil {
			out = append(out, fuzzy...)
		}
	}

	return out, nil
}

func (s *Store) fuzzySearchTitles(ctx context.Context, query string, visibility []string, exclude map[int64]bool, maxResults int) ([]SearchHit, error) {
	if maxResults <= 0 {
		return []SearchHit{}, nil
	}
	args := []any{}
	clause := ""
	if len(visibility) > 0 {
		placeholders := strings.TrimRight(strings.Repeat("?,", len(visibility)), ",")
		clause = " WHERE visibility IN (" + placeholders + ")"
		for _, v := range visibility {
			args = append(args, v)
		}
	}
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, entity_type_id, title, slug, COALESCE(summary, '')
		 FROM entities`+clause+` LIMIT 5000`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	qLower := strings.ToLower(query)
	qRunes := []rune(qLower)
	type scored struct {
		hit   SearchHit
		score float64
	}
	cands := make([]scored, 0)
	for rows.Next() {
		var hit SearchHit
		if err := rows.Scan(&hit.ID, &hit.EntityTypeID, &hit.Title, &hit.Slug, &hit.Summary); err != nil {
			return nil, err
		}
		if exclude[hit.ID] {
			continue
		}
		titleLower := strings.ToLower(hit.Title)
		titleRunes := []rune(titleLower)
		if len(titleRunes) == 0 {
			continue
		}
		score := similarity(qLower, titleLower, qRunes, titleRunes)
		if strings.Contains(titleLower, qLower) {
			score = maxFloat(score, 0.85)
		}
		if score < 0.5 {
			continue
		}
		hit.Rank = -score
		cands = append(cands, scored{hit: hit, score: score})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	sort.Slice(cands, func(i, j int) bool {
		return cands[i].score > cands[j].score
	})
	if len(cands) > maxResults {
		cands = cands[:maxResults]
	}
	out := make([]SearchHit, 0, len(cands))
	for _, c := range cands {
		out = append(out, c.hit)
	}
	return out, nil
}

func similarity(a, b string, ar, br []rune) float64 {
	maxLen := len(ar)
	if len(br) > maxLen {
		maxLen = len(br)
	}
	if maxLen == 0 {
		return 0
	}
	d := levenshtein(ar, br)
	return 1.0 - float64(d)/float64(maxLen)
}

func levenshtein(a, b []rune) int {
	la, lb := len(a), len(b)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}
	if la > lb {
		a, b = b, a
		la, lb = lb, la
	}
	prev := make([]int, la+1)
	curr := make([]int, la+1)
	for i := 0; i <= la; i++ {
		prev[i] = i
	}
	for j := 1; j <= lb; j++ {
		curr[0] = j
		for i := 1; i <= la; i++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			curr[i] = min3(curr[i-1]+1, prev[i]+1, prev[i-1]+cost)
		}
		prev, curr = curr, prev
	}
	return prev[la]
}

func min3(a, b, c int) int {
	if a > b {
		a = b
	}
	if a > c {
		return c
	}
	return a
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func ftsQuery(q string) string {
	parts := strings.Fields(q)
	for i, p := range parts {
		p = strings.Map(func(r rune) rune {
			switch r {
			case '"', '*', '(', ')':
				return -1
			}
			return r
		}, p)
		if p == "" {
			continue
		}
		parts[i] = "\"" + p + "\"*"
	}
	return strings.Join(parts, " ")
}
