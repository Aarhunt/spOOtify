package model

import (
	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"
)

type ItemType int

const (
	PlaylistItem ItemType = iota
	Artist 
	Album
	Track
)

var itemType = map[ItemType]string{
	PlaylistItem: 	"playlist",
	Artist: 		"connected",
	Album:			"album",
	Track:			"track",
}

type IdItem struct {
	gorm.Model
	SpotifyID	spotify.ID	`json:"spotid"`
	ItemType 	ItemType	`json:"ItemType"`
	PlaylistID 	uint
}
