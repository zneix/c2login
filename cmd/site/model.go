package main

type TokenResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    uint64   `json:"expires_in"`
	Scope        []string `json:"scope,omitempty"`
	TokenType    string   `json:"token_type"`
}
