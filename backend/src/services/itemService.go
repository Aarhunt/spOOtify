package services

import (
	"github.com/aarhunt/autistify/src"
	"github.com/aarhunt/autistify/src/model"
	"gorm.io/gorm"
)

func IncludeItem(req model.ItemRequest) (*model.PlaylistResponse, error) {
	dbConn := src.GetDbConn()
	ctx, db := dbConn.Ctx, dbConn.Db

	playlist, err := gorm.G[model.Playlist](db).Where("id = ?", req.PlaylistID).First(ctx)

	playlist.AddItem(model.IdItem{
		SpotifyID: req.ItemSpotifyID,	
		ItemType: req.ItemType,
		PlaylistID: req.PlaylistID,
	})	

	db.Save(&playlist)

	return playlist.ToResponse(), err
}


func ExcludeItem(req model.ItemRequest) (*model.PlaylistResponse, error) {
	dbConn := src.GetDbConn()
	ctx, db := dbConn.Ctx, dbConn.Db

	playlist, err := gorm.G[model.Playlist](db).Where("id = ?", req.PlaylistID).First(ctx)

	playlist.ExcludeItem(model.IdItem{
		SpotifyID: req.ItemSpotifyID,	
		ItemType: req.ItemType,
		PlaylistID: req.PlaylistID,
	})	

	db.Save(&playlist)

	return playlist.ToResponse(), err
}

func IncludePlaylist(req model.ItemPlaylistRequest) (*model.PlaylistResponse, error) {
	dbConn := src.GetDbConn()
	ctx, db := dbConn.Ctx, dbConn.Db

	parentPlaylist, err := gorm.G[model.Playlist](db).Where("id = ?", req.ParentSpotifyID).First(ctx)
	childPlaylist, err := gorm.G[model.Playlist](db).Where("id = ?", req.ChildSpotifyID).First(ctx)

	playlist.AddPlaylist()

	db.Save(&playlist)

	return playlist.ToResponse(), err
}
