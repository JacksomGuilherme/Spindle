package handlers

import (
	"net/http"

	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/configs"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/infra/database"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/entity"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/services"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/utils"
)

type HomeHandler struct {
	Config *configs.Config
	UserDB database.UserInterface
}

func NewHomeHandler(config *configs.Config, userDB database.UserInterface) *HomeHandler {
	return &HomeHandler{
		Config: config,
		UserDB: userDB,
	}
}

func (h *HomeHandler) Home(w http.ResponseWriter, r *http.Request) {
	tab := r.URL.Query().Get("tab")
	if tab == "" {
		tab = "playlists"
	}

	cookie, _ := utils.LerCookie(r)

	if cookie["session_id"] == "" {
		http.Redirect(w, r, "/login", 302)
		return
	}
	user, err := h.UserDB.FindBySpotifyUserId(cookie["session_id"])
	if user == nil || err != nil {
		utils.DeletarCookie(w)
		http.Redirect(w, r, "/login", 302)
		return
	}

	userDevices, err := services.GetUserDevices(user, h.UserDB, h.Config)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var content interface{}

	content = getContent(tab, user, h.UserDB, h.Config)

	utils.ExecutarTemplate(w, "home", map[string]interface{}{
		"Content":         content,
		"Devices":         userDevices,
		"DeviceConnected": true,
		"ActiveTab":       tab,
	})
}

func getContent(tab string, user *entity.User, userDB database.UserInterface, config *configs.Config) interface{} {
	switch tab {
	case "playlists":
		return services.GetUserPlaylists(user, userDB, config)
	case "artists":
		return services.GetUserFollowedArtists(user, userDB, config)
	case "albums":
		return services.GetUserSavedAlbums(user, userDB, config)
	}
	return nil
}
