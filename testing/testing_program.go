package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
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

// Returns the data to be sent to the server
func prepare_data() Library {
	var data = new(Library)

	// opens genres.json and dumps contents into list_of_genres
	genre_file_contents, err := ioutil.ReadFile("genres.json")
	if err != nil {
	}
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
		song.Length = rand.Int() % 720
		song.Genre = random_genre(list_of_genres.Name)
		log.Printf("%+v", song)
		data.Data = append(data.Data, *song)
	}
	return *data
}

// Returns the data to be sent to the server, includes only one song
func prepare_one_song() Library {
	var data = new(Library)

	// opens genres.json and dumps contents into list_of_genres
	genre_file_contents, err := ioutil.ReadFile("genres.json")
	if err != nil {
	}
	list_of_genres := &Genre{}
	json.Unmarshal(genre_file_contents, &list_of_genres)

	// opens songnames.json and dumps contents into list_of_titles
	song_names_file_contents, err := ioutil.ReadFile("song_names.json")
	list_of_titles := &Song_Titles{}
	json.Unmarshal(song_names_file_contents, &list_of_titles)

	data.User = 1
	rand.Seed(time.Now().UTC().UnixNano())
	var song = new(Song)
	song.Name = random_name(rand.Intn(1)+1, list_of_titles.Name)
	song.Artist = random_name(rand.Intn(5)+1, list_of_titles.Name)
	song.Length = rand.Int() % 7200
	song.Genre = random_genre(list_of_genres.Name)
	log.Printf("%+v", song)
	data.Data = append(data.Data, *song)
	return *data
}

// Returns the data the server received
func send_data(data Library) Library {
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

	// get the songs that was sent to the server
	var server_data = new(Library)
	log.Printf("%v", server_data)
	for i := 0; i < len(data.Data); i++ {
		// id number is the num of songs in the database before sending stuff + i + 1
		db.Where("id = ?", i+int(num_of_songs.Length)+1).Find(&server_data.Data)
	}
	return *server_data
}

// EVERYTHING BELOW HERE SENDS 10 SONGS TO THE SERVER AND CHECKS IF THE
// SERVER RECEIVED ALL OF THE SONGS CORRECTLY

func main() {
	data := prepare_data()
	server_data := send_data(data)

	fmt.Println("Sent data")
	fmt.Println(data)
	fmt.Println(len(data.Data))
	fmt.Println()
	fmt.Println("Server received data")
	fmt.Println(server_data)
	fmt.Println(len(server_data.Data))
	fmt.Println()

	var same_data = true
	if len(server_data.Data) == 0 {
		same_data = false
		fmt.Println("The server did not receive any data.")
	} else {
		// use a loop to go through each song and check if data matches
		for i := 0; i < len(data.Data) && same_data && len(server_data.Data) > 0; i++ {
			if data.Data[i] != server_data.Data[i] {
				same_data = false
				fmt.Println("Data does not match!")
				fmt.Println("From sent data, song #", i+1, data.Data[i])
				fmt.Println("From server data", server_data.Data[i])
			}
		}
	}
	if same_data {
		fmt.Println("Data matches!")
	}
}
