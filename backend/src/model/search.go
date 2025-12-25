package model

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aarhunt/autistify/src"
	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify/v2"
)

func Search(c *gin.Context) {	
	conn := src.GetSpotifyConn()
	ctx, client := conn.Ctx, conn.Client
	query := c.Param("query")

	results, err := client.Search(ctx, query, spotify.SearchTypePlaylist|spotify.SearchTypeAlbum|spotify.SearchTypeArtist)
	if err != nil {
		log.Fatal(err)
	}

	// handle album results
	if results.Artists != nil {
		fmt.Println("Artists:")
		for _, item := range results.Artists.Artists {
			fmt.Println("   ", item.Name)
		}
	}

	// handle album results
	if results.Albums != nil {
		fmt.Println("Albums:")
		for _, item := range results.Albums.Albums {
			fmt.Println("   ", item.Name)
		}
	}

	// handle playlist results
	if results.Playlists != nil {
		fmt.Println("Playlists:")
		for _, item := range results.Playlists.Playlists {
			fmt.Println("   ", item.Name)
		}
	}

	c.IndentedJSON(http.StatusOK, results.Artists)
}
