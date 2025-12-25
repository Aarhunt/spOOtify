package services

import (
	"log"
	"net/http"

	"github.com/aarhunt/autistify/src"
	"github.com/aarhunt/autistify/src/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetPlaylists godoc
// @Summary      Get all playlists
// @Description  Responds with the list of all playlists as JSON.
// @Tags         playlists
// @Produce      json
// @Success      200  {array}   model.Playlist
// @Failure      500  {object}  map[string]string
// @Router       /playlist [get]
func GetPlaylists(c *gin.Context) {
	var playlists []model.Playlist
	db := src.GetDbConn().Db

	db.Find(&playlists)

	c.IndentedJSON(http.StatusOK, playlists)
}

// PostPlaylist godoc
// @Summary      Create new playlist
// @Description  Create a new playlist locally and on Spotify
// @Tags         playlists
// @Accept       json
// @Produce      json
// @Param        playlist body model.PlaylistCreateRequest true "Playlist name"
// @Success      201 {object} model.Playlist
// @Failure      400 {object} gin.H
// @Router       /playlist/create [post]
func PostPlaylist(c *gin.Context) {
	spotiConn := src.GetSpotifyConn()
	ctx, client, user := spotiConn.Ctx, spotiConn.Client, spotiConn.UserID
	db := src.GetDbConn().Db

	var req model.PlaylistCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusCreated, localPlaylist)
}
