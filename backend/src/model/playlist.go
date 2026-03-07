package model

type Playlist struct {
	SpotifyID   string `gorm:"primaryKey;type:varchar(255);not null" json:"id" example:"37i9dQZF1DXcBWIGoYBM3M"`
	Name              string `json:"name"`
	Inclusions        []IdItem   `gorm:"many2many:playlist_inclusions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	IncludedPlaylists []*Playlist `gorm:"many2many:playlist_nested_playlists;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Exclusions        []IdItem   `gorm:"many2many:playlist_exclusions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type PlaylistCreateRequest struct {
	Name string `json:"name" binding:"required" example:"My Playlist"`
}

type PlaylistPublishRequest struct {
	SpotifyID         string `json:"spotifyID"`
}

type PlaylistResponse struct {
	Name              string `json:"name"`
	SpotifyID         string `json:"spotifyID"`
}

func (p Playlist) ToResponse() *PlaylistResponse {
	return &PlaylistResponse{
		Name:              p.Name,
		SpotifyID:         p.SpotifyID,
	}
}

