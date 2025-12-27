package services

import (
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

func GetExclusionMap(playlistID spotify.ID, ids []spotify.ID) map[spotify.ID]bool {
	var matchedIDs []spotify.ID
	src.GetDbConn().Db.Table("playlist_exclusions").
        Where("playlist_spotify_id = ? AND id_item_spotify_id IN ?", playlistID, ids).
        Pluck("id_item_spotify_id", &matchedIDs)

	m := make(map[spotify.ID]bool)
	for _, id := range matchedIDs { m[id] = true }
	return m
}

func GetIncludedPlaylistsFromPlaylist(p *model.Playlist) ([]model.Playlist) {
	dbConn := src.GetDbConn()
	ctx, db := dbConn.Ctx, dbConn.Db

	var includedPlaylists = []model.Playlist{}
	db.Model(&p).Association("IncludedPlaylists").Find(ctx, &includedPlaylists)

	return includedPlaylists
}
