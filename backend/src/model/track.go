package model

import (
	"github.com/aarhunt/spootify/src/utils"
	"github.com/zmb3/spotify/v2"
)

type Track struct {
	ID       string     `json:"id"`
	Name     string `json:"name"`
	// The number of the track.  If an album has several
	// discs, the track number is the number on the specified
	// DiscNumber.
	TrackNumber int `json:"track_number"`
	Album   Album    `json:"album"`
}

func ToTrack[T *spotify.FullTrack | spotify.FullTrack | spotify.SimpleTrack](t T) Track {
	var id, name string
	var trackNum int
	var album spotify.SimpleAlbum

	// Inline type switch to extract common fields
	switch track := any(t).(type) {
	case *spotify.FullTrack:
		id, name, trackNum, album = track.ID.String(), track.Name, int(track.TrackNumber), track.Album
	case spotify.SimpleTrack:
		id, name, trackNum, album = track.ID.String(), track.Name, int(track.TrackNumber), track.Album
	}

	return Track{
		ID:          id,
		Name:        name,
		TrackNumber: trackNum,
		Album:       ToAlbum(album),
	}
}

func ToTracks[T *spotify.FullTrack | spotify.FullTrack | spotify.SimpleTrack](ts []T) []Track {
    return utils.Map(ts, ToTrack[T])
}
