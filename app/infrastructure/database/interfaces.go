package database

import (
	"dynamic-user-segmentation-service/core"

	"github.com/google/uuid"
)

type IUser interface {
	GetUserByUUID(user_uuid uuid.UUID) (userOut core.User, err error)
	isUserExist(user_uuid uuid.UUID) (bool, error)
	AddSlugToUser(user_uuid uuid.UUID, add_slugs []string)
	DeleteSlugToUser(user_uuid uuid.UUID, delete_slugs []core.Slug)
	SelectSlugsIdsForUser(user_uuid uuid.UUID) (slugs_ids []uint32, err error)
	GetAll() (users []core.User, amount int, err error)
	CreateUser(user *core.UserRequestCreate) (user_uuid uuid.UUID, err error)
	DeleteByUUID(user_uuid uuid.UUID) (err error)
}

type ISlug interface {
	isSlugsExist(sulgs []string) (bool, error)
	DeleteSlugByName(slugName string) (err error)
	GetSlugsIds(slugs []string) (ids []uint32, err error)
	GetByUserUUID(user_uuid uuid.UUID) (slugOut core.Slug, err error)
	GetByID(slug_id uint32) (core.Slug, error)
	InsertSlugsForUser(user_uuid uuid.UUID, slugs_ids []uint32) (err error)
	DeleteSlugsForUser(user_uuid uuid.UUID, slugs_ids []uint32) (err error)
	CreateSlug(slug *core.SlugRequestAdd) (slugID uint32, err error)
}

type IRepository interface {
	InsertSlugsForUser(user_uuid uuid.UUID, slugs_ids []uint32) (err error)
	DeleteSlugsForUser(user_uuid uuid.UUID, delete_slugs []string) error
}
