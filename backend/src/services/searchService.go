package services

import (
	"fmt"

	"github.com/aarhunt/autistify/src"
	"github.com/aarhunt/autistify/src/model"
	"github.com/zmb3/spotify/v2"
)

func Search(query string) ([]model.ItemResponse, error) {	
	conn := src.GetSpotifyConn()
	ctx, client := conn.Ctx, conn.Client

	results, err := client.Search(ctx, query, spotify.SearchTypePlaylist|spotify.SearchTypeAlbum|spotify.SearchTypeArtist, spotify.Limit(5))

	responses := []model.ItemResponse{}

	// handle album results
	if results.Artists != nil {
		fmt.Println("Artists:")
		for _, item := range results.Artists.Artists {
			responses = append(responses, model.ItemResponse{
				Icon: item.Images,
				Name: item.Name,
				SpotifyID: item.ID,
				ItemType: model.Artist,
			})
		}
	}

	// handle album results
	if results.Albums != nil {
		fmt.Println("Albums:")
		for _, item := range results.Albums.Albums {
			responses = append(responses, model.ItemResponse{
				Icon: item.Images,
				Name: item.Name,
				SpotifyID: item.ID,
				ItemType: model.Artist,
			})
		}
	}

	// handle playlist results
	if results.Playlists != nil {
		fmt.Println("Playlists:")
		for _, item := range results.Playlists.Playlists {
			responses = append(responses, model.ItemResponse{
				Icon: item.Images,
				Name: item.Name,
				SpotifyID: item.ID,
				ItemType: model.PlaylistItem,
			})
		}
	}

	return responses, err
}

