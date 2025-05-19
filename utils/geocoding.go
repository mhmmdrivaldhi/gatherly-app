package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type Geocode struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

type LocationIQResponse []struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

type IPAPIResponse struct {
	Status      string  `json:"status"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Query       string  `json:"query"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
}

func GetCoordinateFromIP(ipAddress string) (*Geocode, error) {
	// Jika IP kosong atau localhost, gunakan endpoint tanpa parameter IP
	endpoint := "http://ip-api.com/json/"
	if ipAddress != "" && ipAddress != "127.0.0.1" && ipAddress != "::1" {
		endpoint = fmt.Sprintf("http://ip-api.com/json/%s", ipAddress)
	}

	response, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("gagal melakukan request ke ip-api: %w", err)
	}
	defer response.Body.Close()

	var ipResponse IPAPIResponse
	err = json.NewDecoder(response.Body).Decode(&ipResponse)
	if err != nil {
		return nil, fmt.Errorf("gagal parsing response dari ip-api: %w", err)
	}

	if ipResponse.Status != "success" {
		return nil, fmt.Errorf("tidak bisa mendapatkan lokasi dari IP")
	}

	return &Geocode{
		Latitude:  ipResponse.Lat,
		Longitude: ipResponse.Lon,
	}, nil
}

func GetCoordinatesFromAddress(address string) (*Geocode, error) {
	apiKey := os.Getenv("LOCATIONIQ_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API KEY LocationIQ tidak ditemukan")
	}

	endpoint := "https://us1.locationiq.com/v1/search.php"
	params := url.Values{}
	params.Add("key", apiKey)
	params.Add("q", address)
	params.Add("format", "json")

	requestURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	response, err := http.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var responseLocation LocationIQResponse
	err = json.NewDecoder(response.Body).Decode(&responseLocation)
	if err != nil {
		return nil, err
	}

	if len(responseLocation) == 0 {
		return nil, fmt.Errorf("lokasi tidak ditemukan")
	}

	var latitude, longitude float64

	fmt.Sscanf(responseLocation[0].Lat, "%f", &latitude)
	fmt.Sscanf(responseLocation[0].Lon, "%f", &longitude)

	return &Geocode{
		Latitude:  latitude,
		Longitude: longitude,
	}, nil
}
