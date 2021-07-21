package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/zneix/c2login/pkg/config"
)

var (
	clientID     string
	clientSecret string
	redirectURI  string
)

func index(w http.ResponseWriter, r *http.Request) {
	// "code" is required for us to receive code instead of token
	const responseType = "code"
	encodedScopes := url.PathEscape(strings.Join(scopes, " "))

	link := fmt.Sprintf("<a href=\"https://id.twitch.tv/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=%s&scope=%s\">login</a>", clientID, redirectURI, responseType, encodedScopes)

	w.Write([]byte(link))
}

func code(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no code provided..."))
		return
	}
	_, err := url.PathUnescape(r.URL.Query().Get("scope"))
	if err != nil {
		log.Printf("Error while unescaping scope list: %s\n", err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Malformed scopes provided..."))
		return
	}
	//scopeList := strings.Split(scopeString, " ")

	// Forming HTTP request to obtain access token
	const grantType = "authorization_code"
	url := fmt.Sprintf("https://id.twitch.tv/oauth2/token?client_id=%s&client_secret=%s&code=%s&grant_type=%s&redirect_uri=%s", clientID, clientSecret, code, grantType, redirectURI)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Printf("Error while forming token request: %v\n", err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error..."))
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Executing HTTP request
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Error while executing token request: %v\n", err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error..."))
		return
	}
	defer resp.Body.Close()

	// Reading body into a string
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error while reading token request body: %v\n", err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error..."))
		return
	}

	// Received non-200 response from id.twitch.tv/oauth2/token
	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusMultipleChoices {
		log.Printf("Non-200 response received from id.twitch.tv/oauth2/token: %v\n", resp.StatusCode)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error..."))
		return
	}

	// Unmarshal magic
	var root TokenResponse
	if err := json.Unmarshal(body, &root); err != nil {
		log.Printf("Error while unmarshaling token request body: %v\n", err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error..."))
		return
	}

	accountData := fmt.Sprintf("oauth_token=%s;refresh_token=%s;username=%s;user_id=%s;client_id=%s", root.AccessToken, root.RefreshToken, "...", "...", clientID)

	w.Write([]byte(accountData))
}

func handleMainRoutes(router *chi.Mux, cfg config.SiteConfig) {
	if cfg.ClientID == "" || cfg.ClientSecret == "" || cfg.RedirectURI == "" {
		log.Fatal("Client secret, ID and redirect URI are necessary to resolve codes to tokens!")
	}

	clientID = cfg.ClientID
	clientSecret = cfg.ClientSecret
	redirectURI = cfg.RedirectURI

	router.Get("/", index)
	router.Get("/code", code)
}
