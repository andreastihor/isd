package postgres

import (
	"bytes"
	"context"
	"fmt"

	"github.com/andreastihor/isd/isdsvc/backend/storage"
	"github.com/pkg/errors"
)

// SQL queries for athlete management
const (
	createAthleteSQL = `
        INSERT INTO athlete (
            id,
            club_id,
            name,
            dob,
            phone_number,
            gender,
            email,
            register_date,
            active
        ) VALUES (
            :id,
            :club_id,
            :name,
            :dob,
            :phone_number,
            :gender,
            :email,
            :register_date,
            :active
        ) RETURNING id`

	retrieveAthleteSQL = `
        SELECT
            id,
            club_id,
            name,
            dob,
            phone_number,
            gender,
            email,
            register_date,
            active
        FROM athlete`

	updateAthleteSQL = `
        UPDATE athlete SET
            club_id = :club_id,
            name = :name,
            dob = :dob,
            phone_number = :phone_number,
            gender = :gender,
            email = :email,
            register_date = :register_date,
            active = :active
        WHERE
            id = :id`

	deleteAthleteSQL = "DELETE FROM athlete WHERE id = :id"
)

// CreateAthlete adds one athlete record to the store.
func (s *Storage) CreateAthlete(ctx context.Context, athlete *storage.Athlete) (string, error) {
	var id string
	// Prepare statement for creating athlete data
	nstmt, err := s.Db.PrepareNamedContext(ctx, createAthleteSQL)
	if err != nil {
		s.Logger.WithError(err).Error("failed to prepare statement for creating athlete data")
		return "", errors.Wrap(err, "failed to prepare statement for creating athlete data")
	}
	defer nstmt.Close()
	// Execute query
	args := map[string]interface{}{
		"id":            athlete.ID,
		"club_id":       athlete.ClubID,
		"name":          athlete.Name,
		"dob":           athlete.DOB,
		"phone_number":  athlete.PhoneNumber,
		"gender":        athlete.Gender,
		"email":         athlete.Email,
		"register_date": athlete.RegisterDate,
		"active":        athlete.Active,
	}
	if err := nstmt.QueryRowContext(ctx, args).Scan(&id); err != nil {
		s.Logger.WithField("args", args).WithError(err).Error("failed to create athlete data")
		return "", errors.Wrap(err, "failed to create athlete data")
	}
	return id, nil
}

// GetAthletes retrieves athlete records from the store.
func (s *Storage) GetAthletes(ctx context.Context, IDs ...string) ([]storage.Athlete, error) {
	athletes := []storage.Athlete{}
	var query string
	var params map[string]interface{}
	if len(IDs) == 0 {
		query = retrieveAthleteSQL
		params = map[string]interface{}{}
	} else {
		query, params = s.buildAthleteQuery(IDs...)
	}
	// Prepare statement for retrieving athlete data
	nstmt, err := s.Db.PrepareNamedContext(ctx, query)
	if err != nil {
		s.Logger.WithError(err).Error("failed to prepare statement for retrieving athlete data")
		return nil, errors.Wrap(err, "failed to prepare statement for retrieving athlete data")
	}
	defer nstmt.Close()
	// Execute query
	if err = nstmt.SelectContext(ctx, &athletes, params); err != nil {
		s.Logger.Error("failed to retrieve athlete data")
		return nil, errors.Wrap(err, "failed to retrieve athlete data")
	}
	return athletes, nil
}

// UpdateAthlete updates athlete information in the store.
func (s *Storage) UpdateAthlete(ctx context.Context, athlete *storage.Athlete) error {
	// Prepare statement for updating athlete data
	nstmt, err := s.Db.PrepareNamedContext(ctx, updateAthleteSQL)
	if err != nil {
		s.Logger.WithError(err).Error("failed to prepare statement for updating athlete data")
		return errors.Wrap(err, "failed to prepare statement for updating athlete data")
	}
	defer nstmt.Close()
	// Execute query
	args := map[string]interface{}{
		"id":            athlete.ID,
		"club_id":       athlete.ClubID,
		"name":          athlete.Name,
		"dob":           athlete.DOB,
		"phone_number":  athlete.PhoneNumber,
		"gender":        athlete.Gender,
		"email":         athlete.Email,
		"register_date": athlete.RegisterDate,
		"active":        athlete.Active,
	}
	if _, err := nstmt.ExecContext(ctx, args); err != nil {
		s.Logger.WithField("args", args).WithError(err).Error("failed to update athlete data")
		return errors.Wrap(err, "failed to update athlete data")
	}
	return nil
}

// DeleteAthlete deletes athlete records from the store.
func (s *Storage) DeleteAthlete(ctx context.Context, id string) error {
	// Prepare query for deleting athlete data
	stmt, err := s.Db.PrepareNamed(deleteAthleteSQL)
	if err != nil {
		return fmt.Errorf("error preparing query deleteAthleteSQL: %w", err)
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

func (s *Storage) buildAthleteQuery(IDs ...string) (string, map[string]interface{}) {
	query := bytes.NewBufferString(retrieveAthleteSQL)
	params := map[string]interface{}{}
	var condVal bytes.Buffer
	fLen := len(IDs)
	for k, v := range IDs {
		key := fmt.Sprintf("AthleteID_%d", k)
		condVal.WriteString(fmt.Sprintf(":%s", key))
		params[key] = v
		if k != fLen-1 {
			condVal.WriteString(",")
		}
	}
	query.WriteString(fmt.Sprintf(" WHERE id IN (%s)", condVal.String()))
	return query.String(), params
}
