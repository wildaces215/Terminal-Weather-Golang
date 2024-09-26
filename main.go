package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type PointsResponse struct {
	Properties struct {
		ForecastHourly string `json:"forecastHourly"`
	} `json:"properties"`
}

type HourlyForecastResponse struct {
	Properties struct {
		Periods []struct {
			StartTime       time.Time `json:"startTime"`
			Temperature     int       `json:"temperature"`
			TemperatureUnit string    `json:"temperatureUnit"`
			ShortForecast   string    `json:"shortForecast"`
		} `json:"periods"`
	} `json:"properties"`
}

type GeoCoordinates struct {
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	City        string  `json:"city"`
	Postal      string  `json:"postal"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	IPv4        string  `json:"IPv4"`
	State       string  `json:"state"`
}

func getIpAddress() (string, error) {
	url := "https://api.ipify.org"
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// Trim any whitespace and return the IP
	return strings.TrimSpace(string(ip)), nil
}

func getGeoCoordinates(ipAddr string) (float64, float64, error) {
	url1 := "https://geolocation-db.com/json/"
	url2 := "&position=true"
	combinedUrl := url1 + ipAddr + url2
	resp, err := http.Get(combinedUrl)
	if err != nil {
		fmt.Println(err)
	}
	body, err := io.ReadAll(resp.Body)

	//Process the response in the api call
	var geo GeoCoordinates
	errr := json.Unmarshal(body, &geo)

	if errr != nil {
		fmt.Println("Error:", errr)
	}

	return geo.Latitude, geo.Longitude, err
}

func getHourlyWeather(latitude, longitude float64) (string, error) {
	// Step 1: Get forecast office and grid coordinates
	pointsURL := fmt.Sprintf("https://api.weather.gov/points/%.4f,%.4f", latitude, longitude)
	pointsResp, err := http.Get(pointsURL)
	if err != nil {
		return "", fmt.Errorf("error fetching points data: %v", err)
	}
	defer pointsResp.Body.Close()

	if pointsResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: unable to fetch data. Status code: %d", pointsResp.StatusCode)
	}

	var pointsData PointsResponse
	if err := json.NewDecoder(pointsResp.Body).Decode(&pointsData); err != nil {
		return "", fmt.Errorf("error decoding points response: %v", err)
	}

	// Step 2: Get the hourly forecast
	forecastResp, err := http.Get(pointsData.Properties.ForecastHourly)
	if err != nil {
		return "", fmt.Errorf("error fetching hourly forecast data: %v", err)
	}
	defer forecastResp.Body.Close()

	if forecastResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: unable to fetch hourly forecast. Status code: %d", forecastResp.StatusCode)
	}

	var forecastData HourlyForecastResponse
	if err := json.NewDecoder(forecastResp.Body).Decode(&forecastData); err != nil {
		return "", fmt.Errorf("error decoding hourly forecast response: %v", err)
	}

	if len(forecastData.Properties.Periods) == 0 {
		return "", fmt.Errorf("no hourly forecast data available")
	}

	var result string
	for i, period := range forecastData.Properties.Periods {
		if i >= 12 { // Limit to next 12 hours
			break
		}
		result += fmt.Sprintf("Time: %s\nTemperature: %d°%s\nForecast: %s\n\n",
			period.StartTime.Format("2006-01-02 15:04 MST"),
			period.Temperature,
			period.TemperatureUnit,
			period.ShortForecast)
	}

	return result, nil
}

func main() {
	addr, err := getIpAddress()
	if err != nil {
		fmt.Println(err)
	}
	lat, long, err2 := getGeoCoordinates(addr)
	if err2 != nil {
		fmt.Println(err2)
	}

	weather, err := getHourlyWeather(lat,long)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Hourly Weather Forecast:")
	fmt.Println(weather)
}
