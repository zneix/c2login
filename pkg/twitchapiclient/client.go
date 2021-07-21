package twitchapiclient

import (
	"errors"
	"net/http"

	"github.com/nicklaw5/helix"
	"github.com/zneix/c2login/pkg/config"
)

// New returns a helix.Client that will be used for requesting user access tokens
func New(cfg config.SiteConfig, httpClient *http.Client) (*helix.Client, error) {
	if cfg.ClientID == "" {
		return nil, errors.New("client-id is missing, can't make Twitch requests")
	}
	if cfg.ClientSecret == "" {
		return nil, errors.New("client-secret is missing, can't make Twitch requests")
	}
	if cfg.RedirectURI == "" {
		return nil, errors.New("redirect-uri is missing, can't make Twitch requests")
	}

	apiClient, err := helix.NewClient(&helix.Options{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURI:  cfg.RedirectURI,
		HTTPClient:   httpClient,
	})
	if err != nil {
		return nil, err
	}

	return apiClient, nil
}
