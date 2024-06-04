package go_sql

import (
	"database/sql"

	"github.com/pkg/errors"
)

func init() {
	addMigration(up_00004, down_00004)
}

func up_00004(tx *sql.Tx) error {
	// Create table
	_, err := tx.Exec(`
	CREATE TABLE account (
		id              VARCHAR(36) PRIMARY KEY,
		name            VARCHAR(255) NOT NULL,
		email		VARCHAR(255) NOT NULL,
		password		TEXT NOT NULL
	);	
	`)
	if err != nil {
		return errors.Wrap(err, "failed to create account table")
	}

	// Index id and name
	_, err = tx.Exec(`
	CREATE INDEX idx_account_id ON account (id);
	CREATE INDEX idx_password_id ON account (password);
	`)
	if err != nil {
		return errors.Wrap(err, "failed to create indexes")
	}

	return nil
}

func down_00004(tx *sql.Tx) error {
	// Drop table
	_, err := tx.Exec("DROP TABLE IF EXISTS account;")
	if err != nil {
		return errors.Wrap(err, "failed to drop account table")
	}

	return nil
}
