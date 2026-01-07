package services

import (
	"github.com/aarhunt/spootify/src"
	"github.com/zmb3/spotify/v2"
)

func getTracks(ids []spotify.ID) []*spotify.FullTrack {
	spotiConn := src.GetSpotifyConn()
	ctx, client := spotiConn.Ctx, spotiConn.Client

	tracks, _ := client.GetTracks(ctx, ids)
	return tracks
}
