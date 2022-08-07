package main

import (
	"context"
	"spotifies-be/config"
	"spotifies-be/externalapi/httpspotify"
	"spotifies-be/util/log"
	"spotifies-be/web"

	"github.com/gorilla/sessions"
	spotifyAuth "github.com/zmb3/spotify/v2/auth"
)

func main() {
	log.InitLogger()
	defer log.Logger.Sync()

	cfg := config.LoadConfig()

	// FIXME: Learn about cookie path
	sessionStore := sessions.NewCookieStore([]byte(cfg.SessionKey))

	spotifyAuthClient := spotifyAuth.New(
		spotifyAuth.WithClientID(cfg.Spotify.ClientID),
		spotifyAuth.WithRedirectURL(cfg.Spotify.RedirectURI),
		spotifyAuth.WithScopes(spotifyAuth.ScopeUserReadPrivate),
	)

	spotifyClient := httpspotify.NewHTTPSpotifyClient(cfg)

	err := web.NewHTTPServer(cfg, sessionStore, *spotifyAuthClient, spotifyClient).Start()
	if err != nil {
		log.Error(context.Background(), "Failed to start HTTP server", err)
		return
	}
}
