package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zneix/c2login/pkg/config"
)

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("here users will be authenticating"))
}

func code(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("here we will translate codes to tokens AlienPls\nafter successfull translation we will send users back to root with their token and stuff passed via headers"))
}

func handleMainRoutes(router *chi.Mux, cfg config.SiteConfig) {
	router.Get("/", index)
	router.Get("/code", code)
}
