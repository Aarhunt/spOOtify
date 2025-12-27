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

		if *req.Include {
			tx.Model(&playlist).Association("Exclusions").Delete(&newItem)
			err = tx.Model(&playlist).Association("Inclusions").Append(&newItem)
			returnItem.Included = 1;
		} else {
			tx.Model(&playlist).Association("Inclusions").Delete(&newItem)
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

	return trackToResponse(tracks, playlist, req.ParentID), err
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

func albumToResponse(albums []spotify.SimpleAlbum, playlist *model.Playlist) []model.ItemResponse {
	albumIDs := utils.Map(albums, func(a spotify.SimpleAlbum) spotify.ID { return a.ID })
	artistIDs := utils.Map(albums, func(a spotify.SimpleAlbum) spotify.ID { return a.Artists[0].ID })

	incMap := GetInclusionMap(playlist.SpotifyID, append(albumIDs, artistIDs...))
	excMap := GetExclusionMap(playlist.SpotifyID, append(albumIDs, artistIDs...))

	return utils.Map(albums, func(a spotify.SimpleAlbum) model.ItemResponse {

		var included model.InclusionType = model.Nothing;
		if incMap[a.ID] {
			included = model.Included
		} else if excMap[a.ID] {
			included = model.Excluded
		} else if incMap[a.Artists[0].ID] {
			included = model.Included
		} else if excMap[a.Artists[0].ID] {
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


func trackToResponse(tracks []spotify.SimpleTrack, playlist *model.Playlist, album spotify.ID) []model.ItemResponse {
	trackIDs := utils.Map(tracks, func(a spotify.SimpleTrack) spotify.ID { return a.ID })

	incMap := GetInclusionMap(playlist.SpotifyID, append(trackIDs, album))
	excMap := GetExclusionMap(playlist.SpotifyID, append(trackIDs, album))

	return utils.Map(tracks, func(a spotify.SimpleTrack) model.ItemResponse {

		var included model.InclusionType = model.Nothing;
		if incMap[a.ID] {
			included = model.Included
		} else if excMap[a.ID] {
			included = model.Excluded
		} else if incMap[album] {
			included = model.Included
		} else if excMap[album] {
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

func getInclusionsExclusions(playlist *model.Playlist, itemIDs []spotify.ID) ([]spotify.ID, []spotify.ID) {
    incSet := make(map[spotify.ID]bool)
    excSet := make(map[spotify.ID]bool)

    processPlaylist := func(p *model.Playlist) {
        for _, id := range GetIncludedItemsFromPlaylist(p, itemIDs) {
            incSet[id] = true
        }
        for _, id := range GetExcludedItemsFromPlaylist(p, itemIDs) {
            excSet[id] = true
        }
    }
    processPlaylist(playlist)
    currentPlaylists := GetIncludedPlaylistsFromPlaylist(playlist)
    
    visited := make(map[spotify.ID]bool)
    visited[playlist.SpotifyID] = true

    for len(currentPlaylists) > 0 {
        var nextLevel = []model.Playlist{}
        for _, p := range currentPlaylists {
            if visited[p.SpotifyID] {
                continue
            }
            visited[p.SpotifyID] = true
            processPlaylist(&p)
            nextLevel = append(nextLevel, GetIncludedPlaylistsFromPlaylist(&p)...)
        }
        currentPlaylists = nextLevel
    }
    return mapToSlice(incSet), mapToSlice(excSet)
}

func mapToSlice(m map[spotify.ID]bool) []spotify.ID {
    slice := make([]spotify.ID, 0, len(m))
    for k := range m {
        slice = append(slice, k)
    }
    return slice
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
