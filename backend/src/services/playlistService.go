package services

import (
	"log"
	"net/http"

	"github.com/aarhunt/autistify/src"
	"github.com/aarhunt/autistify/src/model"
	"github.com/aarhunt/autistify/src/utils"
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

func DeletePlaylist(id uint) *gorm.DB {
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
		Name:              req.Name,
		SpotifyID:         spotPlaylist.ID,
		Inclusions:        []model.IdItem{},
		IncludedPlaylists: []model.Playlist{},
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
