package model

import (
	"context"
	"github.com/zmb3/spotify/v2"
	"log"
)

func getTracksFromArtist(ctx context.Context, client spotify.Client, idItem IdItem) []IdItem{
	albums := getAlbumsFromArtist(ctx, client, idItem)
	var result []IdItem = []IdItem{}

	// handle album results
	if albums != nil {
		for _, album := range albums {
			tracks := getTracksFromAlbum(ctx, client, album)
			result = append(result, tracks...)
		}
	}

	return result
}

func getAlbumsFromArtist(ctx context.Context, client spotify.Client, idItem IdItem) []IdItem{
	albums, err := client.GetArtistAlbums(ctx, idItem.SpotifyID, []spotify.AlbumType{spotify.AlbumTypeAlbum, spotify.AlbumTypeSingle})
	var result []IdItem = []IdItem{}

	if err != nil {
		log.Fatal(err)
	}

	// handle album results
	if albums != nil {
		for _, album := range albums.Albums {
			result = append(result, IdItem{ID: "1", SpotifyID: album.ID, ItemType: Album})
		}
	}

	return result
}
