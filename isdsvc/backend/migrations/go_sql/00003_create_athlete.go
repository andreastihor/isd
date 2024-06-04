package go_sql

import (
	"database/sql"

	"github.com/pkg/errors"
)

func init() {
	addMigration(up_00003, down_00003)
}

func up_00003(tx *sql.Tx) error {
	// Create table
	_, err := tx.Exec(`
	CREATE TABLE athlete (
		id              VARCHAR(36) PRIMARY KEY,
		club_id VARCHAR(36) NOT NULL,
		name            VARCHAR(255) NOT NULL,
		register_date  DATE NOT NULL,
		email		VARCHAR(255) NOT NULL,
		phone_number    VARCHAR(20) NOT NULL,
		gender VARCHAR(1) NOT NULL,
		date_of_birth  DATE NOT NULL,
		active VARCHAR(10) NOT NULL
	);	
	`)
	if err != nil {
		return errors.Wrap(err, "failed to create athlete table")
	}

	// Index id and name
	_, err = tx.Exec(`
	CREATE INDEX idx_athlete_id ON athlete (id);
	CREATE INDEX idx_athlete_name ON athlete (name);
	`)
	if err != nil {
		return errors.Wrap(err, "failed to create indexes")
	}

	return nil
}

func down_00003(tx *sql.Tx) error {
	// Drop table
	_, err := tx.Exec("DROP TABLE IF EXISTS athlete;")
	if err != nil {
		return errors.Wrap(err, "failed to drop athlete table")
	}

	return nil
}
