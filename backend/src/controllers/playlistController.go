package controllers

import (
	"net/http"

	"github.com/aarhunt/spootify/src/model"
	"github.com/aarhunt/spootify/src/services"
	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify/v2"
)

// GetPlaylists godoc
// @Summary      Get all playlists
// @Description  Responds with the list of all playlists as JSON.
// @Tags         playlist
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
// @Tags         playlist
// @Param        id   path      string  true  "Playlist ID"
// @Produce      json
// @Success      200  {array}   model.ItemResponse
// @Failure      500  {object}  map[string]string
// @Router       /playlist/{id}/playlists [get]
func GetPlaylistsById(c *gin.Context) {
    id := c.Param("id")

	playlists := services.SearchPlaylist(model.SearchRequest{Query: "", PlaylistID: spotify.ID(id), ItemType: model.PlaylistItem})

	c.IndentedJSON(http.StatusOK, playlists)
}

// DeletePlaylist godoc
// @Summary      Delete a playlist
// @Description  Delete a specific playlist by its ID
// @Tags         playlist
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

// RenamePlaylistgodoc
// @Summary      Rename a playlist
// @Description  Updates the display name of a playlist in the local database.
// @Tags         playlist
// @Accept       json
// @Produce      json
// @Param        id    path      string         true  "Spotify Playlist ID"
// @Param        body  body      model.PlaylistCreateRequest  true  "New name for the playlist"
// @Success      200   {object}  map[string]interface{} "message: Success"
// @Failure      400   {object}  map[string]string "error: Invalid input"
// @Failure      404   {object}  map[string]string "error: Playlist not found"
// @Failure      500   {object}  map[string]string "error: Database error"
// @Router       /playlist/{id}/rename [put]
func RenamePlaylist(c *gin.Context) {
    playlistID := c.Param("id")
	var req model.PlaylistCreateRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "New name is required"})
        return
    }

    rowsAffected, err := services.RenamePlaylist(spotify.ID(playlistID), req.Name)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update database", "details": err.Error()})
        return
    }

    if rowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "No playlist found with that ID"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Playlist renamed successfully",
        "newName": req.Name,
    })
}

// PostPlaylist godoc
// @Summary      Create new playlist
// @Description  Create a new playlist locally and on Spotify
// @Tags         playlist
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
// @Tags         playlist
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
// @Tags         playlist
// @Accept       json
// @Produce      json
// @Param        request  body      model.PlaylistPublishRequest  true  "Playlist Publish Request"
// @Success      200      {object}  map[string]string "message: Success"
// @Failure      400      {object}  map[string]string "error: Bad Request"
// @Failure      500      {object}  map[string]string "error: Internal Server Error"
// @Router       /playlist/publish [post]
func PublishPlaylist(c *gin.Context) {
    var req model.PlaylistPublishRequest

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

// PublishPlaylist handles the synchronization of the local playlists to Spotify.
// @Summary      Publish all playlists to Spotify
// @Description  Calculates the current tracklist based on inclusions/exclusions and replaces the Spotify playlist content.
// @Tags         playlist
// @Accept       json
// @Produce      json
// @Success      200      {object}  map[string]string "message: Success"
// @Failure      400      {object}  map[string]string "error: Bad Request"
// @Router       /playlist/publishall [post]
func PublishAllPlaylists(c *gin.Context) {
	playlists, err := services.GetPlaylists()

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    }

	for _, p := range playlists {
		err := services.PublishPlaylist(model.PlaylistPublishRequest{SpotifyID: p.SpotifyID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to sync with Spotify",
				"details": err.Error(),
			})
			return
		}
	}

    c.JSON(http.StatusOK, gin.H{
        "message":    "Playlists successfully synchronized",
    })
}

// GetPlaylistInclusions godoc
// @Summary      Get all included items for a playlist
// @Description  Fetches the list of all Playlists, Artists, Albums, and Tracks manually included in a specific playlist.
// @Tags         playlist
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

    inclusions := services.GetAllInclusions(spotify.ID(id))

    c.JSON(http.StatusOK, inclusions)
}

// GetPlaylistExclusions godoc
// @Summary      Get all excluded items for a playlist
// @Description  Fetches the list of all Artists, Albums, and Tracks manually excluded in a specific playlist.
// @Tags         playlist
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

    exclusions := services.GetAllExclusions(spotify.ID(id))

    c.JSON(http.StatusOK, exclusions)
}

// IncludePlaylist godoc
// @Summary      Nest a Playlist
// @Description  Includes one playlist inside another parent playlist
// @Tags         playlist
// @Accept       json
// @Produce      json
// @Param        request  body      model.ItemPlaylistRequest  true  "Playlist Linking Details"
// @Success      200      {object}  model.PlaylistResponse
// @Failure      500      {object}  map[string]string
// @Router       /playlist/include [post]
func IncludePlaylist(c *gin.Context) {
	var req model.ItemPlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ChildSpotifyID == req.ParentSpotifyID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot include playlist in itself"})
		return
	}

	res, err := services.IncludePlaylist(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// UndoIncludePlaylist godoc
// @Summary      Undo a nested Playlist
// @Description  Undoes inclusion of one playlist inside another parent playlist
// @Tags         playlist
// @Accept       json
// @Produce      json
// @Param        request  body      model.ItemPlaylistRequest  true  "Playlist Linking Details"
// @Success      200      {object}  model.PlaylistResponse
// @Failure      500      {object}  map[string]string
// @Router       /playlist/include/undo [post]
func UndoIncludePlaylist(c *gin.Context) { //TODO
	var req model.ItemPlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := services.UndoIncludePlaylist(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}


