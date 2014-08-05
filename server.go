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

type ImageResult struct {
	Faces string
}

func analysis(ws *websocket.Conn, path string) {
	out, err := exec.Command("python2.7", "analysis/analysis.py", path).Output()
	if err != nil {
		log.Fatal("error:", err)
	}
	ws.Write(out)
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

func main() {
	log.Println("Starting Lyra Server")
	http.Handle("/image", websocket.Handler(imageHandler))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("Something went bad with the server: %s", err)
	}
}
