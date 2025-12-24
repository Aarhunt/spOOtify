package model

import (
	"github.com/zmb3/spotify/v2"
)

type Artist struct {
	ID string `json:"id"`
	SpotifyID	spotify.ID	`json:"spotid"`
	Name	string	`json:"name"`
	Albums	[]Album `json:"albums"`
}

func (a Artist) getSpotifyID() spotify.ID {
	return a.SpotifyID;
}


