package admin

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"worldwright/backend/internal/store"
)

type Service struct {
	store     *store.Store
	assetsDir string
	log       *slog.Logger
}

func New(st *store.Store, assetsDir string, log *slog.Logger) *Service {
	if log == nil {
		log = slog.Default()
	}
	return &Service{store: st, assetsDir: assetsDir, log: log}
}

type OnboardingStatus struct {
	Completed       bool `json:"completed"`
	SeedDataPresent bool `json:"seed_data_present"`
}

func (s *Service) Status(ctx context.Context) (OnboardingStatus, error) {
	settings, err := s.store.GetSystemSettings(ctx)
	if err != nil {
		return OnboardingStatus{}, err
	}
	return OnboardingStatus{
		Completed:       settings.OnboardingCompleted,
		SeedDataPresent: settings.SeedDataPresent,
	}, nil
}

// KeepSeed marks onboarding as completed but leaves seed data intact.
func (s *Service) KeepSeed(ctx context.Context) error {
	return s.store.SetSystemSetting(ctx, "onboarding_completed", "true")
}

// PruneSeed removes all rows flagged is_seed=1 (entities, events, maps, ...)
// and marks onboarding as completed. Safe to call any time; user-created
// content is preserved.
func (s *Service) PruneSeed(ctx context.Context) error {
	paths, err := s.store.PruneContent(ctx, true)
	if err != nil {
		return err
	}
	s.removeAssetFiles(paths)
	if _, err := s.ensureDefaultCalendar(ctx); err != nil {
		return err
	}
	if err := s.store.SetSystemSetting(ctx, "seed_data_present", "false"); err != nil {
		return err
	}
	return s.store.SetSystemSetting(ctx, "onboarding_completed", "true")
}

// PruneAll wipes ALL content (regardless of is_seed). Used as preparation
// for an import flow that fully replaces the codex.
func (s *Service) PruneAll(ctx context.Context) error {
	paths, err := s.store.PruneContent(ctx, false)
	if err != nil {
		return err
	}
	s.removeAssetFiles(paths)
	if err := s.store.SetSystemSetting(ctx, "seed_data_present", "false"); err != nil {
		return err
	}
	return nil
}

// CompleteOnboarding marks onboarding done without touching content.
// Called after a successful import.
func (s *Service) CompleteOnboarding(ctx context.Context) error {
	return s.store.SetSystemSetting(ctx, "onboarding_completed", "true")
}

func (s *Service) AssetsDir() string { return s.assetsDir }
func (s *Service) Store() *store.Store { return s.store }
func (s *Service) Logger() *slog.Logger { return s.log }

func (s *Service) removeAssetFiles(paths []string) {
	for _, p := range paths {
		if p == "" {
			continue
		}
		if err := os.Remove(filepath.Join(s.assetsDir, p)); err != nil && !os.IsNotExist(err) {
			s.log.Warn("remove asset file", "path", p, "err", err)
		}
	}
}
