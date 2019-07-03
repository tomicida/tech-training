package main

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

//Source of help
//https://medium.com/@rafaelacioly/construindo-uma-api-restful-com-go-d6007e4faff6

func CityAPICall(cityID string){
	//TODO - get API data from
	//https://openweathermap.org/forecast5

	//Build JSON
	//return JSON
}

func GetCityWeather(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	cityID := params["cityID"]

	//TODO - Call CityAPICall
	//TODO - Build JSON
	//TODO - Return JSON
}

func GetWeather(w http.ResponseWriter, r *http.Request) {
	//TODO - get city IDs from file
	//TODO - Build JSON
	
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/city/{cityID}",GetCityWeather).Methods("GET")
	router.HandleFunc("/city",GetWeather).Methods("GET")
	fmt.Println("Deploying Weather API Server")
	log.Fatal(http.ListenAndServe(":8000",router))
}
