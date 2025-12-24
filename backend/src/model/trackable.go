package model

import (
	"github.com/zmb3/spotify/v2"
)

type Trackable interface {
	getSpotifyID()		spotify.ID
}
