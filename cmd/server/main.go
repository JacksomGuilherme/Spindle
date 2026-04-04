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
	utils.ConfigurarCookies(config)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	utils.CarregarTemplates()

	db, err := database.GetConnection(config)
	if err != nil {
		panic(err)
	}

	userDB := database.NewUserRepository(db)

	pairingStore := utils.NewPairingStore()
	loginHandler := handlers.NewLoginHandler(pairingStore, config, userDB)

	homeHandler := handlers.NewHomeHandler(config, userDB)

	r.Route("/", func(r chi.Router) {
		r.Get("/", homeHandler.Home)
		r.Get("/tab/content", homeHandler.TabContent)
		r.Get("/login", loginHandler.Login)
		r.Get("/logout", loginHandler.Logout)
		r.Get("/login/qr", loginHandler.GetQRCode)
	})

	spotifyLogin := handlers.NewSpotifyLoginHandler(pairingStore, config, userDB)
	r.Route("/auth", func(r chi.Router) {
		r.Get("/spotify/login", spotifyLogin.Auth)
		r.Get("/spotify/callback", spotifyLogin.Callback)
		r.Get("/spotify/login/status", spotifyLogin.Status)
	})

	playbackHander := handlers.NewPlaybackHandler(config, userDB)

	r.Route("/playback", func(r chi.Router) {
		r.Post("/play", playbackHander.Play)
		r.Post("/pause", playbackHander.Pause)
		r.Post("/next", playbackHander.Next)
		r.Post("/previous", playbackHander.Previous)
		r.Get("/state", playbackHander.State)
	})

	fs := http.FileServer(http.Dir("../../website/assets"))
	r.Handle("/assets/*", http.StripPrefix("/assets/", fs))

	go pairingStore.Cleanup()
	http.ListenAndServeTLS(
		fmt.Sprintf(":%s", config.WebServerPort),
		"localhost.pem",
		"localhost-key.pem",
		r,
	)
}
