package services

import (
	"log"
	"slices"

	"github.com/aarhunt/spootify/src"
	"github.com/aarhunt/spootify/src/model"
	"github.com/aarhunt/spootify/src/utils"
	"github.com/zmb3/spotify/v2"
)

// Get the tracks given by the ids
func getTracksByIds(ids []string) []model.Track {
	spotiConn := src.GetSpotifyConn()
	ctx, client := spotiConn.Ctx, spotiConn.Client

	spotifyIDs := utils.Map(ids, model.ToSpotifyID) 

	chunks := slices.Chunk(spotifyIDs, 50)
	tracks := []*spotify.FullTrack{}

	for chunk := range chunks {
		res, err := client.GetTracks(ctx, chunk)
		if err != nil {
			log.Fatal(err)
		}
		tracks = append(tracks, res...)
	}

	return model.ToTracks(tracks)
}
