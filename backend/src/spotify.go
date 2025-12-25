package src

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const redirectURI = "http://127.0.0.1:3000/callback"

var (
	lockSpotifyConn      = &sync.Mutex{}
	spotifyConnInstance  *SpotifyConn
	auth                *spotifyauth.Authenticator
	state               = "abc123" // replace with random per session in prod
	authClientChan      = make(chan *spotify.Client)
)

// SpotifyConn holds an authenticated Spotify client
type SpotifyConn struct {
	Ctx    context.Context
	Client *spotify.Client
	UserID string
}


func initSpotifyAuth() {
	if auth != nil {
		return
	}

	if os.Getenv("SPOTIFY_ID") == "" || os.Getenv("SPOTIFY_SECRET") == "" {
		log.Fatal("SPOTIFY_ID and SPOTIFY_SECRET must be set")
	}

	auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopePlaylistModifyPublic,
			spotifyauth.ScopePlaylistModifyPrivate,
		),
	)
}


func startAuthServer() {
	http.HandleFunc("/callback", completeAuth)

	go func() {
		log.Println("Starting auth server on :3000")
		if err := http.ListenAndServe(":3000", nil); err != nil {
			log.Fatal(err)
		}
	}()
}


func authenticateUser() *spotify.Client {
	initSpotifyAuth()
	startAuthServer()

	url := auth.AuthURL(state)
	fmt.Println("Login to Spotify:", url)

	client := <-authClientChan
	return client
}


func completeAuth(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != state {
		http.Error(w, "State mismatch", http.StatusBadRequest)
		return
	}

	token, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Println(err)
		return
	}

	client := spotify.New(auth.Client(r.Context(), token))
	fmt.Fprintln(w, "Login successful!")

	authClientChan <- client
}


func createSpotifyConn() *SpotifyConn {
	ctx := context.Background()

	client := authenticateUser()

	user, err := client.CurrentUser(ctx)
	if err != nil {
		log.Fatal("Failed to get user:", err)
	}

	return &SpotifyConn{
		Ctx:    ctx,
		Client: client,
		UserID: user.ID,
	}
}


func GetSpotifyConn() *SpotifyConn {
	if spotifyConnInstance == nil {
		lockSpotifyConn.Lock()
		defer lockSpotifyConn.Unlock()

		if spotifyConnInstance == nil {
			log.Println("Creating Spotify connection")
			spotifyConnInstance = createSpotifyConn()
		}
	}

	return spotifyConnInstance
}
