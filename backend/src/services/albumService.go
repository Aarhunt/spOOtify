package services

import (
	"log"
	"slices"

	"github.com/aarhunt/spootify/src"
	"github.com/zmb3/spotify/v2"
)

func getAlbumsByIds(ids []spotify.ID) []*spotify.FullAlbum {
	spotiConn := src.GetSpotifyConn()
	ctx, client := spotiConn.Ctx, spotiConn.Client

	chunks := slices.Chunk(ids, 20)
	albums := []*spotify.FullAlbum{}

	for chunk := range chunks {
		res, err := client.GetAlbums(ctx, chunk)
		if err != nil {
			log.Fatal(err)
		}
		albums = append(albums, res...)	
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

