package main

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"os"
	"bufio"
)

func check(e error){
	if e != nil{
		panic(e)
	}
}

//Source of help
//https://medium.com/@rafaelacioly/construindo-uma-api-restful-com-go-d6007e4faff6

func CityAPICall(cityID string){
	//TODO - get API data from
	//https://openweathermap.org/forecast5

	//Build JSON
	//return JSON
}

func GetCityWeather(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)
	//cityID := params["cityID"]

	//TODO - Call CityAPICall
	//TODO - Build JSON
	//TODO - Return JSON
}

func GetCityList() (array []string) {
	//Define Filepath
	//TODO - Improve reusability
	f, err := os.Open("/home/user/Desktop/cityIDs.txt")
	check(err)

	defer func() {
        	if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	//Read Lines
	s := bufio.NewScanner(f)

	//Store Each Line
	for s.Scan() {
		array = append(array, s.Text())	  
		//fmt.Println(s.Text())
    	}
	check(s.Err())
	return
}

func GetWeather(w http.ResponseWriter, r *http.Request) {
	//get city IDs from file
	array := GetCityList()
	for i:= range array {
		fmt.Println(array[i])	
	}

	//TODO - Call CityAPICall
	//TODO - Build JSON
	json.NewEncoder(w).Encode(array)
	
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/city/{cityID}",GetCityWeather).Methods("GET")
	router.HandleFunc("/city",GetWeather).Methods("GET")
	fmt.Println("Deploying Weather API Server")
	log.Fatal(http.ListenAndServe(":8000",router))
}
