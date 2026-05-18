package store

import (
	"context"
	"database/sql"
	"errors"
)

func (s *Store) ListCalendars(ctx context.Context) ([]Calendar, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, COALESCE(epoch_name,''), current_year, current_month_index, current_day, COALESCE(leap_rule,'')
		 FROM calendars ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Calendar, 0)
	for rows.Next() {
		var c Calendar
		if err := rows.Scan(&c.ID, &c.Name, &c.EpochName, &c.CurrentYear, &c.CurrentMonthIndex, &c.CurrentDay, &c.LeapRule); err != nil {
			return nil, err
		}
		c.Months = []CalendarMonth{}
		c.Weekdays = []CalendarWeekday{}
		c.Eras = []CalendarEra{}
		c.Moons = []CalendarMoon{}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (s *Store) CalendarByID(ctx context.Context, id int64) (*Calendar, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, name, COALESCE(epoch_name,''), current_year, current_month_index, current_day, COALESCE(leap_rule,'')
		 FROM calendars WHERE id = ?`, id)
	var c Calendar
	if err := row.Scan(&c.ID, &c.Name, &c.EpochName, &c.CurrentYear, &c.CurrentMonthIndex, &c.CurrentDay, &c.LeapRule); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	c.Months = []CalendarMonth{}
	c.Weekdays = []CalendarWeekday{}
	c.Eras = []CalendarEra{}
	c.Moons = []CalendarMoon{}

	monthRows, err := s.db.QueryContext(ctx,
		`SELECT id, calendar_id, sort_order, name, days FROM calendar_months
		 WHERE calendar_id = ? ORDER BY sort_order`, id)
	if err != nil {
		return nil, err
	}
	defer monthRows.Close()
	for monthRows.Next() {
		var m CalendarMonth
		if err := monthRows.Scan(&m.ID, &m.CalendarID, &m.SortOrder, &m.Name, &m.Days); err != nil {
			return nil, err
		}
		c.Months = append(c.Months, m)
	}

	wRows, err := s.db.QueryContext(ctx,
		`SELECT id, calendar_id, sort_order, name FROM calendar_weekdays
		 WHERE calendar_id = ? ORDER BY sort_order`, id)
	if err != nil {
		return nil, err
	}
	defer wRows.Close()
	for wRows.Next() {
		var w CalendarWeekday
		if err := wRows.Scan(&w.ID, &w.CalendarID, &w.SortOrder, &w.Name); err != nil {
			return nil, err
		}
		c.Weekdays = append(c.Weekdays, w)
	}

	eraRows, err := s.db.QueryContext(ctx,
		`SELECT id, calendar_id, name, start_year, COALESCE(suffix,'') FROM calendar_eras
		 WHERE calendar_id = ? ORDER BY start_year`, id)
	if err != nil {
		return nil, err
	}
	defer eraRows.Close()
	for eraRows.Next() {
		var e CalendarEra
		if err := eraRows.Scan(&e.ID, &e.CalendarID, &e.Name, &e.StartYear, &e.Suffix); err != nil {
			return nil, err
		}
		c.Eras = append(c.Eras, e)
	}

	moonRows, err := s.db.QueryContext(ctx,
		`SELECT id, calendar_id, name, cycle_days, offset_days FROM calendar_moons
		 WHERE calendar_id = ? ORDER BY cycle_days`, id)
	if err != nil {
		return nil, err
	}
	defer moonRows.Close()
	moonByID := map[int64]int{}
	for moonRows.Next() {
		var m CalendarMoon
		if err := moonRows.Scan(&m.ID, &m.CalendarID, &m.Name, &m.CycleDays, &m.OffsetDays); err != nil {
			return nil, err
		}
		m.Phases = []CalendarMoonPhase{}
		moonByID[m.ID] = len(c.Moons)
		c.Moons = append(c.Moons, m)
	}

	if len(c.Moons) > 0 {
		phaseRows, err := s.db.QueryContext(ctx,
			`SELECT p.id, p.moon_id, p.sort_order, p.name, p.cycle_position, p.randomness, COALESCE(p.icon, '')
			 FROM calendar_moon_phases p
			 JOIN calendar_moons m ON m.id = p.moon_id
			 WHERE m.calendar_id = ?
			 ORDER BY p.moon_id, p.sort_order, p.cycle_position`, id)
		if err != nil {
			return nil, err
		}
		defer phaseRows.Close()
		for phaseRows.Next() {
			var p CalendarMoonPhase
			if err := phaseRows.Scan(&p.ID, &p.MoonID, &p.SortOrder, &p.Name, &p.CyclePosition, &p.Randomness, &p.Icon); err != nil {
				return nil, err
			}
			if idx, ok := moonByID[p.MoonID]; ok {
				c.Moons[idx].Phases = append(c.Moons[idx].Phases, p)
			}
		}
	}

	return &c, nil
}

func (s *Store) UpdateCalendar(ctx context.Context, c Calendar) (*Calendar, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx,
		`UPDATE calendars
		 SET name = ?, epoch_name = ?, current_year = ?, current_month_index = ?, current_day = ?, leap_rule = ?
		 WHERE id = ?`,
		c.Name, nullable(c.EpochName), c.CurrentYear, c.CurrentMonthIndex, c.CurrentDay, nullable(c.LeapRule), c.ID); err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM calendar_months WHERE calendar_id = ?`, c.ID); err != nil {
		return nil, err
	}
	for i, m := range c.Months {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO calendar_months (calendar_id, sort_order, name, days) VALUES (?, ?, ?, ?)`,
			c.ID, i+1, m.Name, m.Days); err != nil {
			return nil, err
		}
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM calendar_weekdays WHERE calendar_id = ?`, c.ID); err != nil {
		return nil, err
	}
	for i, w := range c.Weekdays {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO calendar_weekdays (calendar_id, sort_order, name) VALUES (?, ?, ?)`,
			c.ID, i+1, w.Name); err != nil {
			return nil, err
		}
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM calendar_eras WHERE calendar_id = ?`, c.ID); err != nil {
		return nil, err
	}
	for _, e := range c.Eras {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO calendar_eras (calendar_id, name, start_year, suffix) VALUES (?, ?, ?, ?)`,
			c.ID, e.Name, e.StartYear, nullable(e.Suffix)); err != nil {
			return nil, err
		}
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM calendar_moons WHERE calendar_id = ?`, c.ID); err != nil {
		return nil, err
	}
	for _, m := range c.Moons {
		res, err := tx.ExecContext(ctx,
			`INSERT INTO calendar_moons (calendar_id, name, cycle_days, offset_days) VALUES (?, ?, ?, ?)`,
			c.ID, m.Name, m.CycleDays, m.OffsetDays)
		if err != nil {
			return nil, err
		}
		moonID, _ := res.LastInsertId()
		for i, p := range m.Phases {
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO calendar_moon_phases (moon_id, sort_order, name, cycle_position, randomness, icon)
				 VALUES (?, ?, ?, ?, ?, ?)`,
				moonID, i+1, p.Name, p.CyclePosition, p.Randomness, nullable(p.Icon)); err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.CalendarByID(ctx, c.ID)
}
