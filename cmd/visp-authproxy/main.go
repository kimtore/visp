package main

import (
	"net/http"
	"os"

	"github.com/ambientsound/visp/spotify/proxyserver"
	"github.com/ambientsound/visp/version"
	log "github.com/sirupsen/logrus"
)

// Simple HTTP server that lets users authenticate with Spotify.
// Access tokens are sent back to the client.

func main() {
	clientID := os.Getenv("VISP_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("VISP_OAUTH_CLIENT_SECRET")
	redirectURL := os.Getenv("VISP_OAUTH_REDIRECT_URL")
	renderMode := os.Getenv("VISP_OAUTH_RENDER_MODE") // render in text or JSON
	listenAddr := os.Getenv("VISP_OAUTH_LISTEN_ADDR")

	var renderer spotify_proxyserver.Renderer
	if renderMode == "json" {
		renderer = &spotify_proxyserver.JSONRenderer{}
	} else {
		renderer = &spotify_proxyserver.TextRenderer{}
	}

	server := spotify_proxyserver.New(clientID, clientSecret, redirectURL, renderer)

	log.Infof("Visp-authproxy %s starting", version.Version)
	log.Infof("Listening for connections on %s...\n", listenAddr)
	handler := spotify_proxyserver.Router(server)
	err := http.ListenAndServe(listenAddr, handler)

	if err != nil {
		log.Errorf("Fatal error: %s", err)
		os.Exit(1)
	}

	log.Errorf("Visp oauth proxy terminated")
}
