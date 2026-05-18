package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"worldwright/backend/internal/auth"
	"worldwright/backend/internal/store"
)

type timelineRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	EventIDs    []int64 `json:"event_ids"`
}

func listTimelines(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tls, err := st.ListTimelines(c.UserContext())
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"timelines": tls})
	}
}

func getTimeline(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		tl, err := st.TimelineByID(c.UserContext(), id, auth.VisibilityFor(c))
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"timeline": tl})
	}
}

func createTimeline(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req timelineRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		if req.Name == "" {
			return fiber.NewError(fiber.StatusBadRequest, "name is required")
		}
		tl, err := st.CreateTimeline(c.UserContext(), store.NewTimeline{
			Name:        req.Name,
			Description: req.Description,
			EventIDs:    req.EventIDs,
		})
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"timeline": tl})
	}
}

func updateTimeline(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		var req timelineRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		tl, err := st.UpdateTimeline(c.UserContext(), id, store.NewTimeline{
			Name:        req.Name,
			Description: req.Description,
			EventIDs:    req.EventIDs,
		})
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"timeline": tl})
	}
}

func deleteTimeline(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		if err := st.DeleteTimeline(c.UserContext(), id); err != nil {
			return err
		}
		return c.JSON(fiber.Map{"ok": true})
	}
}
