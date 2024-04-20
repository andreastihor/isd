package go_sql

import (
	"database/sql"
	"runtime"
	"sync"

	"github.com/pkg/errors"
	"github.com/pressly/goose"
)

var IsProd bool

func init() {
	addMigration(up_00001, down_00001)
}

type Migration struct {
	filename string
	up       func(*sql.Tx) error
	down     func(*sql.Tx) error
}

type syncMap struct {
	sync.Mutex
	mig map[string]Migration
}

var m = syncMap{Mutex: sync.Mutex{}, mig: make(map[string]Migration, 0)}

func RegisterMigrations() {
	m.Lock()
	defer m.Unlock()
	for _, v := range m.mig {
		goose.AddNamedMigration(v.filename, v.up, v.down)
	}
}

func addMigration(up func(*sql.Tx) error, down func(*sql.Tx) error) {
	m.Lock()
	defer m.Unlock()
	_, filename, _, _ := runtime.Caller(1)
	m.mig[filename] = Migration{
		filename: filename,
		up:       up,
		down:     down,
	}
}

var IsLive bool // current app env is live or not

func up_00001(tx *sql.Tx) error {
	// Create table
	_, err := tx.Exec(`
	CREATE TABLE club (
		id              UUID PRIMARY KEY,
		name            VARCHAR(255) NOT NULL,
		country         VARCHAR(255) NOT NULL,
		province        VARCHAR(255) NOT NULL,
		district        VARCHAR(255) NOT NULL,
		establish_date  DATE NOT NULL,
		logo            VARCHAR(255) NOT NULL,
		address         TEXT NOT NULL,
		email_pic		VARCHAR(255) NOT NULL,
		pic             VARCHAR(255) NOT NULL,
		discipline       VARCHAR(255) NOT NULL,
		phone_number    VARCHAR(20) NOT NULL,
		active VARCHAR(10) NOT NULL
	);	
	`)
	if err != nil {
		return errors.Wrap(err, "failed to create club table")
	}

	// Index id and name
	_, err = tx.Exec(`
	CREATE INDEX idx_club_id ON club (id);
	CREATE INDEX idx_club_name ON club (name);
	`)
	if err != nil {
		return errors.Wrap(err, "failed to create indexes")
	}

	return nil
}

func down_00001(tx *sql.Tx) error {
	// Drop table
	_, err := tx.Exec("DROP TABLE IF EXISTS club;")
	if err != nil {
		return errors.Wrap(err, "failed to drop club table")
	}

	return nil
}
