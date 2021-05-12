package spotify_proxyserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// https://developer.spotify.com/documentation/general/guides/scopes/
var scopes = []string{
	"playlist-modify-private",
	"playlist-modify-public",
	"playlist-read-collaborative",
	"playlist-read-private",
	"user-follow-modify",
	"user-follow-read",
	"user-library-modify",
	"user-library-read",
	"user-modify-playback-state",
	"user-read-currently-playing",
	"user-read-playback-state",
	"user-read-recently-played",
	"user-top-read",
}

const (
	cookieName  = "token"
	loginURL    = "/oauth/login"
	callbackURL = "/oauth/callback"
	RefreshURL  = "/oauth/refresh"
)

type Handler struct {
	auth     spotify.Authenticator
	frontend Renderer
	json     Renderer
}

func New(clientID, clientSecret, redirectURL string, frontendRenderer Renderer) *Handler {
	authenticator := spotify.NewAuthenticator(redirectURL, scopes...)
	authenticator.SetAuthInfo(clientID, clientSecret)

	return &Handler{
		auth:     authenticator,
		frontend: frontendRenderer,
		json:     &JSONRenderer{},
	}
}

// First step of oauth2 client credentials flow.
// Store a cookie in the user's browser with the XSRF protection token,
// then redirect to Spotify's authentication URL.
func (h *Handler) ServeLogin(w http.ResponseWriter, r *http.Request) {
	u, err := uuid.NewRandom()
	if err != nil {
		h.frontend.Render(w, http.StatusServiceUnavailable, err, nil)
		log.Errorf("generate uuid: %s", err)
		return
	}

	cookie := &http.Cookie{
		Name:    cookieName,
		Value:   u.String(),
		Expires: time.Now().Add(time.Hour),
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, h.auth.AuthURL(u.String()), http.StatusFound)

	return
}

// Callback URL of oauth2 client credentials flow.
// Check the XSRF cookie and exchange the authentication code from the URL with
// an access token using Spotify's oauth2 API.
// The token is returned back to the user, for use in the client.
func (h *Handler) ServeCallback(w http.ResponseWriter, r *http.Request) {
	// Get state parameter from cookie
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		h.frontend.Render(w, http.StatusBadRequest, err, nil)
		log.Errorf("get cookie: %s", err)
		return
	}

	// Exchange credentials into Spotify token
	token, err := h.auth.Token(cookie.Value, r)
	if err != nil {
		h.frontend.Render(w, http.StatusForbidden, err, nil)
		log.Errorf("token exchange: %s", err)
		return
	}

	// Return token to client
	h.frontend.Render(w, http.StatusOK, nil, token)
}

// Token refresh helper endpoint. Takes a valid oauth2 token containing an
// access code and a refresh token, and returns a fresh token in return.
func (h *Handler) RefreshCallback(w http.ResponseWriter, r *http.Request) {
	token := &oauth2.Token{}
	err := json.NewDecoder(r.Body).Decode(token)

	if err != nil {
		h.json.Render(w, http.StatusBadRequest, err, nil)
		log.Errorf("decode oauth2 token: %s", err)
		return
	}

	// force token expiration
	token.Expiry = time.Now().Add(-time.Hour)

	// retrieve new token through automatic refresh
	cli := h.auth.NewClient(token)
	token, err = cli.Token()

	if err != nil {
		h.json.Render(w, http.StatusInternalServerError, err, nil)
		log.Errorf("refresh oauth2 token: %s", err)
		return
	}

	h.json.Render(w, http.StatusOK, nil, token)
}

func Router(handler *Handler) chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.NoCache)

	router.Get(loginURL, handler.ServeLogin)
	router.Get(callbackURL, handler.ServeCallback)
	router.Post(RefreshURL, handler.RefreshCallback)

	return router
}
