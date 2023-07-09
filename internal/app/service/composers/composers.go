package composers

import (
	"context"

	"github.com/ttrueno/rl2-final/internal/app/models"
)

type storage interface {
	GetAll(ctx context.Context, offset int, limit int) ([]models.Composer, error)
}

type service struct {
	storage storage
	Cfg     config
}

type config struct {
	pageLength int
}

func NewConfig(pageLength int) *config {
	return &config{
		pageLength: pageLength,
	}
}

func NewService(storage storage, config config) *service {
	return &service{
		storage: storage,
		Cfg:     config,
	}
}

func (s *service) GetAll(ctx context.Context, pageNumber int) ([]models.Composer, error) {
	return s.storage.GetAll(ctx, (pageNumber-1)*s.Cfg.pageLength, s.Cfg.pageLength)
}

func (s *service) PageLength() int {
	return s.Cfg.pageLength
}
