package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"worldwright/backend/internal/store"
)

func shareEntity(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		slug := c.Params("slug")
		e, err := st.EntityBySlug(c.UserContext(), slug, []string{"public"})
		if err != nil {
			return err
		}
		if e.EntityTypeID != 0 {
			t, _ := st.EntityTypeByID(c.UserContext(), e.EntityTypeID)
			e.EntityType = t
		}
		bls, _ := st.Backlinks(c.UserContext(), e.ID, []string{"public"})
		return c.JSON(fiber.Map{"entity": e, "backlinks": bls})
	}
}

func shareMap(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		m, err := st.MapByID(c.UserContext(), id, []string{"public"})
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"map": m})
	}
}
