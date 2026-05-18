package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"worldwright/backend/internal/store"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrSessionExpired     = errors.New("session expired")
	ErrUnauthorized       = errors.New("unauthorized")
)

const (
	SessionCookieName = "worldwright_session"
	SessionTTL        = 30 * 24 * time.Hour
)

type Service struct {
	store        *store.Store
	cookieSecure bool
}

func New(s *store.Store, cookieSecure bool) *Service {
	return &Service{store: s, cookieSecure: cookieSecure}
}

func (s *Service) CookieSecure() bool { return s.cookieSecure }

func (s *Service) EnsureAdmin(ctx context.Context, email, password string) error {
	count, err := s.store.CountUsers(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.store.CreateUser(ctx, store.NewUser{
		Email:        strings.ToLower(strings.TrimSpace(email)),
		PasswordHash: string(hash),
		Role:         "admin",
		Locale:       "en",
	})
	return err
}

func (s *Service) ChangePassword(ctx context.Context, userID int64, currentPassword, newPassword string) error {
	user, err := s.store.UserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user.PasswordHash != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
			return ErrInvalidCredentials
		}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.store.UpdateUserPassword(ctx, userID, string(hash))
}

func (s *Service) Store() *store.Store { return s.store }

func (s *Service) LoginAsUser(ctx context.Context, user *store.User) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}
	if err := s.store.CreateSession(ctx, token, user.ID, time.Now().Add(SessionTTL)); err != nil {
		return "", err
	}
	return token, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, *store.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	user, err := s.store.UserByEmail(ctx, email)
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}
	token, err := generateToken()
	if err != nil {
		return "", nil, err
	}
	expiresAt := time.Now().Add(SessionTTL)
	if err := s.store.CreateSession(ctx, token, user.ID, expiresAt); err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	return s.store.DeleteSession(ctx, token)
}

func (s *Service) UserBySession(ctx context.Context, token string) (*store.User, error) {
	if token == "" {
		return nil, ErrSessionExpired
	}
	userID, expiresAt, err := s.store.SessionUser(ctx, token)
	if err != nil {
		return nil, ErrSessionExpired
	}
	if time.Now().After(expiresAt) {
		_ = s.store.DeleteSession(ctx, token)
		return nil, ErrSessionExpired
	}
	return s.store.UserByID(ctx, userID)
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
