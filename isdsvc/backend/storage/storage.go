//go:generate mockgen -source=storage.go -destination=mock_storage/mock_storage.go
//go:generate gofumpt -w mock_storage/mock_storage.go
package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/andreastihor/isd/isdsvc/backend/util"
	_ "github.com/lib/pq"
)

// Club represents the data structure for a club in storage.
type Club struct {
	ID            string            `json:"id" db:"id"`
	Name          string            `json:"name" db:"name"`
	Country       string            `json:"country" db:"country"`
	Province      string            `json:"province" db:"province"`
	District      string            `json:"district" db:"district"` // kabupaten
	EstablishDate time.Time         `json:"establish_date" db:"establish_date"`
	Logo          string            `json:"logo" db:"logo"`
	Address       string            `json:"address" db:"address"`
	EmailPIC      string            `json:"email_pic" db:"email_pic"`
	Pic           string            `json:"pic" db:"pic"`
	Discipline    string            `json:"discipline" db:"discipline"`
	PhoneNumber   string            `json:"phone_number" db:"phone_number"`
	Active        util.OptionalBool `json:"active" db:"active"`
}

// Storage provides storage operations for direct
type Storage interface {
	GetDBConn() *sql.DB

	ClubStore
}

// ClubStore is an interface for club storage operations.
type ClubStore interface {
	CreateClub(ctx context.Context, club *Club) (string, error)
	GetClubs(ctx context.Context, clubIDs ...string) ([]Club, error)
	UpdateClub(ctx context.Context, club *Club) error
	DeleteClub(ctx context.Context, clubID string) error
}
