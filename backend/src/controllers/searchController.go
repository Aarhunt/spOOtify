package controllers

import (
	"log"
	"net/http"

	"github.com/aarhunt/autistify/src/services"
	"github.com/gin-gonic/gin"
)

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
