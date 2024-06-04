package go_sql

import (
	"database/sql"

	"github.com/pkg/errors"
)

func init() {
	addMigration(up_00005, down_00005)
}

func up_00005(tx *sql.Tx) error {
	// Create table
	_, err := tx.Exec(`
	CREATE TABLE token (
		account_id              VARCHAR(36) PRIMARY KEY,
		token		TEXT NOT NULL,
		expired date NOT NULL
	);	
	`)
	if err != nil {
		return errors.Wrap(err, "failed to create token table")
	}

	// Index id and name
	_, err = tx.Exec(`
	CREATE INDEX idx_token_account_id ON token (account_id);
	`)
	if err != nil {
		return errors.Wrap(err, "failed to create indexes")
	}

	return nil
}

func down_00005(tx *sql.Tx) error {
	// Drop table
	_, err := tx.Exec("DROP TABLE IF EXISTS token;")
	if err != nil {
		return errors.Wrap(err, "failed to drop token table")
	}

	return nil
}
