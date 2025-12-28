package services

import (
	"log"

	"github.com/aarhunt/autistify/src"
	"github.com/aarhunt/autistify/src/model"
	"github.com/zmb3/spotify/v2"
)

func getAlbums(ids []spotify.ID) []*spotify.FullAlbum {
	spotiConn := src.GetSpotifyConn()
	ctx, client := spotiConn.Ctx, spotiConn.Client

	albums, _ := client.GetAlbums(ctx, ids)
	return albums
}


func getTracksFromAlbum(id spotify.ID) []model.IdItem{
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

