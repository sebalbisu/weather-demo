package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

// server address
const address = "localhost:5000"

// Temperature type of json websocket response
type Temperature struct {
	Id    string
	Value float64
}

// Humidity type of json websocket response
type Humidity struct {
	Id    string
	Value float64
}

// Response type of websocket
type Response struct {
	Temperature []Temperature
	Humidity    []Humidity
}

// getData from raw data source and parse it
func getData() (response Response, err error) {
	file, err := ioutil.ReadFile("./server/data.json")
	if err != nil {
		log.Fatal(err)
	}

	type WeatherJson struct {
		Weather []struct {
			Id          string
			Temperature float64
			Humidity    float64
		}
	}
	var weatherData = new(WeatherJson)

	if err = json.Unmarshal([]byte(file), &weatherData); err != nil {
		return response, err
	}

	temps := make([]Temperature, 0)
	hum := make([]Humidity, 0)
	for _, item := range weatherData.Weather {
		temps = append(temps, Temperature{Id: item.Id, Value: item.Temperature})
		hum = append(hum, Humidity{Id: item.Id, Value: item.Humidity})
	}

	return Response{Temperature: temps, Humidity: hum}, nil
}

// Updates the client with new weathers
func Updates(responseCh chan Response) {
	for {
		if data, err := getData(); err != nil {
			log.Printf("error: %s", err)
		} else {
			responseCh <- data
		}
		time.Sleep(time.Second)
	}
}

// Server is the websocket server handler
func Server(ws *websocket.Conn) {

	responseCh := make(chan Response, 0)
	go Updates(responseCh)

	for {
		select {
		case response := <-responseCh:
			websocket.JSON.Send(ws, response)
		}
	}
}

func main() {
	http.Handle("/weather", websocket.Handler(Server))

	fmt.Printf("listening on %s/weather\n", address)

	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %s", err.Error())
	}
}
