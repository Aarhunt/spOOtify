package model

import (
	"github.com/aarhunt/spootify/src/utils"
	"github.com/zmb3/spotify/v2"
)

// Image identifies an image associated with an item.
type Image struct {
	// The image height, in pixels.
	Height int `json:"height"`
	// The image width, in pixels.
	Width int `json:"width"`
	// The source URL of the image.
	URL string `json:"url"`
}

func ToImage(i spotify.Image) Image {
	return Image{
		Height: int(i.Height),
		Width:  int(i.Width),
		URL:    i.URL,
	}
}

func ToImages(is []spotify.Image) []Image {
	return utils.Map(is, ToImage)
}
