package main

import (
	"code.google.com/p/go.net/websocket"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type ImageUpload struct {
	User string
	Data string
}

type Library struct {
	User int64
	Data []Song
}

type Person struct {
	Id       int64
	Name     string `sql:"not null"`
	Gender   bool
	Location string //validate it?
	Password string //gotta hash that shit
	Email    string
	Songs    []Song `gorm:"many2many:person_library;"`
}

type Song struct {
	Id     int64
	Name   string `sql:"not null"`
	Genre  string
	Length int
	Artist string
	Nsfw   bool
	Bpm    int
	Volume int
}

type Listen struct {
	Id        int64
	CreatedAt time.Time `sql:"not null"`
	SId       int64     `sql:"not null"`
	PId       int64     `sql:"not null"`
	Location  Loc
	Skip      bool
	Faces     int
	Rgb       Colour
	Hue       int
	Lighting  int
	Volume    int
}

type Loc struct {
	Latitude  float64
	Longitude float64
}

type Colour struct {
	Red   float64
	Green float64
	Blue  float64
}

func errorcheck(err error, msg string) {
	if err != nil {
		log.Print(msg)
	}
}

func analysis(ws *websocket.Conn, path string) {
	_, err := exec.Command("python2.7", "analysis/analysis.py", path).Output()
	errorcheck(err, "something went bad with analysis")
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
	Image, err := base64.StdEncoding.DecodeString(data.Data)
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

func libraryHandler(db gorm.DB) websocket.Handler {
	return func(ws *websocket.Conn) {
		var data = new(Library)
		if err := websocket.JSON.Receive(ws, &data); err != nil {
			log.Printf("Error in the library handler %s", err)
		}
		//confirm data.User
		var user = Person{}
		log.Printf("Connection from %d", data.User)
		db.Table("persons").Where("id = ?", data.User).First(&user)
		log.Printf("%+v", user)
		//if user is unknown
		if user.Id != 0 {
			db.Model(&user).Association("Songs").Clear()
			//Go through song list.
			for i := 0; i < len(data.Data); i++ {
				var songs = Song{}
				//search for song, if none found
				db.Table("songs").Where("name = ?", data.Data[i].Name).First(&songs)
				log.Printf("%d", songs.Id)
				if songs.Id == 0 {
					//Associate song with user
					db.Model(&user).Association("Songs").Append(data.Data[i])
				} else {
					db.Model(&user).Association("Songs").Append(songs)
				}
			}
			user.Songs = data.Data
		} else {
			websocket.Message.Send(ws, "WHO IS THIS?")
		}
		//If firsttime, add all songs
		//else delete songs they may no longer have
	}
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

func initialize(db gorm.DB) {
	log.Println("Initilaizing tables")
	//drop any existing tables
	db.DropTableIfExists(Person{})
	db.DropTableIfExists(Listen{})
	db.DropTableIfExists(Song{})
	//add tables.
	db.CreateTable(Person{})
	db.CreateTable(Listen{})
	db.CreateTable(Song{})
	//add indexes
	db.Model(Song{}).AddIndex("idx_song_name", "name")
}

func main() {
	db := initDb()
	if len(os.Args) > 1 {
		if os.Args[1] == "-i" {
			initialize(db)
		}
	}
	log.Println("Starting Lyra Server")
	http.Handle("/image", websocket.Handler(imageHandler))
	http.Handle("/library", websocket.Handler(libraryHandler(db)))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("Something went bad with the server: %s", err)
	}
}
