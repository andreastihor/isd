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
	createVerifCodeSQL = `
        INSERT INTO verif_code (
            account_id,
            code,
           
        ) VALUES (
            :id,
            :code,
        ) `

	retrieveVerifCodeSQL = `
        SELECT
           code
        FROM verif_code`

	deleteVerifCodeSQL = "DELETE FROM code WHERE code = :code"
)

// CreateVerifCode adds a verification code record to the store.
func (s *Storage) CreateVerifCode(ctx context.Context, verifCode *storage.VerifCode) error {
	// Prepare statement for creating verification code data
	nstmt, err := s.Db.PrepareNamedContext(ctx, createVerifCodeSQL)
	if err != nil {
		s.Logger.WithError(err).Error("failed to prepare statement for creating verification code data")
		return errors.Wrap(err, "failed to prepare statement for creating verification code data")
	}
	defer nstmt.Close()

	// Execute query
	args := map[string]interface{}{
		"account_id": verifCode.AccountID,
		"code":       verifCode.Code,
	}
	if _, err := nstmt.ExecContext(ctx, args); err != nil {
		s.Logger.WithField("args", args).WithError(err).Error("failed to create verification code data")
		return errors.Wrap(err, "failed to create verification code data")
	}
	return nil
}

// GetVerifCodes retrieves verification code records from the store.
func (s *Storage) GetVerifCodes(ctx context.Context, accountIDs ...string) ([]storage.VerifCode, error) {
	verifCodes := []storage.VerifCode{}
	var query string
	var params map[string]interface{}
	if len(accountIDs) == 0 {
		query = retrieveVerifCodeSQL
		params = map[string]interface{}{}
	} else {
		query, params = s.buildVerifCodeQuery(accountIDs...)
	}

	// Prepare statement for retrieving verification code data
	nstmt, err := s.Db.PrepareNamedContext(ctx, query)
	if err != nil {
		s.Logger.WithError(err).Error("failed to prepare statement for retrieving verification code data")
		return nil, errors.Wrap(err, "failed to prepare statement for retrieving verification code data")
	}
	defer nstmt.Close()

	// Execute query
	if err = nstmt.SelectContext(ctx, &verifCodes, params); err != nil {
		s.Logger.Error("failed to retrieve verification code data")
		return nil, errors.Wrap(err, "failed to retrieve verification code data")
	}
	return verifCodes, nil
}

// DeleteVerifCode deletes a verification code record from the store.
func (s *Storage) DeleteVerifCode(ctx context.Context, code string) error {
	// Prepare query for deleting verification code data
	stmt, err := s.Db.PrepareNamed(deleteVerifCodeSQL)
	if err != nil {
		s.Logger.WithError(err).Error("failed to prepare statement for deleting verification code data")
		return errors.Wrap(err, "failed to prepare statement for deleting verification code data")
	}
	defer stmt.Close()

	// Execute query
	if _, err := stmt.ExecContext(ctx, map[string]interface{}{
		"code": code,
	}); err != nil {
		s.Logger.WithField("code", code).WithError(err).Error("failed to delete verification code data")
		return errors.Wrap(err, "failed to delete verification code data")
	}
	return nil
}

func (s *Storage) buildVerifCodeQuery(accountIDs ...string) (string, map[string]interface{}) {
	query := bytes.NewBufferString(retrieveVerifCodeSQL)
	params := map[string]interface{}{}
	var condVal bytes.Buffer
	fLen := len(accountIDs)
	for k, v := range accountIDs {
		key := fmt.Sprintf("AccountID_%d", k)
		condVal.WriteString(fmt.Sprintf(":%s", key))
		params[key] = v
		if k != fLen-1 {
			condVal.WriteString(",")
		}
	}
	query.WriteString(fmt.Sprintf(" WHERE account_id IN (%s)", condVal.String()))
	return query.String(), params
}
