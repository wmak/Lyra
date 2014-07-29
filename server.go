package main

import (
	"code.google.com/p/go.net/websocket"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type ImageUpload struct {
	User string
	Image string
}

func imageHandler(ws *websocket.Conn) {
	var data = new(ImageUpload)
	if err := websocket.JSON.Receive(ws, &data); err != nil {
		log.Printf("Error in the image handler %s", err)
	}
	//confirm data.User
	Image, err := base64.StdEncoding.DecodeString(data.Image)
	if err != nil {
		log.Fatal("error:", err)
	}

	hasher := md5.New()
	hasher.Write([]byte(Image))
	Sum := hex.EncodeToString(hasher.Sum(nil))
	ioutil.WriteFile("images/" + Sum + ".jpg", Image, 0644)
	log.Printf("Saved new image at %s", "images/" + Sum + ".jpg")
}

func main() {
	fmt.Println("Starting Lyra Server")
	http.Handle("/image", websocket.Handler(imageHandler))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("Something went bad with the server: %s", err)
	}
}
