package auth

import (
	"github.com/gofiber/fiber/v2"

	"worldwright/backend/internal/store"
)

const userContextKey = "worldwright_user"

func (s *Service) Middleware(required bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies(SessionCookieName)
		user, err := s.UserBySession(c.UserContext(), token)
		if err == nil {
			c.Locals(userContextKey, user)
		} else if required {
			return fiber.NewError(fiber.StatusUnauthorized, "authentication required")
		}
		return c.Next()
	}
}

func UserFrom(c *fiber.Ctx) *store.User {
	v := c.Locals(userContextKey)
	if v == nil {
		return nil
	}
	u, _ := v.(*store.User)
	return u
}

func VisibilityFor(c *fiber.Ctx) []string {
	user := UserFrom(c)
	if user == nil {
		return []string{"public"}
	}
	if user.Role == "admin" {
		return []string{"secret", "player", "public"}
	}
	return []string{"player", "public"}
}
