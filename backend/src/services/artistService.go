package services

import (
	"log"
	"slices"

	"github.com/aarhunt/spootify/src"
	"github.com/aarhunt/spootify/src/model"
	"github.com/aarhunt/spootify/src/utils"
	"github.com/zmb3/spotify/v2"
)

// Get a list of artists by their IDs.
func getArtistsByIds(ids []string) []model.Artist {
	spotiConn := src.GetSpotifyConn()
	ctx, client := spotiConn.Ctx, spotiConn.Client

	spotifyIDs := utils.Map(ids, model.ToSpotifyID)

	chunks := slices.Chunk(spotifyIDs, 50)
	artists := []model.Artist{}

	for chunk := range chunks {
		res, err := client.GetArtists(ctx, chunk...)
		if err != nil {
			log.Fatal(err)
		}
		artists = append(artists, model.ToArtists(res)...)	
	}
	return artists 
}

// Get all albums from an artist given its id.
func GetAlbumsFromArtistById(id string, includeSingles bool) []model.Album{
	spotiConn := src.GetSpotifyConn()
	ctx, client := spotiConn.Ctx, spotiConn.Client

	types := []spotify.AlbumType{spotify.AlbumTypeAlbum}
	if includeSingles {types = append(types, spotify.AlbumTypeSingle)}

	albums, err := client.GetArtistAlbums(ctx, model.ToSpotifyID(id), types, spotify.Limit(50))

	if err != nil {
		log.Fatal(err)
	}

	return model.ToAlbums(albums.Albums)
}

func GetGenresFromArtist(id string) ([]string, error) {
    spotiConn := src.GetSpotifyConn()
    ctx, client := spotiConn.Ctx, spotiConn.Client

    artist, err := client.GetArtist(ctx, model.ToSpotifyID(id))
    if err != nil {
        return nil, err
    }

    return model.ToArtist(artist).Genres, nil
}


// Get all tracks from an artist given its id.
func getTracksFromArtistById(id string) []model.Track{
	albums := GetAlbumsFromArtistById(id, false)
	var result []model.Track = []model.Track{}

	// handle album results
	for _, album := range albums {
		tracks := GetTracksFromAlbumById(album.ID)
		result = append(result, tracks...)
	}

	return result
}
