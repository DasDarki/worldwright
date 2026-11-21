package http

import (
	"context"
	"errors"
	"log/slog"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"

	"worldwright/backend/internal/auth"
	"worldwright/backend/internal/onyx"
	"worldwright/backend/internal/store"
)

const (
	onyxProvider    = "onyx"
	onyxStateCookie = "worldwright_onyx_state"
	onyxStatePath   = "/api/auth/onyx"
	onyxStateTTL    = 10 * time.Minute
	onyxClaimTTL    = 15 * time.Minute
	onyxModeLogin   = "login"
	onyxModeLink    = "link"
)

var (
	errOnyxTaken         = errors.New("onyx account already linked to another user")
	errOnyxOtherIdentity = errors.New("user already linked to another onyx account")
)

type onyxFlow struct {
	mode     string
	verifier string
	userID   int64
	redirect string
	profile  *onyx.Profile
	expires  time.Time
}

type onyxFlowCache struct {
	mutex sync.Mutex
	items map[string]onyxFlow
}

func newOnyxFlowCache() *onyxFlowCache {
	return &onyxFlowCache{items: map[string]onyxFlow{}}
}

func (c *onyxFlowCache) put(key string, flow onyxFlow) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	now := time.Now()
	for existing, item := range c.items {
		if now.After(item.expires) {
			delete(c.items, existing)
		}
	}
	c.items[key] = flow
}

func (c *onyxFlowCache) take(key string) (onyxFlow, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	flow, ok := c.items[key]
	if !ok {
		return onyxFlow{}, false
	}
	delete(c.items, key)
	if time.Now().After(flow.expires) {
		return onyxFlow{}, false
	}
	return flow, true
}

type onyxRoutes struct {
	svc       *auth.Service
	client    *onyx.Client
	flows     *onyxFlowCache
	publicURL string
}

type userView struct {
	*store.User
	OnyxEnabled bool `json:"onyx_enabled"`
	OnyxLinked  bool `json:"onyx_linked"`
}

func (o *onyxRoutes) view(ctx context.Context, user *store.User) userView {
	_, err := o.svc.Store().OAuthSubjectForUser(ctx, user.ID, onyxProvider)
	return userView{User: user, OnyxEnabled: o.client.Enabled(), OnyxLinked: err == nil}
}

func safeRedirect(value string) string {
	if !strings.HasPrefix(value, "/") || strings.HasPrefix(value, "//") || strings.HasPrefix(value, "/\\") {
		return "/"
	}
	return value
}

func (o *onyxRoutes) redirect(c *fiber.Ctx, path string, values url.Values) error {
	target := o.publicURL + path
	if len(values) > 0 {
		target += "?" + values.Encode()
	}
	return c.Redirect(target, fiber.StatusFound)
}

func (o *onyxRoutes) start(c *fiber.Ctx, flow onyxFlow) string {
	state := onyx.RandomString()
	flow.verifier = onyx.RandomString()
	flow.expires = time.Now().Add(onyxStateTTL)
	o.flows.put("state:"+state, flow)
	c.Cookie(&fiber.Cookie{
		Name: onyxStateCookie, Value: state,
		HTTPOnly: true, Secure: o.svc.CookieSecure(), SameSite: "Lax",
		Path: onyxStatePath, Expires: time.Now().Add(onyxStateTTL),
	})
	return o.client.AuthorizeURL(state, flow.verifier)
}

func (o *onyxRoutes) signIn(c *fiber.Ctx, user *store.User, redirect string) error {
	token, err := o.svc.LoginAsUser(c.UserContext(), user)
	if err != nil {
		slog.Error("onyx: create session", "err", err)
		return o.redirect(c, "/login", url.Values{"onyx": {"failed"}})
	}
	setSessionCookie(c, o.svc, token)
	return o.redirect(c, safeRedirect(redirect), nil)
}

func (o *onyxRoutes) login(c *fiber.Ctx) error {
	if !o.client.Enabled() {
		return o.redirect(c, "/login", url.Values{"onyx": {"failed"}})
	}
	target := o.start(c, onyxFlow{mode: onyxModeLogin, redirect: safeRedirect(c.Query("redirect"))})
	return c.Redirect(target, fiber.StatusFound)
}

func (o *onyxRoutes) link(c *fiber.Ctx) error {
	if !o.client.Enabled() {
		return fiber.NewError(fiber.StatusNotFound, "onyx is not configured")
	}
	user := auth.UserFrom(c)
	target := o.start(c, onyxFlow{mode: onyxModeLink, userID: user.ID})
	return c.JSON(fiber.Map{"url": target})
}

func (o *onyxRoutes) callback(c *fiber.Ctx) error {
	state := c.Query("state")
	cookie := c.Cookies(onyxStateCookie)
	c.Cookie(&fiber.Cookie{
		Name: onyxStateCookie, Value: "", HTTPOnly: true, Secure: o.svc.CookieSecure(),
		SameSite: "Lax", Path: onyxStatePath, Expires: time.Now().Add(-time.Hour), MaxAge: -1,
	})

	flow, ok := o.flows.take("state:" + state)
	if state == "" || !ok {
		return o.redirect(c, "/login", url.Values{"onyx": {"expired"}})
	}
	failurePath := "/login"
	if flow.mode == onyxModeLink {
		failurePath = "/"
	}
	if cookie != state {
		return o.redirect(c, failurePath, url.Values{"onyx": {"expired"}})
	}
	if c.Query("error") != "" || c.Query("code") == "" {
		return o.redirect(c, failurePath, url.Values{"onyx": {"cancelled"}})
	}
	profile, err := o.client.Identify(c.UserContext(), c.Query("code"), flow.verifier)
	if err != nil {
		slog.Warn("onyx: identify failed", "err", err)
		return o.redirect(c, failurePath, url.Values{"onyx": {"failed"}})
	}

	if flow.mode == onyxModeLink {
		return o.redirect(c, "/", url.Values{"onyx": {o.linkResult(c.UserContext(), flow.userID, profile)}})
	}

	user, err := o.svc.Store().UserByOAuth(c.UserContext(), onyxProvider, profile.ID)
	if errors.Is(err, store.ErrNotFound) {
		claim := onyx.RandomString()
		o.flows.put("claim:"+claim, onyxFlow{profile: profile, redirect: flow.redirect, expires: time.Now().Add(onyxClaimTTL)})
		return o.redirect(c, "/login", url.Values{"onyx_claim": {claim}, "onyx_user": {profile.Username}})
	}
	if err != nil {
		slog.Error("onyx: user lookup", "err", err)
		return o.redirect(c, "/login", url.Values{"onyx": {"failed"}})
	}
	if user.DisplayName != profile.Username {
		if err := o.svc.Store().UpdateUserDisplayName(c.UserContext(), user.ID, profile.Username); err == nil {
			user.DisplayName = profile.Username
		}
	}
	return o.signIn(c, user, flow.redirect)
}

func (o *onyxRoutes) claim(c *fiber.Ctx) error {
	var body struct {
		Claim string `json:"claim"`
	}
	if err := c.BodyParser(&body); err != nil || body.Claim == "" {
		return fiber.NewError(fiber.StatusBadRequest, "claim required")
	}
	flow, ok := o.flows.take("claim:" + body.Claim)
	if !ok || flow.profile == nil {
		return c.JSON(fiber.Map{"result": "expired"})
	}
	user := auth.UserFrom(c)
	result := o.linkResult(c.UserContext(), user.ID, flow.profile)
	response := fiber.Map{"result": result}
	if updated, err := o.svc.Store().UserByID(c.UserContext(), user.ID); err == nil {
		response["user"] = o.view(c.UserContext(), updated)
	}
	return c.JSON(response)
}

func (o *onyxRoutes) register(c *fiber.Ctx) error {
	flow, ok := o.flows.take("claim:" + c.Query("claim"))
	if !ok || flow.profile == nil {
		return o.redirect(c, "/login", url.Values{"onyx": {"expired"}})
	}
	ctx := c.UserContext()
	user, err := o.svc.Store().UserByOAuth(ctx, onyxProvider, flow.profile.ID)
	if errors.Is(err, store.ErrNotFound) {
		user, err = o.svc.Store().CreateUser(ctx, store.NewUser{
			Role:        "player",
			Locale:      "en",
			DisplayName: flow.profile.Username,
		})
		if err == nil {
			err = o.svc.Store().LinkOAuthIdentity(ctx, user.ID, onyxProvider, flow.profile.ID, "")
		}
	}
	if err != nil {
		slog.Error("onyx: register", "err", err)
		return o.redirect(c, "/login", url.Values{"onyx": {"failed"}})
	}
	return o.signIn(c, user, flow.redirect)
}

func (o *onyxRoutes) linkResult(ctx context.Context, userID int64, profile *onyx.Profile) string {
	err := o.linkOnyx(ctx, userID, profile)
	switch {
	case err == nil:
		return "linked"
	case errors.Is(err, errOnyxTaken):
		return "taken"
	case errors.Is(err, errOnyxOtherIdentity):
		return "other"
	}
	slog.Error("onyx: link", "err", err)
	return "failed"
}

func (o *onyxRoutes) linkOnyx(ctx context.Context, userID int64, profile *onyx.Profile) error {
	st := o.svc.Store()
	owner, err := st.UserByOAuth(ctx, onyxProvider, profile.ID)
	if err == nil && owner.ID != userID {
		return errOnyxTaken
	}
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return err
	}
	subject, err := st.OAuthSubjectForUser(ctx, userID, onyxProvider)
	if err == nil && subject != profile.ID {
		return errOnyxOtherIdentity
	}
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return err
	}
	if err := st.LinkOAuthIdentity(ctx, userID, onyxProvider, profile.ID, ""); err != nil {
		return err
	}
	return st.UpdateUserDisplayName(ctx, userID, profile.Username)
}
