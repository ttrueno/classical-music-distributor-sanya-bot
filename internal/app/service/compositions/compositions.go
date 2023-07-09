package compositions

import (
	"context"

	"github.com/ttrueno/rl2-final/internal/app/models"
)

type config struct {
	MaxCompositions int
}

func NewConfig(maxCompositions int) *config {
	return &config{
		MaxCompositions: maxCompositions,
	}
}

type compositionsStorage interface {
	GetAllByComposerID(ctx context.Context, ComposerID string, Offset, Limit int) ([]models.Composition, error)
}

type service struct {
	storage compositionsStorage
	config  config
}

func New(
	compositionsStorage compositionsStorage,
	config config,
) *service {
	return &service{
		storage: compositionsStorage,
		config:  config,
	}
}

func (s *service) GetAllByComposerID(ctx context.Context, ComposerID string) ([]models.Composition, error) {
	return s.storage.GetAllByComposerID(ctx, ComposerID, 0, s.config.MaxCompositions)
}
