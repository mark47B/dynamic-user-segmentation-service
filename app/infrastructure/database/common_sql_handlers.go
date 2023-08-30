package database

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

type Repository struct {
	S *SlugRepository
	U *UserRepository
	IRepository
}

func (r *Repository) AddSlugToUser(user_uuid uuid.UUID, add_slugs []string) (err error) {
	const op = "interfaces.database.AddSlugToUser"

	// Check exising user
	isExist, err := r.U.isUserExist(user_uuid)
	if err != nil {
		log.Println("Error while trying to call 'isUserExist'", err)
		return fmt.Errorf("%s: %w", op, err)
	}
	if !isExist {
		log.Println("There are no the user in the database", err)
		return fmt.Errorf("user (uuid %d) not found", user_uuid)
	}

	// Check existing slugs
	isExist, err = r.S.isSlugsExist(add_slugs)

	if err != nil {
		log.Println("Error while trying to call 'isSlugsExist'", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	if !isExist {
		log.Println("There are no slugs (like add_slugs) in the database")
		return fmt.Errorf("%s: %s", op, "Create slugs!")
	}

	// Select slugs ids
	slugs_ids, err := r.S.GetSlugsIds(add_slugs)
	if err != nil {
		log.Println("Error while getting slugss' ids", err)
		return fmt.Errorf("%s: %s", op, err)
	}

	// Inserting slugs
	err = r.S.InsertSlugsForUser(user_uuid, slugs_ids)
	if err != nil {
		log.Println("Error while inserting slugs for user:", err)
		return fmt.Errorf("%s: %s", op, err)
	}

	return
}

func (r *Repository) DeleteSlugsForUser(user_uuid uuid.UUID, delete_slugs []string) error {
	const op = "interfaces.database.DeleteSlugsForUser"

	fmt.Println(user_uuid, delete_slugs, op)
	// Check exising user
	isExist, err := r.U.isUserExist(user_uuid)
	if err != nil {
		log.Println("Error while checking exising user ", err)
		return fmt.Errorf("%s: %w", op, err)
	}
	if !isExist {
		log.Println("There are no the user in the database", err)
		return fmt.Errorf("user (uuid %d) not found", user_uuid)
	}

	// Check existing slugs
	isExist, err = r.S.isSlugsExist(delete_slugs)

	if err != nil {
		log.Println("Error while trying checking existing slugs", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	if !isExist {
		log.Println("There are no slugs (like add_slugs) in the database")
		return fmt.Errorf("%s: %s", op, "There are no your slugs")
	}

	// Select slugs ids
	slugs_ids, err := r.S.GetSlugsIds(delete_slugs)
	if err != nil {
		log.Println("Error while gettings slugs' ids:", err)
		return fmt.Errorf("%s: %s", op, err)
	}

	// Delete slugs
	err = r.S.DeleteSlugsForUser(user_uuid, slugs_ids)
	if err != nil {
		log.Println("Error while deleting slug for users", err)
		return fmt.Errorf("%s: %s", op, err)
	}

	return nil

}
