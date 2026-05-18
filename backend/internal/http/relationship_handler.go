package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"worldwright/backend/internal/auth"
	"worldwright/backend/internal/store"
)

type relationshipRequest struct {
	FromEntityID       int64  `json:"from_entity_id"`
	ToEntityID         int64  `json:"to_entity_id"`
	RelationshipTypeID int64  `json:"relationship_type_id"`
	Description        string `json:"description"`
}

func listRelationshipTypes(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		types, err := st.ListRelationshipTypes(c.UserContext())
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"relationship_types": types})
	}
}

func listRelationships(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		edges, err := st.ListRelationshipsForEntity(c.UserContext(), id, auth.VisibilityFor(c))
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"relationships": edges})
	}
}

func createRelationship(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req relationshipRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		if req.FromEntityID == 0 || req.ToEntityID == 0 || req.RelationshipTypeID == 0 {
			return fiber.NewError(fiber.StatusBadRequest, "from, to and type are required")
		}
		if req.FromEntityID == req.ToEntityID {
			return fiber.NewError(fiber.StatusBadRequest, "from and to must differ")
		}
		exists, err := st.RelationshipExists(c.UserContext(), req.FromEntityID, req.ToEntityID, req.RelationshipTypeID)
		if err != nil {
			return err
		}
		if exists {
			return fiber.NewError(fiber.StatusConflict, "relationship already exists")
		}
		rel, err := st.CreateRelationship(c.UserContext(), store.NewRelationship{
			FromEntityID:       req.FromEntityID,
			ToEntityID:         req.ToEntityID,
			RelationshipTypeID: req.RelationshipTypeID,
			Description:        req.Description,
		})
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"relationship": rel})
	}
}

func deleteRelationship(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		if err := st.DeleteRelationship(c.UserContext(), id); err != nil {
			return err
		}
		return c.JSON(fiber.Map{"ok": true})
	}
}
