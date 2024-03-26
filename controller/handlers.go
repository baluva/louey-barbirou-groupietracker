package controller

import (
	"fmt"
	"log"
	"net/http"
	"projet/backend"
	"projet/templates"
	"strconv"
)

func AccueilPage(w http.ResponseWriter, r *http.Request) {
	// Retrieve the numCoins value from the form
	numCoins := 10

	// Now that numCoins is validated as an integer greater than 0 or set to a default, you can pass it to backend.GetCoin
	items, err := backend.GetCoin(numCoins)
	if err != nil {
		fmt.Println("Error fetching coin data:", err)
		return
	}

	// Use the items in your template
	if err := templates.Temp.ExecuteTemplate(w, "accueil", items); err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
	}
}

func FilterSubmit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filter := r.FormValue("numCoins")

	numCoins, err := strconv.Atoi(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	items, err := backend.GetCoin(numCoins)
	if err != nil {
		fmt.Println("Error fetching coin data:", err)
		return
	}

	if err := templates.Temp.ExecuteTemplate(w, "accueil", items); err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
	}
}
