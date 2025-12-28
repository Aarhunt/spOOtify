package services

import (
	"github.com/aarhunt/autistify/src"
	"github.com/aarhunt/autistify/src/model"
	"github.com/aarhunt/autistify/src/utils"
	"github.com/zmb3/spotify/v2"
)

func getTracks(items []model.IdItem) []*spotify.FullTrack {
	spotiConn := src.GetSpotifyConn()
	ctx, client := spotiConn.Ctx, spotiConn.Client

	tracks, _ := client.GetTracks(ctx, utils.Map(items, func(i model.IdItem) spotify.ID { return i.SpotifyID }))
	return tracks
}
