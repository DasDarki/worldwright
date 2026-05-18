package http

import (
	"github.com/gofiber/fiber/v2"

	"worldwright/backend/internal/admin"
)

func adminStatus(svc *admin.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		s, err := svc.Status(c.UserContext())
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"onboarding": s})
	}
}

func adminKeepSeed(svc *admin.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := svc.KeepSeed(c.UserContext()); err != nil {
			return err
		}
		return c.JSON(fiber.Map{"ok": true})
	}
}

func adminPruneSeed(svc *admin.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := svc.PruneSeed(c.UserContext()); err != nil {
			return err
		}
		return c.JSON(fiber.Map{"ok": true})
	}
}

func adminImportLegendKeeper(svc *admin.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fh, err := c.FormFile("file")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "missing file field")
		}
		if fh.Size > 200*1024*1024 {
			return fiber.NewError(fiber.StatusRequestEntityTooLarge, "file too large")
		}
		src, err := fh.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		summary, err := svc.ImportLegendKeeper(c.UserContext(), src)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		if err := svc.CompleteOnboarding(c.UserContext()); err != nil {
			return err
		}
		return c.JSON(fiber.Map{"summary": summary})
	}
}
