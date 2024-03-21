package controller

import (
	"log"
	"net/http"
	"projet/backend"
	"projet/templates"
)

func AccueilPage(w http.ResponseWriter, r *http.Request) {
	datas, err := backend.Getcoin()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error fetching coin data: %v", err)
		return
	}

	if err := templates.Temp.ExecuteTemplate(w, "accueil", datas); err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
	}
}
