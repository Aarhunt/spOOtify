package model

import (
	"github.com/zmb3/spotify/v2"
)

type ItemType int
type InclusionType int

const (
	PlaylistItem ItemType = iota
	Artist 
	Album
	Track
)

const (
	Nothing InclusionType = iota
	Included
	Excluded
)

var itemType = map[ItemType]string{
	PlaylistItem: 	"playlist",
	Artist: 		"artist",
	Album:			"album",
	Track:			"track",
}

var inclusionType = map[InclusionType]string{
	Nothing: 	"nothing",
	Included: 	"included",
	Excluded: 	"excluded",
}

type IdItem struct {
	SpotifyID   spotify.ID `gorm:"primaryKey;type:varchar(255);not null" json:"id" example:"37i9dQZF1DXcBWIGoYBM3M"`
	ItemType 	ItemType
	Playlists 	[]Playlist `gorm:"many2many:playlist_inclusions;" json:"included_in"`
}

type ItemInclusionRequest struct {
	ItemSpotifyID	spotify.ID `json:"spotid" binding:"required"`
	ItemType 	ItemType `json:"type" binding:"required"`
	PlaylistID 	spotify.ID `json:"playlistid" binding:"required"`
	Include 	bool `json:"include" binding:"required"`
}

type ItemRequest struct {
	ArtistID 	spotify.ID	`json:"artistid" binding:"required"`
	PlaylistID 	spotify.ID `json:"playlistid" binding:"required"`
	ItemType 	ItemType `json:"type" binding:"required"`
}

type SearchRequest struct {
	Query 		string `json:"query" binding:"required"`
	PlaylistID 	spotify.ID `json:"playlistid" binding:"required"`
	ItemType 	ItemType `json:"type" binding:"required"`
}

type ItemPlaylistRequest struct {
	ParentSpotifyID spotify.ID	`json:"pspotid" binding:"required"`
	ChildSpotifyID spotify.ID`json:"cspotid" binding:"required"`
}

type ItemResponse struct {
	SpotifyID spotify.ID
	Name      string `json:"name"`
	Icon      []spotify.Image
	ItemType 	ItemType
	Included 	InclusionType
}
