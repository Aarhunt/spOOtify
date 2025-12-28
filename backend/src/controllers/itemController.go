package controllers

import (
	"log"
	"net/http"

	"github.com/aarhunt/autistify/src/model"
	"github.com/aarhunt/autistify/src/services"
	"github.com/gin-gonic/gin"
)

// IncludeExcludeItem godoc
// @Summary      Include or Exclude an Item
// @Description  Adds an IdItem to either the inclusions or exclusions of a playlist
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        request  body      model.ItemInclusionRequest  true  "Inclusion/Exclusion Details"
// @Success      200      {object}  model.InclusionResponse
// @Failure      500      {object}  map[string]string
// @Router       /playlist/item [post]
func IncludeExcludeItem(c *gin.Context) {
	var req model.ItemInclusionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Print(c.Get("include"))
		return
	}

	res, err := services.IncludeExcludeItem(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// UndoIncludeExcludeItem godoc
// @Summary      Undo Include or Exclude
// @Description  Removes an IdItem from both inclusions and exclusions of a playlist
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        request  body      model.ItemInclusionRequest  true  "Item to remove"
// @Success      200      {object}  model.InclusionResponse
// @Failure      500      {object}  map[string]string
// @Router       /playlist/item/undo [post]
func UndoIncludeExcludeItem(c *gin.Context) {
    var req model.ItemInclusionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    res, err := services.UndoIncludeExcludeItem(req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Returns the InclusionResponse with Included: 0 (Nothing)
    c.JSON(http.StatusOK, res)
}

// IncludePlaylist godoc
// @Summary      Nest a Playlist
// @Description  Includes one playlist inside another parent playlist
// @Tags         playlists
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

	res, err := services.IncludePlaylist(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetAlbumFromArtist godoc
// @Summary      Get Albums by Artist
// @Description  Fetches albums and singles from Spotify for a specific artist
// @Tags         spotify
// @Accept       json
// @Produce      json
// @Param        request  body      model.ItemRequest  true  "Artist and Playlist Context"
// @Success      200      {array}   model.ItemResponse
// @Failure      500      {object}  map[string]string
// @Router       /spotify/artist/albums [post]
func GetAlbumsFromArtist(c *gin.Context) {
	var req model.ItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := services.GetAlbumFromArtist(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetTracksFromAlbum godoc
// @Summary      Get Tracks by Album
// @Description  Fetches all tracks for a specific Spotify album
// @Tags         spotify
// @Accept       json
// @Produce      json
// @Param        request  body      model.ItemRequest  true  "Album and Playlist Context"
// @Success      200      {array}   model.ItemResponse
// @Failure      500      {object}  map[string]string
// @Router       /spotify/album/tracks [post]
func GetTracksFromAlbum(c *gin.Context) {
	var req model.ItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := services.GetTracksFromAlbum(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// Search godoc
// @Summary      Unified Spotify Search
// @Description  Search Spotify for Artists, Albums, or Tracks based on the provided ItemType.
// @Tags         search
// @Accept       json
// @Produce      json
// @Param        request  body      model.SearchRequest  true  "Search Query and Item Type"
// @Success      200      {array}   model.ItemResponse
// @Failure      400      {object}  map[string]string "Invalid Item Type"
// @Router       /search [post]
func Search(c *gin.Context) {
    var req model.SearchRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var results []model.ItemResponse

    // Route to the specific service based on ItemType
    switch req.ItemType {
	case model.PlaylistItem:
		results, _ = services.GetPlaylists()
    case model.Artist:
        results = services.SearchArtist(req)
    case model.Album:
        results = services.SearchAlbum(req)
    case model.Track:
        results = services.SearchTrack(req)
    default:
        c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported search item type"})
        return
    }

    c.JSON(http.StatusOK, results)
}
