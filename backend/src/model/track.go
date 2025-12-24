package model

type Track struct {
	ID string `json:"id"`
	SpotifyID	string	`json:"spotid"`
	Title	string	`json:"title"`
	Artist	string	`json:"artist"`
}

func (t Track) getSpotifyID() string {
	return t.SpotifyID;
}


