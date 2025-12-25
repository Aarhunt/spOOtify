package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aarhunt/autistify/src/model"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
"github.com/swaggo/gin-swagger" // gin-swagger middleware
"github.com/aarhunt/autistify/docs"
"github.com/swaggo/files" // swagger embed files
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
	err := godotenv.Load() 
    if err != nil {
    	log.Fatal(err)
    }

	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router := gin.Default()
	// router.GET("/albums", src.GetAlbums)
	// router.GET("/albums/:id", src.GetAlbumByID)
	// router.POST("/albums", src.PostAlbums)
	router.GET("/search/:query", model.Search)
	router.POST("/playlist/create/", model.PostPlaylist)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run("localhost:8080")
}

func handler(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprint(w, "Hello, World!")
}
