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

type Club struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Country       string            `json:"country"`
	Province      string            `json:"province"`
	District      string            `json:"district"` // kabupaten
	EstablishDate time.Time         `json:"establish_date"`
	Logo          string            `json:"logo"`
	Address       string            `json:"address"`
	EmailPIC      string            `json:"email_pic"`
	Pic           string            `json:"pic"`
	Discipline    string            `json:"discipline"`
	PhoneNumber   string            `json:"phone_number"`
	Active        util.OptionalBool `json:"active"`
}

type CreateClubRequest struct {
	Name          string `json:"name"`
	Country       string `json:"country"`
	Province      string `json:"province"`
	District      string `json:"district"` // kabupaten
	EstablishDate string `json:"establish_date"`
	Logo          string `json:"logo"`
	Address       string `json:"address"`
	Pic           string `json:"pic"`
	EmailPIC      string `json:"email_pic"`
	Discipline    string `json:"discipline"`
	PhoneNumber   string `json:"phone_number"`
	Active        string `json:"active"`
}

type CreateClubResponse struct {
	ID string `json:"id"`
}

type GetClubRequest struct{}
type GetClubResponse struct {
	Clubs []*Club `json:"clubs"`
	Total int     `json:"total_count"`
}

type DeleteClubRequest struct {
	ID string `json:"id"`
}

type DeleteClubResponse struct{}

type UpdateClubRequest struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Country       string `json:"country"`
	Province      string `json:"province"`
	District      string `json:"district"` // kabupaten
	EstablishDate string `json:"establish_date"`
	Logo          string `json:"logo"`
	Address       string `json:"address"`
	Pic           string `json:"pic"`
	EmailPIC      string `json:"email_pic"`
	Discipline    string `json:"discipline"`
	PhoneNumber   string `json:"phone_number"`
	Active        string `json:"active"`
}

type UpdateClubResponse struct {
	Club *storage.Club `json:"club"`
}

var methodLogField = "method"

// CreateClub creates a new club.
func (h *Handler) CreateClub(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	const methodName = "CreateClub"
	ctx = util.SetCallerMethodToCtx(ctx, methodName)
	l := logrus.New()
	logger := l.WithContext(ctx)
	logger = logger.WithField(methodLogField, methodName)
	logger.Infof("%v ...", methodName)

	var reqBody CreateClubRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		logger.Errorf("error when decoding request %v", err)
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate the request body
	if err := validateClubRequest(reqBody); err != nil {
		logger.Errorf("error when validating request %v", err)
		util.HandleError(w, http.StatusBadRequest, err.Message)
		return
	}

	establishedDate, err := time.Parse("2006-01-02", reqBody.EstablishDate)
	if err != nil {
		logger.Errorf("error when converting time %v", err)
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Call storage method to create club
	club := &storage.Club{
		ID:            uuid.NewString(),
		Name:          reqBody.Name,
		Country:       reqBody.Country,
		Province:      reqBody.Province,
		District:      reqBody.District,
		EstablishDate: establishedDate,
		Logo:          reqBody.Logo,
		Address:       reqBody.Address,
		Pic:           reqBody.Pic,
		EmailPIC:      reqBody.EmailPIC,
		Discipline:    reqBody.Discipline,
		PhoneNumber:   reqBody.PhoneNumber,
		Active:        requestGetVal(reqBody.Active),
	}

	id, err := h.clubStore.CreateClub(ctx, club)
	if err != nil {
		logger.Errorf("error when creating club %v", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := CreateClubResponse{
		ID: id,
	}

	json.NewEncoder(w).Encode(resp)
}

// validateClubRequest validates the fields of ClubRequest.
func validateClubRequest(req CreateClubRequest) *util.Error {
	var missingFields []string
	var errFields []string

	if req.Name == "" {
		missingFields = append(missingFields, "name")
	}

	if req.Country == "" {
		missingFields = append(missingFields, "country")
	}

	if req.Province == "" {
		missingFields = append(missingFields, "province")
	}

	if req.District == "" {
		missingFields = append(missingFields, "district")
	}

	if req.EmailPIC == "" {
		missingFields = append(missingFields, "email_pic")
	}

	// Validate EstablishedDate
	if req.EstablishDate == "" {
		missingFields = append(missingFields, "establish_date")
	} else {
		_, err := time.Parse("2006-01-02", req.EstablishDate)
		if err != nil {
			errFields = append(errFields, "establish_date (format: yyyy-mm-dd)")
		}
	}
	if req.Logo == "" {
		missingFields = append(missingFields, "logo")
	}
	if req.Address == "" {
		missingFields = append(missingFields, "address")
	}
	if req.Pic == "" {
		missingFields = append(missingFields, "pic")
	}
	if req.Discipline == "" {
		missingFields = append(missingFields, "discipline")
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

// GetClub creates a new club.
func (h *Handler) GetClub(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	const methodName = "GetClub"
	ctx = util.SetCallerMethodToCtx(ctx, methodName)
	l := logrus.New()
	logger := l.WithContext(ctx)
	logger = logger.WithField(methodLogField, methodName)
	logger.Infof("%v ...", methodName)

	// Call storage method to create club
	clubIDs := []string{}
	clubs, err := h.clubStore.GetClubs(ctx, clubIDs...)
	if err != nil {
		logger.Errorf("error when getting clubs : %v", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := &GetClubResponse{
		Total: len(clubs),
	}

	for _, c := range clubs {
		resp.Clubs = append(resp.Clubs, convertStorageToClub(c))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// DeleteClub creates a new club.
func (h *Handler) DeleteClub(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	const methodName = "DeleteClub"
	ctx = util.SetCallerMethodToCtx(ctx, methodName)
	l := logrus.New()
	logger := l.WithContext(ctx)
	logger = logger.WithField(methodLogField, methodName)
	logger.Infof("%v ...", methodName)

	var reqBody DeleteClubRequest
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

	// Call storage method to delete club
	clubIDs := []string{reqBody.ID}

	clubs, err := h.clubStore.GetClubs(ctx, clubIDs...)
	if err != nil {
		logger.Errorf("failed to get club :  %v", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(clubs) == 0 {
		logger.Infof("wrong uuid provided")
		util.HandleError(w, http.StatusBadRequest, "wrong uuid provided")
		return
	}

	if err = h.clubStore.DeleteClub(ctx, reqBody.ID); err != nil {
		logger.Errorf("error when deleting club : %v", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nil)
}

// UpdateClub creates a new club.
func (h *Handler) UpdateClub(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	const methodName = "UpdateClub"
	ctx = util.SetCallerMethodToCtx(ctx, methodName)
	l := logrus.New()
	logger := l.WithContext(ctx)
	logger = logger.WithField(methodLogField, methodName)
	logger.Infof("%v ...", methodName)

	var reqBody UpdateClubRequest
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

	// Call storage method to create club
	clubIDs := []string{reqBody.ID}

	clubs, err := h.clubStore.GetClubs(ctx, clubIDs...)
	if err != nil {
		logger.Errorf("failed to get club", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(clubs) == 0 {
		logger.Info("wrong uuid provided")
		util.HandleError(w, http.StatusBadRequest, "wrong uuid provided")
		return
	}

	club := clubs[0]

	// validating
	if reqBody.Active != "" {
		if reqBody.Active != string(util.OptionalBool_TRUE) && reqBody.Active != string(util.OptionalBool_FALSE) {
			logger.Info("wrong active value given")
			util.HandleError(w, http.StatusBadRequest, fmt.Sprintf("wrong value for active , should be %v or %v", util.OptionalBool_FALSE, util.OptionalBool_TRUE))
			return
		}

		club.Active = requestGetVal(reqBody.Active)
	}

	if reqBody.Name != "" {
		club.Name = reqBody.Name
	}

	if reqBody.Country != "" {
		club.Country = reqBody.Country
	}

	if reqBody.Province != "" {
		club.Province = reqBody.Province
	}

	if reqBody.District != "" {
		club.District = reqBody.District
	}

	if reqBody.EstablishDate != "" {
		establishedDate, err := time.Parse("2006-01-02", reqBody.EstablishDate)
		if err != nil {
			logger.Errorf("error when converting time %v", err)
			util.HandleError(w, http.StatusBadRequest, err.Error())
			return
		}

		club.EstablishDate = establishedDate
	}

	if reqBody.Logo != "" {
		club.Logo = reqBody.Logo
	}

	if reqBody.Address != "" {
		club.Address = reqBody.Address
	}

	if reqBody.Pic != "" {
		club.Pic = reqBody.Pic
	}

	if reqBody.EmailPIC != "" {
		club.EmailPIC = reqBody.EmailPIC
	}

	if reqBody.Discipline != "" {
		club.Discipline = reqBody.Discipline
	}

	if reqBody.PhoneNumber != "" {
		club.PhoneNumber = reqBody.PhoneNumber
	}

	if err = h.clubStore.UpdateClub(ctx, &club); err != nil {
		logger.Errorf("error when updating club : %v", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res, err := h.clubStore.GetClubs(ctx, clubIDs...)
	if err != nil {
		logger.Errorf("failed to get club :  %v", err)
		util.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := UpdateClubResponse{
		Club: &res[0],
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func convertStorageToClub(str storage.Club) *Club {
	return &Club{
		ID:            str.ID,
		Name:          str.Name,
		Country:       str.Country,
		Province:      str.Province,
		District:      str.District,
		EstablishDate: str.EstablishDate,
		Logo:          str.Logo,
		Address:       str.Address,
		Pic:           str.Pic,
		EmailPIC:      str.EmailPIC,
		Discipline:    str.Discipline,
		PhoneNumber:   str.PhoneNumber,
		Active:        str.Active,
	}
}

func requestGetVal(val string) util.OptionalBool {
	if val == "TRUE" {
		return util.OptionalBool_TRUE
	} else if val == "FALSE" {
		return util.OptionalBool_FALSE
	}

	return util.OptionalBool_UNKNOWN_OptionalBool
}
