package isd

import (
	"net/http"
)

// RegisterRoutes registers all routes for the application.
func RegisterRoutes(handler *Handler) {

	http.HandleFunc("/v1/club", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateClub(w, r)
		} else if r.Method == http.MethodGet {
			handler.GetClub(w, r)
		} else if r.Method == http.MethodDelete {
			handler.DeleteClub(w, r)
		} else if r.Method == http.MethodPut {
			handler.UpdateClub(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/v1/organizer", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateOrganizer(w, r)
		} else if r.Method == http.MethodGet {
			handler.GetOrganizer(w, r)
		} else if r.Method == http.MethodDelete {
			handler.DeleteOrganizer(w, r)
		} else if r.Method == http.MethodPut {
			handler.UpdateOrganizer(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/v1/athlete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateAthlete(w, r)
		} else if r.Method == http.MethodGet {
			handler.GetAthlete(w, r)
		} else if r.Method == http.MethodDelete {
			handler.DeleteAthlete(w, r)
		} else if r.Method == http.MethodPut {
			handler.UpdateAthlete(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
