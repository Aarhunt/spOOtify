package controllers

import (
	"net/http"
	"strconv"

	"github.com/aarhunt/autistify/src/services"
	"github.com/gin-gonic/gin"
)

// DeletePlaylist godoc
// @Summary      Delete a playlist
// @Description  Delete a specific playlist by its ID
// @Tags         playlists
// @Param        id   path      string  true  "Playlist ID"
// @Success      204  {object}  nil
// @Failure      500  {object}  map[string]string
// @Router       /playlist/{id} [delete]
func DeletePlaylist(c *gin.Context) {
    idStr := c.Param("id")

    // 1. Convert string to uint
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format. ID must be a positive integer"})
        return
    }

    // 2. Pass the uint to the service
    result := services.DeletePlaylist(uint(id))
    
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
