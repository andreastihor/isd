package isd

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/andreastihor/isd/isdsvc/backend/storage"
	"github.com/andreastihor/isd/isdsvc/backend/util"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Account represents account information.
type Account struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateAccountRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateAccountResponse struct {
	ID string `json:"id"`
}

type UpdateAccountRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateAccountResponse struct {
	Account *Account `json:"account"`
}

type DeleteAccountRequest struct {
	ID string `json:"id"`
}

type GetAccountResponse struct {
	Total    int        `json:"total"`
	Accounts []*Account `json:"accounts"`
}

// CreateAccount creates a new account.
func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	const methodName = "CreateAccount"
	ctx = util.SetCallerMethodToCtx(ctx, methodName)
	l := logrus.New()
	logger := l.WithContext(ctx)
	logger = logger.WithField(methodLogField, methodName)
	logger.Infof("%v ...", methodName)

	var reqBody CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		logger.Errorf("error when decoding request %v", err)
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate the request body
	if err := validateAccountRequest(reqBody); err != nil {
		logger.Errorf("error when validating request %v", err)
		util.HandleError(w, http.StatusBadRequest, err.Message)
		return
	}

	// Call storage method to create account
	account := &storage.Account{
		ID:       uuid.NewString(),
		Name:     reqBody.Name,
		Email:    reqBody.Email,
		Password: reqBody.Password,
	}

	id, err := h.accountStore.CreateAccount(ctx, account)
	if err != nil {
		logger.Errorf("error when creating account %v", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := CreateAccountResponse{
		ID: id,
	}

	json.NewEncoder(w).Encode(resp)
}

// validateAccountRequest validates the fields of AccountRequest.
func validateAccountRequest(req CreateAccountRequest) *util.Error {
	var missingFields []string

	if req.Name == "" {
		missingFields = append(missingFields, "name")
	}

	if req.Email == "" {
		missingFields = append(missingFields, "email")
	}

	if req.Password == "" {
		missingFields = append(missingFields, "password")
	}

	if len(missingFields) > 0 {
		return util.NewError(http.StatusBadRequest, fmt.Sprintf("Missing required fields: [%v]", strings.Join(missingFields, ",")))
	}

	return nil
}

// GetAccount retrieves all accounts.
func (h *Handler) GetAccount(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	const methodName = "GetAccount"
	ctx = util.SetCallerMethodToCtx(ctx, methodName)
	l := logrus.New()
	logger := l.WithContext(ctx)
	logger = logger.WithField(methodLogField, methodName)
	logger.Infof("%v ...", methodName)

	// Extract query parameters
	queryParams := r.URL.Query()
	idParam := queryParams.Get("id")

	// Prepare the list of account IDs
	var accountIDs []string
	if idParam != "" {
		accountIDs = append(accountIDs, idParam)
	}

	filter := storage.GetAccountParams{}
	if len(accountIDs) > 0 {
		filter.AccountIDs = accountIDs
	}

	// Call storage method to get accounts
	accounts, err := h.accountStore.GetAccounts(ctx, filter)
	if err != nil {
		logger.Errorf("error when getting accounts : %v", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := &GetAccountResponse{
		Total: len(accounts),
	}

	for _, a := range accounts {
		resp.Accounts = append(resp.Accounts, convertStorageToAccount(&a))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// DeleteAccount deletes an existing account.
func (h *Handler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	const methodName = "DeleteAccount"
	ctx = util.SetCallerMethodToCtx(ctx, methodName)
	l := logrus.New()
	logger := l.WithContext(ctx)
	logger = logger.WithField(methodLogField, methodName)
	logger.Infof("%v ...", methodName)

	var reqBody DeleteAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		logger.Errorf("error when decoding request %v", err)
		util.HandleError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if reqBody.ID == "" {
		logger.Info("no uuid provided")
		util.HandleError(w, http.StatusBadRequest, "no uuid provided")
		return
	}

	// Call storage method to delete account
	accountIDs := []string{reqBody.ID}

	filter := storage.GetAccountParams{}
	if len(accountIDs) > 0 {
		filter.AccountIDs = accountIDs
	}

	accounts, err := h.accountStore.GetAccounts(ctx, filter)
	if err != nil {
		logger.Errorf("failed to get account :  %v", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(accounts) == 0 {
		logger.Infof("wrong uuid provided")
		util.HandleError(w, http.StatusBadRequest, "wrong uuid provided")
		return
	}

	if err = h.accountStore.DeleteAccount(ctx, reqBody.ID); err != nil {
		logger.Errorf("error when deleting account : %v", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nil)
}

// UpdateAccount updates an existing account.
func (h *Handler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	const methodName = "UpdateAccount"
	ctx = util.SetCallerMethodToCtx(ctx, methodName)
	l := logrus.New()
	logger := l.WithContext(ctx)
	logger = logger.WithField(methodLogField, methodName)
	logger.Infof("%v ...", methodName)

	var reqBody UpdateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		logger.Errorf("error when decoding request %v", err)
		util.HandleError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if reqBody.ID == "" {
		logger.Info("no uuid provided")
		util.HandleError(w, http.StatusBadRequest, "no uuid provided")
		return
	}

	account := &storage.Account{
		ID:       reqBody.ID,
		Name:     reqBody.Name,
		Email:    reqBody.Email,
		Password: reqBody.Password,
	}

	filter := storage.GetAccountParams{}

	filter.AccountIDs = append(filter.AccountIDs, account.ID)

	// Call storage method to update account
	accounts, err := h.accountStore.GetAccounts(ctx, filter)
	if err != nil {
		logger.Errorf("failed to get account :  %v", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(accounts) == 0 {
		logger.Infof("wrong uuid provided")
		util.HandleError(w, http.StatusBadRequest, "wrong uuid provided")
		return
	}

	if err = h.accountStore.UpdateAccount(ctx, account); err != nil {
		logger.Errorf("error when updating account : %v", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := &UpdateAccountResponse{
		Account: convertStorageToAccount(account),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// convertStorageToAccount converts storage Account to Account for response.
func convertStorageToAccount(a *storage.Account) *Account {
	return &Account{
		ID:       a.ID,
		Name:     a.Name,
		Email:    a.Email,
		Password: a.Password,
	}
}

// tihor

// SignInRequest represents the request body for sign-in.
type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignInResponse represents the response body for sign-in.
type SignInResponse struct {
	Token string `json:"token"`
}

// SignOut handles the sign-out process.
func (h *Handler) SignOut(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	const methodName = "SignOut"
	ctx = util.SetCallerMethodToCtx(ctx, methodName)
	l := logrus.New()
	logger := l.WithContext(ctx)
	logger = logger.WithField(methodLogField, methodName)
	logger.Infof("%v ...", methodName)

	var reqBody SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		logger.Errorf("error when decoding request %v", err)
		util.HandleError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate the request body
	if reqBody.Email == "" || reqBody.Password == "" {
		logger.Error("missing email or password")
		util.HandleError(w, http.StatusBadRequest, "missing email or password")
		return
	}

	filter := storage.GetAccountParams{}

	filter.Email = reqBody.Email

	// Retrieve account from storage
	accounts, err := h.accountStore.GetAccounts(ctx, filter)
	if err != nil {
		logger.Errorf("error when retrieving account: %v", err)
		util.HandleError(w, http.StatusInternalServerError, "error retrieving account")
		return
	}

	if accounts == nil {
		logger.Info("accounts not found")
		util.HandleError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	acc := accounts[0]

	if acc.Password != reqBody.Password {
		logger.Infof("invalid password for account: %s", reqBody.Email)
		util.HandleError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	if err = h.accountStore.DeleteToken(ctx, acc.ID); err != nil {
		logger.Infof("failed to create token")
		util.HandleError(w, http.StatusInternalServerError, "failed to delete token")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(nil)
}

// SignIn handles the sign-out process.
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	const methodName = "SignIn"
	ctx = util.SetCallerMethodToCtx(ctx, methodName)
	l := logrus.New()
	logger := l.WithContext(ctx)
	logger = logger.WithField(methodLogField, methodName)
	logger.Infof("%v ...", methodName)

	var reqBody SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		logger.Errorf("error when decoding request %v", err)
		util.HandleError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate the request body
	if reqBody.Email == "" || reqBody.Password == "" {
		logger.Error("missing email or password")
		util.HandleError(w, http.StatusBadRequest, "missing email or password")
		return
	}

	filter := storage.GetAccountParams{}

	filter.Email = reqBody.Email

	// Retrieve account from storage
	accounts, err := h.accountStore.GetAccounts(ctx, filter)
	if err != nil {
		logger.Errorf("error when retrieving account: %v", err)
		util.HandleError(w, http.StatusInternalServerError, "error retrieving account")
		return
	}

	if accounts == nil {
		logger.Info("accounts not found")
		util.HandleError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	acc := accounts[0]

	// Verify password
	// if err := bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(reqBody.Password)); err != nil {
	// 	logger.Infof("invalid password for account: %s", reqBody.Email)
	// 	util.HandleError(w, http.StatusUnauthorized, "invalid email or password")
	// 	return
	// }

	if acc.Password != reqBody.Password {
		logger.Infof("invalid password for account: %s", reqBody.Email)
		util.HandleError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	// Check if token already exist and not expired yet
	token, err := h.accountStore.GetTokens(ctx, acc.ID)
	if err != nil {
		logger.Errorf("error getting token: %v", err)
		util.HandleError(w, http.StatusInternalServerError, "error getting token")
		return
	}

	returnToken := ""

	if token == nil || isTokenExpired(token.Expired) {
		// Generate token (placeholder function, implement token generation as needed)
		token, err := generateToken()
		if err != nil {
			logger.Errorf("error generating token: %v", err)
			util.HandleError(w, http.StatusInternalServerError, "error generating token")
			return
		}

		if err = h.accountStore.DeleteToken(ctx, acc.ID); err != nil {
			logger.Infof("failed to create token")
			util.HandleError(w, http.StatusInternalServerError, "failed to delete token")
			return
		}

		if err = h.accountStore.CreateToken(ctx, acc.ID, token); err != nil {
			logger.Infof("failed to create token")
			util.HandleError(w, http.StatusInternalServerError, "failed to insert token")
			return
		}

		returnToken = token
	} else {
		returnToken = token.Token
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := SignInResponse{
		Token: returnToken,
	}

	json.NewEncoder(w).Encode(resp)
}

func isTokenExpired(expired time.Time) bool {
	return expired.Before(time.Now())
}

// Placeholder function to generate a token. Replace with actual implementation.
func generateToken() (string, error) {
	// Generate a new UUID
	uuidToken, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Wrap(err, "failed to generate UUID token")
	}

	// Convert UUID to string
	token := uuidToken.String()

	// Hash the token using SHA-256
	hasher := sha256.New()
	_, err = hasher.Write([]byte(token))
	if err != nil {
		return "", errors.Wrap(err, "failed to hash token")
	}
	hashedToken := hex.EncodeToString(hasher.Sum(nil))

	return hashedToken, nil
}
