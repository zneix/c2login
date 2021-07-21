package main

type chatterinoLoginPayload struct {
	OauthToken   string `json:"oauthToken"`
	RefreshToken string `json:"refreshToken"`
	Username     string `json:"username"`
	UserID       string `json:"userID"`
	ClientID     string `json:"clientID"`
}
