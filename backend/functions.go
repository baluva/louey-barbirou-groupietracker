package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

const favoritesFilePath = "favorite.json"

func Getcoin() (ApiResponse, error) {
	var apiResponse ApiResponse
	url := "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return apiResponse, err
	}

	req.Header.Add("X-CMC_PRO_API_KEY", "569a71f5-374c-4e00-b44a-5a908d012347")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return apiResponse, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return apiResponse, err
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return apiResponse, err
	}

	// Convert Symbol to lower case for each item in apiResponse.Data
	for i, item := range apiResponse.Data {
		apiResponse.Data[i].Symbol = strings.ToLower(item.Symbol)
	}
	if len(apiResponse.Data) > 22 {
		// If so, slice the Data to keep only the top 20
		apiResponse.Data = apiResponse.Data[:50]
	}

	return apiResponse, nil
}
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("Symbol")
	if query == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	var response DetailedCryptoResponse
	url := fmt.Sprintf("https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest?symbol=%s", strings.ToUpper(query))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Failed to create request to CoinMarketCap", http.StatusInternalServerError)
		return
	}

	req.Header.Add("X-CMC_PRO_API_KEY", "569a71f5-374c-4e00-b44a-5a908d012347") // Replace with your actual API key
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to get coin data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &response); err != nil {
		http.Error(w, "Failed to unmarshal response", http.StatusInternalServerError)
		return
	}

	// Adjusting Symbol to lowercase for each entry in the response
	for key, data := range response.Data {
		data.Symbol = strings.ToLower(data.Symbol)
		response.Data[key] = data // Update the map entry
	}

	tmpl, err := template.ParseFiles("templates/search.html") // Ensure you handle the error from ParseFiles appropriately
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "search", response)
	if err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		return
	}
}
func UpdateFavorites(symbol string) {
	var favs Favorites
	file, err := os.ReadFile(favoritesFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Error reading favorites file: %v", err)
			return
		}
	} else {
		if err := json.Unmarshal(file, &favs); err != nil {
			log.Printf("Error unmarshaling favorites: %v", err)
			return
		}
	}

	// Check if the symbol is already a favorite
	index := -1
	for i, fav := range favs.Favorites {
		if fav == symbol {
			index = i
			break
		}
	}

	// If the symbol is not found, add it; otherwise, remove it
	if index == -1 {
		favs.Favorites = append(favs.Favorites, symbol)
	} else {
		favs.Favorites = append(favs.Favorites[:index], favs.Favorites[index+1:]...)
	}

	// Save the updated favorites
	updatedFavs, err := json.Marshal(favs)
	if err != nil {
		log.Printf("Error marshaling favorites: %v", err)
		return
	}

	if err := os.WriteFile(favoritesFilePath, updatedFavs, 0644); err != nil {
		log.Printf("Error writing favorites file: %v", err)
		return
	}
}
func FavoritePage(w http.ResponseWriter, r *http.Request) {
	var favs Favorites

	// Read the favorites
	file, err := os.ReadFile(favoritesFilePath)
	if err != nil {
		http.Error(w, "Failed to read favorites", http.StatusInternalServerError)
		log.Printf("Error reading favorites file: %v", err)
		return
	}

	if err := json.Unmarshal(file, &favs); err != nil {
		http.Error(w, "Failed to load favorites", http.StatusInternalServerError)
		log.Printf("Error unmarshaling favorites: %v", err)
		return
	}

	// Logic to render the page with the favorites
	// This could be using template.Execute to generate HTML
	// For simplicity, just listing the symbols:
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, "<html><body><h1>Favorite Cryptocurrencies</h1><ul>")
	for _, symbol := range favs.Favorites {
		fmt.Fprintf(w, "<li>%s</li>", symbol)
	}
	fmt.Fprintln(w, "</ul></body></html>")
}
