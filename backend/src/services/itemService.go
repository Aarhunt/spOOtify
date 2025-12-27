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

func IncludeExcludeItem(req model.ItemInclusionRequest) (*model.InclusionResponse, error) {
	dbConn := src.GetDbConn()
	db := dbConn.Db

	playlist, err := getPlaylist(req.PlaylistID)

	newItem := model.IdItem{
		SpotifyID: req.ItemSpotifyID,
		ItemType:  req.ItemType,
		Playlists: []model.Playlist{},
	}

	returnItem := model.InclusionResponse{
		SpotifyID: req.ItemSpotifyID,
		Included: model.InclusionType(0),
	}

	db.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "spotify_id"}},
			UpdateAll: true,
		}).Omit(clause.Associations).Create(&newItem).Error
		if err != nil {
			return err
		}

		if req.Include {
			err = tx.Model(&playlist).Association("Inclusions").Append(&newItem)
			returnItem.Included = 1;
		} else {
			err = tx.Model(&playlist).Association("Exclusions").Append(&newItem)
			returnItem.Included = 2;
		}
		if err != nil {
			return err
		}

		return nil
	})

	return &returnItem, err
}

func UndoIncludeExcludeItem(req model.ItemInclusionRequest) (*model.InclusionResponse, error) {
    dbConn := src.GetDbConn()
    db := dbConn.Db

    playlist, err := getPlaylist(req.PlaylistID)
    if err != nil {
        return nil, err
    }

    item := model.IdItem{SpotifyID: req.ItemSpotifyID}

    returnItem := model.InclusionResponse{
        SpotifyID: req.ItemSpotifyID,
        Included:  model.InclusionType(0), 
    }

    err = db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Model(&playlist).Association("Inclusions").Delete(&item); err != nil {
            return err
        }

        if err := tx.Model(&playlist).Association("Exclusions").Delete(&item); err != nil {
            return err
        }

        return nil
    })

    return &returnItem, err
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
	results, err := client.GetArtistAlbums(ctx, req.ParentID, []spotify.AlbumType{spotify.AlbumTypeAlbum}, spotify.Limit(50))

	var albums = results.Albums

	return albumToResponse(albums, playlist), err
}

func GetTracksFromAlbum(req model.ItemRequest) ([]model.ItemResponse, error) {
	conn := src.GetSpotifyConn()
	ctx, client := conn.Ctx, conn.Client

	playlist, err := getPlaylist(req.PlaylistID)
	results, err := client.GetAlbumTracks(ctx, req.ParentID, spotify.Limit(50))

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
			Icon:      a.Album.Images,
			ItemType:  model.Track,
			Included:  included,
			SortData:  int(a.TrackNumber),
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
			SortData:  a.ReleaseDateTime().Year(),
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
			ItemType:  model.Artist,
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
	results, err := client.Search(ctx, req.Query, spotify.SearchTypeArtist, spotify.Limit(5))

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
	results, err := client.Search(ctx, req.Query, spotify.SearchTypeAlbum, spotify.Limit(5))

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
	results, err := client.Search(ctx, req.Query, spotify.SearchTypeTrack, spotify.Limit(5))

	for _, item := range results.Tracks.Tracks {
		log.Print(item.Album.Images)
	}

	// handle album results
	if err != nil {
		log.Fatal("help")
	}
	return trackToResponse(fullToSimpleTrack(results.Tracks.Tracks), playlist)
}

func fullToSimpleTrack(tracks []spotify.FullTrack) []spotify.SimpleTrack {
	return utils.Map(tracks, func(t spotify.FullTrack) spotify.SimpleTrack {
		return spotify.SimpleTrack {
			Album: t.Album,	
			ID: t.ID,	
			Name: t.Name,
		}
	})
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
