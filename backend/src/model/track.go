package model

import (
	"github.com/zmb3/spotify/v2"
)

type Track struct {
	ID string `json:"id"`
	SpotifyID	spotify.ID	`json:"spotid"`
	Title	string	`json:"title"`
	Artist	string	`json:"artist"`
}

func (t Track) getSpotifyID() spotify.ID {
	return t.SpotifyID;
}


