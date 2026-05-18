package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
)

type EventFilter struct {
	CalendarID int64
	Visibility []string
}

type NewEvent struct {
	Title         string
	Body          json.RawMessage
	BodyText      string
	CalendarID    int64
	EraID         *int64
	Year          int
	MonthIndex    int
	Day           int
	EndYear       *int
	EndMonthIndex *int
	EndDay        *int
	Importance    int
	Visibility    string
	Participants  []NewEventParticipant
}

type NewEventParticipant struct {
	EntityID int64
	Role     string
}

func (s *Store) ListEvents(ctx context.Context, filter EventFilter) ([]Event, error) {
	clauses := []string{"1=1"}
	args := []any{}
	if filter.CalendarID != 0 {
		clauses = append(clauses, "calendar_id = ?")
		args = append(args, filter.CalendarID)
	}
	if len(filter.Visibility) > 0 {
		placeholders := strings.TrimRight(strings.Repeat("?,", len(filter.Visibility)), ",")
		clauses = append(clauses, "visibility IN ("+placeholders+")")
		for _, v := range filter.Visibility {
			args = append(args, v)
		}
	}
	query := `SELECT id, title, body, calendar_id, era_id, year, month_index, day,
	             end_year, end_month_index, end_day, importance, visibility
	          FROM events WHERE ` + strings.Join(clauses, " AND ") + `
	          ORDER BY year, month_index, day`
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Event, 0)
	for rows.Next() {
		var e Event
		var bodyText string
		var eraID, endY, endM, endD sql.NullInt64
		if err := rows.Scan(&e.ID, &e.Title, &bodyText, &e.CalendarID, &eraID,
			&e.Year, &e.MonthIndex, &e.Day, &endY, &endM, &endD,
			&e.Importance, &e.Visibility); err != nil {
			return nil, err
		}
		e.Body = json.RawMessage(bodyText)
		if eraID.Valid { id := eraID.Int64; e.EraID = &id }
		if endY.Valid { v := int(endY.Int64); e.EndYear = &v }
		if endM.Valid { v := int(endM.Int64); e.EndMonthIndex = &v }
		if endD.Valid { v := int(endD.Int64); e.EndDay = &v }
		e.Participants = []EventEntity{}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (s *Store) EventByID(ctx context.Context, id int64, visibility []string) (*Event, error) {
	clause := ""
	args := []any{id}
	if len(visibility) > 0 {
		placeholders := strings.TrimRight(strings.Repeat("?,", len(visibility)), ",")
		clause = " AND visibility IN (" + placeholders + ")"
		for _, v := range visibility {
			args = append(args, v)
		}
	}
	row := s.db.QueryRowContext(ctx,
		`SELECT id, title, body, calendar_id, era_id, year, month_index, day,
		   end_year, end_month_index, end_day, importance, visibility
		 FROM events WHERE id = ?`+clause, args...)
	var e Event
	var bodyText string
	var eraID, endY, endM, endD sql.NullInt64
	if err := row.Scan(&e.ID, &e.Title, &bodyText, &e.CalendarID, &eraID,
		&e.Year, &e.MonthIndex, &e.Day, &endY, &endM, &endD,
		&e.Importance, &e.Visibility); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	e.Body = json.RawMessage(bodyText)
	if eraID.Valid { v := eraID.Int64; e.EraID = &v }
	if endY.Valid { v := int(endY.Int64); e.EndYear = &v }
	if endM.Valid { v := int(endM.Int64); e.EndMonthIndex = &v }
	if endD.Valid { v := int(endD.Int64); e.EndDay = &v }

	parts, err := s.eventParticipants(ctx, e.ID)
	if err != nil {
		return nil, err
	}
	e.Participants = parts
	return &e, nil
}

func (s *Store) eventParticipants(ctx context.Context, eventID int64) ([]EventEntity, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT ee.event_id, ee.entity_id, ee.role,
		    e.id, e.entity_type_id, e.title, e.slug, COALESCE(e.summary,''), e.parent_id, e.visibility
		 FROM event_entities ee
		 JOIN entities e ON e.id = ee.entity_id
		 WHERE ee.event_id = ?
		 ORDER BY ee.role, e.title`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]EventEntity, 0)
	for rows.Next() {
		var ee EventEntity
		var summary EntitySummary
		var parentID sql.NullInt64
		if err := rows.Scan(&ee.EventID, &ee.EntityID, &ee.Role,
			&summary.ID, &summary.EntityTypeID, &summary.Title, &summary.Slug,
			&summary.Summary, &parentID, &summary.Visibility); err != nil {
			return nil, err
		}
		if parentID.Valid {
			id := parentID.Int64
			summary.ParentID = &id
		}
		s := summary
		ee.Entity = &s
		out = append(out, ee)
	}
	return out, rows.Err()
}

func (s *Store) CreateEvent(ctx context.Context, n NewEvent) (*Event, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	body := string(n.Body)
	if body == "" {
		body = "{}"
	}
	res, err := tx.ExecContext(ctx,
		`INSERT INTO events (title, body, calendar_id, era_id, year, month_index, day,
		   end_year, end_month_index, end_day, importance, visibility)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		n.Title, body, n.CalendarID, nullableInt(n.EraID),
		n.Year, n.MonthIndex, n.Day,
		nullablePtrInt(n.EndYear), nullablePtrInt(n.EndMonthIndex), nullablePtrInt(n.EndDay),
		n.Importance, n.Visibility)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()

	if err := syncEventParticipants(ctx, tx, id, n.Participants); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.EventByID(ctx, id, nil)
}

func (s *Store) UpdateEvent(ctx context.Context, id int64, n NewEvent) (*Event, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	body := string(n.Body)
	if body == "" {
		body = "{}"
	}
	if _, err := tx.ExecContext(ctx,
		`UPDATE events SET title = ?, body = ?, calendar_id = ?, era_id = ?, year = ?, month_index = ?, day = ?,
		   end_year = ?, end_month_index = ?, end_day = ?, importance = ?, visibility = ?
		 WHERE id = ?`,
		n.Title, body, n.CalendarID, nullableInt(n.EraID),
		n.Year, n.MonthIndex, n.Day,
		nullablePtrInt(n.EndYear), nullablePtrInt(n.EndMonthIndex), nullablePtrInt(n.EndDay),
		n.Importance, n.Visibility, id); err != nil {
		return nil, err
	}

	if err := syncEventParticipants(ctx, tx, id, n.Participants); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.EventByID(ctx, id, nil)
}

func (s *Store) DeleteEvent(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM events WHERE id = ?`, id)
	return err
}

func (s *Store) EventsForEntity(ctx context.Context, entityID int64, visibility []string) ([]Event, error) {
	vClause := ""
	args := []any{entityID}
	if len(visibility) > 0 {
		placeholders := strings.TrimRight(strings.Repeat("?,", len(visibility)), ",")
		vClause = " AND ev.visibility IN (" + placeholders + ")"
		for _, v := range visibility {
			args = append(args, v)
		}
	}
	rows, err := s.db.QueryContext(ctx,
		`SELECT DISTINCT ev.id, ev.title, ev.body, ev.calendar_id, ev.era_id, ev.year, ev.month_index, ev.day,
		    ev.end_year, ev.end_month_index, ev.end_day, ev.importance, ev.visibility
		 FROM events ev
		 JOIN event_entities ee ON ee.event_id = ev.id
		 WHERE ee.entity_id = ?`+vClause+`
		 ORDER BY ev.year, ev.month_index, ev.day`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Event, 0)
	for rows.Next() {
		var e Event
		var bodyText string
		var eraID, endY, endM, endD sql.NullInt64
		if err := rows.Scan(&e.ID, &e.Title, &bodyText, &e.CalendarID, &eraID,
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
		out = append(out, e)
	}
	return out, rows.Err()
}

func syncEventParticipants(ctx context.Context, tx *sql.Tx, eventID int64, parts []NewEventParticipant) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM event_entities WHERE event_id = ?`, eventID); err != nil {
		return err
	}
	for _, p := range parts {
		role := strings.TrimSpace(p.Role)
		if role == "" {
			role = "participant"
		}
		if _, err := tx.ExecContext(ctx,
			`INSERT OR IGNORE INTO event_entities (event_id, entity_id, role) VALUES (?, ?, ?)`,
			eventID, p.EntityID, role); err != nil {
			return err
		}
	}
	return nil
}

func nullablePtrInt(p *int) any {
	if p == nil {
		return nil
	}
	return *p
}
