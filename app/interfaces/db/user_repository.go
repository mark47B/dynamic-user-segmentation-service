package db

import (
	"database/sql"
	"dynamic-user-segmentation-service/core"
	"fmt"
	"log"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type UserRepository struct {
	db *sql.DB
	IUser
}

func NewUserRepository(connection *sql.DB) (*UserRepository, error) {
	return &UserRepository{db: connection}, nil
}

// TODO: Decompose
func (u *UserRepository) GetUserByUUID(user_uuid uuid.UUID) (*core.User, error) {
	const op = "interfaces.db.GetUserByUUID"

	// Select user for extracting username and other info
	qry := `SELECT public."user".username FROM public."user" WHERE public."user".user_uuid=$1`
	var username string

	err := u.db.QueryRow(qry, user_uuid).Scan(&username)
	if err != nil {
		log.Println("Error while trying to get user by UUID:", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Select all slugs for user with apropriate user_uuid
	// Check existing of slugs for user
	var slugs []core.Slug

	qry_count := `SELECT COUNT(id) from public.user_slug WHERE user_uuid = $1 `
	var amount int
	err = u.db.QueryRow(qry_count, user_uuid).Scan(&amount)
	if err != nil {
		log.Println("Error while trying to get all users (amount):", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if amount == 0 {
		return &core.User{
			UUID:     user_uuid,
			Username: username,
			Slugs:    slugs,
		}, nil
	}

	qry = `SELECT public.slug.id, public.slug.name 
			FROM public.user
			LEFT JOIN public.user_slug ON public.user.user_uuid = public.user_slug.user_UUID 
			LEFT JOIN public.slug ON public.user_slug.slug_id = public.slug.id
			WHERE public.user.user_uuid=$1`

	rows, err := u.db.Query(qry, user_uuid)
	if err != nil {
		log.Println("Error while trying to get user's slugs:", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		fmt.Println(rows)
		var ID uint32
		var slug_name string
		err := rows.Scan(&ID, &slug_name)
		if err != nil {
			log.Println("Error while scanning rows:", err)
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		NewSlug := core.Slug{
			ID:   ID,
			Name: slug_name,
		}
		slugs = append(slugs, NewSlug)
	}

	return &core.User{
		UUID:     user_uuid,
		Username: username,
		Slugs:    slugs,
	}, nil
}

// func (u *UserRepository) GetAll() (users []core.User, amount int, err error) {
// 	return users, 3, nil
// }

func (u *UserRepository) CreateUser(user core.UserRequestCreate) (user_uuid uuid.UUID, err error) {
	const op = "interfaces.db.GetUserByUUID"

	qry := `INSERT INTO public."user" (username) VALUES ($1) RETURNING user_uuid;`
	_, err = u.db.Prepare(qry)
	if err != nil {
		log.Println("Error preparing inserting user:", err)
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	err = u.db.QueryRow(qry, user.Username).Scan(&user_uuid)
	if err != nil {
		log.Println("Error while executing insert user:", err)

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return
}

func (u *UserRepository) isUserExist(user_uuid uuid.UUID) (bool, error) {
	const op = "interfaces.db.isUserExist"

	qry_count := `SELECT COUNT(user_uuid) FROM public.user WHERE public.user.user_uuid = $1;`

	var amount int
	err := u.db.QueryRow(qry_count, user_uuid).Scan(&amount)
	if err != nil {
		log.Println("Error while trying to add slug to user(checking exisiting user):", err)
		return false, fmt.Errorf("%s: %w", op, err)
	}
	if amount == 0 {
		return false, nil
	}

	return true, nil
}

// func (user_uuid uuid.UUID, delete_slugs []core.Slug)

// func (u *UserRepository) DeleteByUUID(id string) (err error) {
// 	return
// }
