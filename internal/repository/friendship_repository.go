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

type friendshipRepositoryGorm struct {
	ezutil.CRUDRepository[entity.Friendship]
	db *gorm.DB
}

func NewFriendshipRepository(db *gorm.DB) FriendshipRepository {
	return &friendshipRepositoryGorm{
		ezutil.NewCRUDRepository[entity.Friendship](db),
		db,
	}
}

func (fr *friendshipRepositoryGorm) Insert(ctx context.Context, friendship entity.Friendship) (entity.Friendship, error) {
	db, err := fr.getGormInstance(ctx)
	if err != nil {
		return entity.Friendship{}, err
	}

	if err = db.Create(&friendship).Error; err != nil {
		return entity.Friendship{}, eris.Wrap(err, appconstant.ErrDataInsert)
	}

	return friendship, nil
}

func (fr *friendshipRepositoryGorm) FindFirstBySpec(ctx context.Context, spec entity.FriendshipSpecification) (entity.Friendship, error) {
	var friendship entity.Friendship

	db, err := fr.getGormInstance(ctx)
	if err != nil {
		return entity.Friendship{}, err
	}

	query := db.
		Scopes(
			ezutil.WhereBySpec(spec.Model),
			ezutil.PreloadRelations(spec.PreloadRelations),
		).
		Joins("JOIN user_profiles AS up1 ON up1.id = friendships.profile_id1").
		Joins("JOIN user_profiles AS up2 ON up2.id = friendships.profile_id2")

	if spec.Name != "" {
		query = query.Where(
			db.Where("up1.name = ? AND friendships.profile_id1 <> ?", spec.Name, spec.Model.ProfileID1).
				Or("up2.name = ? AND friendships.profile_id2 <> ?", spec.Name, spec.Model.ProfileID1),
		)
	}

	err = query.Take(&friendship).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.Friendship{}, nil
		}
		return entity.Friendship{}, eris.Wrap(err, appconstant.ErrDataSelect)
	}

	return friendship, nil
}

func (fr *friendshipRepositoryGorm) FindAllBySpec(ctx context.Context, spec entity.FriendshipSpecification) ([]entity.Friendship, error) {
	var friendships []entity.Friendship

	db, err := fr.getGormInstance(ctx)
	if err != nil {
		return nil, err
	}

	err = db.
		Where(entity.Friendship{ProfileID1: spec.Model.ProfileID1}).
		Or(entity.Friendship{ProfileID2: spec.Model.ProfileID1}).
		Scopes(
			ezutil.PreloadRelations(spec.PreloadRelations),
			ezutil.DefaultOrder(),
		).
		Find(&friendships).
		Error

	if err != nil {
		return nil, eris.Wrap(err, appconstant.ErrDataSelect)
	}

	return friendships, nil
}

func (fr *friendshipRepositoryGorm) FindByProfileIDs(ctx context.Context, profileID1, profileID2 uuid.UUID) (entity.Friendship, error) {
	db, err := fr.getGormInstance(ctx)
	if err != nil {
		return entity.Friendship{}, err
	}

	var friendship entity.Friendship
	err = db.Where("(profile_id1 = ? AND profile_id2 = ?) OR (profile_id1 = ? AND profile_id2 = ?)", profileID1, profileID2, profileID2, profileID1).
		First(&friendship).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.Friendship{}, nil
		}
		return entity.Friendship{}, eris.Wrap(err, appconstant.ErrDataSelect)
	}

	return friendship, nil
}

func (fr *friendshipRepositoryGorm) getGormInstance(ctx context.Context) (*gorm.DB, error) {
	tx, err := ezutil.GetTxFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if tx != nil {
		return tx, nil
	}

	return fr.db.WithContext(ctx), nil
}
