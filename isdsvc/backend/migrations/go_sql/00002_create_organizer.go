package go_sql

import (
	"database/sql"

	"github.com/pkg/errors"
)

func init() {
	addMigration(up_00002, down_00002)
}

func up_00002(tx *sql.Tx) error {
	// Create table
	_, err := tx.Exec(`
	CREATE TABLE organizer (
		id              VARCHAR(36) PRIMARY KEY,
		club_id VARCHAR(36) NOT NULL,
		name            VARCHAR(255) NOT NULL,
		register_date  DATE NOT NULL,
		email		VARCHAR(255) NOT NULL,
		position		VARCHAR(255) NOT NULL,
		phone_number    VARCHAR(20) NOT NULL,
		active VARCHAR(10) NOT NULL
	);	
	`)
	if err != nil {
		return errors.Wrap(err, "failed to create organizer table")
	}

	// Index id and name
	_, err = tx.Exec(`
	CREATE INDEX idx_organizer_id ON organizer (id);
	CREATE INDEX idx_organizer_name ON organizer (name);
	`)
	if err != nil {
		return errors.Wrap(err, "failed to create indexes")
	}

	return nil
}

func down_00002(tx *sql.Tx) error {
	// Drop table
	_, err := tx.Exec("DROP TABLE IF EXISTS organizer;")
	if err != nil {
		return errors.Wrap(err, "failed to drop organizer table")
	}

	return nil
}
