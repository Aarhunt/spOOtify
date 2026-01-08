package controllers

import (
	"net/http"

	"github.com/aarhunt/spootify/src"
	"github.com/gin-gonic/gin"
)

// GetAuthStatusController godoc
// @Summary      Check authentication status
// @Description  Returns true if the backend has a valid Spotify connection.
// @Tags         auth
// @Produce      json
// @Success      200  {object}  map[string]bool "authenticated: true"
// @Router       /auth/status [get]
func GetAuthStatus(c *gin.Context) {
    isLoggedIn := (src.GetSpotifyConn() != nil)
    
    c.JSON(http.StatusOK, gin.H{
        "authenticated": isLoggedIn,
    })
}
