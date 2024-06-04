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

// Coach represents coach information.
type Coach struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	DOB          string            `json:"dob"`
	PhoneNumber  string            `json:"phone_number"`
	Gender       util.Gender       `json:"gender"`
	Email        string            `json:"email"`
	Discipline   string            `json:"discipline"`
	RegisterDate time.Time         `json:"register_date"`
	Active       util.OptionalBool `json:"active"`
}
type CreateCoachRequest struct {
	Name         string      `json:"name"`
	DOB          string      `json:"dob"`
	PhoneNumber  string      `json:"phone_number"`
	Gender       util.Gender `json:"gender"`
	Email        string      `json:"email"`
	Discipline   string      `json:"discipline"`
	RegisterDate string      `json:"register_date"`
	Active       string      `json:"active"`
}

type CreateCoachResponse struct {
	ID string `json:"id"`
}

type UpdateCoachRequest struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	DOB          string      `json:"dob"`
	PhoneNumber  string      `json:"phone_number"`
	Gender       util.Gender `json:"gender"`
	Email        string      `json:"email"`
	Discipline   string      `json:"discipline"`
	RegisterDate string      `json:"register_date"`
	Active       string      `json:"active"`
}

type UpdateCoachResponse struct {
	Coach Coach `json:"Coach"`
}

type DeleteCoachRequest struct {
	ID string `json:"id"`
}

type DeleteCoachResponse struct{}

type GetCoachRequest struct{}

type GetCoachResponse struct {
	Coach []*Coach `json:"coaches"`
	Total int      `json:"total_count"`
}

// CreateCoach creates a new coach.
func (h *Handler) CreateCoach(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	const methodName = "CreateCoach"
	ctx = util.SetCallerMethodToCtx(ctx, methodName)
	l := logrus.New()
	logger := l.WithContext(ctx)
	logger = logger.WithField(methodLogField, methodName)
	logger.Infof("%v ...", methodName)

	var reqBody CreateCoachRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		logger.Errorf("error when decoding request %v", err)
		util.HandleError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate the request body
	if err := validateCoachRequest(reqBody); err != nil {
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

	// Call storage method to create coach
	coach := &storage.Coach{
		ID:           uuid.NewString(),
		Name:         reqBody.Name,
		DOB:          reqBody.DOB,
		PhoneNumber:  reqBody.PhoneNumber,
		Gender:       reqBody.Gender,
		Email:        reqBody.Email,
		Discipline:   reqBody.Discipline,
		RegisterDate: registerDate,
		Active:       requestGetVal(reqBody.Active),
	}

	id, err := h.coachStore.CreateCoach(ctx, coach)
	if err != nil {
		logger.Errorf("error when creating coach %v", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := CreateCoachResponse{
		ID: id,
	}

	json.NewEncoder(w).Encode(resp)
}

// validateCoachRequest validates the fields of CreateCoachRequest.
func validateCoachRequest(req CreateCoachRequest) *util.Error {
	var missingFields []string
	var errFields []string

	if req.Name == "" {
		missingFields = append(missingFields, "name")
	}

	if req.DOB == "" {
		missingFields = append(missingFields, "dob")
	}

	if req.PhoneNumber == "" {
		missingFields = append(missingFields, "phone_number")
	}

	if req.Email == "" {
		missingFields = append(missingFields, "email")
	}

	if req.Discipline == "" {
		missingFields = append(missingFields, "discipline")
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

// GetCoach retrieves coach records.
// todo WIP
// func (h *Handler) GetCoach(w http.ResponseWriter, r *http.Request) {
// 	ctx := context.Background()
// 	const methodName = "GetCoach"
// 	ctx = util.SetCallerMethodToCtx(ctx, methodName)
// 	l := logrus.New()
// 	logger := l.WithContext(ctx)
// 	logger = logger.WithField(methodLogField, methodName)
// 	logger.Infof("%v ...", methodName)

// 	// Call storage method to retrieve coaches
// 	coaches, err := h.coachStore.GetCoaches(ctx)
// 	if err != nil {
// 		logger.Errorf("error when getting coaches: %v", err)
// 		util.HandleError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	resp := &GetCoachResponse{
// 		Total: len(coaches),
// 	}

// 	for _, c := range coaches {
// 		resp.Coach = append(resp.Coach, c)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(resp)
// }

// DeleteCoach deletes a coach record.
func (h *Handler) DeleteCoach(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	const methodName = "DeleteCoach"
	ctx = util.SetCallerMethodToCtx(ctx, methodName)
	l := logrus.New()
	logger := l.WithContext(ctx)
	logger = logger.WithField(methodLogField, methodName)
	logger.Infof("%v ...", methodName)

	var reqBody DeleteCoachRequest
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

	// Call storage method to delete coach
	if err := h.coachStore.DeleteCoach(ctx, reqBody.ID); err != nil {
		logger.Errorf("error when deleting coach: %v", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nil)
}

// UpdateCoach updates a coach record.
// func (h *Handler) UpdateCoach(w http.ResponseWriter, r *http.Request) {
// 	ctx := context.Background()
// 	const methodName = "UpdateCoach"
// 	ctx = util.SetCallerMethodToCtx(ctx, methodName)
// 	l := logrus.New()
// 	logger := l.WithContext(ctx)
// 	logger = logger.WithField(methodLogField, methodName)
// 	logger.Infof("%v ...", methodName)

// 	var reqBody UpdateCoachRequest
// 	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
// 		logger.Errorf("error when decoding request %v", err)
// 		util.HandleError(w, http.StatusBadRequest, "invalid request body")
// 		return
// 	}

// 	if reqBody.ID == "" {
// 		logger.Info("no uuid provided")
// 		util.HandleError(w, http.StatusBadRequest, "no uuid provided")
// 		return
// 	}

// 	// Call storage method to update coach
// 	coaches, err := h.coachStore.GetCoaches(ctx, reqBody.ID)
// 	if err != nil {
// 		logger.Errorf("failed to get coach: %v", err)
// 		util.HandleError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	if len(coaches) == 0 {
// 		logger.Info("wrong uuid provided")
// 		util.HandleError(w, http.StatusBadRequest, "wrong uuid provided")
// 		return
// 	}

// 	coach := coaches[0]

// 	// Update fields if provided in the request
// 	if reqBody.Name != "" {
// 		coach.Name = reqBody.Name
// 	}

// 	if reqBody.DOB != "" {
// 		coach.DOB = reqBody.DOB
// 	}

// 	if reqBody.PhoneNumber != "" {
// 		coach.PhoneNumber = reqBody.PhoneNumber
// 	}

// 	if reqBody.Gender != "" {
// 		coach.Gender = reqBody.Gender
// 	}

// 	if reqBody.Email != "" {
// 		coach.Email = reqBody.Email
// 	}

// 	if reqBody.Discipline != "" {
// 		coach.Discipline = reqBody.Discipline
// 	}

// 	if reqBody.RegisterDate != "" {
// 		registerDate, err := time.Parse("2006-01-02", reqBody.RegisterDate)
// 		if err != nil {
// 			logger.Errorf("error when converting time: %v", err)
// 			util.HandleError(w, http.StatusBadRequest, err.Error())
// 			return
// 		}
// 		coach.RegisterDate = registerDate
// 	}

// 	if reqBody.Active != "" {
// 		if reqBody.Active != string(util.OptionalBool_TRUE) && reqBody.Active != string(util.OptionalBool_FALSE) {
// 			logger.Info("wrong active value given")
// 			util.HandleError(w, http.StatusBadRequest, fmt.Sprintf("wrong value for active, should be %v or %v", util.OptionalBool_FALSE, util.OptionalBool_TRUE))
// 			return
// 		}

// 		coach.Active = requestGetVal(reqBody.Active)
// 	}

// 	if err := h.coachStore.UpdateCoach(ctx, &coach); err != nil {
// 		logger.Errorf("error when updating coach: %v", err)
// 		util.HandleError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	response := UpdateCoachResponse{
// 		Coach: coach,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(response)
// }
