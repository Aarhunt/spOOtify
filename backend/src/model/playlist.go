package model

import (
	"context"
	"log"

	// "log"
	"net/http"
	"slices"

	"github.com/aarhunt/autistify/src"
	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"
)

type Playlist struct {
	gorm.Model
	Name              string
	SpotifyID         spotify.ID
	Inclusions        []IdItem   `gorm:"many2many:playlist_playlists;"`
	IncludedPlaylists []Playlist `gorm:"many2many:playlist_playlists;"`
	Exclusions        []IdItem   `gorm:"many2many:playlist_playlists;"`
}

// PlaylistCreateRequest represents the payload to create a playlist
type PlaylistCreateRequest struct {
	Name string `json:"name" binding:"required" example:"My Playlist"`
}

func GetPlaylists(c *gin.Context) {
	var playlists []Playlist
	db := src.GetDbConn().Db

	results := db.Find(&playlists)

	c.IndentedJSON(http.StatusOK, results)
}

// PostPlaylist godoc
// @Summary      Create new playlist
// @Description  Create a new playlist locally and on Spotify
// @Tags         playlists
// @Accept       json
// @Produce      json
// @Param        playlist body PlaylistCreateRequest true "Playlist name"
// @Success      201 {object} Playlist
// @Failure      400 {object} gin.H
// @Router       /playlist/create [post]
func PostPlaylist(c *gin.Context) {
	spotiConn := src.GetSpotifyConn()
	ctx, client, user := spotiConn.Ctx, spotiConn.AuthClient, spotiConn.User
	db := src.GetDbConn().Db

	var req PlaylistCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	spotPlaylist, err := client.CreatePlaylistForUser(ctx, user, req.Name, "", false, false)

	if err != nil {
		log.Fatal(err)
	}

	localPlaylist := Playlist{
		Name:              req.Name,
		SpotifyID:         spotPlaylist.ID,
		Inclusions:        []IdItem{},
		IncludedPlaylists: []Playlist{},
		Exclusions:        []IdItem{},
	}

	err = gorm.G[Playlist](db).Create(ctx, &localPlaylist)

	if err != nil {
		log.Fatal(err)
	}


	c.IndentedJSON(http.StatusCreated, localPlaylist)
}

func (p *Playlist) AddItem(item IdItem) {
	p.Inclusions = append(p.Inclusions, item)
}

func (p *Playlist) ExcludeItem(item IdItem) {
	p.Exclusions = append(p.Exclusions, item)
}

func (p *Playlist) AddPlaylist(p1 Playlist) {
	p.IncludedPlaylists = append(p.IncludedPlaylists, p1)
}

func (p Playlist) getSpotifyID() spotify.ID {
	return p.SpotifyID
}

func (p Playlist) getTracksFromPlaylist(ctx context.Context, client spotify.Client) []IdItem {
	var excludeTracks []IdItem
	var resultTracks []IdItem

	for _, p1 := range p.IncludedPlaylists {
		playlistTracks := p1.getTracksFromPlaylist(ctx, client)
		resultTracks = append(resultTracks, playlistTracks...)
	}

	for _, v := range p.Exclusions {
		switch v.ItemType {
		case Artist:
			tracks := getTracksFromArtist(ctx, client, v)
			excludeTracks = append(excludeTracks, tracks...)
		case Album:
			tracks := getTracksFromAlbum(ctx, client, v)
			excludeTracks = append(excludeTracks, tracks...)
		case Track:
			excludeTracks = append(excludeTracks, v)
		}
	}

	for _, v := range p.Inclusions {
		switch v.ItemType {
		case Artist:
			tracks := getTracksFromArtist(ctx, client, v)
			for _, track := range tracks {
				if !slices.Contains(excludeTracks, track) {
					resultTracks = append(tracks, track)
				}
			}
		case Album:
			tracks := getTracksFromAlbum(ctx, client, v)
			for _, track := range tracks {
				if !slices.Contains(excludeTracks, track) {
					resultTracks = append(tracks, track)
				}
			}
		case Track:
			if !slices.Contains(excludeTracks, v) {
				resultTracks = append(resultTracks, v)
			}
		}
	}

	return resultTracks
}
