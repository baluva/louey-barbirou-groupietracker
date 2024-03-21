package backend

type ApiResponse struct {
	Data []struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
		Quote  map[string]struct {
			Price            float64 `json:"price"`
			PercentChange24h float64 `json:"percent_change_24h"`
			MarketCap        float64 `json:"market_cap"`
		} `json:"quote"`
	} `json:"data"`
}
type QuoteDetail struct {
	Price            float64 `json:"price"`
	PercentChange24h float64 `json:"percent_change_24h"`
	MarketCap        float64 `json:"market_cap"`
	// Add other quote details as necessary
}

type CryptoData struct {
	ID     int64                  `json:"id"`
	Name   string                 `json:"name"`
	Symbol string                 `json:"symbol"`
	Quote  map[string]QuoteDetail `json:"quote"`
	// Include other fields as necessary...
}

type DetailedCryptoResponse struct {
	Data map[string]CryptoData `json:"data"`
}

type Favorites struct {
	Favorites []string `json:"favorites"`
}
