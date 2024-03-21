package routeur

import (
	"log"
	"net/http"
	"projet/backend"
	"projet/controller"
)

func AccueilHandler(w http.ResponseWriter, r *http.Request) {
	// Example condition: differentiate based on HTTP method
	if r.Method == "GET" {
		controller.AccueilPage(w, r)
	} else if r.Method == "POST" {
		backend.SearchHandler(w, r)
	} else {
		// Handle unsupported methods
		http.Error(w, "Unsupported request method.", http.StatusMethodNotAllowed)
	}
}

func Initserv() {

	css := http.FileServer(http.Dir("./assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", css))

	// Register the new combined handler for /accueil
	http.HandleFunc("/accueil", AccueilHandler)
	http.HandleFunc("/search", backend.SearchHandler)

	addr := ":8000"
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
