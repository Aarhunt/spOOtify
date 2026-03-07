package model

import (
	"github.com/aarhunt/spootify/src/utils"
	"github.com/zmb3/spotify/v2"
)

type Album struct {
	// The name of the album.
	Name string `json:"name"`
	// The type of the album: one of "album",
	// "single", or "compilation".
	AlbumType string `json:"album_type"`
	// The [Spotify ID] for the album.
	ID string `json:"id"`
	// The cover art for the album in various sizes,
	// widest first.
	Images []Image `json:"images"`
	// The date the album was first released.  For example, "1981-12-15".
	// Depending on the ReleaseDatePrecision, it might be shown as
	// "1981" or "1981-12". You can use [SimpleAlbum.ReleaseDateTime] to convert
	// this to a [time.Time] value.
	ReleaseYear int `json:"release_date"`
	// The precision with which ReleaseDate value is known: "year", "month", or "day"
	ReleaseDatePrecision string `json:"release_date_precision"`
	// The number of tracks on the album.
	TotalTracks int `json:"total_tracks"`

	ArtistID string `json:"artistid"`
}

func ToAlbum[A *spotify.FullAlbum | spotify.SimpleAlbum](a A) Album {
	var id, name, albumType, releaseDatePrecision, artistID string
	var images []Image
	var releaseYear, totalTracks int



	switch album := any(a).(type) {
		case *spotify.FullAlbum:
		id, name, albumType, releaseDatePrecision, images, releaseYear, totalTracks = album.ID.String(), album.Name, album.AlbumType, album.ReleaseDatePrecision, ToImages(album.Images), album.ReleaseDateTime().Year(), int(album.TotalTracks)
		if album.Artists != nil {
			artistID = album.Artists[0].ID.String()
		} 
		case spotify.SimpleAlbum:
		id, name, albumType, releaseDatePrecision, images, releaseYear, totalTracks = album.ID.String(), album.Name, album.AlbumType, album.ReleaseDatePrecision, ToImages(album.Images), album.ReleaseDateTime().Year(), int(album.TotalTracks)
		if album.Artists != nil {
			artistID = album.Artists[0].ID.String()
		} 
	}

	return Album {
		Name:                 name,
		AlbumType:            albumType,
		ID:                   id,
		Images:               images,
		ReleaseYear:          releaseYear,
		ReleaseDatePrecision: releaseDatePrecision,
		TotalTracks:          totalTracks,
		ArtistID: 			  artistID,			
	}
}

func ToAlbums[A *spotify.FullAlbum | spotify.SimpleAlbum](as []A) []Album {
	return utils.Map(as, ToAlbum[A])
}

