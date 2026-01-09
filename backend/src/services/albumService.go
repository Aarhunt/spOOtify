package services

import (
	"log"

	"github.com/aarhunt/spootify/src"
	"github.com/zmb3/spotify/v2"
)

func getAlbumsByIds(ids []spotify.ID) []*spotify.FullAlbum {
	spotiConn := src.GetSpotifyConn()
	ctx, client := spotiConn.Ctx, spotiConn.Client

	albums, err := client.GetAlbums(ctx, ids)

	if err != nil {
		log.Print(err)
	}

	return albums
}


func getTracksFromAlbumById(id spotify.ID) []spotify.SimpleTrack{
	spotiConn := src.GetSpotifyConn()
	ctx, client := spotiConn.Ctx, spotiConn.Client

	results, err := client.GetAlbumTracks(ctx, id, spotify.Limit(50))

	if err != nil {
		log.Fatal(err)
	}

	return results.Tracks
}

