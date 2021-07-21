package config

import (
	"os"
	"reflect"
	"strings"
)

const (
	envPrefix = "C2LOGIN_"
)

var (
	envReplacer = strings.NewReplacer("-", "_")
)

type SiteConfig struct {
	// Core

	// Base URL which is meant to be the root endpoint of the site accessible from outside
	BaseURL string `key:"base-url"`
	// Address on which application will be listening, localhost
	BindAddress string `key:"bind-address"`

	// Secrets

	// Twitch client id used to identify the client secret
	ClientID string `key:"client-id"`
	// Twitch client secret used for obtaining tokens from received codes
	ClientSecret string `key:"client-secret"`
	// Redirect URI to which Twitch will send us the code after user successfully authenticates
	RedirectURI string `key:"redirect-uri"`
}

func getEnv(envSuffix string) (value string, exists bool) {

	key := strings.ToUpper(envReplacer.Replace(envPrefix + envSuffix))
	return os.LookupEnv(key)
}

func New() (cfg SiteConfig) {
	// Default config
	cfg = SiteConfig{
		BindAddress: ":1234",
	}

	v := reflect.ValueOf(&cfg).Elem()

	for i := 0; i < v.NumField(); i++ {
		key := v.Type().Field(i).Tag.Get("key")
		value, exists := getEnv(key)
		if exists && len(value) > 0 {
			v.Field(i).SetString(value)
		}
	}

	//fmt.Printf("%# v\n", cfg) // uncomment for debugging purposes
	return
}
