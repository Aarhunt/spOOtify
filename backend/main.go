package main

import (
	"fmt"
	"net/http"

	"github.com/aarhunt/spootify/docs"
	"github.com/aarhunt/spootify/src"
	"github.com/aarhunt/spootify/src/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

	_ = godotenv.Load()
	_ = src.GetSpotifyConn;

	docs.SwaggerInfo.Title = "Spootify API"
	docs.SwaggerInfo.Description = "Spoooootify"
	docs.SwaggerInfo.Version = "1.0"
	// docs.SwaggerInfo.Host = "http://localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"

	router := gin.Default()


	// Apply CORS middleware before your routes
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"}, // Your React URL
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	{
		v1 := router.Group("/api/v1")
		v1.POST("/search", controllers.Search)

		{
			play := v1.Group("/playlist")
			play.GET("", controllers.GetPlaylists)
			play.POST("", controllers.PostPlaylist)
			play.DELETE("/:id", controllers.DeletePlaylist)
			play.DELETE("", controllers.ClearPlaylists)
			play.POST("/item", controllers.IncludeExcludeItem)
			play.POST("/item/undo", controllers.UndoIncludeExcludeItem)
			play.POST("/include", controllers.IncludePlaylist)
			play.POST("/include/undo", controllers.UndoIncludePlaylist)
			play.POST("/publish", controllers.PublishPlaylist)
			play.POST("/publishall", controllers.PublishAllPlaylists)
			play.GET("/:id/inclusions", controllers.GetPlaylistInclusions)
			play.GET("/:id/exclusions", controllers.GetPlaylistExclusions)
			play.GET("/:id/playlists", controllers.GetPlaylistsById)
			play.PUT("/:id/rename", controllers.RenamePlaylist)
		}

		{
			spot := v1.Group("/spotify")
			spot.POST("/artist/albums", controllers.GetAlbumsFromArtist)
			spot.POST("/album/tracks", controllers.GetTracksFromAlbum)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run("0.0.0.0:8080")
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}
