package controllers

import (
	"log"
	"net/http"

	"github.com/aarhunt/autistify/src/services"
	"github.com/gin-gonic/gin"
)

// Search godoc
// @Summary      Search for items
// @Description  Searches for songs or items based on a query string parameter.
// @Tags         search
// @Produce      json
// @Param        query  path      string  true  "Search Query"
// @Success      200    {array}   model.SearchResponse  "Successful search"
// @Failure      404    {object}  map[string]string   "No elements found"
// @Failure      500    {object}  map[string]string   "Internal server error"
// @Router       /search/{query} [get]
func Search(c *gin.Context) {	
	query := c.Param("query")

	results, err := services.Search(query)

	if err != nil {
		log.Fatal(err)
	}

    if len(results) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"message": "No elements found"})
        return
    }

	c.IndentedJSON(http.StatusOK, results)
}
