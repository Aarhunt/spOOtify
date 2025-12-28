package services

import (
	"log"

	"github.com/aarhunt/autistify/src"
	"github.com/aarhunt/autistify/src/model"
	"github.com/aarhunt/autistify/src/utils"
	"github.com/zmb3/spotify/v2"
)

func getArtists(items []model.IdItem) []*spotify.FullArtist {
	spotiConn := src.GetSpotifyConn()
	ctx, client := spotiConn.Ctx, spotiConn.Client

	artists, _ := client.GetArtists(ctx, utils.Map(items, func(i model.IdItem) spotify.ID { return i.SpotifyID })...)
	return artists
}

func GetAlbumsFromArtist(idItem model.IdItem) []model.IdItem{
	spotiConn := src.GetSpotifyConn()
	ctx, client := spotiConn.Ctx, spotiConn.Client

	albums, err := client.GetArtistAlbums(ctx, idItem.SpotifyID, []spotify.AlbumType{spotify.AlbumTypeAlbum, spotify.AlbumTypeSingle})
	var result []model.IdItem = []model.IdItem{}

	if err != nil {
		log.Fatal(err)
	}

	if albums != nil {
		for _, album := range albums.Albums {
			result = append(result, model.IdItem{SpotifyID: album.ID, ItemType: model.Album})
		}
	}

	return result
}

func getTracksFromArtist(idItem model.IdItem) []model.IdItem{
	albums := GetAlbumsFromArtist(idItem)
	var result []model.IdItem = []model.IdItem{}

	// handle album results
	if albums != nil {
		for _, album := range albums {
			tracks := getTracksFromAlbum(album)
			result = append(result, tracks...)
		}
	}

	return result
}
