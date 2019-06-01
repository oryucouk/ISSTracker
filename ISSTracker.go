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

func handleMyErrors(e error, text string) {
	if e != nil {
	log.Fatalf(text, e)
    }
}

func getCoordinates(location string) (float64, float64) {
	//Takes a location and returns latitude, longitude coordinates as float by querying API
	const siteAPI = "https://geocode.xyz/"
	const siteAPISuffix = "?geoit=csv"
	var lat, lon float64
	urlFull := strings.Join([]string{siteAPI, location, siteAPISuffix},"")
	
	
	response, err := http.Get(urlFull)
	handleMyErrors(err, "http.Get issue in getCoordinates function")

	responseData, err := ioutil.ReadAll(response.Body)
	handleMyErrors(err, "ioutil.ReadAll issue in getCoordinates function ")
	response.Body.Close()
    
	latLonArray := strings.Split(string(responseData), ",")
	//fmt.Println(latLonArray)
	lat,_ = strconv.ParseFloat(latLonArray[2],64)
	lon,_ = strconv.ParseFloat(latLonArray[3],64)
	return lat, lon
}




func main() {

	const urlISSCurLoc = "http://api.open-notify.org/iss-now.json"
	const urlISS = "http://api.open-notify.org/iss-pass.json"
	lat, lon  := getCoordinates(os.Args[1])
	urlFull := strings.Join([]string{urlISS, "?", "lat=", strconv.FormatFloat(lat, 'f', 4, 64), "&lon=", strconv.FormatFloat(lon, 'f', 4, 64)}, "")

	response, err := http.Get(urlFull)
	handleMyErrors(err, "http.Get issue in main function ")
	

	responseData, err := ioutil.ReadAll(response.Body)
	handleMyErrors(err, "ioutil.ReadAll issue in main function ")
	response.Body.Close()


	var responseObject BigResponse
	json.Unmarshal(responseData, &responseObject)

    fmt.Println("Timestamp     : ",time.Now())
    fmt.Println("Your Location : ", os.Args[1])
	fmt.Println("Your Longitude: ", responseObject.Request.Longitude)
	fmt.Println("Your Latitude : ", responseObject.Request.Latitude)
	fmt.Println("Number of ISS sightings requested: ", responseObject.Request.Passes,"\n")
	for k, _ := range responseObject.Response {
		fmt.Println("Next ISS sighting for ", responseObject.Response[k].Duration/60, "mins at ", time.Unix(responseObject.Response[k].Risetime, 0))
	}
    //Find out current position of ISS, reverse geocode it and print the city/place.
    responseNow, err := http.Get(urlISSCurLoc)
    handleMyErrors(err, "http.Get error for current ISS position ")

    responseNowData, err := ioutil.ReadAll(responseNow.Body)
    handleMyErrors(err,"ioutil.ReadAll issue for current ISS position ")
    responseNow.Body.Close()
    //fmt.Println(string(responseNowData))

    var responseNowObject BigResponseNow
    json.Unmarshal(responseNowData, &responseNowObject)
    //fmt.Println(responseNowObject)

    fmt.Println("\n\nISS Current Longitude:", responseNowObject.ISSPositionObj.Longitude)
    fmt.Println("ISS Current Latitude :", responseNowObject.ISSPositionObj.Latitude)



}

type BigResponseNow struct{
	ISSPositionObj ISSPosition `json:"iss_position"`
	Timestamp int64 `json:"timestamp"`
	Message string `json:"message"`

}

type ISSPosition struct{
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
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
