package src

import (
	"context"
	"fmt"
	"golang.org/x/oauth2/clientcredentials"
	"os"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"github.com/zmb3/spotify/v2"
	"net/http"
	"log"
)

func init() {

	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)
	results, err := client.Search(ctx, "Rammstein", spotify.SearchTypeArtist)
	if err != nil {
		log.Fatal(err)
	}

	// handle album results
	if results.Artists != nil {
		fmt.Println("Artists:")
		for _, item := range results.Artists.Artists {
			fmt.Println("   ", item.Name)
		}
	}
}

