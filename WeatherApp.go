package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"
)

type WeatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func initConfigLoader() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
		panic(err)
	}
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status:1","message":"Testing API"}`))
}

func query(city string) (error, WeatherData) {

	getKey := os.Getenv("OPEN_WEATHER_API_KEYS")
	weatherURL := os.Getenv("OPEN_WEATHER_BASEURL") + getKey + "&q=" + city

	response, err := http.Get(weatherURL)

	if err != nil {
		return err, WeatherData{}
	}

	defer response.Body.Close()
	var weatherData WeatherData

	if err := json.NewDecoder(response.Body).Decode(&weatherData); err != nil {
		return err, weatherData
	}
	return nil, weatherData
}

func getWeather(w http.ResponseWriter, r *http.Request) {
	city := strings.SplitN(r.URL.Path, "/", 3)[2]
	err, data := query(city)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func main() {

	initConfigLoader()
	port := os.Getenv("SERVER_PORT")

	http.HandleFunc("/testApi", testHandler)
	http.HandleFunc("/weather/", getWeather)

	fmt.Println("Server started on Port :2000")
	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		panic("Server is running of the Specified Port")
	}
}
