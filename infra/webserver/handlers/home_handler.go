package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/configs"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/infra/database"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/dao"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/entity"
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
	cookie, _ := utils.LerCookie(r)

	if cookie["session_id"] == "" {
		http.Redirect(w, r, "/login", 302)
		return
	}
	user, err := h.UserDB.FindBySpotifyUserId(cookie["session_id"])
	if user == nil || err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	userDevices, err := getUserDevices(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userPlaylists, err := getUserPlaylists(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.ExecutarTemplate(w, "home", map[string]interface{}{
		"Playlists":       userPlaylists,
		"Devices":         userDevices,
		"DeviceConnected": len(userDevices) > 0,
	})
}

func getUserPlaylists(user *entity.User) ([]dao.PlaylistsResponse, error) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/playlists", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+user.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResponse dao.PlaylistsAPIResponse

	json.NewDecoder(resp.Body).Decode(&apiResponse)

	return apiResponse.Items, nil
}

func getUserDevices(user *entity.User) ([]dao.Device, error) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/devices", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+user.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResponse dao.DevicesAPIResponse

	json.NewDecoder(resp.Body).Decode(&apiResponse)

	return apiResponse.Devices, nil
}
