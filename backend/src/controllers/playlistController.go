package controllers

import (
	"net/http"

	"github.com/aarhunt/autistify/src/model"
	"github.com/aarhunt/autistify/src/services"
	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify/v2"
)

// GetPlaylists godoc
// @Summary      Get all playlists
// @Description  Responds with the list of all playlists as JSON.
// @Tags         playlists
// @Produce      json
// @Success      200  {array}   model.PlaylistResponse
// @Failure      500  {object}  map[string]string
// @Router       /playlist [get]
func GetPlaylists(c *gin.Context) {
	playlists, err := services.GetPlaylists()

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    }

	c.IndentedJSON(http.StatusOK, playlists)
}

// DeletePlaylist godoc
// @Summary      Delete a playlist
// @Description  Delete a specific playlist by its ID
// @Tags         playlists
// @Param        id   path      string  true  "Playlist ID"
// @Success      204  {object}  nil
// @Failure      500  {object}  map[string]string
// @Router       /playlist/{id} [delete]
func DeletePlaylist(c *gin.Context) {
    id := c.Param("id")

    result := services.DeletePlaylist(spotify.ID(id))
    
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    if result.RowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"message": "Playlist not found"})
        return
    }

    c.Status(http.StatusNoContent)
}

// PostPlaylist godoc
// @Summary      Create new playlist
// @Description  Create a new playlist locally and on Spotify
// @Tags         playlists
// @Accept       json
// @Produce      json
// @Param        playlist body model.PlaylistCreateRequest true "Playlist name"
// @Success      201 {object} model.PlaylistResponse
// @Failure      400 {object} gin.H
// @Router       /playlist [post]
func PostPlaylist(c *gin.Context) {

	var req model.PlaylistCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.PostPlaylist(req)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result})
        return
    }

	c.IndentedJSON(http.StatusCreated, result)
}

// ClearPlaylists godoc
// @Summary      Clear all playlists
// @Description  Deletes every playlist record in the database
// @Tags         playlists
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /playlist [delete]
func ClearPlaylists(c *gin.Context) {
    result, err := services.ClearPlaylists()

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.IndentedJSON(http.StatusOK, result)
}
