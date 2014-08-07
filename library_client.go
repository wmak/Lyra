package main

import (
	"code.google.com/p/go.net/websocket"
	"io"
	"log"
	"math/rand"
)

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

func random_genre() string{
	genres := []string{
		"Alternative Rock",
		"College Rock",
		"Experimental Rock",
		"Goth Rock",
		"Grunge",
		"Hardcore Punk",
		"Hard Rock",
		"Indie Rock",
		"New Wave",
		"Progressive Rock",
		"Punk",
		"Acoustic Blues",
		"Chicago Blues",
		"Classic Blues",
		"Contemporary Blues",
		"Country Blues",
		"Delta Blues",
		"Electric Blues",
		"Lullabies",
		"Sing-Along",
		"Stories",
		"Avant-Garde",
		"Baroque",
		"Chamber Music",
		"Chant",
		"Choral",
		"Classical Crossover",
		"Early Music",
		"High Classical",
		"Impressionist",
		"Medieval",
		"Minimalism",
		"Modern Composition",
		"Opera",
		"Orchestral",
		"Renaissance",
		"Romantic",
		"Wedding Music",
		"Novelty",
		"Standup Comedy",
		"Alternative Country",
		"Americana",
		"Bluegrass",
		"Contemporary Bluegrass",
		"Contemporary Country",
		"Country Gospel",
		"Honky Tonk",
		"Outlaw Country",
		"Traditional Bluegrass",
		"Traditional Country",
		"Urban Cowboy",
		"Breakbeat",
		"Dubstep",
		"Exercise",
		"Garage",
		"Hardcore",
		"Hard Dance",
		"Hi-NRG / Eurodance",
		"House",
		"Jungle/Drum’n’bass",
		"Techno",
		"Trance",
		"Bop",
		"Lounge",
		"Swing",
		"Ambient",
		"Downtempo",
		"Electro",
		"Electronica",
		"Electronic Rock",
		"IDM/Experimental",
		"Industrial",
		"Alternative Rap",
		"Bounce",
		"Dirty South",
		"East Coast Rap",
		"Gangsta Rap",
		"Hardcore Rap",
		"Hip-Hop",
		"Latin Rap",
		"Old School Rap",
		"Rap",
		"Underground Rap",
		"West Coast Rap",
		"Chanukah",
		"Christmas",
		"Christmas: Children’s",
		"Christmas: Classic",
		"Christmas: Classical",
		"Christmas: Jazz",
		"Christmas: Modern",
		"Christmas: Pop",
		"Christmas: R&B",
		"Christmas: Religious",
		"Christmas: Rock",
		"Easter",
		"Halloween",
		"Holiday: Other",
		"Thanksgiving",
		"CCM",
		"Christian Metal",
		"Christian Pop",
		"Christian Rap",
		"Christian Rock",
		"Classic Christian",
		"Contemporary Gospel",
		"Gospel",
		"Christian & Gospel",
		"Praise & Worship",
		"Qawwali (with thx to Jillian Edwards)",
		"Southern Gospel",
		"Traditional Gospel",
		"March (Marching Band)",
	}
	return genres[rand.Intn(len(genres))]
}


func random_name(x int) string{
	song_names := []string{
		"Gimme",
		"Love",
		"All you need is",
		"The",
		"Cheeseits",
		"Ipsum",
		"Pop!",
		"ChaCha",
		"Mocha",
		"Power",
		"Hate",
		"Lyra",
		"Rude",
		"Am I",
		"Wrong",
		"Stay With",
		"Me",
		"Problem",
		"Latch",
		"Summer",
		"Fancy",
	}
	var name string
	for i := 0; i < x; i++ {
		name += song_names[rand.Intn(len(song_names))]
		name += " "
	}
	return name
}

func main() {
	var data = new(Library)
	data.User = "Bob"
	for i := 0; i < 1000; i++ {
		var song = new(Song)
		song.Name = random_name(rand.Intn(5) + 1)
		song.Artist = random_name(rand.Intn(5) + 1)
		song.Length = rand.Int() % 720
		song.Genre = random_genre()
		data.Songs = append(data.Songs, *song)
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
