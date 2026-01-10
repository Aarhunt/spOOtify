package services

import (
	"log"
	"slices"

	"github.com/aarhunt/spootify/src"
	"github.com/zmb3/spotify/v2"
)

func getArtistsByIds(ids []spotify.ID) []*spotify.FullArtist {
	spotiConn := src.GetSpotifyConn()
	ctx, client := spotiConn.Ctx, spotiConn.Client

	chunks := slices.Chunk(ids, 50)
	artists := []*spotify.FullArtist{}

	for chunk := range chunks {
		res, err := client.GetArtists(ctx, chunk...)
		if err != nil {
			log.Fatal(err)
		}
		artists = append(artists, res...)	
	}
	return artists
}

func GetAlbumsFromArtistById(id spotify.ID) []spotify.SimpleAlbum{
	spotiConn := src.GetSpotifyConn()
	ctx, client := spotiConn.Ctx, spotiConn.Client

	albums, err := client.GetArtistAlbums(ctx, id, []spotify.AlbumType{spotify.AlbumTypeAlbum}, spotify.Limit(50))

	if err != nil {
		log.Fatal(err)
	}

	return albums.Albums
}

func GetAlbumsAndSinglesFromArtistById(id spotify.ID) []spotify.SimpleAlbum{
	spotiConn := src.GetSpotifyConn()
	ctx, client := spotiConn.Ctx, spotiConn.Client

	albums, err := client.GetArtistAlbums(ctx, id, []spotify.AlbumType{spotify.AlbumTypeAlbum, spotify.AlbumTypeSingle}, spotify.Limit(50))

	if err != nil {
		log.Fatal(err)
	}

	return albums.Albums
}

func getTracksFromArtistById(id spotify.ID) []spotify.SimpleTrack{
	albums := GetAlbumsFromArtistById(id)
	var result []spotify.SimpleTrack = []spotify.SimpleTrack{}

	// handle album results
	if albums != nil {
		for _, album := range albums {
			tracks := getTracksFromAlbumById(album.ID)
			result = append(result, tracks...)
		}
	}

	return result
}
