package services

import (
	"fmt"
	"log"
	"regexp"
	"slices"
	"strings"

	"github.com/aarhunt/spootify/src"
	"github.com/aarhunt/spootify/src/model"
	"github.com/aarhunt/spootify/src/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetPlaylists() ([]model.PlaylistResponse, error) {
	var playlists []model.Playlist
	db := src.GetDbConn().Db

	err := db.Find(&playlists)
	return utils.Map(playlists, func (p model.Playlist) model.PlaylistResponse { return  *p.ToResponse() }), err.Error
}

func SearchPlaylist(req model.SearchRequest) []model.ItemResponse {
	playlist, _ := GetPlaylist(req.PlaylistID)
	playlists, _ := GetPlaylists()
	ids := getPlaylistsRecursive(*playlist, make(map[string]bool));

	playlists = slices.DeleteFunc(playlists, func(p model.PlaylistResponse) bool {
		match, _ := regexp.Match(strings.ToLower(req.Query), []byte(strings.ToLower(p.Name)))
		return p.SpotifyID == playlist.SpotifyID || !match })

	return utils.Map(playlists, func(p model.PlaylistResponse) model.ItemResponse {
		included := model.Nothing
		if (ids[p.SpotifyID] >= 2) {
			included = model.IncludedByProxy
		} else if (ids[p.SpotifyID] == 1) {
			included = model.Included
		}

		return model.ItemResponse{
			SpotifyID: p.SpotifyID,
			Name:      p.Name,
			Icon:      []model.Image{}, 
			ItemType:  model.PlaylistItem,
			Included:  included,
		}
	})
}

func GetPlaylist(id string) (*model.Playlist, error) {
	dbConn := src.GetDbConn()
	ctx, db := dbConn.Ctx, dbConn.Db
	playlist, err := gorm.G[model.Playlist](db).Where("spotify_id = ?", id).First(ctx)

	return &playlist, err
}

func DeletePlaylist(id string) *gorm.DB {
	dbConn := src.GetDbConn()
	ctx, db := dbConn.Ctx, dbConn.Db
	client := src.GetSpotifyConn().Client

	client.UnfollowPlaylist(ctx, model.ToSpotifyID(id))

   	playlist, _ := GetPlaylist(id)

    return db.Select(clause.Associations).Delete(&playlist)
}

func RenamePlaylist(id string, name string) (int, error) {
	dbConn := src.GetDbConn()
	ctx, db := dbConn.Ctx, dbConn.Db

	client := src.GetSpotifyConn().Client

	client.ChangePlaylistName(ctx, model.ToSpotifyID(id), name)

	return gorm.G[model.Playlist](db).Where("spotify_id = ?", id).Update(ctx, "name", name)
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
		SpotifyID:		   spotPlaylist.ID.String(),
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

func GetIncludedIDsFromPlaylist(p *model.Playlist, ids []string) ([]string) {
	db:= src.GetDbConn().Db

	var includedIDs = []string{}
	_ =	db.Table("playlist_inclusions").
        Where("playlist_spotify_id = ?", p.SpotifyID).
        Where("id_item_spotify_id IN ?", ids).
        Pluck("id_item_spotify_id", &includedIDs).Error

	return includedIDs
}

func isItemIncluded(playlistID string, itemID string) bool {
    var count int64
    src.GetDbConn().Db.Table("playlist_inclusions").
        Where("playlist_spotify_id = ? AND id_item_spotify_id = ?", playlistID, itemID).
        Count(&count)
    
    return count > 0
}

func GetExcludedIDsFromPlaylist(p *model.Playlist, ids []string) ([]string) {
	db:= src.GetDbConn().Db

	var excludedIDs = []string{}
	_ =	db.Table("playlist_exclusions").
        Where("playlist_spotify_id = ?", p.SpotifyID).
        Where("id_item_spotify_id IN ?", ids).
        Pluck("id_item_spotify_id", &excludedIDs).Error

	return excludedIDs
}

func isItemExcluded(playlistID string, itemID string) bool {
    var count int64
    src.GetDbConn().Db.Table("playlist_exclusions").
        Where("playlist_spotify_id = ? AND id_item_spotify_id = ?", playlistID, itemID).
        Count(&count)
    
    return count > 0
}

func GetInclusionMap(playlistID string, ids []string) map[string]bool {
	var matchedIDs []string
	src.GetDbConn().Db.Table("playlist_inclusions").
        Where("playlist_spotify_id = ? AND id_item_spotify_id IN ?", playlistID, ids).
        Pluck("id_item_spotify_id", &matchedIDs)

	m := make(map[string]bool)
	for _, id := range matchedIDs { m[id] = true }
	return m
}

func GetAllInclusions(id string) []model.ItemResponse {
    var items []model.IdItem
	var playlists = SearchPlaylist(model.SearchRequest{Query: "", PlaylistID: id, ItemType: model.PlaylistItem})

	playlists = slices.DeleteFunc(playlists, func(req model.ItemResponse) bool {return !(req.Included == model.Included || req.Included == model.IncludedByProxy)})


    err := src.GetDbConn().Db.
        Table("id_items").
        Joins("JOIN playlist_inclusions ON playlist_inclusions.id_item_spotify_id = id_items.spotify_id").
        Where("playlist_inclusions.playlist_spotify_id = ?", id).
        Find(&items).Error

    if err != nil {
        fmt.Printf("Error fetching inclusions: %v\n", err)
        return []model.ItemResponse{}
    }

	itemResponses := ToItemResponse(items, model.Included)

    return append(playlists, itemResponses...)
}

func GetExclusionMap(playlistID string, ids []string) map[string]bool {
	var matchedIDs []string
	src.GetDbConn().Db.Table("playlist_exclusions").
        Where("playlist_spotify_id = ? AND id_item_spotify_id IN ?", playlistID, ids).
        Pluck("id_item_spotify_id", &matchedIDs)


	m := make(map[string]bool)
	for _, id := range matchedIDs { m[id] = true }
	return m
}

func GetAllExclusions(id string) []model.ItemResponse {
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

	itemResponses := ToItemResponse(items, model.Excluded)

    return itemResponses
}

func GetPlaylistParents(p *model.Playlist) ([]model.Playlist) {
	var includedParents = []model.Playlist{}

    _ = src.GetDbConn().Db.
        Table("playlists").
        Joins("JOIN playlist_nested_playlists ON playlist_nested_playlists.playlist_spotify_id = playlists.spotify_id").
        Where("playlist_nested_playlists.included_playlist_spotify_id = ?", p.SpotifyID).
        Find(&includedParents).Error

	return includedParents
}

func getParentsRecursive(p model.Playlist, visited map[string]bool) map[string]bool {
    if visited[p.SpotifyID] {
        return nil
    }

    visited[p.SpotifyID] = true

	parentPlaylists := GetPlaylistParents(&p)

	for _, nested := range parentPlaylists {
		nestedPlaylistsMap := getParentsRecursive(nested, visited) 
		for key := range nestedPlaylistsMap {
			visited[key] = true 
		}
	}

	return visited
}

func GetIncludedPlaylistsFromPlaylist(p *model.Playlist) ([]model.Playlist) {
	var includedPlaylists = []model.Playlist{}

    _ = src.GetDbConn().Db.
        Table("playlists").
        Joins("JOIN playlist_nested_playlists ON playlist_nested_playlists.included_playlist_spotify_id = playlists.spotify_id").
        Where("playlist_nested_playlists.playlist_spotify_id = ?", p.SpotifyID).
        Find(&includedPlaylists).Error

	return includedPlaylists
}

func getPlaylistsRecursive(p model.Playlist, visited map[string]bool) map[string]int {
    if visited[p.SpotifyID] {
        return nil
    }

    visited[p.SpotifyID] = true

    includedPlaylists := make(map[string]int)

	nestedPlaylists := GetIncludedPlaylistsFromPlaylist(&p)

	for _, nested := range nestedPlaylists {
		nestedPlaylistsMap := getPlaylistsRecursive(nested, visited) 
		for key, val := range nestedPlaylistsMap {
			includedPlaylists[key] = val + 1
		}
	}

	for _, nested := range nestedPlaylists {
		includedPlaylists[nested.SpotifyID] = 1
	}

	return includedPlaylists
}

func getTracksFromPlaylist(p model.Playlist) []string {
	inclusions, exclusions := getTracksRecursive(p, make(map[string]bool))

    finalTracks := []string{}
    for id, inc := range inclusions {
		exc := exclusions[id]
		if model.IsIncluded(inc, exc) {finalTracks = append(finalTracks, id)}
    }

	return finalTracks
}

func getTracksRecursive(p model.Playlist, visited map[string]bool) (map[string]int, map[string]int) {
    if visited[p.SpotifyID] {
        return nil, nil
    }
    visited[p.SpotifyID] = true

    excludedMap := make(map[string]int)
	includedMap := make(map[string]int)

	for _, nested := range GetIncludedPlaylistsFromPlaylist(&p) {
		nestedInclusions, nestedExclusions := getTracksRecursive(nested, visited)
		for id, val := range nestedInclusions {
			if val != 0 {
				includedMap[id]	= val + 3
			} else {
				includedMap[id] = val
			}
		}
		for id, val := range nestedExclusions {
			if val != 0 {
				excludedMap[id]	= val + 3
			} else {
				excludedMap[id] = val
			}
		}
	}

	inclusions, exclusions := GetAllInclusions(p.SpotifyID), GetAllExclusions(p.SpotifyID)

    for _, v := range exclusions {
        switch v.ItemType {
        case model.ArtistItem:
            for _, t := range getTracksFromArtistById(v.SpotifyID) {
				if excludedMap[t.ID] == 0 {
					excludedMap[t.ID] = -3 
				}
            }
        case model.AlbumItem:
            for _, t := range GetTracksFromAlbumById(v.SpotifyID) {
				if excludedMap[t.ID] == 0 || excludedMap[t.ID] == -3 {
					excludedMap[t.ID] = -2
				}
            }
        case model.TrackItem:
            excludedMap[v.SpotifyID] = -1
		}
    }

    for _, v := range inclusions {
        switch v.ItemType {
        case model.ArtistItem:
            for _, t := range getTracksFromArtistById(v.SpotifyID) {
				if includedMap[t.ID] == 0 {
					includedMap[t.ID] = 3 
				}
            }
        case model.AlbumItem:
            for _, t := range GetTracksFromAlbumById(v.SpotifyID) {
				if includedMap[t.ID] == 0 || excludedMap[t.ID] == 3 {
					includedMap[t.ID] = 2
				}
            }
        case model.TrackItem:
            includedMap[v.SpotifyID] = 1
		}
    }

    return includedMap, excludedMap
}

func PublishPlaylist(req model.PlaylistPublishRequest) error {
    spotiConn := src.GetSpotifyConn()
    ctx, client := spotiConn.Ctx, spotiConn.Client

    playlist, err := GetPlaylist(req.SpotifyID)
    if err != nil {
        return err
    }

	affected := getParentsRecursive(*playlist, make(map[string]bool))
	affectedPlaylists := []*model.Playlist{playlist}
	for id := range affected {
		parent, err := GetPlaylist(id)
		if err != nil {
			return err
		}
		affectedPlaylists = append(affectedPlaylists, parent)
	}

	for _, p := range affectedPlaylists {

		trackIDs := model.ToSpotifyIDs(getTracksFromPlaylist(*p))

		chunks := slices.Chunk(trackIDs, 100)

		err = client.ReplacePlaylistTracks(ctx, model.ToSpotifyID(p.SpotifyID))
		if err != nil {
			return err
		}

		for chunk := range chunks {
			_, err := client.AddTracksToPlaylist(ctx, model.ToSpotifyID(p.SpotifyID), chunk...)
			if err != nil {
				return err
			}
		}
	}
    return nil
}
