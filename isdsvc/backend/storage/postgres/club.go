package postgres

import (
	"bytes"
	"context"
	"fmt"

	"github.com/andreastihor/isd/isdsvc/backend/storage"
	"github.com/pkg/errors"
)

// SQL queries for club management
const (
	createClubSQL = `
        INSERT INTO club (
            id,
            name,
            country,
            province,
            district,
            establish_date,
            logo,
            address,
            email_pic,
            pic,
            discipline,
            phone_number,
			active
        ) VALUES (
            :id,
            :name,
            :country,
            :province,
            :district,
            :establish_date,
            :logo,
            :address,
            :email_pic,
            :pic,
            :discipline,
            :phone_number,
			:active
        ) RETURNING id`

	retrieveClubSQL = `
        SELECT
            id,
            name,
            country,
            province,
            district,
            establish_date,
            logo,
            address,
            email_pic,
            pic,
            discipline,
            phone_number,
			active
        FROM club`

	updateClubSQL = `
        UPDATE club SET
            name = :name,
            country = :country,
            province = :province,
            district = :district,
            establish_date = :establish_date,
            logo = :logo,
            address = :address,
            email_pic = :email_pic,
            pic = :pic,
            discipline = :discipline,
            phone_number = :phone_number,
			active = :active
        WHERE
            id = :id`

	deleteClubSQL = "DELETE FROM club WHERE id = :id"
)

// CreateClub adds one club record to the store.
func (s *Storage) CreateClub(ctx context.Context, club *storage.Club) (string, error) {
	var id string
	// Prepare statement for creating club data
	nstmt, err := s.Db.PrepareNamedContext(ctx, createClubSQL)
	if err != nil {
		s.Logger.WithError(err).Error("failed to prepare statement for creating club data")
		return "", errors.Wrap(err, "failed to prepare statement for creating club data")
	}
	defer nstmt.Close()
	// Execute query
	args := map[string]interface{}{
		"id":             club.ID,
		"name":           club.Name,
		"country":        club.Country,
		"province":       club.Province,
		"district":       club.District,
		"establish_date": club.EstablishDate,
		"logo":           club.Logo,
		"address":        club.Address,
		"email_pic":      club.EmailPIC,
		"pic":            club.Pic,
		"discipline":     club.Discipline,
		"phone_number":   club.PhoneNumber,
		"active":         club.Active,
	}
	if err := nstmt.QueryRowContext(ctx, args).Scan(&id); err != nil {
		s.Logger.WithField("args", args).WithError(err).Error("failed to create club data")
		return "", errors.Wrap(err, "failed to create club data")
	}
	return id, nil
}

// GetClubs retrieves club records from the store.
func (s *Storage) GetClubs(ctx context.Context, IDs ...string) ([]storage.Club, error) {
	clubs := []storage.Club{}
	var query string
	var params map[string]interface{}
	if len(IDs) == 0 {
		query = retrieveClubSQL
		params = map[string]interface{}{}
	} else {
		query, params = s.buildClubQuery(IDs...)
	}
	// Prepare statement for retrieving club data
	nstmt, err := s.Db.PrepareNamedContext(ctx, query)
	if err != nil {
		s.Logger.WithError(err).Error("failed to prepare statement for retrieving club data")
		return nil, errors.Wrap(err, "failed to prepare statement for retrieving club data")
	}
	defer nstmt.Close()
	// Execute query
	if err = nstmt.SelectContext(ctx, &clubs, params); err != nil {
		s.Logger.Error("failed to retrieve club data")
		return nil, errors.Wrap(err, "failed to retrieve club data")
	}
	return clubs, nil
}

// UpdateClub updates club information in the store.
func (s *Storage) UpdateClub(ctx context.Context, club *storage.Club) error {
	// Prepare statement for updating club data
	nstmt, err := s.Db.PrepareNamedContext(ctx, updateClubSQL)
	if err != nil {
		s.Logger.WithError(err).Error("failed to prepare statement for updating club data")
		return errors.Wrap(err, "failed to prepare statement for updating club data")
	}
	defer nstmt.Close()
	// Execute query
	args := map[string]interface{}{
		"id":             club.ID,
		"name":           club.Name,
		"country":        club.Country,
		"province":       club.Province,
		"district":       club.District,
		"establish_date": club.EstablishDate,
		"logo":           club.Logo,
		"address":        club.Address,
		"email_pic":      club.EmailPIC,
		"pic":            club.Pic,
		"discipline":     club.Discipline,
		"phone_number":   club.PhoneNumber,
		"active":         club.Active,
	}
	if _, err := nstmt.ExecContext(ctx, args); err != nil {
		s.Logger.WithField("args", args).WithError(err).Error("failed to update club data")
		return errors.Wrap(err, "failed to update club data")
	}
	return nil
}

// DeleteClub deletes club records from the store.
func (s *Storage) DeleteClub(ctx context.Context, id string) error {
	// Prepare query for deleting club data
	stmt, err := s.Db.PrepareNamed(deleteClubSQL)
	if err != nil {
		return fmt.Errorf("error preparing query deleteClubSQL: %w", err)
	}
	defer stmt.Close()
	// Execute query
	if _, err := stmt.ExecContext(ctx, map[string]interface{}{
		"id": id,
	}); err != nil {
		return err
	}
	return nil
}

func (s *Storage) buildClubQuery(IDs ...string) (string, map[string]interface{}) {
	query := bytes.NewBufferString(retrieveClubSQL)
	params := map[string]interface{}{}
	var condVal bytes.Buffer
	fLen := len(IDs)
	for k, v := range IDs {
		key := fmt.Sprintf("ClubID_%d", k)
		condVal.WriteString(fmt.Sprintf(":%s", key))
		params[key] = v
		if k != fLen-1 {
			condVal.WriteString(",")
		}
	}
	query.WriteString(fmt.Sprintf(" WHERE id IN (%s)", condVal.String()))
	return query.String(), params
}
