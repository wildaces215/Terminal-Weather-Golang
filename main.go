package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

//types

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

type Weather struct {
	Context  []interface{} `json:"@context"`
	ID       string        `json:"id"`
	Type     string        `json:"type"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		ID                  string `json:"@id"`
		Type                string `json:"@type"`
		Cwa                 string `json:"cwa"`
		ForecastOffice      string `json:"forecastOffice"`
		GridID              string `json:"gridId"`
		GridX               int    `json:"gridX"`
		GridY               int    `json:"gridY"`
		Forecast            string `json:"forecast"`
		ForecastHourly      string `json:"forecastHourly"`
		ForecastGridData    string `json:"forecastGridData"`
		ObservationStations string `json:"observationStations"`
		RelativeLocation    struct {
			Type     string `json:"type"`
			Geometry struct {
				Type        string    `json:"type"`
				Coordinates []float64 `json:"coordinates"`
			} `json:"geometry"`
			Properties struct {
				City     string `json:"city"`
				State    string `json:"state"`
				Distance struct {
					UnitCode string  `json:"unitCode"`
					Value    float64 `json:"value"`
				} `json:"distance"`
				Bearing struct {
					UnitCode string `json:"unitCode"`
					Value    int    `json:"value"`
				} `json:"bearing"`
			} `json:"properties"`
		} `json:"relativeLocation"`
		ForecastZone    string `json:"forecastZone"`
		County          string `json:"county"`
		FireWeatherZone string `json:"fireWeatherZone"`
		TimeZone        string `json:"timeZone"`
		RadarStation    string `json:"radarStation"`
	} `json:"properties"`
}

//This gets the public ip address

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

// Calls to the NOAA api to retrieve the  hourly weather forecast
func getWeatherHourlyURL(lat, long string) (string, error) {
	url := "https://api.weather.gov/points/" + lat + "," + long
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	var weather Weather

	body, err := io.ReadAll(resp.Body)
	err2 := json.Unmarshal(body, &weather)
	if err2 != nil {
		fmt.Println("Error:", err2)
	}
	return weather.Properties.ForecastHourly, err2
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
	hourlyURL, err3 := getWeatherHourlyURL(strconv.FormatFloat(lat, 'f', -1, 64), strconv.FormatFloat(long, 'f', -1, 64))
	if err3 != nil {
		fmt.Println(err3)
	}
	fmt.Println(hourlyURL)
}
