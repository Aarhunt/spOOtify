package model

import (
	"github.com/zmb3/spotify/v2"
)

type SearchResponse struct {
	Icon      []spotify.Image
	Name      string `json:"name"`
	SpotifyID spotify.ID
	ItemType 	ItemType
}
