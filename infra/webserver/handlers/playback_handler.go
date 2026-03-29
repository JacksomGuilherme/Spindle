package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/infra/database"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/dao"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/utils"
)

type PlaybackHandler struct {
	UserDB database.UserInterface
}

func NewPlaybackHandler(userDB database.UserInterface) *PlaybackHandler {
	return &PlaybackHandler{
		UserDB: userDB,
	}
}

func (h *PlaybackHandler) Play(w http.ResponseWriter, r *http.Request) {
	deviceId := r.URL.Query().Get("device_id")

	var playReq dao.PlaybackRequest
	err := json.NewDecoder(r.Body).Decode(&playReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bodyBytes, err := json.Marshal(playReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("https://api.spotify.com/v1/me/player/play?device_id=%s", deviceId)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cookie, _ := utils.LerCookie(r)
	user, err := h.UserDB.FindBySpotifyUserId(cookie["session_id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+user.AccessToken)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
