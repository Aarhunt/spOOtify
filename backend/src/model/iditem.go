package model

type ItemType int
type InclusionType int

const (
	PlaylistItem ItemType = iota
	ArtistItem
	AlbumItem
	TrackItem
)

const (
	Nothing InclusionType = iota
	Included
	Excluded
	IncludedByProxy
	ExcludedByProxy
)

var itemType = map[ItemType]string{
	PlaylistItem: "playlist",
	ArtistItem:   "artist",
	AlbumItem:    "album",
	TrackItem:    "track",
}

var inclusionType = map[InclusionType]string{
	Nothing:         "nothing",
	Included:        "included",
	Excluded:        "excluded",
	IncludedByProxy: "includedbyproxy",
	ExcludedByProxy: "excludedbyproxy",
}

type IdItem struct {
	SpotifyID string `gorm:"primaryKey;type:varchar(255);not null" json:"id" example:"37i9dQZF1DXcBWIGoYBM3M"`
	ItemType  ItemType
	// Playlists 	[]Playlist `gorm:"many2many:playlist_inclusions;" json:"included_in"`
}

type ItemInclusionRequest struct {
	ItemSpotifyID string   `json:"spotid" binding:"required"`
	ItemType      ItemType `json:"type" binding:"required"`
	PlaylistID    string   `json:"playlistid" binding:"required"`
	Include       *bool    `json:"include" binding:"required"`
}

type ItemRequest struct {
	ParentID   string   `json:"parentid" binding:"required"`
	PlaylistID string   `json:"playlistid" binding:"required"`
	ItemType   ItemType `json:"type" binding:"required"`
}

type SearchRequest struct {
	Query      string   `json:"query"`
	PlaylistID string   `json:"playlistid" binding:"required"`
	ItemType   ItemType `json:"type"`
}

type ItemPlaylistRequest struct {
	ParentSpotifyID string `json:"pspotid" binding:"required"`
	ChildSpotifyID  string `json:"cspotid" binding:"required"`
}

type ItemResponse struct {
	SpotifyID        string        `json:"spotifyID" binding:"required"`
	Name             string        `json:"name" binding:"required"`
	Icon             []Image       `json:"icon"`
	ItemType         ItemType      `json:"itemType"`
	Included         InclusionType `json:"included" binding:"required"`
	InclusionByProxy *bool         `json:"inclusionByProxy"`
	SortData         int           `json:"sortdata" binding:"required"`
}

type InclusionResponse struct {
	SpotifyID string        `json:"spotifyID"`
	Included  InclusionType `json:"included"`
}

func IsIncluded(iInc int, iExc int) bool {
	if iInc == 0 {
		return false
	}
	if iExc == 0 {
		return iInc > 0
	}
	return (iInc + iExc) < 0
}
