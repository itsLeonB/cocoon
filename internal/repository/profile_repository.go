package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/appconstant"
	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/itsLeonB/go-crud"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type profileRepositoryGorm struct {
	crud.Repository[entity.UserProfile]
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) UserProfileRepository {
	return &profileRepositoryGorm{
		crud.NewRepository[entity.UserProfile](db),
		db,
	}
}

func (pr *profileRepositoryGorm) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.UserProfile, error) {
	var profiles []entity.UserProfile

	db, err := pr.GetGormInstance(ctx)
	if err != nil {
		return nil, err
	}

	if err = db.Where("id IN ?", ids).Find(&profiles).Error; err != nil {
		return nil, eris.Wrap(err, appconstant.ErrDataSelect)
	}

	return profiles, nil
}
