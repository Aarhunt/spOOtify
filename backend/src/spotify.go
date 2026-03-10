package src

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

var (
    lockSpotifyConn     = &sync.Mutex{}
    spotifyConnInstance *SpotifyConn
    auth                *spotifyauth.Authenticator
    state               = "abc123_random_string" 
)

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

    redirectURI := os.Getenv("SPOTIFY_REDIRECT_URL")
    if redirectURI == "" {
        redirectURI = "http://127.0.0.1:8080/api/v1/callback"
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

func GetSpotifyConn() *SpotifyConn {
    lockSpotifyConn.Lock()
    defer lockSpotifyConn.Unlock()
    return spotifyConnInstance
}

// GetAuthURLController godoc
// @Summary      Get Spotify Authorization URL
// @Description  Returns the URL the user needs to visit to authorize the app on Spotify.
// @Tags         auth
// @Produce      json
// @Success      200  {object}  map[string]string "url: https://accounts.spotify.com/..."
// @Router       /auth/url [get]
func GetAuthURLController(c *gin.Context) {
    initSpotifyAuth()
    newState := generateRandomState(16)

    // Check if we are on a secure connection
    // (Localhost is usually NOT secure, so this must be false)
    isSecure := c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https"

    // Set SameSite to Lax (standard for OAuth)
    c.SetSameSite(http.SameSiteLaxMode)

    // SetCookie: name, value, maxAge, path, domain, secure, httpOnly
    // IMPORTANT: Passing an empty string "" for domain is best for localhost.
    c.SetCookie(
        "spotify_auth_state", 
        newState, 
        3600, 
        "/", 
        "",       // Domain MUST be empty for localhost to work reliably
        isSecure, // Must be FALSE on http://localhost
        true,     // HttpOnly: Yes, for security
    )

    url := auth.AuthURL(newState)
    c.JSON(http.StatusOK, gin.H{"url": url})
}

func generateRandomState(length int) string {
    b := make([]byte, length)
    if _, err := rand.Read(b); err != nil {
        return "fallback_random_state" // Extremely unlikely to fail
    }
    return hex.EncodeToString(b)
}

func CompleteAuthGin(c *gin.Context) {
    initSpotifyAuth() 

    // DEBUG: Print all cookies to the console to see if it's there
    log.Println("All Cookies:", c.Request.Header.Get("Cookie"))

	storedState, err := c.Cookie("spotify_auth_state")
    if err != nil {
        log.Println("Cookie Error:", err) // This will tell you if it's "named cookie not present"
        c.JSON(http.StatusBadRequest, gin.H{"error": "State cookie missing"})
        return
    }

    if c.Query("state") != storedState {
        c.JSON(http.StatusBadRequest, gin.H{"error": "State mismatch"})
        return
    }

    c.SetCookie("spotify_auth_state", "", -1, "/", "", false, true)

    token, err := auth.Token(c.Request.Context(), state, c.Request)
    if err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": "Couldn't get token", "details": err.Error()})
        log.Println("Auth Error:", err)
        return
    }

    client := spotify.New(auth.Client(c.Request.Context(), token))
    user, err := client.CurrentUser(c.Request.Context())
    if err != nil {
         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info"})
         return
    }

    lockSpotifyConn.Lock()
    spotifyConnInstance = &SpotifyConn{
        Ctx:    context.Background(),
        Client: client,
        UserID: user.ID,
    }
    lockSpotifyConn.Unlock()

    c.Data(http.StatusOK, "text/html", []byte(`
        <div style="font-family: sans-serif; text-align: center; margin-top: 50px;">
            <h1 style="color: #1DB954;">Connected!</h1>
            <p>You can close this window now.</p>
            <script>window.close();</script>
        </div>
    `))
}

