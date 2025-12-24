package model

import (
	"github.com/zmb3/spotify/v2"
)

type Album struct {
	ID string `json:"id"`
	SpotifyID	spotify.ID	`json:"spotid"`
	Name	string	`json:"name"`
	Tracks	[]Track	`json:"songs"`
}

func (a Album) getSpotifyID() spotify.ID {
	return a.SpotifyID;
}

func (a Album) getTracks() []Track {
		
}
