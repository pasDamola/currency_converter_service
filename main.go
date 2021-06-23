package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var apis map[int]string

func getDataResponse(API int) {
	url := apis[API]
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			var result map[string]interface{}
			json.Unmarshal([]byte(body), &result)
			switch API {

			case 1:
				if result["success"] == true {
					fmt.Println(result["rates"].(map[string]interface{})["USD"])
				} else {
					fmt.Println(result["error"].(map[string]interface{})["info"])
				}
			case 2: // for the openweathermap.org API
				if result["main"] != nil {
					fmt.Println(result["main"].(map[string]interface{})["temp"])
				} else {
					fmt.Println(result["message"])
				}

			}
		} else {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
}

func main() {
	apis = make(map[int]string)
	apis[1] = "http://data.fixer.io/api/latest?access_key=21ce410f9900f7c07838c2e6d1f80e31"
	apis[2] = "http://api.openweathermap.org/data/2.5/weather?q=London&appid=61b3761a8c87df93dff55c1ca9eb93a4"

	go getDataResponse(1)
	go getDataResponse(2)

	fmt.Scanln()
}
