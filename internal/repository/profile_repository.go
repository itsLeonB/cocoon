package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/appconstant"
	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/itsLeonB/ezutil"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type profileRepositoryGorm struct {
	ezutil.CRUDRepository[entity.UserProfile]
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) UserProfileRepository {
	return &profileRepositoryGorm{
		ezutil.NewCRUDRepository[entity.UserProfile](db),
		db,
	}
}

func (pr *profileRepositoryGorm) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.UserProfile, error) {
	var profiles []entity.UserProfile

	db, err := pr.CRUDRepository.GetGormInstance(ctx)
	if err != nil {
		return nil, err
	}

	if err = db.Where("id IN ?", ids).Find(&profiles).Error; err != nil {
		return nil, eris.Wrap(err, appconstant.ErrDataSelect)
	}

	return profiles, nil
}
