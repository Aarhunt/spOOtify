package services

import (
	"fmt"
	"log"

	"github.com/aarhunt/autistify/src"
	"github.com/aarhunt/autistify/src/model"
	"github.com/aarhunt/autistify/src/utils"
	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"
)

func GetPlaylists() ([]model.Playlist, error) {
	var playlists []model.Playlist
	db := src.GetDbConn().Db

	err := db.Find(&playlists)
	return playlists, err.Error
}

func GetPlaylistsResponse(id spotify.ID) []model.ItemResponse {
	p, _ := getPlaylist(id)
	playlists, _ := GetPlaylists()
	ids := getPlaylistsRecursive(*p, make(map[spotify.ID]bool), 1);

	return utils.Map(playlists, func(p model.Playlist) model.ItemResponse {
		included := model.Nothing
		if (ids[p.SpotifyID] > 0) {
			included = model.Included
		}
		return model.ItemResponse{
			SpotifyID: p.SpotifyID,
			Name:      p.Name,
			Icon:      p.Images,
			ItemType:  model.PlaylistItem,
			Included:  included,
		}
	})
}

func getPlaylist(id spotify.ID) (*model.Playlist, error) {
	dbConn := src.GetDbConn()
	ctx, db := dbConn.Ctx, dbConn.Db
	playlist, err := gorm.G[model.Playlist](db).Where("spotify_id = ?", id).First(ctx)

	return &playlist, err
}

func DeletePlaylist(id spotify.ID) *gorm.DB {
	db := src.GetDbConn().Db

	return db.Delete(&model.Playlist{}, id)
}

func PostPlaylist(req model.PlaylistCreateRequest) (*model.PlaylistResponse, error) {
	spotiConn := src.GetSpotifyConn()
	ctx, client, user := spotiConn.Ctx, spotiConn.Client, spotiConn.UserID
	db := src.GetDbConn().Db

	spotPlaylist, err := client.CreatePlaylistForUser(ctx, user, req.Name, "", false, false)

	if err != nil {
		log.Fatal(err)
	}

	localPlaylist := model.Playlist{
		SpotifyID:		   spotPlaylist.ID,
		Name:              req.Name,
		Inclusions:        []model.IdItem{},
		IncludedPlaylists: []*model.Playlist{},
		Exclusions:        []model.IdItem{},
		Images: 		   spotPlaylist.Images,
	}

	err = gorm.G[model.Playlist](db).Create(ctx, &localPlaylist)

	return localPlaylist.ToResponse(), err
}

func ClearPlaylists() (int, error) {
	dbConn := src.GetDbConn()
	ctx, db := dbConn.Ctx, dbConn.Db

	return gorm.G[model.Playlist](db).Where("true").Delete(ctx)
}

func GetIncludedIDsFromPlaylist(p *model.Playlist, ids []spotify.ID) ([]spotify.ID) {
	db:= src.GetDbConn().Db

	var includedIDs = []spotify.ID{}
	_ =	db.Table("playlist_inclusions").
        Where("playlist_spotify_id = ?", p.SpotifyID).
        Where("id_item_spotify_id IN ?", ids).
        Pluck("id_item_spotify_id", &includedIDs).Error

	return includedIDs
}

func isItemIncluded(playlistID spotify.ID, itemID spotify.ID) bool {
    var count int64
    src.GetDbConn().Db.Table("playlist_inclusions").
        Where("playlist_spotify_id = ? AND id_item_spotify_id = ?", playlistID, itemID).
        Count(&count)
    
    return count > 0
}

func GetExcludedIDsFromPlaylist(p *model.Playlist, ids []spotify.ID) ([]spotify.ID) {
	db:= src.GetDbConn().Db

	var excludedIDs = []spotify.ID{}
	_ =	db.Table("playlist_exclusions").
        Where("playlist_spotify_id = ?", p.SpotifyID).
        Where("id_item_spotify_id IN ?", ids).
        Pluck("id_item_spotify_id", &excludedIDs).Error

	return excludedIDs
}

func isItemExcluded(playlistID spotify.ID, itemID spotify.ID) bool {
    var count int64
    src.GetDbConn().Db.Table("playlist_exclusions").
        Where("playlist_spotify_id = ? AND id_item_spotify_id = ?", playlistID, itemID).
        Count(&count)
    
    return count > 0
}

func GetInclusionMap(playlistID spotify.ID, ids []spotify.ID) map[spotify.ID]bool {
	var matchedIDs []spotify.ID
	src.GetDbConn().Db.Table("playlist_inclusions").
        Where("playlist_spotify_id = ? AND id_item_spotify_id IN ?", playlistID, ids).
        Pluck("id_item_spotify_id", &matchedIDs)

	m := make(map[spotify.ID]bool)
	for _, id := range matchedIDs { m[id] = true }
	return m
}

func GetAllInclusions(id spotify.ID) []model.ItemResponse {
    var items []model.IdItem
	var playlists = GetPlaylistsResponse(id)


    err := src.GetDbConn().Db.
        Table("id_items").
        Joins("JOIN playlist_inclusions ON playlist_inclusions.id_item_spotify_id = id_items.spotify_id").
        Where("playlist_inclusions.playlist_spotify_id = ?", id).
        Find(&items).Error

    if err != nil {
        fmt.Printf("Error fetching inclusions: %v\n", err)
        return []model.ItemResponse{}
    }

	itemResponses := IncludedItemsToResponse(items, model.Included)

    return append(playlists, itemResponses...)
}

func GetExclusionMap(playlistID spotify.ID, ids []spotify.ID) map[spotify.ID]bool {
	var matchedIDs []spotify.ID
	src.GetDbConn().Db.Table("playlist_exclusions").
        Where("playlist_spotify_id = ? AND id_item_spotify_id IN ?", playlistID, ids).
        Pluck("id_item_spotify_id", &matchedIDs)


	m := make(map[spotify.ID]bool)
	for _, id := range matchedIDs { m[id] = true }
	return m
}

func GetAllExclusions(id spotify.ID) []model.ItemResponse {
    var items []model.IdItem

    err := src.GetDbConn().Db.
        Table("id_items").
        Joins("JOIN playlist_exclusions ON playlist_exclusions.id_item_spotify_id = id_items.spotify_id").
        Where("playlist_exclusions.playlist_spotify_id = ?", id).
        Find(&items).Error

    if err != nil {
        fmt.Printf("Error fetching exclusions: %v\n", err)
        return []model.ItemResponse{}
    }

	itemResponses := IncludedItemsToResponse(items, model.Included)

    return itemResponses
}

func GetIncludedPlaylistsFromPlaylist(p *model.Playlist) ([]model.Playlist) {
	dbConn := src.GetDbConn()
	ctx, db := dbConn.Ctx, dbConn.Db

	var includedPlaylists = []model.Playlist{}
	db.Model(&p).Association("IncludedPlaylists").Find(ctx, &includedPlaylists)

	return includedPlaylists
}

func getPlaylistsRecursive(p model.Playlist, visited map[spotify.ID]bool, i int) map[spotify.ID]int {
    if visited[p.SpotifyID] {
        return nil
    }

    visited[p.SpotifyID] = true

    includedPlaylists := make(map[spotify.ID]int)

	// Get all the playlists that are included in p
	for _, nested := range GetIncludedPlaylistsFromPlaylist(&p) {
		nestedPlaylists := getPlaylistsRecursive(nested, visited, i + 1)
		for id, _ := range nestedPlaylists {
			includedPlaylists[id] = i
		}
	}

	return includedPlaylists
}

func getTracksFromPlaylist(p model.Playlist) []spotify.ID {
	inclusions, exclusions := getTracksRecursive(p, make(map[spotify.ID]bool), 0)

    finalTracks := make([]spotify.ID, 0, len(inclusions))
    for id, inc := range inclusions {
		exc := exclusions[id]
		if model.IsIncluded(inc, exc) {finalTracks = append(finalTracks, id)}
    }

	return finalTracks
}

func getTracksRecursive(p model.Playlist, visited map[spotify.ID]bool, i int) (map[spotify.ID]int, map[spotify.ID]int) {
    if visited[p.SpotifyID] {
        return nil, nil
    }
    visited[p.SpotifyID] = true

    excludedMap := make(map[spotify.ID]int)
	includedMap := make(map[spotify.ID]int)

	for _, nested := range GetIncludedPlaylistsFromPlaylist(&p) {
		nestedInclusions, nestedExclusions := getTracksRecursive(nested, visited, i + 1)
		for id, val := range nestedInclusions {
			if val != 0 {
				includedMap[id]	= val + 3 * i
			} else {
				includedMap[id] = val
			}
		}
		for id, val := range nestedExclusions {
			if val != 0 {
				excludedMap[id]	= val + 3 * i
			} else {
				excludedMap[id] = val
			}
		}
	}

	inclusions, exclusions := GetAllInclusions(p.SpotifyID), GetAllExclusions(p.SpotifyID)

    for _, v := range exclusions {
        switch v.ItemType {
        case model.Artist:
            for _, t := range getTracksFromArtist(v.SpotifyID) {
				if excludedMap[t.SpotifyID] == 0 {
					excludedMap[t.SpotifyID] = -3 
				}
            }
        case model.Album:
            for _, t := range getTracksFromAlbum(v.SpotifyID) {
				if excludedMap[t.SpotifyID] == 0 || excludedMap[t.SpotifyID] == -3 {
					excludedMap[t.SpotifyID] = -2
				}
            }
        case model.Track:
            excludedMap[v.SpotifyID] = -1
		}
    }

    for _, v := range inclusions {
        switch v.ItemType {
        case model.Artist:
            for _, t := range getTracksFromArtist(v.SpotifyID) {
				if includedMap[t.SpotifyID] == 0 {
					includedMap[t.SpotifyID] = 3 
				}
            }
        case model.Album:
            for _, t := range getTracksFromAlbum(v.SpotifyID) {
				if includedMap[t.SpotifyID] == 0 || excludedMap[t.SpotifyID] == 3 {
					includedMap[t.SpotifyID] = 2
				}
            }
        case model.Track:
            includedMap[v.SpotifyID] = 1
		}
    }

    return includedMap, excludedMap
}

func PublishPlaylist(req model.PlaylistPublishRequest) error {
    spotiConn := src.GetSpotifyConn()
    ctx, client := spotiConn.Ctx, spotiConn.Client

    playlist, err := getPlaylist(req.SpotifyID)
    if err != nil {
        return err
    }

    trackIDs := getTracksFromPlaylist(*playlist)

    initialBatch := trackIDs[:min(len(trackIDs), 100)]
    err = client.ReplacePlaylistTracks(ctx, req.SpotifyID, initialBatch...)
    if err != nil {
        return err
    }

	log.Print(len(trackIDs))

    if len(trackIDs) > 100 {
        for i := 100; i < len(trackIDs); i += 100 {
            end := min(i + 100, len(trackIDs))
            
            chunk := trackIDs[i:end]
			_, err := client.AddTracksToPlaylist(ctx, req.SpotifyID, chunk...)
            if err != nil {
                return err
            }
        }
    }

    return nil
}
