# Spotify

## Library
Visp uses [zmb3's Spotify Go library](https://github.com/zmb3/spotify) to interface with Spotify.
This is a Go wrapper for working with Spotify's Web API.

## Authentication
Spotify's Web API authenticates users by using [OAuth2 tokens](https://oauth.net/2/). Unfortunately for you (the user)
this means that you cannot simply enter your username and password into Visp's configuration file.

To obtain a OAuth2 token that works with Spotify Web API, one must authenticate through a
[Spotify App](https://developer.spotify.com/dashboard/applications). This mechanism is more secure than using a username
and password, but works best for web applications because the Spotify App also has credentials. It is not allowed
to distribute these credentials, thus you are left with two options:

  1) Authenticate through the official Visp Spotify App:
     This is the default, and easiest option. (FIXME: set up as a public service)

  2) Create your own Spotify App and run your own authentication server:
     The server code is part of this repository, is small enough to understand fairly quickly, and
     guarantees that yours truly will not be able to access your Spotify data.
     [See entry point for the server code](../cmd/visp-oauth/main.go).
     (FIXME: document how to run it)

The OAuth flow is as follows:

  1) The user opens the URL `http://localhost:59999/oauth/login` in the browser.
  2) User is redirected by the authentication server to Spotify's login page.
  3) User logs into Spotify and is redirected back by Spotify to the authentication server together with a code.
  4) Authentication server uses the code to retrieve an access token to Spotify Web API.
  5) Authentication server gives the access token back to the user and discards it, nothing is saved.
  6) User copies the access token from the browser and pastes it into Visp.
  7) Visp can now access Spotify until the access token expires.
  8) When the access token expires, Visp automatically refreshes it through the authentication server.
  9) Step 8 repeats until the token no longer can be refreshed, and the user must start from step 1.
