package model

import (
	"context"
	"slices"
	"gorm.io/gorm"
	"github.com/zmb3/spotify/v2"
)

type Playlist struct {
	gorm.Model
	SpotifyID         spotify.ID `json:"spotid"`
	Inclusions        []IdItem   `json:"inclusions"`
	IncludedPlaylists []Playlist `gorm:"many2many:playlist_playlists;"`
	Exclusions        []IdItem   `json:"exclusions"`
}

func getPlaylist(id spotify.ID) {

}

func (p *Playlist) AddItem(item IdItem) {
	p.Inclusions = append(p.Inclusions, item)
}

func (p *Playlist) ExcludeItem(item IdItem) {
	p.Exclusions = append(p.Exclusions, item)
}

func (p Playlist) getSpotifyID() spotify.ID {
	return p.SpotifyID
}

func (p Playlist) getTracksFromPlaylist(ctx context.Context, client spotify.Client) []IdItem {
	var excludeTracks []IdItem
	var resultTracks []IdItem

	for _, p1 := range p.IncludedPlaylists {
		playlistTracks := p1.getTracksFromPlaylist(ctx, client);
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

// func createPlaylist(ctx context.Context, client spotify.Client, user spotify.ID, name string) (*Playlist, error) {
// 	results, err := client.CreatePlaylistForUser(ctx, user.String(), name, "", false, false);
//
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	if results != nil {
// 		return &Playlist{
// 			ID: "1", //TODO
// 			SpotifyID: results.ID,
// 			Inclusions: []spotify.ID{},
// 			Exclusions: []spotify.ID{},
// 		}, err
// 	}
//
// 	return nil, err
// }
