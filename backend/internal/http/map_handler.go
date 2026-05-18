package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"worldwright/backend/internal/auth"
	"worldwright/backend/internal/store"
)

type mapRequest struct {
	Name           string `json:"name"`
	AssetID        int64  `json:"asset_id"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	ParentEntityID *int64 `json:"parent_entity_id"`
}

type mapPinRequest struct {
	X              float64 `json:"x"`
	Y              float64 `json:"y"`
	Label          string  `json:"label"`
	Icon           string  `json:"icon"`
	TargetEntityID *int64  `json:"target_entity_id"`
	TargetMapID    *int64  `json:"target_map_id"`
	Visibility     string  `json:"visibility"`
}

func listMaps(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		maps, err := st.ListMaps(c.UserContext())
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"maps": maps})
	}
}

func getMap(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		m, err := st.MapByID(c.UserContext(), id, auth.VisibilityFor(c))
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"map": m})
	}
}

func createMap(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req mapRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		if req.Name == "" || req.AssetID == 0 {
			return fiber.NewError(fiber.StatusBadRequest, "name and asset_id are required")
		}
		m, err := st.CreateMap(c.UserContext(), store.NewMap{
			Name: req.Name, AssetID: req.AssetID,
			Width: req.Width, Height: req.Height,
			ParentEntityID: req.ParentEntityID,
		})
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"map": m})
	}
}

func updateMap(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		var req mapRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		m, err := st.UpdateMap(c.UserContext(), id, store.NewMap{
			Name: req.Name, AssetID: req.AssetID,
			Width: req.Width, Height: req.Height,
			ParentEntityID: req.ParentEntityID,
		})
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"map": m})
	}
}

func deleteMap(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		if err := st.DeleteMap(c.UserContext(), id); err != nil {
			return err
		}
		return c.JSON(fiber.Map{"ok": true})
	}
}

func createMapPin(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		mapID, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		var req mapPinRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		if req.Visibility == "" {
			req.Visibility = "secret"
		}
		pin, err := st.CreateMapPin(c.UserContext(), mapID, store.NewMapPin{
			X: req.X, Y: req.Y, Label: req.Label, Icon: req.Icon,
			TargetEntityID: req.TargetEntityID, TargetMapID: req.TargetMapID,
			Visibility: req.Visibility,
		})
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"pin": pin})
	}
}

func updateMapPin(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseInt(c.Params("pin_id"), 10, 64)
		var req mapPinRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		if req.Visibility == "" {
			req.Visibility = "secret"
		}
		pin, err := st.UpdateMapPin(c.UserContext(), id, store.NewMapPin{
			X: req.X, Y: req.Y, Label: req.Label, Icon: req.Icon,
			TargetEntityID: req.TargetEntityID, TargetMapID: req.TargetMapID,
			Visibility: req.Visibility,
		})
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"pin": pin})
	}
}

func deleteMapPin(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseInt(c.Params("pin_id"), 10, 64)
		if err := st.DeleteMapPin(c.UserContext(), id); err != nil {
			return err
		}
		return c.JSON(fiber.Map{"ok": true})
	}
}
