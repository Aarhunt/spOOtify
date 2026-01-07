package services

import (
	"log"

	"github.com/aarhunt/spootify/src"
	"github.com/aarhunt/spootify/src/model"
	"github.com/zmb3/spotify/v2"
)

func getAlbumsByIds(ids []spotify.ID) []*spotify.FullAlbum {
	spotiConn := src.GetSpotifyConn()
	ctx, client := spotiConn.Ctx, spotiConn.Client

	albums, _ := client.GetAlbums(ctx, ids)
	return albums
}


func getTracksFromAlbumById(id spotify.ID) []model.IdItem{
	spotiConn := src.GetSpotifyConn()
	ctx, client := spotiConn.Ctx, spotiConn.Client

	results, err := client.GetAlbumTracks(ctx, id)
	var result []model.IdItem = []model.IdItem{}

	if err != nil {
		log.Fatal(err)
	}

	// handle album results
	if results != nil {
		for _, item := range results.Tracks {
			result = append(result, model.IdItem{
				SpotifyID: item.ID,
				ItemType: model.Track,
			})	
		}
	}

	return result
}

