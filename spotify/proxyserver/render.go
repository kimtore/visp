package spotify_proxyserver

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"
)

type Response struct {
	Error error  `json:"error,omitempty"`
	Token string `json:"token,omitempty"`
}

type Renderer interface {
	Render(w http.ResponseWriter, code int, err error, token *oauth2.Token)
}

type JSONRenderer struct{}

type TextRenderer struct{}

func EncodeTokenString(token *oauth2.Token) (string, error) {
	jsontok, err := json.Marshal(token)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(jsontok), nil
}

func (r *JSONRenderer) Render(w http.ResponseWriter, code int, err error, token *oauth2.Token) {
	jsontok, errenc := EncodeTokenString(token)
	if errenc != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Response{
		Error: err,
		Token: jsontok,
	})
}

func (r *TextRenderer) Render(w http.ResponseWriter, code int, err error, token *oauth2.Token) {
	jsontok, errenc := EncodeTokenString(token)
	if errenc != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/html")
	w.WriteHeader(code)

	if err == nil {
		w.Write([]byte(`<h1>Successfully authorized!</h1><p>Copy and paste the following text into Visp:</p><p style="font-family: monospace; word-break: break-all">`))
		w.Write([]byte(`:auth ` + jsontok))
		w.Write([]byte(`</p>`))
		return
	}

	w.Write([]byte(`<h1>An error occurred while authorizing:</h1><p>`))
	w.Write([]byte(err.Error()))
	w.Write([]byte(`</p><p><a href="/oauth/login">Click here to start over again</a></p>`))
}
