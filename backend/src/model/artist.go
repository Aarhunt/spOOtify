package model

import (
	"github.com/aarhunt/spootify/src/utils"
	"github.com/zmb3/spotify/v2"
)

type Artist struct {
	Name string `json:"name"`
	ID   string     `json:"id"`
	// Images of the artist in various sizes, widest first.
	Images []Image `json:"images"`
}

func ToArtist[A *spotify.FullArtist | spotify.FullArtist](a A) Artist {
	var name, id string
	var images []Image

	switch artist := any(a).(type) {
	case *spotify.FullArtist:
		id, name, images = artist.ID.String(), artist.Name, ToImages(artist.Images)
	case spotify.FullArtist:
		id, name, images = artist.ID.String(), artist.Name, ToImages(artist.Images)
	}

	return Artist{
		Name:   name,
		ID:     id,
		Images: images,
	}
}

func ToArtists[A *spotify.FullArtist | spotify.FullArtist](as []A) []Artist {
	return utils.Map(as, ToArtist)
}
