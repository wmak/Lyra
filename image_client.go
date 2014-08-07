package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type ImageUpload struct {
	User  string
	Image string
}

type ImageResult struct {
	Faces int
}

func main() {
	file, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Printf("Something went bad %s", err)
	}
	var data = new(ImageUpload)
	data.User = "bob"
	ws, err := websocket.Dial("ws://localhost:8080/image", "", "http://localhost")
	if err != nil {
		log.Printf("Something went bad %s", err)
	}
	data.Image = base64.StdEncoding.EncodeToString(file)
	websocket.JSON.Send(ws, &data)
	var out []byte
	for {
		if err := websocket.Message.Receive(ws, &out); err == io.EOF {
			log.Printf("Exiting, %s", err)
			break
		}
		log.Printf("final result %s", out)
	}
}
