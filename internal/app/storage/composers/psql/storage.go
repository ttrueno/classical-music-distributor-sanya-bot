package psql

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ttrueno/rl2-final/internal/app/models"
	"github.com/ttrueno/rl2-final/internal/app/storage/composers/psql/sql/composers"
	"github.com/ttrueno/rl2-final/internal/db/psql"
	"github.com/ttrueno/rl2-final/internal/lib/cond"
	"github.com/ttrueno/rl2-final/internal/lib/conv"
)

type storage struct {
	queries *composers.Queries
}

func New(db composers.DBTX) *storage {
	return &storage{
		queries: composers.New(db),
	}
}

func (s *storage) CreateComposer(ctx context.Context, composer models.Composer) (*models.Composer, error) {
	dbComposer, err := s.queries.InsertComposer(ctx, composers.InsertComposerParams{
		FirstName: composer.FirstName,
		LastName:  composer.LastName,
		ImageLink: pgtype.Text{
			String: cond.Ter(composer.ImageLink != "", composer.ImageLink, ""),
			Valid:  cond.Ter(composer.ImageLink != "", true, false),
		},
		Description: pgtype.Text{
			String: cond.Ter(composer.Description != "", composer.Description, ""),
			Valid:  cond.Ter(composer.Description != "", true, false),
		},
	})

	if err != nil {
		return nil, err
	}

	c := dbToModel(dbComposer)

	return &c, nil
}

func (s *storage) GetAll(ctx context.Context, offset int, limit int) ([]models.Composer, error) {
	dbComposers, err := s.queries.SelectComposers(ctx, composers.SelectComposersParams{
		Offset: int32(offset),
		Limit:  int32(limit),
	})
	if err != nil {
		return nil, err
	}
	if len(dbComposers) < 1 {
		return nil, psql.ErrNoRecords
	}
	var (
		composersLen = len(dbComposers)
		composers    = make([]models.Composer, 0, composersLen)
	)

	for _, dbComposer := range dbComposers {
		composers = append(composers, dbToModel(dbComposer))
	}
	return composers, nil
}

func (s *storage) GetComposer(ctx context.Context, ID string) (*models.Composer, error) {
	id, err := conv.Atoi64(ID)
	if err != nil {
		if strings.Trim(ID, " ") == "" {
			return nil, models.ErrEmptyID
		}
		return nil, err
	}
	dbComposer, err := s.queries.SelectComposer(ctx, id)
	if err != nil {
		return nil, err
	}

	composer := dbToModel(dbComposer)

	return &composer, nil
}

func dbToModel(dbComposer composers.Composer) models.Composer {
	return models.Composer{
		ID:          conv.Itoa64(dbComposer.ID),
		FirstName:   dbComposer.FirstName,
		LastName:    dbComposer.LastName,
		ImageLink:   dbComposer.ImageLink.String,
		Description: dbComposer.Description.String,
		Version:     conv.Itoa64(dbComposer.Version),
	}
}

func dbFromModel(composer models.Composer) (*composers.Composer, bool) {
	valid := true
	id, err := conv.Atoi64(composer.ID)
	if err != nil {
		valid = false
	}

	version, err := conv.Atoi64(composer.Version)
	if err != nil {
		valid = false
	}

	return &composers.Composer{
		ID:        id,
		FirstName: composer.FirstName,
		LastName:  composer.LastName,
		ImageLink: pgtype.Text{
			String: composer.ImageLink,
			Valid:  true,
		},
		Description: pgtype.Text{
			String: composer.Description,
			Valid:  true,
		},
		Version: version,
	}, valid
}
