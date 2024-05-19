package isd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/andreastihor/isd/isdsvc/backend/storage"
	"github.com/andreastihor/isd/isdsvc/backend/util"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Athlete represents athlete information.
type Athlete struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	DOB          time.Time         `json:"dob"`
	PhoneNumber  string            `json:"phone_number"`
	Gender       util.Gender       `json:"gender"`
	Email        string            `json:"email"`
	RegisterDate time.Time         `json:"register_date"`
	Active       util.OptionalBool `json:"active"`
}

type CreateAthleteRequest struct {
	Name         string      `json:"name"`
	DOB          string      `json:"dob"`
	PhoneNumber  string      `json:"phone_number"`
	Gender       util.Gender `json:"gender"`
	Email        string      `json:"email"`
	RegisterDate string      `json:"register_date"`
	Active       string      `json:"active"`
}

type CreateAthleteResponse struct {
	ID string `json:"id"`
}

type UpdateAthleteRequest struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	DOB          string      `json:"dob"`
	PhoneNumber  string      `json:"phone_number"`
	Gender       util.Gender `json:"gender"`
	Email        string      `json:"email"`
	RegisterDate string      `json:"register_date"`
	Active       string      `json:"active"`
}

type UpdateAthleteResponse struct {
	Athlete *Athlete `json:"athlete"`
}

type DeleteAthleteRequest struct {
	ID string `json:"id"`
}

type GetAthleteResponse struct {
	Total    int        `json:"total"`
	Athletes []*Athlete `json:"athletes"`
}

// CreateAthlete creates a new athlete.
func (h *Handler) CreateAthlete(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	const methodName = "CreateAthlete"
	ctx = util.SetCallerMethodToCtx(ctx, methodName)
	l := logrus.New()
	logger := l.WithContext(ctx)
	logger = logger.WithField(methodLogField, methodName)
	logger.Infof("%v ...", methodName)

	var reqBody CreateAthleteRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		logger.Errorf("error when decoding request %v", err)
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate the request body
	if err := validateAthleteRequest(reqBody); err != nil {
		logger.Errorf("error when validating request %v", err)
		util.HandleError(w, http.StatusBadRequest, err.Message)
		return
	}

	registerDate, err := time.Parse("2006-01-02", reqBody.RegisterDate)
	if err != nil {
		logger.Errorf("error when converting time %v", err)
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	dob, err := time.Parse("2006-01-02", reqBody.DOB)
	if err != nil {
		logger.Errorf("error when converting time %v", err)
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Call storage method to create athlete
	athlete := &storage.Athlete{
		ID:           uuid.NewString(),
		Name:         reqBody.Name,
		DOB:          dob,
		PhoneNumber:  reqBody.PhoneNumber,
		Gender:       reqBody.Gender,
		Email:        reqBody.Email,
		RegisterDate: registerDate,
		Active:       requestGetVal(reqBody.Active),
	}

	id, err := h.athleteStore.CreateAthlete(ctx, athlete)
	if err != nil {
		logger.Errorf("error when creating athlete %v", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := CreateAthleteResponse{
		ID: id,
	}

	json.NewEncoder(w).Encode(resp)
}

// validateAthleteRequest validates the fields of AthleteRequest.
func validateAthleteRequest(req CreateAthleteRequest) *util.Error {
	var missingFields []string
	var errFields []string

	if req.Name == "" {
		missingFields = append(missingFields, "name")
	}

	if req.DOB == "" {
		missingFields = append(missingFields, "register_date")
	} else {
		_, err := time.Parse("2006-01-02", req.DOB)
		if err != nil {
			errFields = append(errFields, "register_date (format: yyyy-mm-dd)")
		}
	}

	if req.PhoneNumber == "" {
		missingFields = append(missingFields, "phone_number")
	}

	if req.Email == "" {
		missingFields = append(missingFields, "email")
	}

	if req.RegisterDate == "" {
		missingFields = append(missingFields, "register_date")
	} else {
		_, err := time.Parse("2006-01-02", req.RegisterDate)
		if err != nil {
			errFields = append(errFields, "register_date (format: yyyy-mm-dd)")
		}
	}

	errorMsg := []string{}
	if len(missingFields) > 0 {
		errorMsg = append(errorMsg, fmt.Sprintf("Missing required fields: [%v]", strings.Join(missingFields, ",")))
	}

	if len(errFields) > 0 {
		errorMsg = append(errorMsg, fmt.Sprintf("Wrong Format fields: [%v]", strings.Join(errFields, ",")))
	}

	if len(errorMsg) > 0 {
		if len(errorMsg) == 1 {
			return util.NewError(http.StatusBadRequest, errorMsg[0])
		} else {
			return util.NewError(http.StatusBadRequest, strings.Join(errorMsg, " + "))
		}
	}

	return nil
}

// GetAthlete retrieves all athletes.
func (h *Handler) GetAthlete(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	const methodName = "GetAthlete"
	ctx = util.SetCallerMethodToCtx(ctx, methodName)
	l := logrus.New()
	logger := l.WithContext(ctx)
	logger = logger.WithField(methodLogField, methodName)
	logger.Infof("%v ...", methodName)

	// Call storage method to get athletes
	athleteIDs := []string{}
	athletes, err := h.athleteStore.GetAthletes(ctx, athleteIDs...)
	if err != nil {
		logger.Errorf("error when getting athletes : %v", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := &GetAthleteResponse{
		Total: len(athletes),
	}

	for _, a := range athletes {
		resp.Athletes = append(resp.Athletes, convertStorageToAthlete(&a))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// DeleteAthlete deletes an existing athlete.
func (h *Handler) DeleteAthlete(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	const methodName = "DeleteAthlete"
	ctx = util.SetCallerMethodToCtx(ctx, methodName)
	l := logrus.New()
	logger := l.WithContext(ctx)
	logger = logger.WithField(methodLogField, methodName)
	logger.Infof("%v ...", methodName)

	var reqBody DeleteAthleteRequest
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

	// Call storage method to delete athlete
	athleteIDs := []string{reqBody.ID}

	athletes, err := h.athleteStore.GetAthletes(ctx, athleteIDs...)
	if err != nil {
		logger.Errorf("failed to get athlete :  %v", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(athletes) == 0 {
		logger.Infof("wrong uuid provided")
		util.HandleError(w, http.StatusBadRequest, "wrong uuid provided")
		return
	}

	if err = h.athleteStore.DeleteAthlete(ctx, reqBody.ID); err != nil {
		logger.Errorf("error when deleting athlete : %v", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nil)
}

// UpdateAthlete updates an existing athlete.
func (h *Handler) UpdateAthlete(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	const methodName = "UpdateAthlete"
	ctx = util.SetCallerMethodToCtx(ctx, methodName)
	l := logrus.New()
	logger := l.WithContext(ctx)
	logger = logger.WithField(methodLogField, methodName)
	logger.Infof("%v ...", methodName)

	var reqBody UpdateAthleteRequest
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

	registerDate, err := time.Parse("2006-01-02", reqBody.RegisterDate)
	if err != nil {
		logger.Errorf("error when converting time %v", err)
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	dob, err := time.Parse("2006-01-02", reqBody.DOB)
	if err != nil {
		logger.Errorf("error when converting time %v", err)
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	athlete := &storage.Athlete{
		ID:           reqBody.ID,
		Name:         reqBody.Name,
		DOB:          dob,
		PhoneNumber:  reqBody.PhoneNumber,
		Gender:       reqBody.Gender,
		Email:        reqBody.Email,
		RegisterDate: registerDate,
		Active:       requestGetVal(reqBody.Active),
	}

	// Call storage method to update athlete
	athletes, err := h.athleteStore.GetAthletes(ctx, athlete.ID)
	if err != nil {
		logger.Errorf("failed to get athlete :  %v", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(athletes) == 0 {
		logger.Infof("wrong uuid provided")
		util.HandleError(w, http.StatusBadRequest, "wrong uuid provided")
		return
	}

	if err = h.athleteStore.UpdateAthlete(ctx, athlete); err != nil {
		logger.Errorf("error when updating athlete : %v", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := &UpdateAthleteResponse{
		Athlete: convertStorageToAthlete(athlete),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// convertStorageToAthlete converts storage Athlete to Athlete for response.
func convertStorageToAthlete(a *storage.Athlete) *Athlete {
	return &Athlete{
		ID:           a.ID,
		Name:         a.Name,
		DOB:          a.DOB,
		PhoneNumber:  a.PhoneNumber,
		Gender:       a.Gender,
		Email:        a.Email,
		RegisterDate: a.RegisterDate,
		Active:       a.Active,
	}
}
