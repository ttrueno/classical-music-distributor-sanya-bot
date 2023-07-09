package psql

import (
	"context"
	"time"

	"github.com/ttrueno/rl2-final/internal/app/models"
	"github.com/ttrueno/rl2-final/internal/app/storage/compositions/psql/sql/compositions"
	"github.com/ttrueno/rl2-final/internal/app/storage/compositions/psql/sql/mirrors"
	"github.com/ttrueno/rl2-final/internal/db/psql"
	"github.com/ttrueno/rl2-final/internal/lib/conv"
)

type storage struct {
	compositionsQueries *compositions.Queries
	mirrorsQueries      *mirrors.Queries
}

type db interface {
	compositions.DBTX
	mirrors.DBTX
}

func New(db db) *storage {
	return &storage{
		compositionsQueries: compositions.New(db),
		mirrorsQueries:      mirrors.New(db),
	}
}

func (s *storage) GetAllByComposerID(ctx context.Context, ComposerID string, Offset, Limit int) ([]models.Composition, error) {
	id, err := conv.Atoi64(ComposerID)
	if err != nil {
		return nil, err
	}

	dbCompositions, err := s.compositionsQueries.SelectAllCompositionsByComposerID(ctx, compositions.SelectAllCompositionsByComposerIDParams{
		ComposerID: id,
		Limit:      int32(Limit),
		Offset:     int32(Offset),
	})
	if err != nil {
		return nil, err
	}

	if len(dbCompositions) < 1 {
		return nil, psql.ErrNoRecords
	}

	var res = make([]models.Composition, 0, len(dbCompositions))
	for _, dbComposition := range dbCompositions {
		composition := dbCompositionToModel(dbComposition)
		err = func() error {
			fetchCtx, fetchCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer fetchCancel()
			dbMirrors, err := s.mirrorsQueries.GetCompositionsMirrorsByCompositionID(fetchCtx, dbComposition.ID)
			if err != nil {
				return err
			}
			composition.Mirrors = make([]models.CompositionMirror, 0, len(dbMirrors))
			for _, dbMirror := range dbMirrors {
				composition.Mirrors = append(composition.Mirrors, dbMirrorToModel(dbMirror))
			}

			return nil
		}()

		res = append(res, composition)
	}

	return res, nil
}

func dbCompositionToModel(composition compositions.Composition) models.Composition {
	return models.Composition{
		ID:         conv.Itoa64(composition.ID),
		ComposerID: conv.Itoa64(composition.ComposerID),
		Name:       composition.Name,
		Version:    conv.Itoa64(composition.Version),
	}
}

func dbCompositionFromModel(composition models.Composition) (*compositions.Composition, bool) {
	valid := true
	id, err := conv.Atoi64(composition.ID)
	if err != nil {
		valid = false
	}

	composerID, err := conv.Atoi64(composition.ComposerID)
	if err != nil {
		valid = false
	}

	version, err := conv.Atoi64(composition.Version)
	if err != nil {
		valid = false
	}

	return &compositions.Composition{
		ID:         id,
		ComposerID: composerID,
		Name:       composition.Name,
		Version:    version,
	}, valid
}

func dbMirrorToModel(composition mirrors.CompositionMirror) models.CompositionMirror {
	return models.CompositionMirror{
		ID:            conv.Itoa64(composition.ID),
		CompositionID: conv.Itoa64(composition.CompositionID),
		Link:          composition.Link,
		Version:       conv.Itoa64(composition.Version),
	}
}

func dbMirrorFromModel(composition models.CompositionMirror) (*mirrors.CompositionMirror, bool) {
	valid := true
	id, err := conv.Atoi64(composition.ID)
	if err != nil {
		valid = false
	}

	compositionID, err := conv.Atoi64(composition.CompositionID)
	if err != nil {
		valid = false
	}

	version, err := conv.Atoi64(composition.Version)
	if err != nil {
		valid = false
	}

	return &mirrors.CompositionMirror{
		ID:            id,
		CompositionID: compositionID,
		Link:          composition.Link,
		Version:       version,
	}, valid
}
