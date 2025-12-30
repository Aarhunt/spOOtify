package services

import (
	"log"

	"github.com/aarhunt/autistify/src"
	"github.com/aarhunt/autistify/src/model"
	"github.com/zmb3/spotify/v2"
)

func getArtistsByIds(ids []spotify.ID) []*spotify.FullArtist {
	spotiConn := src.GetSpotifyConn()
	ctx, client := spotiConn.Ctx, spotiConn.Client

	artists, _ := client.GetArtists(ctx, ids...)
	return artists
}

func GetAlbumsFromArtistById(id spotify.ID) []model.IdItem{
	spotiConn := src.GetSpotifyConn()
	ctx, client := spotiConn.Ctx, spotiConn.Client

	albums, err := client.GetArtistAlbums(ctx, id, []spotify.AlbumType{spotify.AlbumTypeAlbum, spotify.AlbumTypeSingle})
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

func getTracksFromArtistById(id spotify.ID) []model.IdItem{
	albums := GetAlbumsFromArtistById(id)
	var result []model.IdItem = []model.IdItem{}

	// handle album results
	if albums != nil {
		for _, album := range albums {
			tracks := getTracksFromAlbumById(album.SpotifyID)
			result = append(result, tracks...)
		}
	}

	return result
}
