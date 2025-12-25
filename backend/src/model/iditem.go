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
	Artist: 		"artist",
	Album:			"album",
	Track:			"track",
}

type IdItem struct {
	gorm.Model
	SpotifyID	spotify.ID	
	ItemType 	ItemType
	PlaylistID 	uint
}

type ItemRequest struct {
	ItemSpotifyID	spotify.ID `json:"spotid" binding:"required"`
	ItemType 	ItemType `json:"type" binding:"required"`
	PlaylistID 	uint `json:"playlistid" binding:"required"`
}

type ItemPlaylistRequest struct {
	ParentSpotifyID spotify.ID	`json:"pspotid" binding:"required"`
	ChildSpotifyID spotify.ID`json:"cspotid" binding:"required"`
}
