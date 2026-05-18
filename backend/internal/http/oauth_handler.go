package http

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"worldwright/backend/internal/auth"
)

func providers(o *auth.OAuth) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"password": true,
			"oauth":    o.Enabled(),
		})
	}
}

func oauthStart(svc *auth.Service, o *auth.OAuth) fiber.Handler {
	return func(c *fiber.Ctx) error {
		provider := c.Params("provider")
		if !o.HasProvider(provider) {
			return fiber.NewError(fiber.StatusNotFound, "provider not enabled")
		}
		authURL, state, verifier, err := o.Start(provider)
		if err != nil {
			return err
		}
		expires := time.Now().Add(auth.OAuthStateTTL)
		c.Cookie(&fiber.Cookie{
			Name: auth.OAuthStateCookie, Value: state,
			HTTPOnly: true, Secure: svc.CookieSecure(), SameSite: "Lax",
			Path: "/api/auth/oauth", Expires: expires,
		})
		c.Cookie(&fiber.Cookie{
			Name: auth.OAuthVerifierCookie, Value: verifier,
			HTTPOnly: true, Secure: svc.CookieSecure(), SameSite: "Lax",
			Path: "/api/auth/oauth", Expires: expires,
		})
		return c.Redirect(authURL, fiber.StatusFound)
	}
}

func oauthCallback(svc *auth.Service, o *auth.OAuth) fiber.Handler {
	return func(c *fiber.Ctx) error {
		provider := c.Params("provider")
		if !o.HasProvider(provider) {
			return fiber.NewError(fiber.StatusNotFound, "provider not enabled")
		}
		code := c.Query("code")
		state := c.Query("state")
		if code == "" || state == "" {
			return fiber.NewError(fiber.StatusBadRequest, "missing code or state")
		}
		cookieState := c.Cookies(auth.OAuthStateCookie)
		verifier := c.Cookies(auth.OAuthVerifierCookie)
		if cookieState == "" || verifier == "" {
			return fiber.NewError(fiber.StatusBadRequest, "missing state cookie")
		}
		if cookieState != state {
			return fiber.NewError(fiber.StatusForbidden, "state mismatch")
		}
		info, err := o.Callback(c.UserContext(), provider, code, verifier)
		if err != nil {
			return err
		}
		user, err := o.ResolveUser(c.UserContext(), provider, info)
		if err != nil {
			return err
		}
		token, err := svc.LoginAsUser(c.UserContext(), user)
		if err != nil {
			return err
		}
		c.ClearCookie(auth.OAuthStateCookie)
		c.ClearCookie(auth.OAuthVerifierCookie)
		setSessionCookie(c, svc, token)
		return c.Redirect(o.PublicURL()+"/", fiber.StatusFound)
	}
}
