package main

import (
	"encoding/json"
	"fmt"
	ai "goweathergo/AI"
	config "goweathergo/Config"
	"net/http"
	"os"
	"strings"
)

type WeatherResponse struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int64   `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int64   `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree int     `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		PrecipMm   float64 `json:"precip_mm"`
		PrecipIn   float64 `json:"precip_in"`
		Humidity   int     `json:"humidity"`
		Cloud      int     `json:"cloud"`
		FeelslikeC float64 `json:"feelslike_c"`
		FeelslikeF float64 `json:"feelslike_f"`
		WindchillC float64 `json:"windchill_c"`
		WindchillF float64 `json:"windchill_f"`
		HeatindexC float64 `json:"heatindex_c"`
		HeatindexF float64 `json:"heatindex_f"`
		DewpointC  float64 `json:"dewpoint_c"`
		DewpointF  float64 `json:"dewpoint_f"`
		VisKm      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		UV         float64 `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
		ShortRad   float64 `json:"short_rad"`
		DiffRad    float64 `json:"diff_rad"`
		DNI        float64 `json:"dni"`
		GTI        float64 `json:"gti"`
	} `json:"current"`
}

func WeatherInString(wr WeatherResponse) string {
	return fmt.Sprintf(
		"Location:\n"+
			"  Name: %s\n  Region: %s\n  Country: %s\n  Lat: %.2f\n  Lon: %.2f\n  TzID: %s\n  LocaltimeEpoch: %d\n  Localtime: %s\n"+
			"Current:\n"+
			"  LastUpdatedEpoch: %d\n  LastUpdated: %s\n  TempC: %.2f\n  TempF: %.2f\n  IsDay: %d\n"+
			"  Condition:\n    Text: %s\n    Icon: %s\n    Code: %d\n"+
			"  WindMph: %.2f\n  WindKph: %.2f\n  WindDegree: %d\n  WindDir: %s\n"+
			"  PressureMb: %.2f\n  PressureIn: %.2f\n  PrecipMm: %.2f\n  PrecipIn: %.2f\n"+
			"  Humidity: %d\n  Cloud: %d\n  FeelslikeC: %.2f\n  FeelslikeF: %.2f\n"+
			"  WindchillC: %.2f\n  WindchillF: %.2f\n  HeatindexC: %.2f\n  HeatindexF: %.2f\n"+
			"  DewpointC: %.2f\n  DewpointF: %.2f\n  VisKm: %.2f\n  VisMiles: %.2f\n"+
			"  UV: %.2f\n  GustMph: %.2f\n  GustKph: %.2f\n  ShortRad: %.2f\n  DiffRad: %.2f\n  DNI: %.2f\n  GTI: %.2f\n",
		wr.Location.Name, wr.Location.Region, wr.Location.Country, wr.Location.Lat, wr.Location.Lon, wr.Location.TzID, wr.Location.LocaltimeEpoch, wr.Location.Localtime,
		wr.Current.LastUpdatedEpoch, wr.Current.LastUpdated, wr.Current.TempC, wr.Current.TempF, wr.Current.IsDay,
		wr.Current.Condition.Text, wr.Current.Condition.Icon, wr.Current.Condition.Code,
		wr.Current.WindMph, wr.Current.WindKph, wr.Current.WindDegree, wr.Current.WindDir,
		wr.Current.PressureMb, wr.Current.PressureIn, wr.Current.PrecipMm, wr.Current.PrecipIn,
		wr.Current.Humidity, wr.Current.Cloud, wr.Current.FeelslikeC, wr.Current.FeelslikeF,
		wr.Current.WindchillC, wr.Current.WindchillF, wr.Current.HeatindexC, wr.Current.HeatindexF,
		wr.Current.DewpointC, wr.Current.DewpointF, wr.Current.VisKm, wr.Current.VisMiles,
		wr.Current.UV, wr.Current.GustMph, wr.Current.GustKph, wr.Current.ShortRad, wr.Current.DiffRad, wr.Current.DNI, wr.Current.GTI,
	)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("❌ Usage goweathergo <city>")
		os.Exit(1)
	}

	city := strings.Join(os.Args[1:], " ")

	APIKey := os.Getenv("WeatherAPIKey")

	url := fmt.Sprintf(config.BaseURL + config.CurrentWeatherEndpoint + "?key=" + APIKey + "&q=" + city)
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("❌ Error fetching weather:", err)
		return
	}
	defer resp.Body.Close() // Ensure the response body is closed after reading

	if resp.StatusCode != http.StatusOK {
		fmt.Println("❌ Error: Received non-OK HTTP status:", resp.Status)
		return
	}

	var data WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println("❌ Error decoding response:", err)
	}

	fmt.Printf("Weather in %s, %s, %s: \n", data.Location.Name, data.Location.Region, data.Location.Country)
	fmt.Printf("Temperature: %.1f°C\n", data.Current.TempC)

	input := WeatherInString(data)

	ai.GetClothingRecommendation(input)

}
