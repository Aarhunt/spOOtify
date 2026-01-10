package services

import (
	"log"
	"slices"

	"github.com/aarhunt/spootify/src"
	"github.com/zmb3/spotify/v2"
)

func getTracks(ids []spotify.ID) []*spotify.FullTrack {
	spotiConn := src.GetSpotifyConn()
	ctx, client := spotiConn.Ctx, spotiConn.Client

	chunks := slices.Chunk(ids, 50)
	tracks := []*spotify.FullTrack{}

	for chunk := range chunks {
		res, err := client.GetTracks(ctx, chunk)
		if err != nil {
			log.Fatal(err)
		}
		tracks = append(tracks, res...)
	}

	return tracks
}
