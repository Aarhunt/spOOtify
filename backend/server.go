package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/aarhunt/autistify/src"
)

func main() {
	router := gin.Default()
	router.GET("/albums", src.GetAlbums)
	router.GET("/albums/:id", src.GetAlbumByID)
	router.POST("/albums", src.PostAlbums)

	router.Run("localhost:8080")
}

func handler(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprint(w, "Hello, World!")
}
