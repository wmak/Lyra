package main

import (
	"code.google.com/p/go.net/websocket"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

type ImageUpload struct {
	User  string
	Image string
}

type Library struct {
	User  string
	Songs []Song
}

type Song struct {
	Name   string
	Artist string
	Length int
	Genre  string
}

func analysis(ws *websocket.Conn, path string) {
	_, err := exec.Command("python2.7", "analysis/analysis.py", path).Output()
	if err != nil {
		log.Fatal("error:", err)
	}
	websocket.Message.Send(ws, "analysis complete")
	log.Printf("Analysis complete on %s", path)
}

func imageHandler(ws *websocket.Conn) {
	var data = new(ImageUpload)
	if err := websocket.JSON.Receive(ws, &data); err != nil {
		log.Printf("Error in the image handler %s", err)
	}
	//confirm data.User
	log.Printf("Connection from %s", data.User)
	Image, err := base64.StdEncoding.DecodeString(data.Image)
	if err != nil {
		log.Fatal("error:", err)
	}
	hasher := md5.New()
	hasher.Write([]byte(Image))
	Sum := hex.EncodeToString(hasher.Sum(nil))
	path := "images/" + Sum + ".jpg"
	ioutil.WriteFile(path, Image, 0644)
	log.Printf("Saved new image at %s", path)
	analysis(ws, path)
}

func libraryHandler(ws *websocket.Conn) {
	var data = new(Library)
	if err := websocket.JSON.Receive(ws, &data); err != nil {
		log.Printf("Error in the library handler %s", err)
	}
	//confirm data.User
	log.Printf("Connection from %s", data.User)
	//Go through song list.
	//Fuzzy search for song, if none found
	//Add new song
	//Associate song with user
	//If firsttime, add all songs
	//else delete songs they may no longer have
}

func main() {
	log.Println("Starting Lyra Server")
	http.Handle("/image", websocket.Handler(imageHandler))
	http.Handle("/library", websocket.Handler(libraryHandler))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("Something went bad with the server: %s", err)
	}
}
