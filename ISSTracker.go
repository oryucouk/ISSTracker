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

func main() {
	const url = "http://api.open-notify.org/iss-pass.json"
	lat := 51.467
	lon := -0.0145

	urlFull := strings.Join([]string{url, "?", "lat=", strconv.FormatFloat(lat, 'f', 4, 64), "&lon=", strconv.FormatFloat(lon, 'f', 4, 64)}, "")
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

	fmt.Println("API fetch status: ", responseObject.Message)
	fmt.Println("Your longitude: ", responseObject.Request.Longitude)
	fmt.Println("Your latitude: ", responseObject.Request.Latitude)
	fmt.Println("Default number of passes: ", responseObject.Request.Passes)
	for k, _ := range responseObject.Response {
		fmt.Println("You will see the ISS next for a total of ", responseObject.Response[k].Duration/60, " mins duration and it will begin to show at ", time.Unix(responseObject.Response[k].Risetime, 0))
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
