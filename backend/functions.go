package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"
)

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
