package model

import (
	"context"
	"github.com/zmb3/spotify/v2"
	"log"
)

type Playlist struct {
	ID string `json:"id"`
	SpotifyID	spotify.ID	`json:"spotid"`
	Inclusions []Trackable `json:"inclusions"`
	Exclusions []Trackable `json:"exclusions"`
}

func getPlaylist(id spotify.ID) {} 

func (p *Playlist) AddItem(item Trackable) {
	p.Inclusions = append(p.Inclusions, item)
}

func (p *Playlist) ExcludeItem(item Trackable) {
	p.Exclusions = append(p.Exclusions, item)
}

func (p Playlist) getSpotifyID() spotify.ID {
	return p.SpotifyID;
}

func (p Playlist) getTracks() []Track {
	var tracks []spotify.ID;

	for _, v := range p.Inclusions {
		switch v.(type) {
		case Artist:
			
		case Album:
			for _, v := range Album.Tracks {
				tracks = append(tracks, v.getSpotifyID())
			}
		case Track:
			tracks = append(tracks, v.getSpotifyID())
		}
	}
}

func createPlaylist(ctx context.Context, client spotify.Client, user spotify.ID, name string) (*Playlist, error) {
	results, err := client.CreatePlaylistForUser(ctx, user.String(), name, "", false, false);

	if err != nil {
		log.Fatal(err)
	}

	if results != nil {
		return &Playlist{
			ID: "1", //TODO
			SpotifyID: results.ID,	
			Inclusions: []Trackable{},
			Exclusions: []Trackable{},
		}, err
	}
	
	return nil, err
}
