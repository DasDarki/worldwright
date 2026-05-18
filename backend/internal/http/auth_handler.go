package http

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"worldwright/backend/internal/auth"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func login(svc *auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req loginRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		token, user, err := svc.Login(c.UserContext(), req.Email, req.Password)
		if err != nil {
			return err
		}
		setSessionCookie(c, svc, token)
		return c.JSON(fiber.Map{"user": user})
	}
}

func logout(svc *auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies(auth.SessionCookieName)
		_ = svc.Logout(c.UserContext(), token)
		clearSessionCookie(c, svc)
		return c.JSON(fiber.Map{"ok": true})
	}
}

func me(c *fiber.Ctx) error {
	u := auth.UserFrom(c)
	if u == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "not authenticated")
	}
	return c.JSON(fiber.Map{"user": u})
}

type updateProfileRequest struct {
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	Locale      string `json:"locale"`
}

func updateMe(svc *auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		u := auth.UserFrom(c)
		if u == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "not authenticated")
		}
		var req updateProfileRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		locale := req.Locale
		if locale == "" {
			locale = u.Locale
		}
		if err := svc.Store().UpdateUserProfile(c.UserContext(), u.ID, req.DisplayName, req.AvatarURL, locale); err != nil {
			return err
		}
		updated, err := svc.Store().UserByID(c.UserContext(), u.ID)
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"user": updated})
	}
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func changePassword(svc *auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		u := auth.UserFrom(c)
		if u == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "not authenticated")
		}
		var req changePasswordRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		if len(req.NewPassword) < 6 {
			return fiber.NewError(fiber.StatusBadRequest, "new password must be at least 6 characters")
		}
		if err := svc.ChangePassword(c.UserContext(), u.ID, req.CurrentPassword, req.NewPassword); err != nil {
			return err
		}
		return c.JSON(fiber.Map{"ok": true})
	}
}

func listIdentities(svc *auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		u := auth.UserFrom(c)
		if u == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "not authenticated")
		}
		ids, err := svc.Store().ListOAuthIdentities(c.UserContext(), u.ID)
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"identities": ids})
	}
}

func setSessionCookie(c *fiber.Ctx, svc *auth.Service, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     auth.SessionCookieName,
		Value:    token,
		HTTPOnly: true,
		Secure:   svc.CookieSecure(),
		SameSite: "Lax",
		Path:     "/",
		Expires:  time.Now().Add(auth.SessionTTL),
	})
}

func clearSessionCookie(c *fiber.Ctx, svc *auth.Service) {
	c.Cookie(&fiber.Cookie{
		Name:     auth.SessionCookieName,
		Value:    "",
		HTTPOnly: true,
		Secure:   svc.CookieSecure(),
		SameSite: "Lax",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		MaxAge:   -1,
	})
}
