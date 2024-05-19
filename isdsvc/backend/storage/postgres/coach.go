package postgres

import (
	"bytes"
	"context"
	"fmt"

	"github.com/andreastihor/isd/isdsvc/backend/storage"
	"github.com/pkg/errors"
)

// SQL queries for coach management
const (
	createCoachSQL = `
        INSERT INTO coach (
            id,
            name,
            dob,
            phone_number,
            gender,
            email,
            discipline,
            register_date,
			active
        ) VALUES (
            :id,
            :name,
            :dob,
            :phone_number,
            :gender,
            :email,
            :discipline,
            :register_date,
			:active
        ) RETURNING id`

	retrieveCoachSQL = `
        SELECT
            id,
            name,
            dob,
            phone_number,
            gender,
            email,
            discipline,
            register_date,
			active
        FROM coach`

	updateCoachSQL = `
        UPDATE coach SET
            name = :name,
            dob = :dob,
            phone_number = :phone_number,
            gender = :gender,
            email = :email,
            discipline = :discipline,
            register_date = :register_date,
			active = :active
        WHERE
            id = :id`

	deleteCoachSQL = "DELETE FROM coach WHERE id = :id"
)

// CreateCoach adds one coach record to the store.
func (s *Storage) CreateCoach(ctx context.Context, coach *storage.Coach) (string, error) {
	var id string
	// Prepare statement for creating coach data
	nstmt, err := s.Db.PrepareNamedContext(ctx, createCoachSQL)
	if err != nil {
		s.Logger.WithError(err).Error("failed to prepare statement for creating coach data")
		return "", errors.Wrap(err, "failed to prepare statement for creating coach data")
	}
	defer nstmt.Close()
	// Execute query
	args := map[string]interface{}{
		"id":            coach.ID,
		"name":          coach.Name,
		"dob":           coach.DOB,
		"phone_number":  coach.PhoneNumber,
		"gender":        coach.Gender,
		"email":         coach.Email,
		"discipline":    coach.Discipline,
		"register_date": coach.RegisterDate,
		"active":        coach.Active,
	}
	if err := nstmt.QueryRowContext(ctx, args).Scan(&id); err != nil {
		s.Logger.WithField("args", args).WithError(err).Error("failed to create coach data")
		return "", errors.Wrap(err, "failed to create coach data")
	}
	return id, nil
}

// GetCoaches retrieves coach records from the store.
func (s *Storage) GetCoaches(ctx context.Context, IDs ...string) ([]storage.Coach, error) {
	coaches := []storage.Coach{}
	var query string
	var params map[string]interface{}
	if len(IDs) == 0 {
		query = retrieveCoachSQL
		params = map[string]interface{}{}
	} else {
		query, params = s.buildCoachQuery(IDs...)
	}
	// Prepare statement for retrieving coach data
	nstmt, err := s.Db.PrepareNamedContext(ctx, query)
	if err != nil {
		s.Logger.WithError(err).Error("failed to prepare statement for retrieving coach data")
		return nil, errors.Wrap(err, "failed to prepare statement for retrieving coach data")
	}
	defer nstmt.Close()
	// Execute query
	if err = nstmt.SelectContext(ctx, &coaches, params); err != nil {
		s.Logger.Error("failed to retrieve coach data")
		return nil, errors.Wrap(err, "failed to retrieve coach data")
	}
	return coaches, nil
}

// UpdateCoach updates coach information in the store.
func (s *Storage) UpdateCoach(ctx context.Context, coach *storage.Coach) error {
	// Prepare statement for updating coach data
	nstmt, err := s.Db.PrepareNamedContext(ctx, updateCoachSQL)
	if err != nil {
		s.Logger.WithError(err).Error("failed to prepare statement for updating coach data")
		return errors.Wrap(err, "failed to prepare statement for updating coach data")
	}
	defer nstmt.Close()
	// Execute query
	args := map[string]interface{}{
		"id":            coach.ID,
		"name":          coach.Name,
		"dob":           coach.DOB,
		"phone_number":  coach.PhoneNumber,
		"gender":        coach.Gender,
		"email":         coach.Email,
		"discipline":    coach.Discipline,
		"register_date": coach.RegisterDate,
		"active":        coach.Active,
	}
	if _, err := nstmt.ExecContext(ctx, args); err != nil {
		s.Logger.WithField("args", args).WithError(err).Error("failed to update coach data")
		return errors.Wrap(err, "failed to update coach data")
	}
	return nil
}

// DeleteCoach deletes coach records from the store.
func (s *Storage) DeleteCoach(ctx context.Context, id string) error {
	// Prepare query for deleting coach data
	stmt, err := s.Db.PrepareNamed(deleteCoachSQL)
	if err != nil {
		return fmt.Errorf("error preparing query deleteCoachSQL: %w", err)
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

func (s *Storage) buildCoachQuery(IDs ...string) (string, map[string]interface{}) {
	query := bytes.NewBufferString(retrieveCoachSQL)
	params := map[string]interface{}{}
	var condVal bytes.Buffer
	fLen := len(IDs)
	for k, v := range IDs {
		key := fmt.Sprintf("CoachID_%d", k)
		condVal.WriteString(fmt.Sprintf(":%s", key))
		params[key] = v
		if k != fLen-1 {
			condVal.WriteString(",")
		}
	}
	query.WriteString(fmt.Sprintf(" WHERE id IN (%s)", condVal.String()))
	return query.String(), params
}
