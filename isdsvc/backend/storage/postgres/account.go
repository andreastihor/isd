package postgres

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/andreastihor/isd/isdsvc/backend/storage"
	"github.com/pkg/errors"
)

// SQL queries for account management
const (
	createAccountSQL = `
        INSERT INTO account (
            id,
            name,
            email,
            password
        ) VALUES (
            :id,
            :name,
            :email,
            :password
        ) RETURNING id`

	retrieveAccountSQL = `
        SELECT
            id,
            name,
            email,
            password
        FROM account
		WHERE TRUE `

	updateAccountSQL = `
        UPDATE account SET
            name = :name,
            email = :email,
            password = :password
        WHERE
            id = :id`

	deleteAccountSQL = "DELETE FROM account WHERE id = :id"

	signInSQL = "SELECT token,expired from bearer_token where account_id = :account_id"
)

// CreateAccount adds one account record to the store.
func (s *Storage) CreateAccount(ctx context.Context, account *storage.Account) (string, error) {
	var id string
	// Prepare statement for creating account data
	nstmt, err := s.Db.PrepareNamedContext(ctx, createAccountSQL)
	if err != nil {
		s.Logger.WithError(err).Error("failed to prepare statement for creating account data")
		return "", errors.Wrap(err, "failed to prepare statement for creating account data")
	}
	defer nstmt.Close()
	// Execute query
	args := map[string]interface{}{
		"id":       account.ID,
		"name":     account.Name,
		"email":    account.Email,
		"password": account.Password,
	}
	if err := nstmt.QueryRowContext(ctx, args).Scan(&id); err != nil {
		s.Logger.WithField("args", args).WithError(err).Error("failed to create account data")
		return "", errors.Wrap(err, "failed to create account data")
	}
	return id, nil
}

// GetAccounts retrieves account records from the store.
func (s *Storage) GetAccounts(ctx context.Context, accountParams storage.GetAccountParams) ([]storage.Account, error) {
	accounts := []storage.Account{}
	var query string
	var params map[string]interface{}
	query, params = s.buildAccountQuery(accountParams)
	// Prepare statement for retrieving account data
	nstmt, err := s.Db.PrepareNamedContext(ctx, query)
	if err != nil {
		s.Logger.WithError(err).Error("failed to prepare statement for retrieving account data")
		return nil, errors.Wrap(err, "failed to prepare statement for retrieving account data")
	}
	defer nstmt.Close()
	// Execute query
	if err = nstmt.SelectContext(ctx, &accounts, params); err != nil {
		s.Logger.Error("failed to retrieve account data")
		return nil, errors.Wrap(err, "failed to retrieve account data")
	}
	return accounts, nil
}

// UpdateAccount updates account information in the store.
func (s *Storage) UpdateAccount(ctx context.Context, account *storage.Account) error {
	// Prepare statement for updating account data
	nstmt, err := s.Db.PrepareNamedContext(ctx, updateAccountSQL)
	if err != nil {
		s.Logger.WithError(err).Error("failed to prepare statement for updating account data")
		return errors.Wrap(err, "failed to prepare statement for updating account data")
	}
	defer nstmt.Close()
	// Execute query
	args := map[string]interface{}{
		"id":       account.ID,
		"name":     account.Name,
		"email":    account.Email,
		"password": account.Password,
	}
	if _, err := nstmt.ExecContext(ctx, args); err != nil {
		s.Logger.WithField("args", args).WithError(err).Error("failed to update account data")
		return errors.Wrap(err, "failed to update account data")
	}
	return nil
}

// DeleteAccount deletes account records from the store.
func (s *Storage) DeleteAccount(ctx context.Context, id string) error {
	// Prepare query for deleting account data
	stmt, err := s.Db.PrepareNamed(deleteAccountSQL)
	if err != nil {
		return fmt.Errorf("error preparing query deleteAccountSQL: %w", err)
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
func (s *Storage) buildAccountQuery(filter storage.GetAccountParams) (string, map[string]interface{}) {
	query := bytes.NewBufferString(retrieveAccountSQL)
	params := map[string]interface{}{}
	conditions := []string{}

	if len(filter.AccountIDs) > 0 {
		var condVal bytes.Buffer
		for k, v := range filter.AccountIDs {
			key := fmt.Sprintf("AccountID_%d", k)
			condVal.WriteString(fmt.Sprintf(":%s", key))
			params[key] = v
			if k != len(filter.AccountIDs)-1 {
				condVal.WriteString(",")
			}
		}
		conditions = append(conditions, fmt.Sprintf("id IN (%s)", condVal.String()))
	}

	if filter.Email != "" {
		conditions = append(conditions, "email = :email")
		params["email"] = filter.Email
	}

	if len(conditions) > 0 {
		query.WriteString(" AND ")
		query.WriteString(strings.Join(conditions, " AND "))
	}

	return query.String(), params
}

// SignIn gives the token back if user_id exist and token hasn't expired
func (s *Storage) SignIn(ctx context.Context, id string) error {
	// Prepare query for deleting account data
	stmt, err := s.Db.PrepareNamed(signInSQL)
	if err != nil {
		return fmt.Errorf("error preparing query deleteAccountSQL: %w", err)
	}
	defer stmt.Close()
	// Execute query
	if _, err := stmt.ExecContext(ctx, map[string]interface{}{
		"account_id": id,
	}); err != nil {
		return err
	}
	return nil
}

// Token

const (
	createTokenSQL = `
		INSERT INTO token (
			account_id,
			token,
			expired
		) VALUES (
			:account_id,
			:token,
			:expired
		) RETURNING account_id`

	retrieveTokenSQL = `
		SELECT
			account_id,
			token,
			expired
		FROM token`

	deleteTokenSQL = "DELETE FROM token WHERE account_id = :account_id"
)

// CreateToken adds one token record to the store.
func (s *Storage) CreateToken(ctx context.Context, accountID, token string) error {
	var id string

	expired := time.Now().Add(24 * time.Hour) // Set token expiration to 24 hours from now

	nstmt, err := s.Db.PrepareNamedContext(ctx, createTokenSQL)
	if err != nil {
		s.Logger.WithError(err).Error("failed to prepare statement for creating token data")
		return errors.Wrap(err, "failed to prepare statement for creating token data")
	}
	defer nstmt.Close()

	args := map[string]interface{}{
		"account_id": accountID,
		"token":      token,
		"expired":    expired,
	}

	if err := nstmt.QueryRowContext(ctx, args).Scan(&id); err != nil {
		s.Logger.WithField("args", args).WithError(err).Error("failed to create token data")
		return errors.Wrap(err, "failed to create token data")
	}
	return nil
}

// GetTokens retrieves token records from the store.
func (s *Storage) GetTokens(ctx context.Context, accountID string) (*storage.BearerToken, error) {
	token := &storage.BearerToken{}
	query, params := s.buildTokenQuery(accountID)

	nstmt, err := s.Db.PrepareNamedContext(ctx, query)
	if err != nil {
		s.Logger.WithError(err).Error("failed to prepare statement for retrieving token data")
		return nil, errors.Wrap(err, "failed to prepare statement for retrieving token data")
	}
	defer nstmt.Close()

	if err = nstmt.GetContext(ctx, token, params); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No token found, return nil without an error
			return nil, nil
		}
		s.Logger.WithField("params", params).WithError(err).Error("failed to retrieve token data")
		return nil, errors.Wrap(err, "failed to retrieve token data")
	}
	return token, nil
}

// DeleteToken deletes token records from the store.
func (s *Storage) DeleteToken(ctx context.Context, accountID string) error {
	stmt, err := s.Db.PrepareNamed(deleteTokenSQL)
	if err != nil {
		return fmt.Errorf("error preparing query deleteTokenSQL: %w", err)
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, map[string]interface{}{
		"account_id": accountID,
	}); err != nil {
		return err
	}
	return nil
}

func (s *Storage) buildTokenQuery(accountID string) (string, map[string]interface{}) {
	query := bytes.NewBufferString(retrieveTokenSQL)
	params := map[string]interface{}{
		"account_id": accountID,
	}
	query.WriteString(" WHERE account_id = :account_id")
	return query.String(), params
}
