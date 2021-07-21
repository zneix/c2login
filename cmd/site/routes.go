package main

import (
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

	// Get details about user tied to the token
	userClient, err := helix.NewClient(&helix.Options{
		ClientID:        clientID,
		UserAccessToken: resp.Data.AccessToken,
	})

	userResp, err := userClient.GetUsers(&helix.UsersParams{})
	if err != nil {
		log.Printf("Error while requesting user details: %v\n", err)

		w.WriteHeader(userResp.StatusCode)
		w.Write([]byte(fmt.Sprintf("Something went wrong while requesting user details, %d", userResp.StatusCode)))
		return
	}

	user := userResp.Data.Users[0]

	// Success
	accountData := fmt.Sprintf("oauth_token=%s;refresh_token=%s;username=%s;user_id=%s;client_id=%s", resp.Data.AccessToken, resp.Data.RefreshToken, user.Login, user.ID, clientID)
	w.Write([]byte(accountData))
}

func handleMainRoutes(router *chi.Mux, helixClient *helix.Client, cfg config.SiteConfig) {
	APIClient = helixClient
	clientID = cfg.ClientID
	redirectURI = cfg.RedirectURI

	router.Get("/", index)
	router.Get("/code", code)
}
