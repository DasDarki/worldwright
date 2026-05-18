package http

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"worldwright/backend/internal/auth"
	"worldwright/backend/internal/content"
	"worldwright/backend/internal/store"
)

type eventParticipantRequest struct {
	EntityID int64  `json:"entity_id"`
	Role     string `json:"role"`
}

type eventRequest struct {
	Title         string                    `json:"title"`
	Body          json.RawMessage           `json:"body"`
	CalendarID    int64                     `json:"calendar_id"`
	EraID         *int64                    `json:"era_id"`
	Year          int                       `json:"year"`
	MonthIndex    int                       `json:"month_index"`
	Day           int                       `json:"day"`
	EndYear       *int                      `json:"end_year"`
	EndMonthIndex *int                      `json:"end_month_index"`
	EndDay        *int                      `json:"end_day"`
	Importance    int                       `json:"importance"`
	Visibility    string                    `json:"visibility"`
	Participants  []eventParticipantRequest `json:"participants"`
}

func listEvents(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		filter := store.EventFilter{Visibility: auth.VisibilityFor(c)}
		if v := c.Query("calendar_id"); v != "" {
			id, _ := strconv.ParseInt(v, 10, 64)
			filter.CalendarID = id
		}
		events, err := st.ListEvents(c.UserContext(), filter)
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"events": events})
	}
}

func getEvent(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		ev, err := st.EventByID(c.UserContext(), id, auth.VisibilityFor(c))
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"event": ev})
	}
}

func createEvent(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req eventRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		if req.Title == "" || req.CalendarID == 0 {
			return fiber.NewError(fiber.StatusBadRequest, "title and calendar_id are required")
		}
		if req.Visibility == "" {
			req.Visibility = "secret"
		}
		bodyText, _, _ := content.ParseBody(req.Body)
		participants := mapParticipants(req.Participants)

		ev, err := st.CreateEvent(c.UserContext(), store.NewEvent{
			Title: req.Title, Body: req.Body, BodyText: bodyText,
			CalendarID: req.CalendarID, EraID: req.EraID,
			Year: req.Year, MonthIndex: req.MonthIndex, Day: req.Day,
			EndYear: req.EndYear, EndMonthIndex: req.EndMonthIndex, EndDay: req.EndDay,
			Importance: req.Importance, Visibility: req.Visibility,
			Participants: participants,
		})
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"event": ev})
	}
}

func updateEvent(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		var req eventRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		bodyText, _, _ := content.ParseBody(req.Body)
		participants := mapParticipants(req.Participants)
		ev, err := st.UpdateEvent(c.UserContext(), id, store.NewEvent{
			Title: req.Title, Body: req.Body, BodyText: bodyText,
			CalendarID: req.CalendarID, EraID: req.EraID,
			Year: req.Year, MonthIndex: req.MonthIndex, Day: req.Day,
			EndYear: req.EndYear, EndMonthIndex: req.EndMonthIndex, EndDay: req.EndDay,
			Importance: req.Importance, Visibility: req.Visibility,
			Participants: participants,
		})
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"event": ev})
	}
}

func deleteEvent(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		if err := st.DeleteEvent(c.UserContext(), id); err != nil {
			return err
		}
		return c.JSON(fiber.Map{"ok": true})
	}
}

func entityEvents(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		events, err := st.EventsForEntity(c.UserContext(), id, auth.VisibilityFor(c))
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"events": events})
	}
}

func mapParticipants(in []eventParticipantRequest) []store.NewEventParticipant {
	out := make([]store.NewEventParticipant, 0, len(in))
	for _, p := range in {
		if p.EntityID == 0 {
			continue
		}
		out = append(out, store.NewEventParticipant{EntityID: p.EntityID, Role: p.Role})
	}
	return out
}
