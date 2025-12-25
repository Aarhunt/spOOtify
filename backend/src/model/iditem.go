package model

import (
	"github.com/zmb3/spotify/v2"
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
	SpotifyID   spotify.ID `gorm:"primaryKey;type:varchar(255);not null" json:"id" example:"37i9dQZF1DXcBWIGoYBM3M"`
	ItemType 	ItemType
	Playlists 	[]Playlist
}

type ItemRequest struct {
	ItemSpotifyID	spotify.ID `json:"spotid" binding:"required"`
	ItemType 	ItemType `json:"type" binding:"required"`
	PlaylistID 	spotify.ID `json:"playlistid" binding:"required"`
	Include 	bool `json:"include" binding:"required"`
}

type ItemPlaylistRequest struct {
	ParentSpotifyID spotify.ID	`json:"pspotid" binding:"required"`
	ChildSpotifyID spotify.ID`json:"cspotid" binding:"required"`
}
