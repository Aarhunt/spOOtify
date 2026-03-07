package model

import (
	"github.com/aarhunt/spootify/src/utils"
	"github.com/zmb3/spotify/v2"
)

func ToSpotifyID(s string) spotify.ID {
	return spotify.ID(s);
}

func ToSpotifyIDs(ss []string) []spotify.ID {
	return utils.Map(ss, ToSpotifyID)
}

func toStrings(ss []spotify.ID) []string {
	return utils.Map(ss, func(s spotify.ID) string {return s.String()})
}
