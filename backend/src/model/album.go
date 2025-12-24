package model

import (
	"context"
	"github.com/zmb3/spotify/v2"
	"log"
)

func getTracksFromAlbum(ctx context.Context, client spotify.Client, idItem IdItem) []IdItem{
	if idItem.ItemType == Album {
		log.Fatal("ID Submitted is not an album")	
	}

	results, err := client.GetAlbumTracks(ctx, idItem.SpotifyID)
	var result []IdItem = []IdItem{}

	if err != nil {
		log.Fatal(err)
	}

	// handle album results
	if results != nil {
		for _, item := range results.Tracks {
			result = append(result, IdItem{
				SpotifyID: item.ID,
				ItemType: Track,
			})	
		}
	}

	return result
}

