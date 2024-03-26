package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"
)

const favoritesFilePath = "favorite.json"

func GetCoin(numCoins int) (ApiResponse, error) {
	var coin ApiResponse
	url := "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return coin, err
	}

	req.Header.Add("X-CMC_PRO_API_KEY", "569a71f5-374c-4e00-b44a-5a908d012347")
	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	q.Add("limit", fmt.Sprintf("%d", numCoins)) // Use numCoins as the limit for the number of results
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return coin, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return coin, err
	}

	if err := json.Unmarshal(body, &coin); err != nil {
		return coin, err
	}

	// Convert Symbol to lower case for each item in coin.Data
	for i, item := range coin.Data {
		coin.Data[i].Symbol = strings.ToLower(item.Symbol)
	}

	return coin, nil
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
func HandleFavorite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Extract the symbol from the form data
	symbol := r.FormValue("symbol")

	// Initialize a variable to hold the current list of favorites
	var favorites Favorite

	// Try to open the existing JSON file
	file, err := os.Open("favorite.json")
	if err == nil {
		// If the file exists, decode its contents into the favorites variable
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&favorites); err != nil {
			http.Error(w, "Error reading favorite", http.StatusInternalServerError)
			return
		}
		file.Close() // Close the file after reading
	} else if !os.IsNotExist(err) {
		// If an error other than "file not found" occurred, report it
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}

	// Check if the symbol already exists in the list of favorite coins
	for _, fav := range favorites.FavoriteCoins {
		if fav == symbol {
			// Symbol already exists, no need to add it again
			http.Redirect(w, r, "/accueil", http.StatusFound)
			return
		}
	}

	// Append the new symbol to the list of favorite coins
	favorites.FavoriteCoins = append(favorites.FavoriteCoins, symbol)

	// Open the file again, this time for writing (create it if it doesn't exist)
	file, err = os.OpenFile("favorite.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		http.Error(w, "Error opening file for writing", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Encode the updated favorites list to JSON and save it to the file
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(favorites); err != nil {
		http.Error(w, "Error saving favorite", http.StatusInternalServerError)
		return
	}

	// Redirect or inform the user of success
	http.Redirect(w, r, "/accueil", http.StatusFound)
}

func readFavorites() ([]string, error) {
	// Initialize a variable to hold the current list of favorites
	var favorites Favorite

	// Open the JSON file
	file, err := os.Open("favorite.json")
	if err != nil {
		// If the file does not exist, return an empty slice and no error
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		// For any other error, return it
		return nil, err
	}
	defer file.Close()

	// Decode the file contents into the favorites struct
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&favorites); err != nil {
		return nil, err
	}

	// Return the slice of favorite coins
	return favorites.FavoriteCoins, nil
}
func AfficheFavorite(w http.ResponseWriter, r *http.Request) {
	//favoriteCoins, err := readFavorites()

}
