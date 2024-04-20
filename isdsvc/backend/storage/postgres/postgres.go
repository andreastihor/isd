package postgres

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
)

// Storage provides storage operations for direct
type Storage struct {
	StorageUtil
}

// NewStorage initializes a new Storage instance.
func NewStorage(logger logrus.FieldLogger, dbstring string, opts ...func(*Storage)) (*Storage, error) {
	db, err := NewDBStorage(logger, dbstring)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to postgres '%s': %v", dbstring, err)
	}
	strg := &Storage{
		*db,
	}

	for _, o := range opts {
		o(strg)
	}
	return strg, nil
}

// // NewStorageWithTracing initializes a new Storage instance with tracing.
// func NewStorageWithTracing(logger logrus.FieldLogger, config *viper.Viper, opts ...func(*Storage)) (*Storage, error) {
// 	s, err := utils.NewStorageWithTracing(logger, config, utils.FuzzConfig{})
// 	if err != nil {
// 		logging.WithError(err, logger).Error("NewStorageWithTracing: failed to create storage from config")
// 		return nil, err
// 	}
// 	strg := &Storage{
// 		*s,
// 		nil,
// 		utils.StorageMetadata{},
// 	}
// 	for _, o := range opts {
// 		o(strg)
// 	}
// 	return strg, err
// }

// NewDBStorage returns a new StorageUtil from the provides psql database string
func NewDBStorage(logger logrus.FieldLogger, dbstring string) (*StorageUtil, error) {
	db, err := otelsqlx.Connect("postgres", dbstring)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to postgres '%s': %v", dbstring, err)
	}

	// TODO: is this a sane default?
	// The current max_connections in postgres is 100.
	db.SetMaxOpenConns(50)
	db.SetConnMaxLifetime(time.Hour)
	return &StorageUtil{Logger: logger, Db: db}, nil
}

// StorageUtil provides a wrapper around an sql database and provides
// required methods for interacting with the database
type StorageUtil struct {
	Logger logrus.FieldLogger
	Db     *sqlx.DB
}
