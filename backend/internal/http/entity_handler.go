package http

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"worldwright/backend/internal/auth"
	"worldwright/backend/internal/content"
	"worldwright/backend/internal/store"
)

type entityRequest struct {
	EntityTypeID int64             `json:"entity_type_id"`
	Title        string            `json:"title"`
	Slug         string            `json:"slug"`
	Summary      string            `json:"summary"`
	Body         json.RawMessage   `json:"body"`
	ParentID     *int64            `json:"parent_id"`
	Visibility   string            `json:"visibility"`
	Tags         []string          `json:"tags"`
	FieldValues  map[string]string `json:"field_values"`
}

func listEntities(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		filter := store.EntityFilter{
			Visibility: auth.VisibilityFor(c),
			Tag:        c.Query("tag"),
		}
		if v := c.Query("entity_type_id"); v != "" {
			id, _ := strconv.ParseInt(v, 10, 64)
			filter.EntityTypeID = id
		}
		if v := c.Query("parent_id"); v != "" {
			if v == "root" {
				zero := int64(0)
				filter.ParentID = &zero
			} else {
				id, _ := strconv.ParseInt(v, 10, 64)
				filter.ParentID = &id
			}
		}
		list, err := st.ListEntities(c.UserContext(), filter)
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"entities": list})
	}
}

func getEntityBySlug(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		slug := c.Params("slug")
		visibility := auth.VisibilityFor(c)
		e, err := st.EntityBySlug(c.UserContext(), slug, visibility)
		if err != nil {
			return err
		}
		if e.EntityTypeID != 0 {
			t, _ := st.EntityTypeByID(c.UserContext(), e.EntityTypeID)
			e.EntityType = t
		}
		// Scrub wikilinks whose target the caller cannot see, so secret
		// entries don't leak through references in pages they CAN see.
		if visible, err := content.VisibleSlugSet(c.UserContext(), st.DB(), visibility); err == nil {
			e.Body = content.ScrubInvisibleWikilinks(e.Body, visible)
		}
		return c.JSON(fiber.Map{"entity": e})
	}
}

func createEntity(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req entityRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		if req.Title == "" || req.Slug == "" || req.EntityTypeID == 0 {
			return fiber.NewError(fiber.StatusBadRequest, "missing title, slug or entity_type_id")
		}
		if req.Visibility == "" {
			req.Visibility = "secret"
		}
		body := autoLinkBody(c, st, req.Body, req.Slug)
		bodyText, slugs, err := content.ParseBody(body)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body json")
		}
		entity, err := st.CreateEntity(c.UserContext(), store.NewEntity{
			EntityTypeID: req.EntityTypeID,
			Title:        req.Title,
			Slug:         req.Slug,
			Summary:      req.Summary,
			Body:         body,
			BodyText:     bodyText,
			ParentID:     req.ParentID,
			Visibility:   req.Visibility,
			Tags:         req.Tags,
			FieldValues:  req.FieldValues,
			Wikilinks:    slugs,
		})
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"entity": entity})
	}
}

func autoLinkBody(c *fiber.Ctx, st *store.Store, body json.RawMessage, selfSlug string) json.RawMessage {
	if len(body) == 0 {
		return body
	}
	all, err := st.ListEntities(c.UserContext(), store.EntityFilter{})
	if err != nil {
		return body
	}
	cands := make([]content.TitleSlug, 0, len(all))
	for _, e := range all {
		if e.Slug == selfSlug {
			continue
		}
		cands = append(cands, content.TitleSlug{Title: e.Title, Slug: e.Slug})
	}
	out, err := content.AutoLink(body, cands, selfSlug)
	if err != nil {
		return body
	}
	return out
}

func updateEntity(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		var req entityRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		body := autoLinkBody(c, st, req.Body, req.Slug)
		bodyText, slugs, err := content.ParseBody(body)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body json")
		}
		entity, err := st.UpdateEntity(c.UserContext(), id, store.NewEntity{
			EntityTypeID: req.EntityTypeID,
			Title:        req.Title,
			Slug:         req.Slug,
			Summary:      req.Summary,
			Body:         body,
			BodyText:     bodyText,
			ParentID:     req.ParentID,
			Visibility:   req.Visibility,
			Tags:         req.Tags,
			FieldValues:  req.FieldValues,
			Wikilinks:    slugs,
		})
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"entity": entity})
	}
}

func deleteEntity(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		if err := st.DeleteEntity(c.UserContext(), id); err != nil {
			return err
		}
		return c.JSON(fiber.Map{"ok": true})
	}
}

type reorderRequest struct {
	Order []struct {
		ID        int64  `json:"id"`
		ParentID  *int64 `json:"parent_id"`
		SortOrder int    `json:"sort_order"`
	} `json:"order"`
}

type graphRequest struct {
	IDs  []int64 `json:"ids"`
	Lang string  `json:"lang"`
}

func relationshipGraph(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req graphRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		if len(req.IDs) > 200 {
			return fiber.NewError(fiber.StatusBadRequest, "too many ids (max 200)")
		}
		g, err := st.RelationshipGraphForIDs(c.UserContext(), req.IDs, auth.VisibilityFor(c), req.Lang)
		if err != nil {
			return err
		}
		return c.JSON(g)
	}
}

func reorderEntities(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req reorderRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid body")
		}
		entries := make([]store.ReorderEntry, 0, len(req.Order))
		for _, o := range req.Order {
			entries = append(entries, store.ReorderEntry{
				ID:        o.ID,
				ParentID:  o.ParentID,
				SortOrder: o.SortOrder,
			})
		}
		if err := st.ReorderEntities(c.UserContext(), entries); err != nil {
			return err
		}
		return c.JSON(fiber.Map{"ok": true})
	}
}

func entityGenealogy(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		depth, _ := strconv.Atoi(c.Query("depth", "3"))
		g, err := st.GenealogyForEntity(c.UserContext(), id, depth, auth.VisibilityFor(c))
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"genealogy": g})
	}
}

func entityBacklinks(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		bl, err := st.Backlinks(c.UserContext(), id, auth.VisibilityFor(c))
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"backlinks": bl})
	}
}

func recentEntities(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		limit, _ := strconv.Atoi(c.Query("limit", "10"))
		list, err := st.RecentEntities(c.UserContext(), auth.VisibilityFor(c), limit)
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"entities": list})
	}
}

func listTags(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Query("with_counts") == "true" {
			tags, err := st.ListTagsWithCounts(c.UserContext(), auth.VisibilityFor(c))
			if err != nil {
				return err
			}
			return c.JSON(fiber.Map{"tags": tags})
		}
		t, err := st.ListTags(c.UserContext())
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"tags": t})
	}
}
