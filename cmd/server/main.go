package main

import (
	"fmt"
	"net/http"

	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/configs"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/infra/database"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/infra/webserver/handlers"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	config := configs.LoadConfig(".")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	utils.CarregarTemplates()

	pairingStore := utils.NewPairingStore()
	loginHandler := handlers.NewLoginHandler(pairingStore, config)

	r.Route("/login", func(r chi.Router) {
		r.Get("/", loginHandler.Login)
		r.Get("/qr", loginHandler.GetQRCode)
	})

	db, err := database.GetConnection(config)
	if err != nil {
		panic(err)
	}

	userDB := database.NewUserRepository(db)
	spotifyLogin := handlers.NewSpotifyLoginHandler(pairingStore, config, userDB)
	r.Route("/auth", func(r chi.Router) {
		r.Get("/spotify/login", spotifyLogin.Auth)
		r.Get("/spotify/callback", spotifyLogin.Callback)
		r.Get("/spotify/login/status", spotifyLogin.Status)
	})

	r.Get("/", handlers.Home)

	go pairingStore.Cleanup()
	http.ListenAndServe(fmt.Sprintf(":%s", config.WebServerPort), r)
}
