package model

import (
	"github.com/zmb3/spotify/v2"
)

type Playlist struct {
	SpotifyID   spotify.ID `gorm:"primaryKey;type:varchar(255);not null" json:"id" example:"37i9dQZF1DXcBWIGoYBM3M"`
	Name              string `json:"name"`
	Inclusions        []IdItem   `gorm:"many2many:playlist_inclusions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	IncludedPlaylists []*Playlist `gorm:"many2many:playlist_nested_playlists;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	// IncludedIn 		  []*Playlist `gorm:"many2many:playlist_nested_playlists;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Exclusions        []IdItem   `gorm:"many2many:playlist_exclusions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Images 			  []spotify.Image `json:"images"`	
}

type PlaylistCreateRequest struct {
	Name string `json:"name" binding:"required" example:"My Playlist"`
}

type PlaylistRequest struct {
	SpotifyID         spotify.ID `json:"spotifyID"`
}

type PlaylistResponse struct {
	Name              string `json:"name"`
	SpotifyID         spotify.ID `json:"spotifyID"`
	Inclusions        []IdItem   `gorm:"many2many:playlist_playlists;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	IncludedPlaylists []*Playlist `gorm:"many2many:playlist_playlists;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Exclusions        []IdItem   `gorm:"many2many:playlist_playlists;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (p Playlist) ToResponse() *PlaylistResponse {
	return &PlaylistResponse{
		Name:              p.Name,
		SpotifyID:         p.SpotifyID,
		Inclusions:        p.Inclusions,
		IncludedPlaylists: p.IncludedPlaylists,
		Exclusions:        p.Exclusions,
	}
}

