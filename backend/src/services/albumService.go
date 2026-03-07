package services

import (
	"log"
	"slices"

	"github.com/aarhunt/spootify/src"
	"github.com/aarhunt/spootify/src/model"
	"github.com/aarhunt/spootify/src/utils"
	"github.com/zmb3/spotify/v2"
)

// Gets a list of albums by their IDs.
func getAlbumsByIds(ids []string) []model.Album {
	spotiConn := src.GetSpotifyConn()
	ctx, client := spotiConn.Ctx, spotiConn.Client

	spotifyIDs := utils.Map(ids, model.ToSpotifyID) 

	chunks := slices.Chunk(spotifyIDs, 20)
	albums := []model.Album{}

	for chunk := range chunks {
		res, err := client.GetAlbums(ctx, chunk)
		if err != nil {
			log.Fatal(err)
		}
		albums = append(albums, model.ToAlbums(res)...)
	}
	
	return albums
}

// Get all tracks from an album given its id.
func GetTracksFromAlbumById(id string) []model.Track {
	spotiConn := src.GetSpotifyConn()
	ctx, client := spotiConn.Ctx, spotiConn.Client

	results, err := client.GetAlbumTracks(ctx, model.ToSpotifyID(id), spotify.Limit(50))

	if err != nil {
		log.Fatal(err)
	}

	return model.ToTracks(results.Tracks)
}

