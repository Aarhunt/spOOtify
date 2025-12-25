package model

import (
	"context"
	"slices"

	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"
)

type Playlist struct {
	gorm.Model
	Name              string `json:"name"`
	SpotifyID         spotify.ID
	Inclusions        []IdItem   `gorm:"many2many:playlist_playlists;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	IncludedPlaylists []Playlist `gorm:"many2many:playlist_playlists;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Exclusions        []IdItem   `gorm:"many2many:playlist_playlists;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type PlaylistCreateRequest struct {
	Name string `json:"name" binding:"required" example:"My Playlist"`
}

type PlaylistResponse struct {
	Name              string `json:"name"`
	SpotifyID         spotify.ID
	Inclusions        []IdItem   `gorm:"many2many:playlist_playlists;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	IncludedPlaylists []Playlist `gorm:"many2many:playlist_playlists;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Exclusions        []IdItem   `gorm:"many2many:playlist_playlists;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
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

func (p Playlist) ToResponse() *PlaylistResponse {
	return &PlaylistResponse{
		Name:              p.Name,
		SpotifyID:         p.SpotifyID,
		Inclusions:        p.Inclusions,
		IncludedPlaylists: p.IncludedPlaylists,
		Exclusions:        p.Exclusions,
	}
}

func (p Playlist) getTracks(ctx context.Context, client spotify.Client) []IdItem {
	var excludeTracks []IdItem
	var resultTracks []IdItem

	for _, p1 := range p.IncludedPlaylists {
		playlistTracks := p1.getTracks(ctx, client)
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
