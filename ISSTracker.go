package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"strconv"
	"strings"
	"time"
)

func curl(location string) (float64, float64) {
	var lat, lon float64
	urlFull := strings.Join([]string{"https://geocode.xyz/", location,"?geoit=csv"},"")
	response, err := http.Get(urlFull)
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Curl function error", err)
	}
    
	stringifyResponse := string(responseData)
	latLonArray := strings.Split(stringifyResponse, ",")
	//fmt.Println(latLonArray)
	lat,_ = strconv.ParseFloat(latLonArray[2],64)
	lon,_ = strconv.ParseFloat(latLonArray[3],64)
	return lat, lon
}


func main() {

	const urlISS = "http://api.open-notify.org/iss-pass.json"
	lat, lon  := curl(os.Args[1])
	urlFull := strings.Join([]string{urlISS, "?", "lat=", strconv.FormatFloat(lat, 'f', 4, 64), "&lon=", strconv.FormatFloat(lon, 'f', 4, 64)}, "")

	response, err := http.Get(urlFull)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(responseData))

	var responseObject BigResponse
	json.Unmarshal(responseData, &responseObject)

	fmt.Println("Your longitude: ", responseObject.Request.Longitude)
	fmt.Println("Your latitude: ", responseObject.Request.Latitude)
	fmt.Println("Number of passes requested: ", responseObject.Request.Passes)
	for k, _ := range responseObject.Response {
		fmt.Println("Next ISS sighting for ", responseObject.Response[k].Duration/60, "mins at ", time.Unix(responseObject.Response[k].Risetime, 0))
	}
}

type BigResponse struct {
	Message  string `json:"message"`
	Request  TypeRequest
	Response []TypeResponse
}

type TypeRequest struct {
	Altitude  int64   `json:"altitude"`
	Datetime  int64   `json:"datetime"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Passes    int8    `json:"passes"`
}

type TypeResponse struct {
	Duration int64 `json:"duration"`
	Risetime int64 `json:"risetime"`
}
