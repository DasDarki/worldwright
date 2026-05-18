package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"worldwright/backend/internal/store"
)

const (
	OAuthStateCookie    = "worldwright_oauth_state"
	OAuthVerifierCookie = "worldwright_oauth_verifier"
	OAuthReturnCookie   = "worldwright_oauth_return"
	OAuthStateTTL       = 10 * time.Minute
)

var ErrProviderUnknown = errors.New("oauth provider not enabled")
var ErrSignupDisabled = errors.New("oauth signup disabled")

type OAuthProvider struct {
	Name        string
	Config      *oauth2.Config
	UserinfoURL string
	ParseUser   func(map[string]any) (subject, email, name, avatar string)
}

type OAuth struct {
	providers   map[string]*OAuthProvider
	store       *store.Store
	allowSignup bool
	publicURL   string
}

type OAuthSettings struct {
	PublicURL           string
	AllowSignup         bool
	GoogleClientID      string
	GoogleClientSecret  string
	DiscordClientID     string
	DiscordClientSecret string
}

func NewOAuth(s *store.Store, settings OAuthSettings) *OAuth {
	o := &OAuth{
		providers:   make(map[string]*OAuthProvider),
		store:       s,
		allowSignup: settings.AllowSignup,
		publicURL:   strings.TrimRight(settings.PublicURL, "/"),
	}
	redirect := func(provider string) string {
		return o.publicURL + "/api/auth/oauth/" + provider + "/callback"
	}
	if settings.GoogleClientID != "" && settings.GoogleClientSecret != "" {
		o.providers["google"] = &OAuthProvider{
			Name: "google",
			Config: &oauth2.Config{
				ClientID:     settings.GoogleClientID,
				ClientSecret: settings.GoogleClientSecret,
				RedirectURL:  redirect("google"),
				Scopes:       []string{"openid", "email", "profile"},
				Endpoint:     google.Endpoint,
			},
			UserinfoURL: "https://openidconnect.googleapis.com/v1/userinfo",
			ParseUser: func(m map[string]any) (string, string, string, string) {
				sub, _ := m["sub"].(string)
				email, _ := m["email"].(string)
				name, _ := m["name"].(string)
				picture, _ := m["picture"].(string)
				return sub, email, name, picture
			},
		}
	}
	if settings.DiscordClientID != "" && settings.DiscordClientSecret != "" {
		o.providers["discord"] = &OAuthProvider{
			Name: "discord",
			Config: &oauth2.Config{
				ClientID:     settings.DiscordClientID,
				ClientSecret: settings.DiscordClientSecret,
				RedirectURL:  redirect("discord"),
				Scopes:       []string{"identify", "email"},
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://discord.com/oauth2/authorize",
					TokenURL: "https://discord.com/api/oauth2/token",
				},
			},
			UserinfoURL: "https://discord.com/api/users/@me",
			ParseUser: func(m map[string]any) (string, string, string, string) {
				id, _ := m["id"].(string)
				email, _ := m["email"].(string)
				username, _ := m["username"].(string)
				globalName, _ := m["global_name"].(string)
				name := globalName
				if name == "" {
					name = username
				}
				avatar, _ := m["avatar"].(string)
				avatarURL := ""
				if id != "" && avatar != "" {
					avatarURL = "https://cdn.discordapp.com/avatars/" + id + "/" + avatar + ".png"
				}
				return id, email, name, avatarURL
			},
		}
	}
	return o
}

func (o *OAuth) Enabled() []string {
	names := make([]string, 0, len(o.providers))
	for name := range o.providers {
		names = append(names, name)
	}
	return names
}

func (o *OAuth) HasProvider(name string) bool {
	_, ok := o.providers[name]
	return ok
}

func (o *OAuth) Start(provider string) (authURL, state, verifier string, err error) {
	p, ok := o.providers[provider]
	if !ok {
		return "", "", "", ErrProviderUnknown
	}
	state, err = randomURLSafe(32)
	if err != nil {
		return "", "", "", err
	}
	verifier = oauth2.GenerateVerifier()
	authURL = p.Config.AuthCodeURL(state,
		oauth2.AccessTypeOnline,
		oauth2.S256ChallengeOption(verifier),
	)
	return authURL, state, verifier, nil
}

type OAuthUser struct {
	Subject string
	Email   string
	Name    string
	Avatar  string
}

func (o *OAuth) Callback(ctx context.Context, provider, code, verifier string) (*OAuthUser, error) {
	p, ok := o.providers[provider]
	if !ok {
		return nil, ErrProviderUnknown
	}
	token, err := p.Config.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, p.UserinfoURL, nil)
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.New("userinfo: " + resp.Status + " " + string(body))
	}
	var raw map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	sub, email, name, avatar := p.ParseUser(raw)
	if sub == "" {
		return nil, errors.New("userinfo: missing subject")
	}
	return &OAuthUser{Subject: sub, Email: email, Name: name, Avatar: avatar}, nil
}

func (o *OAuth) ResolveUser(ctx context.Context, provider string, info *OAuthUser) (*store.User, error) {
	user, err := o.store.UserByOAuth(ctx, provider, info.Subject)
	if err == nil {
		return user, nil
	}
	if !errors.Is(err, store.ErrNotFound) {
		return nil, err
	}
	if info.Email != "" {
		if existing, err := o.store.UserByEmail(ctx, strings.ToLower(info.Email)); err == nil {
			if linkErr := o.store.LinkOAuthIdentity(ctx, existing.ID, provider, info.Subject, info.Email); linkErr != nil {
				return nil, linkErr
			}
			return existing, nil
		}
	}
	if !o.allowSignup {
		return nil, ErrSignupDisabled
	}
	created, err := o.store.CreateUser(ctx, store.NewUser{
		Email:       strings.ToLower(info.Email),
		Role:        "player",
		Locale:      "en",
		DisplayName: info.Name,
		AvatarURL:   info.Avatar,
	})
	if err != nil {
		return nil, err
	}
	if err := o.store.LinkOAuthIdentity(ctx, created.ID, provider, info.Subject, info.Email); err != nil {
		return nil, err
	}
	return created, nil
}

func (o *OAuth) PublicURL() string { return o.publicURL }

func randomURLSafe(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
