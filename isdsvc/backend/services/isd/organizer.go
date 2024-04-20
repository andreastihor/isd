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

type Organizer struct {
	ID           string            `json:"id" db:"id"`
	Name         string            `json:"name" db:"name"`
	Position     string            `json:"position" db:"position"`
	Club         *Club             `json:"club" db:"-"`
	RegisterDate time.Time         `json:"register_date" db:"register_date"`
	PhoneNumber  string            `json:"phone_number" db:"phone_number"`
	Active       util.OptionalBool `json:"active" db:"active"`
	Email        string            `json:"email" db:"email"`
}

type CreateOrganizerRequest struct {
	Name         string `json:"name"`
	Position     string `json:"position"`
	ClubID       string `json:"club_id"`
	RegisterDate string `json:"register_date"`
	PhoneNumber  string `json:"phone_number"`
	Active       string `json:"active"`
	Email        string `json:"email"`
}

type CreateOrganizerResponse struct {
	ID string `json:"id"`
}

type UpdateOrganizerRequest struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Position     string `json:"position"`
	RegisterDate string `json:"register_date"`
	PhoneNumber  string `json:"phone_number"`
	Active       string `json:"active"`
	Email        string `json:"email"`
}

type UpdateOrganizerResponse struct {
	Organizer Organizer `json:"organizer"`
}

type DeleteOrganizerRequest struct {
	ID string `json:"id"`
}

type DeleteOrganizerResponse struct{}

type GetOrganizerRequest struct{}

type GetOrganizerResponse struct {
	Organizer []*Organizer `json:"organizers"`
	Total     int          `json:"total_count"`
}

// CreateOrganizer creates a new organizer.
func (h *Handler) CreateOrganizer(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	const methodName = "CreateOrganizer"
	ctx = util.SetCallerMethodToCtx(ctx, methodName)
	l := logrus.New()
	logger := l.WithContext(ctx)
	logger = logger.WithField(methodLogField, methodName)
	logger.Infof("%v ...", methodName)

	var reqBody CreateOrganizerRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		logger.Errorf("error when decoding request %v", err)
		util.HandleError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate the request body
	if err := validateOrganizerRequest(reqBody); err != nil {
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

	// Call storage method to create organizer
	organizer := &storage.Organizer{
		ID:           uuid.NewString(),
		Name:         reqBody.Name,
		Position:     reqBody.Position,
		ClubID:       reqBody.ClubID,
		RegisterDate: registerDate,
		PhoneNumber:  reqBody.PhoneNumber,
		Active:       requestGetVal(reqBody.Active),
		Email:        reqBody.Email,
	}

	id, err := h.organizerStore.CreateOrganizer(ctx, organizer)
	if err != nil {
		logger.Errorf("error when creating organizer %v", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := CreateOrganizerResponse{
		ID: id,
	}

	json.NewEncoder(w).Encode(resp)
}

// validateOrganizerRequest validates the fields of CreateOrganizerRequest.
func validateOrganizerRequest(req CreateOrganizerRequest) *util.Error {
	var missingFields []string
	var errFields []string

	if req.Name == "" {
		missingFields = append(missingFields, "name")
	}

	if req.Position == "" {
		missingFields = append(missingFields, "position")
	}

	if req.ClubID == "" {
		missingFields = append(missingFields, "club_id")
	}

	if req.RegisterDate == "" {
		missingFields = append(missingFields, "register_date")
	} else {
		_, err := time.Parse("2006-01-02", req.RegisterDate)
		if err != nil {
			errFields = append(errFields, "register_date (format: yyyy-mm-dd)")
		}
	}

	if req.PhoneNumber == "" {
		missingFields = append(missingFields, "phone_number")
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

// GetOrganizer retrieves organizer records.
func (h *Handler) GetOrganizer(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	const methodName = "GetOrganizer"
	ctx = util.SetCallerMethodToCtx(ctx, methodName)
	l := logrus.New()
	logger := l.WithContext(ctx)
	logger = logger.WithField(methodLogField, methodName)
	logger.Infof("%v ...", methodName)

	// Call storage method to retrieve organizers
	organizers, err := h.organizerStore.GetOrganizers(ctx)
	if err != nil {
		logger.Errorf("error when getting organizers: %v", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := &GetOrganizerResponse{
		Total: len(organizers),
	}

	for _, o := range organizers {
		resp.Organizer = append(resp.Organizer, convertStorageToOrganizer(o))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// DeleteOrganizer deletes an organizer record.
func (h *Handler) DeleteOrganizer(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	const methodName = "DeleteOrganizer"
	ctx = util.SetCallerMethodToCtx(ctx, methodName)
	l := logrus.New()
	logger := l.WithContext(ctx)
	logger = logger.WithField(methodLogField, methodName)
	logger.Infof("%v ...", methodName)

	var reqBody DeleteOrganizerRequest
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

	// Call storage method to delete organizer
	if err := h.organizerStore.DeleteOrganizer(ctx, reqBody.ID); err != nil {
		logger.Errorf("error when deleting organizer: %v", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nil)
}

// UpdateOrganizer updates an organizer record.
func (h *Handler) UpdateOrganizer(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	const methodName = "UpdateOrganizer"
	ctx = util.SetCallerMethodToCtx(ctx, methodName)
	l := logrus.New()
	logger := l.WithContext(ctx)
	logger = logger.WithField(methodLogField, methodName)
	logger.Infof("%v ...", methodName)

	var reqBody UpdateOrganizerRequest
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

	// Call storage method to update organizer
	organizers, err := h.organizerStore.GetOrganizers(ctx, reqBody.ID)
	if err != nil {
		logger.Errorf("failed to get organizer: %v", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(organizers) == 0 {
		logger.Info("wrong uuid provided")
		util.HandleError(w, http.StatusBadRequest, "wrong uuid provided")
		return
	}

	organizer := organizers[0]

	// Update fields if provided in the request
	if reqBody.Name != "" {
		organizer.Name = reqBody.Name
	}

	if reqBody.Position != "" {
		organizer.Position = reqBody.Position
	}

	if reqBody.RegisterDate != "" {
		registerDate, err := time.Parse("2006-01-02", reqBody.RegisterDate)
		if err != nil {
			logger.Errorf("error when converting time: %v", err)
			util.HandleError(w, http.StatusBadRequest, err.Error())
			return
		}
		organizer.RegisterDate = registerDate
	}

	if reqBody.PhoneNumber != "" {
		organizer.PhoneNumber = reqBody.PhoneNumber
	}

	if reqBody.Active != "" {
		if reqBody.Active != string(util.OptionalBool_TRUE) && reqBody.Active != string(util.OptionalBool_FALSE) {
			logger.Info("wrong active value given")
			util.HandleError(w, http.StatusBadRequest, fmt.Sprintf("wrong value for active, should be %v or %v", util.OptionalBool_FALSE, util.OptionalBool_TRUE))
			return
		}

		organizer.Active = requestGetVal(reqBody.Active)
	}

	if reqBody.Email != "" {
		organizer.Email = reqBody.Email
	}

	if err = h.organizerStore.UpdateOrganizer(ctx, &organizer); err != nil {
		logger.Errorf("error when updating organizer: %v", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := UpdateOrganizerResponse{
		Organizer: *convertStorageToOrganizer(organizer),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func convertStorageToOrganizer(str storage.Organizer) *Organizer {
	organizer := &Organizer{
		ID:           str.ID,
		Name:         str.Name,
		Position:     str.Position,
		RegisterDate: str.RegisterDate,
		PhoneNumber:  str.PhoneNumber,
		Active:       str.Active,
		Email:        str.Email,
	}

	if str.ClubName != "" {
		organizer.Club = &Club{
			ID:            str.ClubID,
			Name:          str.ClubName,
			Country:       str.ClubCountry,
			Province:      str.ClubProvince,
			District:      str.ClubDistrict,
			EstablishDate: str.ClubEstablishDate,
			Logo:          str.ClubLogo,
			Address:       str.ClubAddress,
			Pic:           str.ClubPic,
			EmailPIC:      str.ClubEmailPIC,
			Discipline:    str.ClubDiscipline,
			PhoneNumber:   str.ClubPhoneNumber,
			Active:        str.ClubActive,
		}
	}

	return organizer
}
