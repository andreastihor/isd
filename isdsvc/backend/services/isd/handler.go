package isd

import "github.com/andreastihor/isd/isdsvc/backend/storage"

type Handler struct {
	clubStore storage.ClubStore
}

// NewHandler creates a new instance of Handler with the given ClubStore.
func NewHandler(clubStore storage.ClubStore) *Handler {
	return &Handler{clubStore: clubStore}
}
