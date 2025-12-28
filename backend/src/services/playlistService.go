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

func GetPlaylists() ([]model.PlaylistResponse, error) {
	var playlists []model.Playlist
	db := src.GetDbConn().Db

	err := db.Find(&playlists)
	return utils.Map(playlists, func(p model.Playlist) model.PlaylistResponse {
		return *p.ToResponse()
	}), err.Error
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
	}

	err = gorm.G[model.Playlist](db).Create(ctx, &localPlaylist)

	return localPlaylist.ToResponse(), err
}

func ClearPlaylists() (int, error) {
	dbConn := src.GetDbConn()
	ctx, db := dbConn.Ctx, dbConn.Db

	return gorm.G[model.Playlist](db).Where("true").Delete(ctx)
}

func GetIncludedItemsFromPlaylist(p *model.Playlist, ids []spotify.ID) ([]spotify.ID) {
	db:= src.GetDbConn().Db

	var includedIDs = []spotify.ID{}
	_ =	db.Table("playlist_inclusions").
        Where("playlist_spotify_id = ?", p.SpotifyID).
        Where("id_item_spotify_id IN ?", ids).
        Pluck("id_item_spotify_id", &includedIDs).Error

	return includedIDs
}

func IsItemIncluded(playlistID spotify.ID, itemID spotify.ID) bool {
    var count int64
    src.GetDbConn().Db.Table("playlist_inclusions").
        Where("playlist_spotify_id = ? AND id_item_spotify_id = ?", playlistID, itemID).
        Count(&count)
    
    return count > 0
}

func GetExcludedItemsFromPlaylist(p *model.Playlist, ids []spotify.ID) ([]spotify.ID) {
	db:= src.GetDbConn().Db

	var excludedIDs = []spotify.ID{}
	_ =	db.Table("playlist_exclusions").
        Where("playlist_spotify_id = ?", p.SpotifyID).
        Where("id_item_spotify_id IN ?", ids).
        Pluck("id_item_spotify_id", &excludedIDs).Error

	return excludedIDs
}

func IsItemExcluded(playlistID spotify.ID, itemID spotify.ID) bool {
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

func getAllInclusions(playlistID spotify.ID) []model.IdItem {
    var items []model.IdItem

    err := src.GetDbConn().Db.
        Table("id_items").
        Joins("JOIN playlist_inclusions ON playlist_inclusions.id_item_spotify_id = id_items.spotify_id").
        Where("playlist_inclusions.playlist_spotify_id = ?", playlistID).
        Find(&items).Error

    if err != nil {
        fmt.Printf("Error fetching inclusions: %v\n", err)
        return []model.IdItem{}
    }

    return items
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

func getAllExclusions(playlistID spotify.ID) []model.IdItem {
    var items []model.IdItem

    err := src.GetDbConn().Db.
        Table("id_items").
        Joins("JOIN playlist_exclusions ON playlist_exclusions.id_item_spotify_id = id_items.spotify_id").
        Where("playlist_exclusions.playlist_spotify_id = ?", playlistID).
        Find(&items).Error

    if err != nil {
        fmt.Printf("Error fetching exclusions: %v\n", err)
        return []model.IdItem{}
    }

    return items
}

func GetIncludedPlaylistsFromPlaylist(p *model.Playlist) ([]model.Playlist) {
	dbConn := src.GetDbConn()
	ctx, db := dbConn.Ctx, dbConn.Db

	var includedPlaylists = []model.Playlist{}
	db.Model(&p).Association("IncludedPlaylists").Find(ctx, &includedPlaylists)

	return includedPlaylists
}

func getTracks(p model.Playlist) []spotify.ID {
    return getTracksRecursive(p, make(map[spotify.ID]bool))
}

func getTracksRecursive(p model.Playlist, visited map[spotify.ID]bool) []spotify.ID {
    if visited[p.SpotifyID] {
        return nil
    }
    visited[p.SpotifyID] = true

    excludedMap := make(map[spotify.ID]model.InclusionType)
	includedMap := make(map[spotify.ID]model.InclusionType)
    resultMap := make(map[spotify.ID]spotify.ID)

    for _, nested := range GetIncludedPlaylistsFromPlaylist(&p) {
        nestedTracks := getTracksRecursive(nested, visited)
        for _, t := range nestedTracks {
            resultMap[t] = t 
		}
    }

	inclusions, exclusions := getAllInclusions(p.SpotifyID), getAllExclusions(p.SpotifyID)

    for _, v := range exclusions {
        switch v.ItemType {
        case model.Artist:
            for _, t := range getTracksFromArtist(v) {
				if excludedMap[t.SpotifyID] == 0 {
					excludedMap[t.SpotifyID] = -3 
				}
            }
        case model.Album:
            for _, t := range getTracksFromAlbum(v) {
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
            for _, t := range getTracksFromArtist(v) {
				if includedMap[t.SpotifyID] == 0 {
					includedMap[t.SpotifyID] = 3 
				}
            }
        case model.Album:
            for _, t := range getTracksFromAlbum(v) {
				if includedMap[t.SpotifyID] == 0 || excludedMap[t.SpotifyID] == 3 {
					includedMap[t.SpotifyID] = 2
				}
            }
        case model.Track:
            includedMap[v.SpotifyID] = 1
		}
    }

    finalTracks := make([]spotify.ID, 0, len(resultMap))
    for id, inc := range includedMap {
		exc := excludedMap[id]
		if inc.IsIncluded(exc) {finalTracks = append(finalTracks, id)}
    }

    return finalTracks
}

func PublishPlaylist(req model.PlaylistPublishRequest) error {
    spotiConn := src.GetSpotifyConn()
    ctx, client := spotiConn.Ctx, spotiConn.Client

    playlist, err := getPlaylist(req.SpotifyID)
    if err != nil {
        return err
    }

    trackIDs := getTracks(*playlist)

    initialBatch := trackIDs[:min(len(trackIDs), 100)]
    err = client.ReplacePlaylistTracks(ctx, req.SpotifyID, initialBatch...)
    if err != nil {
        return err
    }

	log.Print(len(trackIDs))

    if len(trackIDs) > 100 {
        for i := 100; i < len(trackIDs); i += 100 {
            end := i + 100
            if end > len(trackIDs) {
                end = len(trackIDs)
            }
            
            chunk := trackIDs[i:end]
			_, err := client.AddTracksToPlaylist(ctx, req.SpotifyID, chunk...)
            if err != nil {
                return err
            }
        }
    }

    return nil
}

func min(a, b int) int {
    if a < b { return a }
    return b
}
