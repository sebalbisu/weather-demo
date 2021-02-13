package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	"golang.org/x/net/websocket"
)

// server addresses
const serverHost = "http://localhost"
const weatherUrl = "ws://localhost:5000/weather"

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

func main() {
	ws, err := websocket.Dial(weatherUrl, "", serverHost)
	if err != nil {
		log.Fatal(err)
	}

	for {
		var response = new(Response)
		if err := websocket.JSON.Receive(ws, &response); err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("connection closed")
				break
			}
			log.Print(err)
			continue
		}
		out, err := json.Marshal(response)
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Printf("%+v\n", string(out))
	}
}
