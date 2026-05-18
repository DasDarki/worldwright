package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"worldwright/backend/internal/auth"
	"worldwright/backend/internal/store"
)

func searchHandler(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		q := c.Query("q")
		limit, _ := strconv.Atoi(c.Query("limit", "25"))
		hits, err := st.Search(c.UserContext(), q, auth.VisibilityFor(c), limit)
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"hits": hits})
	}
}
