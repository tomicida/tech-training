package main

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"os"
	"bufio"
	"math"
)

type Sample struct {
	Time	string	`json:"time"`
	Temp	float64	`json:"temp"`
	Weather	string	`json:"weather"`	
}	

type Day struct { 
	TempMin	float64	`json:"temp_min"`
	TempMax float64	`json:"temp_max"`
	TempAvg float64	`json:"temp_avg"`
	Date	string	`json:"date"`
	Rain	bool	`json:"rain"`
	Sample	[]Sample	`json:"sample"`
}	

type Prediction struct {
	Cod     string  `json:"cod"`
	Message float64 `json:"message"`
	Cnt     int     `json:"cnt"`
	List    []struct {
		Dt   int `json:"dt"`
		Main struct {
			Temp      float64 `json:"temp"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  float64 `json:"pressure"`
			SeaLevel  float64 `json:"sea_level"`
			GrndLevel float64 `json:"grnd_level"`
			Humidity  int     `json:"humidity"`
			TempKf    float64 `json:"temp_kf"`
		} `json:"main"`
		Weather []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Wind struct {
			Speed float64 `json:"speed"`
			Deg   float64 `json:"deg"`
		} `json:"wind"`
		Sys struct {
			Pod string `json:"pod"`
		} `json:"sys"`
		DtTxt string `json:"dt_txt"`
		Rain  struct {
			ThreeH float64 `json:"3h"`
		} `json:"rain,omitempty"`
	} `json:"list"`
	City struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Coord struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
		Country  string `json:"country"`
		Timezone int    `json:"timezone"`
	} `json:"city"`
}

func check(e error){
	if e != nil{
		panic(e)
	}
}

//Source of help
//https://medium.com/@rafaelacioly/construindo-uma-api-restful-com-go-d6007e4faff6

func CityAPICall(cityID string, key string)(record Prediction){
	//TODO - get API data from
	//https://openweathermap.org/forecast5
	client := &http.Client{}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?id=%s&APPID=%s",cityID,key)

	req, err := http.NewRequest("GET", url, nil)
	check(err)
		
	resp, err := client.Do(req)
	check(err)

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&record)
	check(err)

	return
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

func GetAPIKey(index int) (key string){
	//Define Filepath
	//TODO - Improve reusability
	f, err := os.Open("/home/user/Desktop/APIkeys.txt")
	check(err)

	defer func() {
        	if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	//Read Lines
	s := bufio.NewScanner(f)

	//Search for key
	i := 0
	for s.Scan() {
		
		if i == index {
			key = s.Text()
		}  
		i++
    	}
	check(s.Err())
	return
}

func SendableConverter(p Prediction) (s []Day){
	numberOfDays := 5
	sampleInterval := 3
	samplesPerDay := 24/sampleInterval

	for i := 0; i < numberOfDays ; i++ {
		dayIndex := i*samplesPerDay
		
		willRain := false		
		minTemp := KelvinToCelsius(p.List[dayIndex].Main.Temp)
		maxTemp := KelvinToCelsius(p.List[dayIndex].Main.Temp)
		avgTemp := float64(0)

		s = append(s,Day{})
		s[i].Date = p.List[dayIndex].DtTxt[:10]
		
		//TODO - Fix index/days counters
		for j := 0; j < samplesPerDay; j++{
			currentSample := p.List[dayIndex + j]
			var sample Sample
	
			sample.Time = currentSample.DtTxt[11:]
			sample.Temp = KelvinToCelsius(currentSample.Main.Temp)
			sample.Weather = currentSample.Weather[0].Main
			
			if minTemp > sample.Temp {
				minTemp = sample.Temp
			}
			if maxTemp < sample.Temp {
				maxTemp = sample.Temp		
			}
			avgTemp += sample.Temp
			if sample.Weather == "Rain"{
				willRain = true			
			}

			s[i].Sample = append(s[i].Sample, sample)
		}
		
		s[i].TempAvg = reduce(avgTemp/float64(samplesPerDay))
		s[i].TempMin = minTemp
		s[i].TempMax = maxTemp
		s[i].Rain = willRain
	}
	

	return
}

func KelvinToCelsius (k float64) (c float64) {
	c = reduce(k - 272.15)
	return
}

func reduce (o float64) (n float64) {
	n = math.Floor((o*100))/100
	return 
}

func GetWeather(w http.ResponseWriter, r *http.Request) {
		
	//get city IDs from file
	array := GetCityList()

	//Get API Key
	key := GetAPIKey(0)
		
	//Call CityAPICall
	var resultsArray []Prediction
	for index := range array{
		appendable := CityAPICall(array[index],key)
		resultsArray = append(resultsArray,appendable)	
	}
	
	

	//TODO - Build JSON
	json.NewEncoder(w).Encode(SendableConverter(resultsArray[0]))
	
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/city/{cityID}",GetCityWeather).Methods("GET")
	router.HandleFunc("/city",GetWeather).Methods("GET")
	fmt.Println("Deploying Weather API Server")
	log.Fatal(http.ListenAndServe(":8000",router))
}
