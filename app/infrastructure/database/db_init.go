package database

import (
	"database/sql"
	"dynamic-user-segmentation-service/settings"
	"fmt"
	"log"
)

var Connection *sql.DB = nil

func GetConnectionToDB(storageInfo settings.Database) (connection *sql.DB, err error) {
	op := "app.interface.db.GetConnectionToDB"

	if Connection == nil {
		dbName := storageInfo.POSTGRES_DB_NAME
		dbHost := storageInfo.POSTGRES_HOST
		dbUserName := storageInfo.POSTGRES_USER
		dbPassword := storageInfo.POSTGRES_PASSWORD
		storagePath := fmt.Sprintf("host=%s dbname=%s user=%s password=%s", dbHost, dbName, dbUserName, dbPassword)

		connection, err := sql.Open("pgx", storagePath)
		if err != nil {
			log.Fatal("Connection Error:", err)
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		err = CreateTables(connection)
		if err != nil {
			log.Fatal("Table create Error:", err)
			return nil, err
		}
		err = connection.Ping()
		if err != nil {
			log.Fatal("Connection Error:", err)
			return nil, err
		}
		Connection = connection
	}
	return Connection, err
}

func CreateTables(db *sql.DB) (err error) {
	const op = "interfaces.db.NewUserRepository"

	qry := `
	BEGIN;
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE IF NOT EXISTS public."user"(
		user_uuid UUID DEFAULT uuid_generate_v4(),
		username VARCHAR(255),
		PRIMARY KEY (user_uuid)
	);

	CREATE TABLE IF NOT EXISTS public."slug"(
		id SERIAL,
		name VARCHAR(255) NOT NULL UNIQUE,
		PRIMARY KEY (id)
	);
	CREATE INDEX IF NOT EXISTS idx_slug_id ON public.slug(id);

	CREATE TABLE IF NOT EXISTS public."user_slug"
	(
		id BIGSERIAL PRIMARY KEY,
		user_uuid UUID NOT NULL REFERENCES public.user(user_uuid) ON DELETE CASCADE,
		slug_id INTEGER NOT NULL REFERENCES public.slug(id) ON DELETE CASCADE,
		UNIQUE (user_uuid, slug_id)
	);

	CREATE UNIQUE INDEX IF NOT EXISTS "idx_user_slug"
	ON "user_slug"
	USING btree
	(user_uuid, slug_id);

	COMMIT;
	`

	_, err = db.Exec(qry)
	if err != nil {
		log.Fatal("Error while execution init tables for database:", err)
		return fmt.Errorf("%s: %w", op, err)
	}
	return
}
