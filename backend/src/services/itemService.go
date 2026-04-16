package services

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/aarhunt/spootify/src"
	"github.com/aarhunt/spootify/src/model"
	"github.com/aarhunt/spootify/src/utils"
	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Include or exclude an item from a playlist
func IncludeExcludeItem(req model.ItemInclusionRequest, recurse bool, override bool) (*model.InclusionResponse, error) {
	dbConn := src.GetDbConn()
	db := dbConn.Db

	playlist, err := GetPlaylist(req.PlaylistID)

	// Create the item to be included
	newItem := model.IdItem{
		SpotifyID: req.ItemSpotifyID,
		ItemType:  req.ItemType,
	}

	// Automatically exclude some queries.
	if *req.Include && recurse {
		autoExclude(req, false)
	}

	returnItem := model.InclusionResponse{
		SpotifyID: req.ItemSpotifyID,
		Included:  model.Nothing,
	}

	db.Transaction(func(tx *gorm.DB) error {
		// Add the item if it does not exist yet.
		err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "spotify_id"}},
			UpdateAll: true,
		}).Omit(clause.Associations).Create(&newItem).Error
		if err != nil {
			return err
		}

		var count int64

		// Check whether the item is excluded already
		tx.Table("playlist_exclusions").
			Where("playlist_spotify_id = ? AND id_item_spotify_id = ?", playlist.SpotifyID, newItem.SpotifyID).
			Count(&count)
		isExcluded := count > 0

		// Check whether the item is included already
		tx.Table("playlist_inclusions").
			Where("playlist_spotify_id = ? AND id_item_spotify_id = ?", playlist.SpotifyID, newItem.SpotifyID).
			Count(&count)
		isIncluded := count > 0

		if *req.Include {
			// If we dont override, return
			if !override && isExcluded {
				returnItem.Included = model.Excluded
				return nil
			}

			// Otherwise, delete a potential exclusion and add the inclusion.
			tx.Model(&playlist).Association("Exclusions").Delete(&newItem)
			err = tx.Model(&playlist).Association("Inclusions").Append(&newItem)
			returnItem.Included = model.Included
		} else {
			// If we dont override, return
			if !override && isIncluded {
				returnItem.Included = model.Included
				return nil
			}

			// Otherwise, delete a potential inclusion and add the exclusion.
			tx.Model(&playlist).Association("Inclusions").Delete(&newItem)
			err = tx.Model(&playlist).Association("Exclusions").Append(&newItem)
			returnItem.Included = model.Excluded
		}
		if err != nil {
			return err
		}

		return nil
	})

	return &returnItem, err
}

// Exclude the autoexclusions
func autoExclude(req model.ItemInclusionRequest, undo bool) {
	// Default values
	ex := []string{}
	var itemType model.ItemType
	included := false
	excluded := false

	switch req.ItemType {
	case model.ArtistItem:
		albums := GetAlbumsFromArtistById(req.ItemSpotifyID, false)

		for i := len(albums) - 1; i >= 0; i-- {
			album := albums[i]

			if strings.Contains(strings.ToLower(album.Name), "live") {
				ex = append(ex, album.ID)
			} else if strings.Contains(strings.ToLower(album.Name), "instrumental") {
				ex = append(ex, album.ID)
			} else if strings.Contains(strings.ToLower(album.Name), "acoustic") {
				ex = append(ex, album.ID)
			} else {
				autoExclude(model.ItemInclusionRequest{
					ItemSpotifyID: album.ID,
					ItemType:      model.AlbumItem,
					PlaylistID:    req.PlaylistID,
					Include:       nil,
				}, undo)
			}
		}

		itemType = model.AlbumItem
	case model.AlbumItem:
		tracks := GetTracksFromAlbumById(req.ItemSpotifyID)
		for i := len(tracks) - 1; i >= 0; i-- {
			track := tracks[i]

			//TODO Use the autoexclusions
			if strings.Contains(strings.ToLower(track.Name), "- live") {
				ex = append(ex, track.ID)
			} else if strings.Contains(strings.ToLower(track.Name), "- instrumental") {
				ex = append(ex, track.ID)
			} else if strings.Contains(strings.ToLower(track.Name), "- acoustic") {
				ex = append(ex, track.ID)
			} else if strings.Contains(strings.ToLower(track.Name), "- orchestral") {
				ex = append(ex, track.ID)
			} else if strings.Contains(strings.ToLower(track.Name), "- single") {
				ex = append(ex, track.ID)
			}
		}

		itemType = model.TrackItem
	case model.TrackItem:
		return
	}

	if !undo {
		for _, id := range ex {
			IncludeExcludeItem(model.ItemInclusionRequest{
				ItemSpotifyID: id,
				ItemType:      itemType,
				PlaylistID:    req.PlaylistID,
				Include:       &included,
			}, false, false)
		}
	} else {
		for _, id := range ex {
			UndoIncludeExcludeItem(model.ItemInclusionRequest{
				ItemSpotifyID: id,
				ItemType:      itemType,
				PlaylistID:    req.PlaylistID,
				Include:       &excluded,
			})
		}
	}
}

// Undo the inclusion or exclusion of an item from a playlist
func UndoIncludeExcludeItem(req model.ItemInclusionRequest) (*model.InclusionResponse, error) {
	dbConn := src.GetDbConn()
	db := dbConn.Db

	playlist, err := GetPlaylist(req.PlaylistID)
	if err != nil {
		return nil, err
	}

	// Automatically undo exclusion of some queries.
	if *req.Include {
		autoExclude(req, true)
	}

	// Create the item to be included
	item := model.IdItem{SpotifyID: req.ItemSpotifyID}

	returnItem := model.InclusionResponse{
		SpotifyID: req.ItemSpotifyID,
		Included:  model.InclusionType(0),
	}

	// Delete the item from the correct association table
	err = db.Transaction(func(tx *gorm.DB) error {
		if *req.Include {
			if err := tx.Model(&playlist).Association("Inclusions").Delete(&item); err != nil {
				return err
			}
		} else {
			if err := tx.Model(&playlist).Association("Exclusions").Delete(&item); err != nil {
				return err
			}
		}

		return nil
	})

	return &returnItem, err
}

// Include a playlist into a playlist
func IncludePlaylist(req model.ItemPlaylistRequest) (*model.PlaylistResponse, error) {
	dbConn := src.GetDbConn()
	ctx, db := dbConn.Ctx, dbConn.Db

	parentPlaylist, err := gorm.G[model.Playlist](db).Where("spotify_id = ?", req.ParentSpotifyID).First(ctx)
	childPlaylist, err := gorm.G[model.Playlist](db).Where("spotify_id = ?", req.ChildSpotifyID).First(ctx)

	err = db.Model(&parentPlaylist).Association("IncludedPlaylists").Append(&childPlaylist)
	return parentPlaylist.ToResponse(), err
}

// Include a playlist into a playlist
func UndoIncludePlaylist(req model.ItemPlaylistRequest) (*model.PlaylistResponse, error) {
	dbConn := src.GetDbConn()
	ctx, db := dbConn.Ctx, dbConn.Db

	parentPlaylist, err := gorm.G[model.Playlist](db).Where("spotify_id = ?", req.ParentSpotifyID).First(ctx)
	childPlaylist, err := gorm.G[model.Playlist](db).Where("spotify_id = ?", req.ChildSpotifyID).First(ctx)

	err = db.Model(&parentPlaylist).Association("IncludedPlaylists").Delete(&childPlaylist)
	return parentPlaylist.ToResponse(), err
}

func ToItemResponse(items []model.IdItem, included model.InclusionType) []model.ItemResponse {
	results := []model.ItemResponse{}

	artists := []model.IdItem{}
	albums := []model.IdItem{}
	tracks := []model.IdItem{}

	for _, item := range items {
		switch item.ItemType {
		case model.ArtistItem:
			artists = append(artists, item)
		case model.AlbumItem:
			albums = append(albums, item)
		case model.TrackItem:
			tracks = append(tracks, item)
		}
	}

	toId := func(i model.IdItem) string { return i.SpotifyID }
	fullArtists := getArtistsByIds(utils.Map(artists, toId))
	fullAlbums := getAlbumsByIds(utils.Map(albums, toId))
	fullTracks := getTracksByIds(utils.Map(tracks, toId))

	results = append(results, utils.Map(fullArtists, func(a model.Artist) model.ItemResponse {
		return model.ItemResponse{
			SpotifyID: a.ID,
			Name:      a.Name,
			Icon:      a.Images,
			ItemType:  model.ArtistItem,
			Included:  included,
		}
	})...)

	results = append(results, utils.Map(fullAlbums, func(a model.Album) model.ItemResponse {
		return model.ItemResponse{
			SpotifyID: a.ID,
			Name:      a.Name,
			Icon:      a.Images,
			ItemType:  model.AlbumItem,
			Included:  included,
			SortData:  a.ReleaseYear,
		}
	})...)

	results = append(results, utils.Map(fullTracks, func(a model.Track) model.ItemResponse {
		return model.ItemResponse{
			SpotifyID: a.ID,
			Name:      a.Name,
			Icon:      a.Album.Images,
			ItemType:  model.TrackItem,
			Included:  included,
			SortData:  int(a.TrackNumber),
		}
	})...)

	return results
}

func ArtistToResponse(artists []model.Artist, playlist *model.Playlist) []model.ItemResponse {
	artistIDs := utils.Map(artists, func(a model.Artist) string {
		return a.ID
	})

	includedArtists, excludedArtists := getInclusionsExclusions(playlist, artistIDs)

	return utils.Map(artists, func(a model.Artist) model.ItemResponse {
		included := model.Nothing
		if slices.Contains(includedArtists, a.ID) {
			included = model.Included
		} else if slices.Contains(excludedArtists, a.ID) {
			included = model.Excluded
		}
		return model.ItemResponse {
			SpotifyID: a.ID,
			Name:      a.Name,
			Icon:      a.Images,
			ItemType:  model.ArtistItem,
			Included:  included,
		}
	})
}

func AlbumToResponse(albums []model.Album, playlist *model.Playlist) []model.ItemResponse {
	albumIDs := utils.Map(albums, func(a model.Album) string { return a.ID })
	artistIDs := utils.Map(albums, func(a model.Album) string { return a.ArtistID })

	incMap := GetInclusionMap(playlist.SpotifyID, append(albumIDs, artistIDs...))
	excMap := GetExclusionMap(playlist.SpotifyID, append(albumIDs, artistIDs...))

	fmt.Println(incMap)

	return utils.Map(albums, func(a model.Album) model.ItemResponse {

		var included model.InclusionType = model.Nothing
		if incMap[a.ID] {
			included = model.Included
		} else if excMap[a.ID] {
			included = model.Excluded
		} else if incMap[a.ArtistID] {
			if a.AlbumType == "album" {
				included = model.IncludedByProxy
			}
		} else if excMap[a.ArtistID] {
			included = model.ExcludedByProxy
		}

		inclusion := a.AlbumType == "album"
		return model.ItemResponse{
			SpotifyID:        a.ID,
			Name:             a.Name,
			Icon:             a.Images,
			ItemType:         model.AlbumItem,
			Included:         included,
			InclusionByProxy: &inclusion,
			SortData:         a.ReleaseYear,
		}
	})
}

func singleAlbumTrackToResponse(tracks []model.Track, playlist *model.Playlist, album string) []model.ItemResponse {
	trackIDs := utils.Map(tracks, func(a model.Track) string { return a.ID })

	incMap := GetInclusionMap(playlist.SpotifyID, append(trackIDs, album))
	excMap := GetExclusionMap(playlist.SpotifyID, append(trackIDs, album))

	return utils.Map(tracks, func(a model.Track) model.ItemResponse {

		var included model.InclusionType = model.Nothing
		if incMap[a.ID] {
			included = model.Included
		} else if excMap[a.ID] {
			included = model.Excluded
		} else if incMap[album] {
			included = model.IncludedByProxy
		} else if excMap[album] {
			included = model.ExcludedByProxy
		}

		inclusion := true
		return model.ItemResponse{
			SpotifyID:        a.ID,
			Name:             a.Name,
			Icon:             a.Album.Images,
			ItemType:         model.TrackItem,
			Included:         included,
			InclusionByProxy: &inclusion,
			SortData:         int(a.TrackNumber),
		}
	})
}

func TrackToResponse(tracks []model.Track, playlist *model.Playlist) []model.ItemResponse {
	trackIDs := utils.Map(tracks, func(a model.Track) string { return a.ID })
	albumIDs := utils.Map(tracks, func(a model.Track) string { return a.Album.ID })

	incMap := GetInclusionMap(playlist.SpotifyID, append(trackIDs, albumIDs...))
	excMap := GetExclusionMap(playlist.SpotifyID, append(trackIDs, albumIDs...))

	return utils.Map(tracks, func(a model.Track) model.ItemResponse {

		var included model.InclusionType = model.Nothing
		if incMap[a.ID] {
			included = model.Included
		} else if excMap[a.ID] {
			included = model.Excluded
		} else if incMap[a.Album.ID] {
			included = model.IncludedByProxy
		} else if excMap[a.Album.ID] {
			included = model.ExcludedByProxy
		}

		inclusion := true
		return model.ItemResponse{
			SpotifyID:        a.ID,
			Name:             a.Name,
			Icon:             a.Album.Images,
			ItemType:         model.TrackItem,
			Included:         included,
			InclusionByProxy: &inclusion,
			SortData:         int(a.TrackNumber),
		}
	})
}

func getInclusionsExclusions(playlist *model.Playlist, itemIDs []string) ([]string, []string) {
	incSet := make(map[string]bool)
	excSet := make(map[string]bool)

	processPlaylist := func(id string) {
		for id := range GetInclusionMap(id, itemIDs) {
			incSet[id] = true
		}
		for id := range GetExclusionMap(id, itemIDs) {
			excSet[id] = true
		}
	}
	processPlaylist(playlist.SpotifyID)
	currentPlaylists := GetPlaylistChildren(playlist.SpotifyID)

	visited := make(map[string]bool)
	visited[playlist.SpotifyID] = true

	for len(currentPlaylists) > 0 {
		var nextLevel = []string{}
		for _, p := range currentPlaylists {
			if visited[p] {
				continue
			}
			visited[p] = true
			processPlaylist(p)
			nextLevel = append(nextLevel, GetPlaylistChildren(p)...)
		}
		currentPlaylists = nextLevel
	}
	return mapToSlice(incSet), mapToSlice(excSet)
}

func mapToSlice(m map[string]bool) []string {
	slice := make([]string, 0, len(m))
	for k := range m {
		slice = append(slice, k)
	}
	return slice
}

func SearchArtist(req model.SearchRequest) []model.ItemResponse {
	conn := src.GetSpotifyConn()
	ctx, client := conn.Ctx, conn.Client

	playlist, err := GetPlaylist(req.PlaylistID)
	results, err := client.Search(ctx, req.Query, spotify.SearchTypeArtist, spotify.Limit(5))

	// handle album results
	if err != nil {
		log.Fatal("help")
	}
	return ArtistToResponse(model.ToArtists(results.Artists.Artists), playlist)
}

func SearchAlbum(req model.SearchRequest) []model.ItemResponse {
	conn := src.GetSpotifyConn()
	ctx, client := conn.Ctx, conn.Client

	playlist, err := GetPlaylist(req.PlaylistID)
	results, err := client.Search(ctx, req.Query, spotify.SearchTypeAlbum, spotify.Limit(5))

	// handle album results
	if err != nil {
		log.Fatal("help")
	}
	return AlbumToResponse(model.ToAlbums(results.Albums.Albums), playlist)
}

func SearchTrack(req model.SearchRequest) []model.ItemResponse {
	conn := src.GetSpotifyConn()
	ctx, client := conn.Ctx, conn.Client

	playlist, err := GetPlaylist(req.PlaylistID)
	results, err := client.Search(ctx, req.Query, spotify.SearchTypeTrack, spotify.Limit(5))

	// handle album results
	if err != nil {
		log.Fatal("help")
	}
	return TrackToResponse(model.ToTracks(results.Tracks.Tracks), playlist)
}

