package http

import (
	"github.com/gofiber/fiber/v2"

	"worldwright/backend/internal/store"
)

func listEntityTypes(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		types, err := st.ListEntityTypes(c.UserContext())
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"entity_types": types})
	}
}
