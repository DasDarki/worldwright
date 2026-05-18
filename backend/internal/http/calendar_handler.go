package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"worldwright/backend/internal/store"
)

func listCalendars(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cals, err := st.ListCalendars(c.UserContext())
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"calendars": cals})
	}
}

func getCalendar(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		cal, err := st.CalendarByID(c.UserContext(), id)
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"calendar": cal})
	}
}

func updateCalendar(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		var req store.Calendar
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		req.ID = id
		cal, err := st.UpdateCalendar(c.UserContext(), req)
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"calendar": cal})
	}
}
