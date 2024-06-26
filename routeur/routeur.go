package routeur

import (
	"log"
	"net/http"
	"projet/backend"
	"projet/controller"
	"projet/templates"
)

func AccueilHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		controller.AccueilPage(w, r)
	} else if r.Method == "POST" {
		// Keep existing logic for handling POST requests.
		backend.SearchHandler(w, r)
	} else {
		// Handle unsupported HTTP methods.
		http.Error(w, "Unsupported request method.", http.StatusMethodNotAllowed)
	}
}
func errorHandler(w http.ResponseWriter, r *http.Request) {
	templates.Temp.ExecuteTemplate(w, "error", nil)
}

func Initserv() {

	css := http.FileServer(http.Dir("./assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", css))
	http.HandleFunc("/accueil", AccueilHandler)
	http.HandleFunc("/search", backend.SearchHandler)
	http.HandleFunc("/favorite", backend.HandleFavorite)
	http.HandleFunc("/filter_submit", controller.FilterSubmit)
	addr := ":8000"
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)

	}
}
