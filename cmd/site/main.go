package main

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/zneix/c2login/pkg/config"
	"github.com/zneix/c2login/pkg/twitchapiclient"
)

var (
	listenPrefix string

	httpClient = &http.Client{
		Timeout: 15 * time.Second,
	}
)

func mountRouter(r *chi.Mux, cfg config.SiteConfig) *chi.Mux {
	if cfg.BaseURL == "" {
		log.Printf("Listening on %s (Prefix=%s, BaseURL=%s)\n", cfg.BindAddress, listenPrefix, cfg.BaseURL)
		return r
	}

	// figure out prefix from address
	u, err := url.Parse(cfg.BaseURL)
	if err != nil {
		log.Fatal(err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		log.Fatal("Scheme must be included in base url")
	}

	listenPrefix = u.Path
	// Empty prefix can't be passed to chi
	if listenPrefix == "" {
		listenPrefix = "/"
	}
	ur := chi.NewRouter()
	ur.Mount(listenPrefix, r)

	log.Printf("Listening on %s (Prefix=%s, BaseURL=%s)\n", cfg.BindAddress, listenPrefix, cfg.BaseURL)
	return ur
}

func listen(bind string, router *chi.Mux) {
	srv := &http.Server{
		Handler:      router,
		Addr:         bind,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func main() {
	cfg := config.New()
	router := chi.NewRouter()

	helixClient, err := twitchapiclient.New(cfg, httpClient)
	if err != nil {
		log.Fatalf("[Twitch] Error initializing new client: %v\n", err)
	}

	handleMainRoutes(router, helixClient, cfg)
	listen(cfg.BindAddress, mountRouter(router, cfg))
}
