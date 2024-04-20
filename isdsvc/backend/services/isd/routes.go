package isd

import (
	"net/http"
)

// RegisterRoutes registers all routes for the application.
func RegisterRoutes(handler *Handler) {
	// Define route for POST method
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
}
