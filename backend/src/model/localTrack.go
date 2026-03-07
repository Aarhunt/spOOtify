package model

import (
	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"
)

type LocalTrack struct {
	*gorm.Model
	Length    int        `json:"length"`
	Start     int        `json:"start"`
	End       int        `json:"end"`
	Name      string     `json:"name"`
	SpotifyId spotify.ID `json:"spotify_id"`
}

type LocalTrackCreateRequest struct {

}

type LocalTrackResponse struct {

}
