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

// GetPlaylistsResponse godoc
// @Summary      Get all playlists as itemresponse items
// @Description  Responds with the list of all playlists as JSON.
// @Tags         playlists
// @Param        id   path      string  true  "Playlist ID"
// @Produce      json
// @Success      200  {array}   model.ItemResponse
// @Failure      500  {object}  map[string]string
// @Router       /playlist/{id}/playlists [get]
func GetPlaylistsResponse(c *gin.Context) {
    id := c.Param("id")

	req := model.PlaylistRequest{SpotifyID: spotify.ID(id)}

	playlists := services.GetPlaylistsResponse(req)

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

	req := model.PlaylistRequest{SpotifyID: spotify.ID(id)}

    result := services.DeletePlaylist(req)
    
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

// PublishPlaylist handles the synchronization of the local playlist state to Spotify.
// @Summary      Publish a playlist to Spotify
// @Description  Calculates the current tracklist based on inclusions/exclusions and replaces the Spotify playlist content.
// @Tags         playlists
// @Accept       json
// @Produce      json
// @Param        request  body      model.PlaylistPublishRequest  true  "Playlist Publish Request"
// @Success      200      {object}  map[string]string "message: Success"
// @Failure      400      {object}  map[string]string "error: Bad Request"
// @Failure      500      {object}  map[string]string "error: Internal Server Error"
// @Router       /playlist/publish [post]
func PublishPlaylist(c *gin.Context) {
    var req model.PlaylistRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }

    if req.SpotifyID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Playlist Spotify ID is required"})
        return
    }

    err := services.PublishPlaylist(req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Failed to sync with Spotify",
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":    "Playlist successfully synchronized",
        "playlistId": req.SpotifyID,
    })
}

// GetPlaylistInclusions godoc
// @Summary      Get all included items for a playlist
// @Description  Fetches the list of all Artists, Albums, and Tracks manually included in a specific playlist.
// @Tags         playlists
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Spotify Playlist ID"
// @Success      200  {array}   model.ItemResponse
// @Failure      400  {object}  map[string]string "error: Invalid ID"
// @Failure      500  {object}  map[string]string "error: Database error"
// @Router       /playlist/{id}/inclusions [get]
func GetPlaylistInclusions(c *gin.Context) {
    id := c.Param("id")
    
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Playlist ID is required"})
        return
    }

	req := model.PlaylistRequest{SpotifyID: spotify.ID(id)}

    inclusions := services.GetAllInclusions(req)

    if inclusions == nil {
        inclusions = []model.IdItem{}
    }

	inclusionResponses := services.IncludedItemsToResponse(inclusions, model.Included)

    c.JSON(http.StatusOK, inclusionResponses)
}

// GetPlaylistExclusions godoc
// @Summary      Get all excluded items for a playlist
// @Description  Fetches the list of all Artists, Albums, and Tracks manually excluded in a specific playlist.
// @Tags         playlists
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Spotify Playlist ID"
// @Success      200  {array}   model.ItemResponse
// @Failure      400  {object}  map[string]string "error: Invalid ID"
// @Failure      500  {object}  map[string]string "error: Database error"
// @Router       /playlist/{id}/exclusions [get]
func GetPlaylistExclusions(c *gin.Context) {
    id := c.Param("id")
    
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Playlist ID is required"})
        return
    }

	req := model.PlaylistRequest{SpotifyID: spotify.ID(id)}

    exclusions := services.GetAllExclusions(req)

    if exclusions == nil {
        exclusions = []model.IdItem{}
    }

	exclusionResponses := services.IncludedItemsToResponse(exclusions, model.Excluded)

    c.JSON(http.StatusOK, exclusionResponses)
}
