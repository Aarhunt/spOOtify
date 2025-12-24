package model

type Artist struct {
	ID string `json:"id"`
	SpotifyID	string	`json:"spotid"`
	Name	string	`json:"name"`
	Albums	[]Album `json:"albums"`
}

func (a Artist) getSpotifyID() string {
	return a.SpotifyID;
}


