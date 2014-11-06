package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

type Library struct {
	User int64
	Data []Song
}

type Song struct {
	Name   string
	Artist string
	Length int
	Genre  string
}

type Genre struct {
	Name []string `json:"Genre"`
}

type Song_Titles struct {
	Name []string `json:"Name"`
}

type Table_Length struct {
	Length uint8
}

func random_genre(genres []string) string {
	return genres[rand.Intn(len(genres))]
}

func random_name(x int, song_names []string) string {
	var name string
	for i := 0; i < x; i++ {
		name += song_names[rand.Intn(len(song_names))]
		name += " "
	}
	return name[0 : len(name)-1]
}

func initDb() gorm.DB {
	//Create a sql connection to the database
	db, err := gorm.Open("mysql",
		"root:password@/Lyra?charset=utf8&parseTime=True")
	if err != nil {
		log.Printf("yup not working: %s", err)
	}
	//Setup logging
	db.LogMode(true)
	return db
}

func main() {
	//for i := 0; i < 50; i++ {
	var data = new(Library)

	// opens genres.json and dumps contents into list_of_genres
	genre_file_contents, err := ioutil.ReadFile("genres.json")
	list_of_genres := &Genre{}
	json.Unmarshal(genre_file_contents, &list_of_genres)

	// opens songnames.json and dumps contents into list_of_titles
	song_names_file_contents, err := ioutil.ReadFile("song_names.json")
	list_of_titles := &Song_Titles{}
	json.Unmarshal(song_names_file_contents, &list_of_titles)

	data.User = 1
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < 10; i++ {
		var song = new(Song)
		song.Name = random_name(rand.Intn(1)+1, list_of_titles.Name)
		song.Artist = random_name(rand.Intn(5)+1, list_of_titles.Name)
		song.Length = rand.Int() % 7200
		song.Genre = random_genre(list_of_genres.Name)
		log.Printf("%+v", song)
		data.Data = append(data.Data, *song)
	}

	var db = initDb()
	// get the number of rows in the songs table
	var num_of_songs = new(Table_Length)
	db.Table("songs").Count(&num_of_songs.Length)

	// create connection to server
	ws, err := websocket.Dial("ws://localhost:8080/library", "", "http://localhost")
	if err != nil {
		log.Printf("Something went bad %s", err)
	}
	// send data to server
	websocket.JSON.Send(ws, &data)
	var out []byte
	for {
		if err := websocket.Message.Receive(ws, &out); err == io.EOF {
			log.Printf("Exiting, %s", err)
			break
		}
		log.Printf("final result %s", out)
	}

	// GET THE ROWS THAT ARE IN THE TABLE AFTER THE NUMBER OF ROWS
	// THAT WAS GRABBED FROM THE SERVER AND CHECK IF IT IS THE SAME
	// AS WHAT WAS SENT

	// get the songs that was sent to the server
	var server_data = new(Library)
	log.Printf("%v", server_data)
	//for i := 0; i < len(data.Data); i++ {
	//db.Where("name = ?", data.Data[i].Name).Find(&server_data.Data)
	//db.First(&server_data.Data, i + int(num_of_songs.Length) + 1)
	//}
	db.Last(&server_data.Data)

	fmt.Println("Sent data")
	fmt.Println(data)
	fmt.Println()
	fmt.Println("Server data")
	fmt.Println(server_data)
	fmt.Println()
	// use a loop to go through each song and check if data matches
	var same_data = true
	for i := 0; i < len(data.Data) && same_data; i++ {
		if data.Data[i] != server_data.Data[i] {
			same_data = false
			fmt.Println("Data does not match!")
			fmt.Println("From sent data, song #", i+1, data.Data[i])
			fmt.Println("From server data", server_data.Data[i])
		}
	}
	// if !same_data {
	// 	break
	// }
	//}
}
