package database

import (
	"database/sql"
	"dynamic-user-segmentation-service/core"
	"fmt"
	"log"
	"strconv"
	"strings"

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

func (u *UserRepository) SelectSlugsIdsForUser(user_uuid uuid.UUID) (slugs_ids []uint32, err error) {
	op := "infrastucture.database.GetSlugsIdsForUser"

	qry := `SELECT id from public.user_slug WHERE user_uuid = $1 `

	rows, err := u.db.Query(qry, user_uuid)
	defer rows.Close()
	if err != nil {
		log.Println("Error while trying to get user's slugs_ids:", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		var ID uint32
		err := rows.Scan(&ID)
		if err != nil {
			log.Println("Error while scanning rows:", err)
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		slugs_ids = append(slugs_ids, ID)
	}
	return
}

func (u *UserRepository) SelectUserSlugsByUUID(user_uuid uuid.UUID) (user_slugs []core.Slug, err error) {
	op := "infrastructure.database.SelectUserSlugsByUUID"

	user_slugs_ids, err := u.SelectSlugsIdsForUser(user_uuid)
	user_slugs_ids_str := *new([]string)
	for _, val := range user_slugs_ids {
		user_slugs_ids_str = append(user_slugs_ids_str, strconv.FormatUint(uint64(val), 10))
	}

	if err != nil {
		log.Println("Error while trying Select Slugs Ids For User", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := u.db.Query(`SELECT id, name from public.user_slug WHERE id IN ('` + strings.Join(user_slugs_ids_str, "','") + `')`)
	defer rows.Close()
	if err != nil {
		log.Println("Error while trying to get slug name", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
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
		user_slugs = append(user_slugs, NewSlug)
	}
	return
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

func (u *UserRepository) GetAll() (users []core.User, err error) {
	op := "infrastructure.database.GetApp"

	qry := `SELECT public.user.user_uuid, public.user.username, public.slug.id, public.slug.name 
			FROM public.user
			LEFT JOIN public.user_slug ON public.user.user_uuid = public.user_slug.user_UUID 
			LEFT JOIN public.slug ON public.user_slug.slug_id = public.slug.id 
			ORDER BY public.user.user_uuid`

	rows, err := u.db.Query(qry)
	defer rows.Close()
	if err != nil {
		log.Println("Error while extracting all users:", err)

		return users, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		var (
			user_uuid      uuid.UUID
			username, name string
			ID             uint32
		)
		rows.Scan(&user_uuid, &username, &ID, &name)
		if len(users) != 0 && users[len(users)-1].UUID == user_uuid {
			users[len(users)-1].Slugs = append(users[len(users)-1].Slugs, core.Slug{ID: ID, Name: name})
		} else {
			users = append(users, core.User{UUID: user_uuid, Username: username, Slugs: []core.Slug{{ID: ID, Name: name}}})
		}
	}
	return
}

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
