package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
)

func (s *Store) ListTimelines(ctx context.Context) ([]Timeline, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT t.id, t.name, COALESCE(t.description, ''), COUNT(te.event_id)
		 FROM timelines t
		 LEFT JOIN timeline_events te ON te.timeline_id = t.id
		 GROUP BY t.id
		 ORDER BY t.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Timeline, 0)
	for rows.Next() {
		var t Timeline
		var count int
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &count); err != nil {
			return nil, err
		}
		t.EventIDs = make([]int64, 0, count)
		out = append(out, t)
	}
	return out, rows.Err()
}

func (s *Store) TimelineByID(ctx context.Context, id int64, visibility []string) (*Timeline, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, name, COALESCE(description, '') FROM timelines WHERE id = ?`, id)
	var t Timeline
	if err := row.Scan(&t.ID, &t.Name, &t.Description); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	clause := ""
	args := []any{id}
	if len(visibility) > 0 {
		placeholders := strings.TrimRight(strings.Repeat("?,", len(visibility)), ",")
		clause = " AND ev.visibility IN (" + placeholders + ")"
		for _, v := range visibility {
			args = append(args, v)
		}
	}

	idRows, err := s.db.QueryContext(ctx,
		`SELECT event_id FROM timeline_events WHERE timeline_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer idRows.Close()
	t.EventIDs = []int64{}
	for idRows.Next() {
		var eid int64
		if err := idRows.Scan(&eid); err != nil {
			return nil, err
		}
		t.EventIDs = append(t.EventIDs, eid)
	}

	evRows, err := s.db.QueryContext(ctx,
		`SELECT ev.id, ev.title, ev.body, ev.calendar_id, ev.era_id, ev.year, ev.month_index, ev.day,
		    ev.end_year, ev.end_month_index, ev.end_day, ev.importance, ev.visibility
		 FROM timeline_events te
		 JOIN events ev ON ev.id = te.event_id
		 WHERE te.timeline_id = ?`+clause+`
		 ORDER BY ev.year, ev.month_index, ev.day`, args...)
	if err != nil {
		return nil, err
	}
	defer evRows.Close()
	t.Events = make([]Event, 0)
	for evRows.Next() {
		var e Event
		var bodyText string
		var eraID, endY, endM, endD sql.NullInt64
		if err := evRows.Scan(&e.ID, &e.Title, &bodyText, &e.CalendarID, &eraID,
			&e.Year, &e.MonthIndex, &e.Day, &endY, &endM, &endD,
			&e.Importance, &e.Visibility); err != nil {
			return nil, err
		}
		e.Body = json.RawMessage(bodyText)
		if eraID.Valid { v := eraID.Int64; e.EraID = &v }
		if endY.Valid { v := int(endY.Int64); e.EndYear = &v }
		if endM.Valid { v := int(endM.Int64); e.EndMonthIndex = &v }
		if endD.Valid { v := int(endD.Int64); e.EndDay = &v }
		e.Participants = []EventEntity{}
		t.Events = append(t.Events, e)
	}
	return &t, evRows.Err()
}

func (s *Store) CreateTimeline(ctx context.Context, n NewTimeline) (*Timeline, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx,
		`INSERT INTO timelines (name, description) VALUES (?, ?)`,
		n.Name, nullable(n.Description))
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()

	if err := syncTimelineEvents(ctx, tx, id, n.EventIDs); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.TimelineByID(ctx, id, nil)
}

func (s *Store) UpdateTimeline(ctx context.Context, id int64, n NewTimeline) (*Timeline, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx,
		`UPDATE timelines SET name = ?, description = ? WHERE id = ?`,
		n.Name, nullable(n.Description), id); err != nil {
		return nil, err
	}

	if err := syncTimelineEvents(ctx, tx, id, n.EventIDs); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.TimelineByID(ctx, id, nil)
}

func (s *Store) DeleteTimeline(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM timelines WHERE id = ?`, id)
	return err
}

func syncTimelineEvents(ctx context.Context, tx *sql.Tx, timelineID int64, eventIDs []int64) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM timeline_events WHERE timeline_id = ?`, timelineID); err != nil {
		return err
	}
	for _, eid := range eventIDs {
		if _, err := tx.ExecContext(ctx,
			`INSERT OR IGNORE INTO timeline_events (timeline_id, event_id) VALUES (?, ?)`,
			timelineID, eid); err != nil {
			return err
		}
	}
	return nil
}
