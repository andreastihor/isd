package postgres

import (
	"bytes"
	"context"
	"fmt"

	"github.com/andreastihor/isd/isdsvc/backend/storage"
	"github.com/pkg/errors"
)

// SQL queries for organizer management
const (
	retrieveOrganizerSQL = `
        SELECT
            o.id,
            o.name,
            o.position,
            o.register_date,
            o.phone_number,
            o.active,
            o.email,
            c.id as club_id,
            c.name as club_name,
            c.country as club_country,
            c.province as club_province,
            c.district as club_district,
            c.establish_date as club_establish_date,
            c.logo as club_logo,
            c.address as club_address,
            c.email_pic as club_email_pic,
            c.pic as club_pic,
            c.discipline as club_discipline,
            c.phone_number as club_phone_number,
            c.active as club_active
        FROM organizer o
        JOIN club c ON o.club_id = c.id
    `

	createOrganizerSQL = `
        INSERT INTO organizer (
            id,
            name,
            position,
            club_id,
            register_date,
            phone_number,
            active,
            email
        ) VALUES (
            :id,
            :name,
            :position,
            :club_id,
            :register_date,
            :phone_number,
            :active,
            :email
        ) RETURNING id
    `

	updateOrganizerSQL = `
    UPDATE organizer SET
        name = :name,
        position = :position,
        register_date = :register_date,
        phone_number = :phone_number,
        active = :active,
        email = :email
    WHERE
        id = :id
`

	deleteOrganizerSQL = "DELETE FROM organizer WHERE id = :id"
)

// CreateOrganizer adds one organizer record to the store.
func (s *Storage) CreateOrganizer(ctx context.Context, organizer *storage.Organizer) (string, error) {
	var id string
	// Prepare statement for creating organizer data
	nstmt, err := s.Db.PrepareNamedContext(ctx, createOrganizerSQL)
	if err != nil {
		s.Logger.WithError(err).Error("failed to prepare statement for creating organizer data")
		return "", errors.Wrap(err, "failed to prepare statement for creating organizer data")
	}
	defer nstmt.Close()
	// Execute query
	args := map[string]interface{}{
		"id":            organizer.ID,
		"name":          organizer.Name,
		"position":      organizer.Position,
		"club_id":       organizer.ClubID,
		"register_date": organizer.RegisterDate,
		"phone_number":  organizer.PhoneNumber,
		"active":        organizer.Active,
		"email":         organizer.Email,
	}
	if err := nstmt.QueryRowContext(ctx, args).Scan(&id); err != nil {
		s.Logger.WithField("args", args).WithError(err).Error("failed to create organizer data")
		return "", errors.Wrap(err, "failed to create organizer data")
	}
	return id, nil
}

// GetOrganizers retrieves organizer records from the store.
func (s *Storage) GetOrganizers(ctx context.Context, IDs ...string) ([]storage.Organizer, error) {
	organizers := []storage.Organizer{}
	var query string
	var params map[string]interface{}
	if len(IDs) == 0 {
		query = retrieveOrganizerSQL
		params = map[string]interface{}{}
	} else {
		query, params = s.buildOrganizerQuery(IDs...)
	}
	// Prepare statement for retrieving organizer data
	nstmt, err := s.Db.PrepareNamedContext(ctx, query)
	if err != nil {
		s.Logger.WithError(err).Error("failed to prepare statement for retrieving organizer data")
		return nil, errors.Wrap(err, "failed to prepare statement for retrieving organizer data")
	}
	defer nstmt.Close()
	// Execute query
	if err = nstmt.SelectContext(ctx, &organizers, params); err != nil {
		s.Logger.Error("failed to retrieve organizer data")
		return nil, errors.Wrap(err, "failed to retrieve organizer data")
	}
	return organizers, nil
}

// UpdateOrganizer updates organizer information in the store.
func (s *Storage) UpdateOrganizer(ctx context.Context, organizer *storage.Organizer) error {
	// Prepare statement for updating organizer data
	nstmt, err := s.Db.PrepareNamedContext(ctx, updateOrganizerSQL)
	if err != nil {
		s.Logger.WithError(err).Error("failed to prepare statement for updating organizer data")
		return errors.Wrap(err, "failed to prepare statement for updating organizer data")
	}
	defer nstmt.Close()
	// Execute query
	args := map[string]interface{}{
		"id":            organizer.ID,
		"name":          organizer.Name,
		"position":      organizer.Position,
		"register_date": organizer.RegisterDate,
		"phone_number":  organizer.PhoneNumber,
		"active":        organizer.Active,
		"email":         organizer.Email,
	}
	if _, err := nstmt.ExecContext(ctx, args); err != nil {
		s.Logger.WithField("args", args).WithError(err).Error("failed to update organizer data")
		return errors.Wrap(err, "failed to update organizer data")
	}
	return nil
}

// DeleteOrganizer deletes organizer records from the store.
func (s *Storage) DeleteOrganizer(ctx context.Context, id string) error {
	// Prepare query for deleting organizer data
	stmt, err := s.Db.PrepareNamed(deleteOrganizerSQL)
	if err != nil {
		return fmt.Errorf("error preparing query deleteOrganizerSQL: %w", err)
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

func (s *Storage) buildOrganizerQuery(IDs ...string) (string, map[string]interface{}) {
	query := bytes.NewBufferString(retrieveOrganizerSQL)
	params := map[string]interface{}{}
	var condVal bytes.Buffer
	fLen := len(IDs)
	for k, v := range IDs {
		key := fmt.Sprintf("OrganizerID_%d", k)
		condVal.WriteString(fmt.Sprintf(":%s", key))
		params[key] = v
		if k != fLen-1 {
			condVal.WriteString(",")
		}
	}
	query.WriteString(fmt.Sprintf(" WHERE o.id IN (%s)", condVal.String()))
	return query.String(), params
}
