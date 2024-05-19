package isd

import "github.com/andreastihor/isd/isdsvc/backend/storage"

type Handler struct {
	clubStore      storage.ClubStore
	organizerStore storage.OrganizerStore
	coachStore     storage.CoachStore
	athleteStore   storage.AthleteStore
}

// NewHandler creates a new instance of Handler with the given ClubStore.
func NewHandler(clubStore storage.ClubStore, organizerStore storage.OrganizerStore, athleteStore storage.AthleteStore) *Handler {
	return &Handler{
		clubStore:      clubStore,
		organizerStore: organizerStore,
		athleteStore:   athleteStore,
	}
}
