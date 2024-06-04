package go_sql

import (
	"database/sql"

	"github.com/pkg/errors"
)

func init() {
	addMigration(up_00006, down_00006)
}

func up_00006(tx *sql.Tx) error {
	// Create table
	_, err := tx.Exec(`
	CREATE TABLE verif_code (
		account_id              VARCHAR(36) PRIMARY KEY,
		code		TEXT NOT NULL
	);	
	`)
	if err != nil {
		return errors.Wrap(err, "failed to create verif_code table")
	}

	// Index id and name
	_, err = tx.Exec(`
	CREATE INDEX idx_verif_code_account_id ON verif_code (account_id);
	`)
	if err != nil {
		return errors.Wrap(err, "failed to create indexes")
	}

	return nil
}

func down_00006(tx *sql.Tx) error {
	// Drop table
	_, err := tx.Exec("DROP TABLE IF EXISTS verif_code;")
	if err != nil {
		return errors.Wrap(err, "failed to drop toverif_codeken table")
	}

	return nil
}
