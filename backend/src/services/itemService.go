package services

import (
	"github.com/aarhunt/autistify/src"
	"github.com/aarhunt/autistify/src/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func IncludeItem(req model.ItemRequest) (*model.PlaylistResponse, error) {
	dbConn := src.GetDbConn()
	ctx, db := dbConn.Ctx, dbConn.Db

	playlist, err := gorm.G[model.Playlist](db).Where("id = ?", req.PlaylistID).First(ctx)

	newItem := model.IdItem{
		SpotifyID: req.ItemSpotifyID,	
		ItemType: req.ItemType,
		Playlists: []model.Playlist{},
	}
	
	db.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
            UpdateAll: true, 
        }).Omit(clause.Associations).Create(&newItem).Error
        if err != nil {
            return err
        }

		if req.Include {
			err = tx.Model(&playlist).Association("Inclusions").Append(&newItem)
		} else {
			err = tx.Model(&playlist).Association("Exclusions").Append(&newItem)
		}
        if err != nil {
            return err
        }

        return nil
    })

	return playlist.ToResponse(), err
}

func IncludePlaylist(req model.ItemPlaylistRequest) (*model.PlaylistResponse, error) {
	dbConn := src.GetDbConn()
	ctx, db := dbConn.Ctx, dbConn.Db

	parentPlaylist, err := gorm.G[model.Playlist](db).Where("id = ?", req.ParentSpotifyID).First(ctx)
	childPlaylist, err := gorm.G[model.Playlist](db).Where("id = ?", req.ChildSpotifyID).First(ctx)

	err = db.Model(&parentPlaylist).Association("IncludedPlaylists").Append(&childPlaylist)
	return parentPlaylist.ToResponse(), err
}
