package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

type ImageUpload struct {
	Auth  Authentication
	Image string
}

type Authentication struct {
	Id        int64
	Key       []byte
	User      int64
	ExpiredBy time.Time
}

type PersonUpload struct {
	New  bool
	User Person
}

type Person struct {
	Name     string
	Gender   bool
	Location string
	Password string
	Email    string
}

type Library struct {
	Auth Authentication
	Data []Song
}

type Song struct {
	Name   string
	Artist string
	Length int
	Genre  string
}

func register() Authentication {
	var data = new(PersonUpload)
	data.User.Name = "bob"
	data.User.Gender = true
	data.User.Location = "Markham"
	data.User.Password = "passwordpasswordpasswordpasswordpassword"
	data.User.Email = "whalelord@email.com"
	data.New = true
	ws, err := websocket.Dial("ws://localhost:8080/user", "", "http://localhost")
	if err != nil {
		log.Printf("Something went bad %s", err)
	}
	websocket.JSON.Send(ws, &data)
	var out Authentication
	for {
		if err := websocket.JSON.Receive(ws, &out); err == io.EOF {
			break
		}
	}
	return out
}

func login(email, password string) Authentication {
	var data = new(PersonUpload)
	data.User.Password = "passwordpasswordpasswordpasswordpassword"
	data.User.Email = "whalelord@email.com"
	data.New = false
	ws, err := websocket.Dial("ws://localhost:8080/user", "", "http://localhost")
	if err != nil {
		log.Printf("Something went bad %s", err)
	}
	websocket.JSON.Send(ws, &data)
	var out Authentication
	for {
		if err := websocket.JSON.Receive(ws, &out); err == io.EOF {
			break
		}
	}
	return out
}

func upload_image(auth Authentication, path string) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Something went bad %s", err)
	}
	var data = new(ImageUpload)
	data.Auth = auth
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

func upload_library(auth Authentication) {
	var data = new(Library)
	data.Auth = auth
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < 10; i++ {
		var song = new(Song)
		song.Name = "Something something song name"
		song.Artist = "Something something artist"
		song.Length = rand.Int() % 720
		song.Genre = "Somethign something genre"
		data.Data = append(data.Data, *song)
	}
	ws, err := websocket.Dial("ws://localhost:8080/library", "", "http://localhost")
	if err != nil {
		log.Printf("Something went bad %s", err)
	}
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

func main() {
	auth := login("whalelord@email.com", "passwordpasswordpasswordpassword")
}
