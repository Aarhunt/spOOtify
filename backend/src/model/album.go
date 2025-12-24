package model

type Album struct {
	ID string `json:"id"`
	SpotifyID	string	`json:"spotid"`
	Name	string	`json:"name"`
	Tracks	[]Track	`json:"songs"`
}

func (a Album) getSpotifyID() string {
	return a.SpotifyID;
}
