package http

import (
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"worldwright/backend/internal/auth"
	"worldwright/backend/internal/store"
)

type createShareTokenRequest struct {
	TTLSeconds int64 `json:"ttl_seconds"`
}

const (
	defaultShareTTL = 24 * time.Hour
	maxShareTTL     = 30 * 24 * time.Hour
	minShareTTL     = 5 * time.Minute
)

func createShareToken(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil || id <= 0 {
			return fiber.NewError(fiber.StatusBadRequest, "invalid entity id")
		}
		// make sure the entity exists (admins can see anything)
		if _, err := st.EntityByID(c.UserContext(), id, []string{"secret", "player", "public"}); err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return fiber.NewError(fiber.StatusNotFound, "entity not found")
			}
			return err
		}

		var req createShareTokenRequest
		_ = c.BodyParser(&req) // body is optional, default TTL applies
		ttl := defaultShareTTL
		if req.TTLSeconds > 0 {
			ttl = time.Duration(req.TTLSeconds) * time.Second
		}
		if ttl < minShareTTL {
			ttl = minShareTTL
		}
		if ttl > maxShareTTL {
			ttl = maxShareTTL
		}

		user := auth.UserFrom(c)
		var byID int64
		if user != nil {
			byID = user.ID
		}
		token, err := st.CreateShareToken(c.UserContext(), id, byID, ttl)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"share_token": token})
	}
}

func listShareTokens(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil || id <= 0 {
			return fiber.NewError(fiber.StatusBadRequest, "invalid entity id")
		}
		tokens, err := st.ListShareTokensFor(c.UserContext(), id)
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"share_tokens": tokens})
	}
}

func revokeShareToken(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Params("token")
		if token == "" {
			return fiber.NewError(fiber.StatusBadRequest, "missing token")
		}
		if err := st.DeleteShareToken(c.UserContext(), token); err != nil {
			return err
		}
		return c.JSON(fiber.Map{"ok": true})
	}
}
