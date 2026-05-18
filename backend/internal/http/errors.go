package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"worldwright/backend/internal/auth"
	"worldwright/backend/internal/store"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	var fe *fiber.Error
	if errors.As(err, &fe) {
		return c.Status(fe.Code).JSON(fiber.Map{"error": fe.Message})
	}
	switch {
	case errors.Is(err, store.ErrNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	case errors.Is(err, auth.ErrInvalidCredentials):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	case errors.Is(err, auth.ErrSessionExpired):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "session expired"})
	case errors.Is(err, auth.ErrProviderUnknown):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "provider not enabled"})
	case errors.Is(err, auth.ErrSignupDisabled):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "signup disabled"})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
}
