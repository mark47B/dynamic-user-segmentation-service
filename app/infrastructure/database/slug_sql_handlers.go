package database

import (
	"database/sql"
	"dynamic-user-segmentation-service/core"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
)

type SlugRepository struct {
	db *sql.DB
	ISlug
}

func NewSlugRepository(connection *sql.DB) (*SlugRepository, error) {
	return &SlugRepository{db: connection}, nil
}

func (s *SlugRepository) CreateSlug(slug *core.SlugRequestAdd) (slugID uint32, err error) {
	const op = "interfaces.db.CreateSlug"

	qry := `INSERT INTO public.slug (name) VALUES ($1) RETURNING id;`
	_, err = s.db.Prepare(qry)
	if err != nil {
		log.Println("Error preparing inserting slug:", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	err = s.db.QueryRow(qry, slug.Name).Scan(&slugID)
	if err != nil {
		log.Println("Error while executing insert user:", err)

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return
}

func (s *SlugRepository) DeleteSlugByName(slugName string) (err error) {
	const op = "interfaces.db.Delete"

	qry := `DELETE FROM slug WHERE name = $1;`
	_, err = s.db.Prepare(qry)
	if err != nil {
		log.Println("Error preparing deleting slug:", err)
		return fmt.Errorf("%s: %w", op, err)
	}
	err = s.db.QueryRow(qry, slugName).Err()
	if err != nil {
		log.Println("Error while executing delete slug:", err)
		return fmt.Errorf("%s: %w", op, err)
	}
	return
}

func (s *SlugRepository) isSlugsExist(sulgs []string) (bool, error) {
	const op = "interfaces.db.isSlugsExist"

	var amount uint32
	rows, err := s.db.Query(
		`SELECT COUNT("name") FROM public."slug" WHERE "name" IN ('` + strings.Join(sulgs, "','") + `')`)
	defer rows.Close()
	rows.Next()
	rows.Scan(&amount)
	if err != nil {
		log.Println("Error while trying to add slug to user(checking amount):", err)
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if amount == 0 {
		return false, nil
	}
	return true, nil
}

func (s *SlugRepository) GetSlugsIds(slugs []string) (slugs_ids []uint32, err error) {
	const op = "interfaces.db.GetSlugsIds"

	rows, err := s.db.Query(`SELECT id FROM public."slug" WHERE name IN ('` + strings.Join(slugs, "','") + `')`)
	defer rows.Close()
	if err != nil {
		log.Println("Error while trying to get slugs ids by slug's names:", err)
		return slugs_ids, fmt.Errorf("%s: %w", op, err)
	}
	for rows.Next() {
		var slug_id uint32
		err := rows.Scan(&slug_id)
		if err != nil {
			log.Println("Error while scanning rows:", err)
			return slugs_ids, fmt.Errorf("%s: %w", op, err)
		}

		slugs_ids = append(slugs_ids, slug_id)
	}
	return
}

func (s *SlugRepository) InsertSlugsForUser(user_uuid uuid.UUID, slugs_ids []uint32) (err error) {
	const op = "interfaces.db.InsertSlugsForUser"

	qry := `INSERT INTO public.user_slug (user_uuid, slug_id) VALUES ($1, $2);`
	_, err = s.db.Prepare(qry)
	if err != nil {
		log.Println("Error preparing inserting slugs for user:", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, val := range slugs_ids {
		_, err = s.db.Exec(qry, user_uuid, val)
		if err != nil {
			log.Println("Error while executing insert slug:", err)
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return
}

func (s *SlugRepository) DeleteSlugsForUser(user_uuid uuid.UUID, slugs_ids []uint32) (err error) {
	const op = "interfaces.db.InsertSlugsForUser"

	qry := `DELETE FROM public.user_slug WHERE user_uuid = $1 AND slug_id = $2;`
	_, err = s.db.Prepare(qry)
	if err != nil {
		log.Println("Error preparing delete slugs for user:", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, val := range slugs_ids {
		fmt.Println(user_uuid, val, op)
		_, err = s.db.Exec(qry, user_uuid, val)
		if err != nil {
			log.Println("Error while executing delete slug:", err)
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return
}

// GetByUserUUID(user_uuid uuid.UUID) (slugOut core.Slug, err error)
// GetByID(slug_id uint32)
// GetAll() (slugs []core.Slug, amount int, err error)
