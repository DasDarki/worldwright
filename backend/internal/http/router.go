package http

import (
	"github.com/gofiber/fiber/v2"

	"worldwright/backend/internal/admin"
	"worldwright/backend/internal/auth"
	"worldwright/backend/internal/store"
)

func RegisterRoutes(app *fiber.App, st *store.Store, authSvc *auth.Service, oauthSvc *auth.OAuth, adminSvc *admin.Service, assetsDir string) {
	app.Get("/api/healthz", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"ok": true}) })

	api := app.Group("/api")
	api.Use(authSvc.Middleware(false))

	api.Get("/share/entities/:slug", shareEntity(st))
	api.Get("/share/maps/:id", shareMap(st))

	api.Post("/auth/login", login(authSvc))
	api.Post("/auth/logout", logout(authSvc))
	api.Get("/auth/providers", providers(oauthSvc))
	api.Get("/auth/oauth/:provider/start", oauthStart(authSvc, oauthSvc))
	api.Get("/auth/oauth/:provider/callback", oauthCallback(authSvc, oauthSvc))

	// Read-only endpoints: anonymous callers see public entries only, players
	// see public+player, admins see everything. Visibility is enforced inside
	// each handler via auth.VisibilityFor(c).
	api.Get("/recent", recentEntities(st))
	api.Get("/entity-types", listEntityTypes(st))
	api.Get("/entities", listEntities(st))
	api.Get("/entities/by-slug/:slug", getEntityBySlug(st))
	api.Get("/entities/:id/backlinks", entityBacklinks(st))
	api.Get("/entities/:id/relationships", listRelationships(st))
	api.Get("/entities/:id/genealogy", entityGenealogy(st))
	api.Get("/entities/:id/events", entityEvents(st))
	api.Post("/entities/relationship-graph", relationshipGraph(st))
	api.Get("/relationship-types", listRelationshipTypes(st))
	api.Get("/calendars", listCalendars(st))
	api.Get("/calendars/:id", getCalendar(st))
	api.Get("/events", listEvents(st))
	api.Get("/events/:id", getEvent(st))
	api.Get("/timelines", listTimelines(st))
	api.Get("/timelines/:id", getTimeline(st))
	api.Get("/maps", listMaps(st))
	api.Get("/maps/:id", getMap(st))
	api.Get("/assets/:id", serveAsset(st, assetsDir))
	api.Get("/tags", listTags(st))

	authed := api.Group("", requireAuth)
	authed.Get("/auth/me", me)
	authed.Patch("/auth/me", updateMe(authSvc))
	authed.Post("/auth/me/password", changePassword(authSvc))
	authed.Get("/auth/me/identities", listIdentities(authSvc))

	authed.Post("/entities", requireAdmin, createEntity(st))
	authed.Post("/entities/reorder", requireAdmin, reorderEntities(st))
	authed.Patch("/entities/:id", requireAdmin, updateEntity(st))
	authed.Delete("/entities/:id", requireAdmin, deleteEntity(st))

	authed.Post("/relationships", requireAdmin, createRelationship(st))
	authed.Delete("/relationships/:id", requireAdmin, deleteRelationship(st))

	authed.Patch("/calendars/:id", requireAdmin, updateCalendar(st))

	authed.Post("/events", requireAdmin, createEvent(st))
	authed.Patch("/events/:id", requireAdmin, updateEvent(st))
	authed.Delete("/events/:id", requireAdmin, deleteEvent(st))

	authed.Post("/timelines", requireAdmin, createTimeline(st))
	authed.Patch("/timelines/:id", requireAdmin, updateTimeline(st))
	authed.Delete("/timelines/:id", requireAdmin, deleteTimeline(st))

	authed.Get("/assets", listAssets(st))
	authed.Post("/assets", requireAdmin, uploadAsset(st, assetsDir))
	authed.Delete("/assets/:id", requireAdmin, deleteAsset(st, assetsDir))

	authed.Post("/maps", requireAdmin, createMap(st))
	authed.Patch("/maps/:id", requireAdmin, updateMap(st))
	authed.Delete("/maps/:id", requireAdmin, deleteMap(st))
	authed.Post("/maps/:id/pins", requireAdmin, createMapPin(st))
	authed.Patch("/maps/:id/pins/:pin_id", requireAdmin, updateMapPin(st))
	authed.Delete("/maps/:id/pins/:pin_id", requireAdmin, deleteMapPin(st))

	authed.Get("/search", searchHandler(st))

	authed.Get("/admin/onboarding", requireAdmin, adminStatus(adminSvc))
	authed.Post("/admin/onboarding/keep-seed", requireAdmin, adminKeepSeed(adminSvc))
	authed.Post("/admin/onboarding/prune-seed", requireAdmin, adminPruneSeed(adminSvc))
	authed.Post("/admin/onboarding/import/legendkeeper", requireAdmin, adminImportLegendKeeper(adminSvc))
	authed.Post("/admin/seed/prune", requireAdmin, adminPruneSeed(adminSvc))
}

func requireAuth(c *fiber.Ctx) error {
	if auth.UserFrom(c) == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "authentication required")
	}
	return c.Next()
}

func requireAdmin(c *fiber.Ctx) error {
	u := auth.UserFrom(c)
	if u == nil || u.Role != "admin" {
		return fiber.NewError(fiber.StatusForbidden, "forbidden")
	}
	return c.Next()
}
