package main

import (
	"reflect"
	"testing"
)

func TestSendTenRandomSongs(t *testing.T) {
	data := prepare_data()
	server_data := send_data(data)
	if !(reflect.DeepEqual(server_data.Data, data.Data)) {
		t.Error("The server did not receive all of the songs.")
	}
}

func TestSendOneSong(t *testing.T) {
	data := prepare_one_song()
	server_data := send_data(data)
	if !(reflect.DeepEqual(server_data.Data, data.Data)) {
		t.Error("The server did not receive the only song that was sent.")
	}
}

func TestSongWithNullEverything(t *testing.T) {
	data := new(Library)
	data.User = 1
	var song = new(Song)
	data.Data = append(data.Data, *song)
	server_data := send_data(*data)
	// new(Library).Data is just an empty slice, []
	if !(reflect.DeepEqual(server_data.Data, new(Library).Data)) {
		t.Error("Failed when sending song with null information.")
	}
}

func TestSongWithNullName(t *testing.T) {
	data := new(Library)
	data.User = 1
	var song = new(Song)
	song.Artist = "Null Name Test Artist"
	song.Length = 6900
	song.Genre = "Null Name Test Genre"
	data.Data = append(data.Data, *song)
	server_data := send_data(*data)
	// new(Library).Data is just an empty slice, []
	if !(reflect.DeepEqual(server_data.Data, new(Library).Data)) {
		t.Error("Failed when sending song with null name.")
	}
}

func TestSongWithNullArtist(t *testing.T) {
	data := new(Library)
	data.User = 1
	var song = new(Song)
	song.Name = "Null Artist Test Song"
	song.Length = 6900
	song.Genre = "Null Artist Test Genre"
	data.Data = append(data.Data, *song)
	server_data := send_data(*data)
	if !(reflect.DeepEqual(server_data.Data, data.Data)) {
		t.Error("Failed when sending song with null artist.")
	}
}

func TestSongWithNullLength(t *testing.T) {
	data := new(Library)
	data.User = 1
	var song = new(Song)
	song.Name = "Null Length Test Song"
	song.Artist = "Null Length Test Artist"
	song.Genre = "Null Length Test Genre"
	data.Data = append(data.Data, *song)
	server_data := send_data(*data)
	if !(reflect.DeepEqual(server_data.Data, data.Data)) {
		t.Error("Failed when sending song with null length.")
	}
}

func TestSongWithNullGenre(t *testing.T) {
	data := new(Library)
	data.User = 1
	var song = new(Song)
	song.Name = "Null Genre Test Song"
	song.Artist = "Null Genre Test Artist"
	song.Length = 6900
	data.Data = append(data.Data, *song)
	server_data := send_data(*data)
	if !(reflect.DeepEqual(server_data.Data, data.Data)) {
		t.Error("Failed when sending song with null genre.")
	}
}

func TestSendSameSongTwice(t *testing.T) {
	data := new(Library)
	data.User = 1
	correct_data := new(Library)
	correct_data.User = 1
	for i := 0; i < 2; i++ {
		var song = new(Song)
		song.Name = "Duplicate Song Test"
		song.Artist = "Duplicate Artist"
		song.Length = 6900
		song.Genre = "Duplicate Genre"
		data.Data = append(data.Data, *song)
		if i == 0 {
			correct_data.Data = append(correct_data.Data, *song)
		}
	}
	server_data := send_data(*data)
	if !(reflect.DeepEqual(server_data.Data, correct_data.Data)) {
		t.Error("Failed when sending the same song twice.")
	}
}

func TestTwoSongsSameName(t *testing.T) {
	data := new(Library)
	data.User = 1
	for i := 0; i < 2; i++ {
		var song = new(Song)
		song.Name = "Same Name Test"
		if i == 0 {
			song.Artist = "Same Name Test Artist"
			song.Length = 5000
			song.Genre = "Blues"
		} else {
			song.Artist = "Same Name Test Artist 2"
			song.Length = 5500
			song.Genre = "Pop"
		}
		data.Data = append(data.Data, *song)
	}
	server_data := send_data(*data)
	if !(reflect.DeepEqual(server_data.Data, data.Data)) {
		t.Error("Failed when sending two songs with the same title.")
	}
}

func TestTwowSongsSameArtist(t *testing.T) {
	data := new(Library)
	data.User = 1
	for i := 0; i < 2; i++ {
		var song = new(Song)
		song.Artist = "Same Artist Test"
		if i == 0 {
			song.Name = "Same Artist Test Song"
			song.Length = 4750
			song.Genre = "Hip Hop"
		} else {
			song.Name = "Same Artist Test Song 2"
			song.Length = 5200
			song.Genre = "Rock"
		}
		data.Data = append(data.Data, *song)
	}
	server_data := send_data(*data)
	if !(reflect.DeepEqual(server_data.Data, data.Data)) {
		t.Error("Failed when sending two songs with the same artist.")
	}
}

func TestTwoSongsSameLength(t *testing.T) {
	data := new(Library)
	data.User = 1
	for i := 0; i < 2; i++ {
		var song = new(Song)
		song.Length = 3900
		if i == 0 {
			song.Name = "Same Length Test Song"
			song.Artist = "Same Length Test Artist"
			song.Genre = "Metal"
		} else {
			song.Name = "Same Length Test Song 2"
			song.Artist = "Same Length Test Artist 2"
			song.Genre = "Rap"
		}
		data.Data = append(data.Data, *song)
	}
	server_data := send_data(*data)
	if !(reflect.DeepEqual(server_data.Data, data.Data)) {
		t.Error("Failed when sending two songs with the same length.")
	}
}

func TestTwoSongsSameGenre(t *testing.T) {
	data := new(Library)
	data.User = 1
	for i := 0; i < 2; i++ {
		var song = new(Song)
		song.Genre = "Electronic"
		if i == 0 {
			song.Name = "Same Genre Test Song"
			song.Artist = "Same Genre Test Artist"
			song.Length = 3000
		} else {
			song.Name = "Same Genre Test Song 2"
			song.Artist = "Same Genre Test Artist 2"
			song.Length = 5300
		}
		data.Data = append(data.Data, *song)
	}
	server_data := send_data(*data)
	if !(reflect.DeepEqual(server_data.Data, data.Data)) {
		t.Error("Failed when sending two songs with the same genre.")
	}
}
