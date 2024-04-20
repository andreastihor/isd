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

// Organizer represents organizer information.
type Organizer struct {
	ID                string            `json:"id" db:"id"`
	Name              string            `json:"name" db:"name"`
	Position          string            `json:"position" db:"position"`
	Club              Club              `json:"club" db:"-"`
	RegisterDate      time.Time         `json:"register_date" db:"register_date"`
	PhoneNumber       string            `json:"phone_number" db:"phone_number"`
	Active            util.OptionalBool `json:"active" db:"active"`
	Email             string            `json:"email" db:"email"`
	ClubID            string            `json:"club_id" db:"club_id"`
	ClubName          string            `json:"club_name" db:"club_name"`
	ClubCountry       string            `json:"club_country" db:"club_country"`
	ClubProvince      string            `json:"club_province" db:"club_province"`
	ClubDistrict      string            `json:"club_district" db:"club_district"` // kabupaten
	ClubEstablishDate time.Time         `json:"club_establish_date" db:"club_establish_date"`
	ClubLogo          string            `json:"club_logo" db:"club_logo"`
	ClubAddress       string            `json:"club_address" db:"club_address"`
	ClubEmailPIC      string            `json:"club_email_pic" db:"club_email_pic"`
	ClubPic           string            `json:"club_pic" db:"club_pic"`
	ClubDiscipline    string            `json:"club_discipline" db:"club_discipline"`
	ClubPhoneNumber   string            `json:"club_phone_number" db:"club_phone_number"`
	ClubActive        util.OptionalBool `json:"club_active" db:"club_active"`
}

// Storage provides storage operations for direct
type Storage interface {
	GetDBConn() *sql.DB

	ClubStore
	OrganizerStore
}

// ClubStore is an interface for club storage operations.
type ClubStore interface {
	CreateClub(ctx context.Context, club *Club) (string, error)
	GetClubs(ctx context.Context, clubIDs ...string) ([]Club, error)
	UpdateClub(ctx context.Context, club *Club) error
	DeleteClub(ctx context.Context, clubID string) error
}

type OrganizerStore interface {
	CreateOrganizer(ctx context.Context, organizer *Organizer) (string, error)
	GetOrganizers(ctx context.Context, organizerIDs ...string) ([]Organizer, error)
	UpdateOrganizer(ctx context.Context, organizer *Organizer) error
	DeleteOrganizer(ctx context.Context, organizerID string) error
}
