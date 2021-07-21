package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/nicklaw5/helix"
	"github.com/zneix/c2login/pkg/config"
)

var (
	APIClient   *helix.Client
	clientID    string
	redirectURI string
)

// index page - user authentication
func index(w http.ResponseWriter, r *http.Request) {
	// "code" is required for us to receive code instead of token
	const responseType = "code"
	encodedScopes := url.PathEscape(strings.Join(scopes, " "))

	link := fmt.Sprintf("https://id.twitch.tv/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=%s&scope=%s&force_verify=true", clientID, redirectURI, responseType, encodedScopes)

	http.Redirect(w, r, link, http.StatusFound)
}

// code page - callback called by Twitch, where we process Twitch's request and redirect to the final page
func code(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No code provided..."))
		return
	}

	// Request access token for given code
	resp, err := APIClient.RequestUserAccessToken(code)
	if err != nil {
		log.Printf("Error while requesting user access token: %v\n", err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error..."))
		return
	}

	// Requesting token was unsuccessful
	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		w.Write([]byte(fmt.Sprintf("Something went wrong while requesting access token, %d", resp.StatusCode)))
		return
	}

	// Validate obtained token for extra data (no way it can be invalid)
	_, validResp, err := APIClient.ValidateToken(resp.Data.AccessToken)

	// Success
	payloadData := chatterinoLoginPayload{
		ClientID:     validResp.Data.ClientID,
		OauthToken:   resp.Data.AccessToken,
		RefreshToken: resp.Data.RefreshToken,
		UserID:       validResp.Data.UserID,
		Username:     validResp.Data.Login,
	}
	payload, err := json.Marshal(payloadData)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
}

func handleMainRoutes(router *chi.Mux, helixClient *helix.Client, cfg config.SiteConfig) {
	APIClient = helixClient
	clientID = cfg.ClientID
	redirectURI = cfg.RedirectURI

	router.Get("/", index)
	router.Get("/code", code)
}
