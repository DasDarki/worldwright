package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"worldwright/backend/internal/content"
	"worldwright/backend/internal/store"
)

func shareEntity(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		slug := c.Params("slug")
		token := c.Query("token")

		// A valid token grants access to one specific entity regardless of
		// its visibility. The entity must match the slug in the URL — that
		// prevents someone with a token for entity A from using it to view
		// arbitrary other entities by swapping slugs.
		bypass := false
		if token != "" {
			if eid, err := st.ResolveShareToken(c.UserContext(), token); err == nil {
				if tEntity, terr := st.EntityByID(c.UserContext(), eid, []string{"secret", "player", "public"}); terr == nil {
					if tEntity.Slug == slug {
						bypass = true
					}
				}
			}
		}

		viewVisibility := []string{"public"}
		if bypass {
			viewVisibility = []string{"secret", "player", "public"}
		}

		e, err := st.EntityBySlug(c.UserContext(), slug, viewVisibility)
		if err != nil {
			return err
		}
		if e.EntityTypeID != 0 {
			t, _ := st.EntityTypeByID(c.UserContext(), e.EntityTypeID)
			e.EntityType = t
		}
		// Wikilink scrubbing always uses public visibility: the token grants
		// access to one entity, not transitively to every entity it mentions.
		if visible, err := content.VisibleSlugSet(c.UserContext(), st.DB(), []string{"public"}); err == nil {
			e.Body = content.ScrubInvisibleWikilinks(e.Body, visible)
		}
		// Always strip secret-vault content for share-view callers: even with
		// a valid token, the reader is anonymous and isn't the vault's author.
		e.Body = content.ScrubSecretVaults(e.Body, 0)
		bls, _ := st.Backlinks(c.UserContext(), e.ID, []string{"public"})
		return c.JSON(fiber.Map{"entity": e, "backlinks": bls})
	}
}

func shareMap(st *store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		m, err := st.MapByID(c.UserContext(), id, []string{"public"})
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"map": m})
	}
}
