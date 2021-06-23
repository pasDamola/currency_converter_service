package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var apis map[int]string
var c chan map[int]interface{}

func getDataResponse(API int) {
	url := apis[API]
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			var result map[string]interface{}
			json.Unmarshal([]byte(body), &result)
			var re = make(map[int]interface{})
			switch API {

			case 1:
				if result["success"] == true {
					re[API] = result["rates"].(map[string]interface{})["USD"]
				} else {
					re[API] = result["rates"].(map[string]interface{})["info"]
				}
				// store the result into the channel
				c <- re
				fmt.Println("Result for API 1 stored")
			case 2: // for the openweathermap.org API
				if result["main"] != nil {
					re[API] = result["main"].(map[string]interface{})["temp"]
				} else {
					re[API] = result["message"]
				}
				c <- re
				fmt.Println("Result for API 2 stored")

			}
		} else {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
}

func main() {
	// creates a channel to store the results from the
	// API calls

	c = make(chan map[int]interface{})
	apis = make(map[int]string)
	apis[1] = "http://data.fixer.io/api/latest?access_key=21ce410f9900f7c07838c2e6d1f80e31"
	apis[2] = "http://api.openweathermap.org/data/2.5/weather?q=London&appid=61b3761a8c87df93dff55c1ca9eb93a4"

	go getDataResponse(1)
	go getDataResponse(2)

	// we expect two results in the channel
	for i := 0; i < 2; i++ {
		fmt.Println(<-c)
	}

	fmt.Println("Done!")
	fmt.Scanln()
}
