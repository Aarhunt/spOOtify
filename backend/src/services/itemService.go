package services

import (
	"log"
	"slices"

	"github.com/aarhunt/autistify/src"
	"github.com/aarhunt/autistify/src/model"
	"github.com/aarhunt/autistify/src/utils"
	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func IncludeExcludeItem(req model.ItemInclusionRequest) (*model.PlaylistResponse, error) {
	dbConn := src.GetDbConn()
	db := dbConn.Db

	playlist, err := getPlaylist(req.PlaylistID)

	newItem := model.IdItem{
		SpotifyID: req.ItemSpotifyID,
		ItemType:  req.ItemType,
		Playlists: []model.Playlist{},
	}

	db.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Omit(clause.Associations).Create(&newItem).Error
		if err != nil {
			return err
		}

		if req.Include {
			err = tx.Model(&playlist).Association("Inclusions").Append(&newItem)
		} else {
			err = tx.Model(&playlist).Association("Exclusions").Append(&newItem)
		}
		if err != nil {
			return err
		}

		return nil
	})

	return playlist.ToResponse(), err
}

func IncludePlaylist(req model.ItemPlaylistRequest) (*model.PlaylistResponse, error) {
	dbConn := src.GetDbConn()
	ctx, db := dbConn.Ctx, dbConn.Db

	parentPlaylist, err := gorm.G[model.Playlist](db).Where("spotify_id = ?", req.ParentSpotifyID).First(ctx)
	childPlaylist, err := gorm.G[model.Playlist](db).Where("spotify_id = ?", req.ChildSpotifyID).First(ctx)

	err = db.Model(&parentPlaylist).Association("IncludedPlaylists").Append(&childPlaylist)
	return parentPlaylist.ToResponse(), err
}

func GetAlbumFromArtist(req model.ItemRequest) ([]model.ItemResponse, error) {
	conn := src.GetSpotifyConn()
	ctx, client := conn.Ctx, conn.Client

	playlist, err := getPlaylist(req.PlaylistID)
	results, err := client.GetArtistAlbums(ctx, req.ArtistID, []spotify.AlbumType{spotify.AlbumTypeAlbum, spotify.AlbumTypeSingle})

	var albums = results.Albums

	return albumToResponse(albums, playlist), err
}

func GetTracksFromAlbum(req model.ItemRequest) ([]model.ItemResponse, error) {
	conn := src.GetSpotifyConn()
	ctx, client := conn.Ctx, conn.Client

	playlist, err := getPlaylist(req.PlaylistID)
	results, err := client.GetAlbumTracks(ctx, req.ArtistID)

	var tracks = results.Tracks

	return trackToResponse(tracks, playlist), err
}

func trackToResponse(tracks []spotify.SimpleTrack, playlist *model.Playlist) []model.ItemResponse {
	trackIDs := utils.Map(tracks, func(a spotify.SimpleTrack) spotify.ID {
		return a.ID
	})

	includedTracks, excludedTracks := getInclusionsExclusions(playlist, trackIDs)

	return utils.Map(tracks, func(a spotify.SimpleTrack) model.ItemResponse {
		included := model.Nothing
		if slices.Contains(includedTracks, a.ID) {
			included = model.Included
		} else if slices.Contains(excludedTracks, a.ID) {
			included = model.Excluded
		}
		return model.ItemResponse{
			SpotifyID: a.ID,
			Name:      a.Name,
			Icon:      []spotify.Image{},
			ItemType:  model.Album,
			Included:  included,
		}
	})
}

func albumToResponse(albums []spotify.SimpleAlbum, playlist *model.Playlist) []model.ItemResponse {
	albumIDs := utils.Map(albums, func(a spotify.SimpleAlbum) spotify.ID {
		return a.ID
	})

	includedAlbums, excludedAlbums := getInclusionsExclusions(playlist, albumIDs)

	return utils.Map(albums, func(a spotify.SimpleAlbum) model.ItemResponse {
		included := model.Nothing
		if slices.Contains(includedAlbums, a.ID) {
			included = model.Included
		} else if slices.Contains(excludedAlbums, a.ID) {
			included = model.Excluded
		}
		return model.ItemResponse{
			SpotifyID: a.ID,
			Name:      a.Name,
			Icon:      a.Images,
			ItemType:  model.Album,
			Included:  included,
		}
	})
}

func artistToResponse(artists []spotify.FullArtist, playlist *model.Playlist) []model.ItemResponse {
	artistIDs := utils.Map(artists, func(a spotify.FullArtist) spotify.ID {
		return a.ID
	})

	includedArtists, excludedArtists := getInclusionsExclusions(playlist, artistIDs)

	return utils.Map(artists, func(a spotify.FullArtist) model.ItemResponse {
		included := model.Nothing
		if slices.Contains(includedArtists, a.ID) {
			included = model.Included
		} else if slices.Contains(excludedArtists, a.ID) {
			included = model.Excluded
		}
		return model.ItemResponse{
			SpotifyID: a.ID,
			Name:      a.Name,
			Icon:      a.Images,
			ItemType:  model.Album,
			Included:  included,
		}
	})
}

func getInclusionsExclusions(playlist *model.Playlist, itemIDs []spotify.ID) ([]spotify.ID, []spotify.ID) {
	includedItems := GetIncludedItemsFromPlaylist(playlist, itemIDs)
	excludedItems := GetExcludedItemsFromPlaylist(playlist, itemIDs)
	includedPlaylists := GetIncludedPlaylistsFromPlaylist(playlist)
	for len(includedPlaylists) > 0 {
		var newIncludedPlaylists = []model.Playlist{}
		for _, item := range includedPlaylists {
			includedItems = append(includedItems, GetIncludedItemsFromPlaylist(&item, itemIDs)...)
			excludedItems = append(excludedItems, GetExcludedItemsFromPlaylist(&item, itemIDs)...)
			newIncludedPlaylists = append(newIncludedPlaylists, GetIncludedPlaylistsFromPlaylist(&item)...)
		}
		includedPlaylists = newIncludedPlaylists
	}
	return includedItems, excludedItems
}

func SearchArtist(req model.SearchRequest) []model.ItemResponse {
	conn := src.GetSpotifyConn()
	ctx, client := conn.Ctx, conn.Client

	playlist, err := getPlaylist(req.PlaylistID)
	results, err := client.Search(ctx, req.Query, spotify.SearchTypeArtist)

	// handle album results
	if err != nil {
		log.Fatal("help")
	}
	return artistToResponse(results.Artists.Artists, playlist)
}

func SearchAlbum(req model.SearchRequest) []model.ItemResponse {
	conn := src.GetSpotifyConn()
	ctx, client := conn.Ctx, conn.Client

	playlist, err := getPlaylist(req.PlaylistID)
	results, err := client.Search(ctx, req.Query, spotify.SearchTypeAlbum)

	// handle album results
	if err != nil {
		log.Fatal("help")
	}
	return albumToResponse(results.Albums.Albums, playlist)
}

func SearchTrack(req model.SearchRequest) []model.ItemResponse {
	conn := src.GetSpotifyConn()
	ctx, client := conn.Ctx, conn.Client

	playlist, err := getPlaylist(req.PlaylistID)
	results, err := client.Search(ctx, req.Query, spotify.SearchTypeAlbum)

	// handle album results
	if err != nil {
		log.Fatal("help")
	}
	return albumToResponse(results.Albums.Albums, playlist)
}

func getFullPlaylist(id spotify.ID) (model.Playlist, error) { //TODO
	var p model.Playlist
	db := src.GetDbConn().Db

	err := db.Preload("Inclusions").
		Preload("Exclusions").
		Preload("IncludedPlaylists").
		First(&p, "spotify_id = ?", id).Error

	return p, err
}
