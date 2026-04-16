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

// GetPlaylists
// Get all playlists 
// Return []model.PlaylistResponse as playlists, and error for potential error that occured
func GetPlaylists() ([]model.PlaylistResponse, error) {
	var playlists []model.Playlist
	db := src.GetDbConn().Db

	err := db.Find(&playlists)
	return utils.Map(playlists, func (p model.Playlist) model.PlaylistResponse { return  *p.ToResponse() }), err.Error
}

// SearchPlaylist
// Search for a playlist by query. Does not return the current playlist.
// Returns ItemResponses for all found playlists.
func SearchPlaylist(req model.SearchRequest) []model.ItemResponse {
	playlist, _ := GetPlaylist(req.PlaylistID)
	playlists, _ := GetPlaylists()
	children := getPlaylistSubtree(req.PlaylistID, make(map[string]bool));

	playlists = slices.DeleteFunc(playlists, func(p model.PlaylistResponse) bool {
		match, _ := regexp.Match(strings.ToLower(req.Query), []byte(strings.ToLower(p.Name)))
		return p.SpotifyID == playlist.SpotifyID || !match })

	newVar := func(p model.PlaylistResponse) model.ItemResponse {
		included := model.Nothing
		if children[p.SpotifyID] >= 2 {
			included = model.IncludedByProxy
		} else if children[p.SpotifyID] == 1 {
			included = model.Included
		}
	
		return model.ItemResponse{
			SpotifyID: p.SpotifyID,
			Name:      p.Name,
			Icon:      []model.Image{},
			ItemType:  model.PlaylistItem,
			Included:  included,
		}
	}
	return utils.Map(playlists, newVar)
}

// GetPlaylist
// Gets a single playlist by id
func GetPlaylist(id string) (*model.Playlist, error) {
	dbConn := src.GetDbConn()
	ctx, db := dbConn.Ctx, dbConn.Db
	playlist, err := gorm.G[model.Playlist](db).Where("spotify_id = ?", id).First(ctx)

	return &playlist, err
}

// DeletePlaylist
// Deletes a single playlist by id
func DeletePlaylist(id string) *gorm.DB {
	dbConn := src.GetDbConn()
	ctx, db := dbConn.Ctx, dbConn.Db
	client := src.GetSpotifyConn().Client

	client.UnfollowPlaylist(ctx, model.ToSpotifyID(id))

   	playlist, _ := GetPlaylist(id)

    return db.Select(clause.Associations).Delete(&playlist)
}

// RenamePlaylist
// Renames a single playlist by id
// Returns amount of lists updated, and error
func RenamePlaylist(id string, name string) (int, error) {
	dbConn := src.GetDbConn()
	ctx, db := dbConn.Ctx, dbConn.Db

	client := src.GetSpotifyConn().Client

	client.ChangePlaylistName(ctx, model.ToSpotifyID(id), name)

	return gorm.G[model.Playlist](db).Where("spotify_id = ?", id).Update(ctx, "name", name)
}

// CreatePlaylist
// Create a playlist
// Returns the playlist, and potential error
func CreatePlaylist(req model.PlaylistCreateRequest) (*model.PlaylistResponse, error) {
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

// ClearPlaylists
// Delete all playlists.
// Returns amount of effected playlists.
func ClearPlaylists() (int, error) {
	dbConn := src.GetDbConn()
	ctx, db := dbConn.Ctx, dbConn.Db

	return gorm.G[model.Playlist](db).Where("true").Delete(ctx)
}

// isItemIncluded
// Checks a playlists by id includes an item by id
func isItemIncluded(playlistID string, itemID string) bool {
    var count int64
    src.GetDbConn().Db.Table("playlist_inclusions").
        Where("playlist_spotify_id = ? AND id_item_spotify_id = ?", playlistID, itemID).
        Count(&count)
    
    return count > 0
}

// isItemExcluded
// Checks a playlists by id excludes an item by id
func isItemExcluded(playlistID string, itemID string) bool {
    var count int64
    src.GetDbConn().Db.Table("playlist_exclusions").
        Where("playlist_spotify_id = ? AND id_item_spotify_id = ?", playlistID, itemID).
        Count(&count)
    
    return count > 0
}

// GetInclusionMap
// From the given IDs, create a map with all IDs that are included in the playlist
func GetInclusionMap(playlistID string, ids []string) map[string]bool {
	var matchedIDs []string
	src.GetDbConn().Db.Table("playlist_inclusions").
        Where("playlist_spotify_id = ? AND id_item_spotify_id IN ?", playlistID, ids).
        Pluck("id_item_spotify_id", &matchedIDs)

	m := make(map[string]bool)
	for _, id := range matchedIDs { m[id] = true }
	return m
}

// GetAllIncludedItems
// Get all included items of a given playlist.
// Return all included items as ItemResponse.
func GetAllIncludedItems(id string) []model.ItemResponse {
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

// GetExclusionMap
// From the given IDs, create a map with all IDs that are excluded in the playlist
func GetExclusionMap(playlistID string, ids []string) map[string]bool {
	var matchedIDs []string
	src.GetDbConn().Db.Table("playlist_exclusions").
        Where("playlist_spotify_id = ? AND id_item_spotify_id IN ?", playlistID, ids).
        Pluck("id_item_spotify_id", &matchedIDs)


	m := make(map[string]bool)
	for _, id := range matchedIDs { m[id] = true }
	return m
}


// GetAllExcludedItems
// Get all excluded items of a given playlist.
// Return all excluded items as ItemResponse.
func GetAllExcludedItems(id string) []model.ItemResponse {
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

// getPlaylistParents
// Get the parents of a playlist
// Returns the playlists as model.Playlist objects
func getPlaylistParents(p *model.Playlist) ([]model.Playlist) {
	var includedParents = []model.Playlist{}

    _ = src.GetDbConn().Db.
        Table("playlists").
        Joins("JOIN playlist_nested_playlists ON playlist_nested_playlists.playlist_spotify_id = playlists.spotify_id").
        Where("playlist_nested_playlists.included_playlist_spotify_id = ?", p.SpotifyID).
        Find(&includedParents).Error

	return includedParents
}

// getPlaylistParentsRecursive
// Returns a map of visited playlists, with a map that is true when the node is a root. .
func getPlaylistParentsRecursive(p model.Playlist, root map[string]bool) map[string]bool {
    if !root[p.SpotifyID] {
        return nil
    }

    root[p.SpotifyID] = false

	parentPlaylists := getPlaylistParents(&p)

	for _, nested := range parentPlaylists {
		nestedPlaylistsMap := getPlaylistParentsRecursive(nested, root) 
		for key := range nestedPlaylistsMap {
			root[key] = false 
		}
	}

	if (len(parentPlaylists) == 0) {root[p.SpotifyID] = false}

	return root
}

func GetPlaylistChildren(id string) ([]string) {
	var includedPlaylists = []string{}

	src.GetDbConn().Db.Table("playlist_nested_playlists").
		Where("playlist_nested_playlists.playlist_spotify_id = ?", id).
		Pluck("included_playlist_spotify_id", &includedPlaylists)


	return includedPlaylists
}

// Get the subtree below the given playlist id.
func getPlaylistSubtree(id string, visited map[string]bool) map[string]int {
    if visited[id] {
        return nil
    }

    visited[id] = true

    includedPlaylists := make(map[string]int)

	nestedPlaylists := GetPlaylistChildren(id)

	for _, nested := range nestedPlaylists {
		nestedPlaylistsMap := getPlaylistSubtree(nested, visited) 
		for key, val := range nestedPlaylistsMap {
			includedPlaylists[key] = val + 1
		}
	}

	for _, nested := range nestedPlaylists {
		includedPlaylists[nested] = 1
	}

	return includedPlaylists
}

// getTracksFromSinglePlaylist()
// Get the tracks from a single playlist 
func getTracksFromSinglePlaylist(id string) (map[string]int, map[string]int) {
    excludedMap := make(map[string]int)
	includedMap := make(map[string]int)

	
	inclusions, exclusions := GetAllIncludedItems(id), GetAllExcludedItems(id)

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

func getTracksFromPlaylist(p model.Playlist) []string {
	inclusions, exclusions := getTracksRecursive(p.SpotifyID, make(map[string]bool))

    finalTracks := []string{}
    for id, inc := range inclusions {
		exc := exclusions[id]
		if model.IsIncluded(inc, exc) {finalTracks = append(finalTracks, id)}
    }


	fmt.Println("%d tracks in playlist %s", len(finalTracks), p.Name)

	return finalTracks
}

func getTracksRecursive(id string, visited map[string]bool) (map[string]int, map[string]int) {
    if visited[id] {
        return nil, nil
    }
    visited[id] = true

    excludedMap := make(map[string]int)
	includedMap := make(map[string]int)

	for _, nested := range GetPlaylistChildren(id) {
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

	inclusions, exclusions := GetAllIncludedItems(id), GetAllExcludedItems(id)

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

func PublishPlaylistOld(req model.PlaylistPublishRequest) error {
    playlist, err := GetPlaylist(req.SpotifyID)
    if err != nil {
        return err
    }

	affected := getPlaylistParentsRecursive(*playlist, make(map[string]bool))
	affectedPlaylists := []*model.Playlist{playlist}
	for id := range affected {
		parent, err := GetPlaylist(id)
		if err != nil {
			return err
		}
		affectedPlaylists = append(affectedPlaylists, parent)
	}

	fmt.Println("%d playlists affected", len(affectedPlaylists))

	for _, p := range affectedPlaylists {

		trackIDs := getTracksFromPlaylist(*p)

		err := publishPlaylist(trackIDs,  p.SpotifyID)
		if err != nil {
			return err
		}
	}
    return nil
}

func publishPlaylist(trackStringIDs []string, id string) error {
    spotiConn := src.GetSpotifyConn()
    ctx, client := spotiConn.Ctx, spotiConn.Client

	trackIDs := model.ToSpotifyIDs(trackStringIDs)

	chunks := slices.Chunk(trackIDs, 100)

	err := client.ReplacePlaylistTracks(ctx, model.ToSpotifyID(id))
	if err != nil {
		return err
	}

	for chunk := range chunks {
		_, err := client.AddTracksToPlaylist(ctx, model.ToSpotifyID(id), chunk...)
		if err != nil {
			return err
		}
	}
	return nil
}

func PublishPlaylistEfficient(req model.PlaylistPublishRequest) error {
    playlist, err := GetPlaylist(req.SpotifyID)
    if err != nil {
        return err
    }

	path := getPlaylistParentsRecursive(*playlist, make(map[string]bool))

	fmt.Println(path)

	for id, val := range path {
		if val {
			if err != nil {
				return err
			}
			PublishPlaylistEfficientRecursive(id, path)

		}
	}

	return nil
	// start from thee root, recurse downwards, update playlist if in the path.
}

func PublishPlaylistEfficientRecursive(id string, path map[string]bool) []string {
	children := GetPlaylistChildren(id)

	totalTracks := []string{}

	for _, c := range children {
		totalTracks = append(totalTracks, PublishPlaylistEfficientRecursive(c, path)...)
	}

	currentInclusions, currentExclusions := getTracksFromSinglePlaylist(id)

	for id,inc := range currentInclusions {
		if exc, ok := currentExclusions[id]; ok {
			if model.IsIncluded(inc, exc) {totalTracks = append(totalTracks, id)}
		} else {
			totalTracks = append(totalTracks, id)
		}
	}
	
	if _, ok := path[id]; ok {
		publishPlaylist(totalTracks, id)
	}

	return totalTracks 
}

// Get all playlist child-parent connections.
func GetPlaylistTreeEdges() ([]model.PlaylistInclusionResponse, error) { 
	var edges []model.PlaylistInclusionResponse
    
    err := src.GetDbConn().Db.Table("playlist_nested_playlists").
        Select("playlist_spotify_id, included_playlist_spotify_id").
        Scan(&edges).Error
        
    return edges, err
}
