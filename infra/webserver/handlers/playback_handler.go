package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/configs"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/infra/database"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/dao"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/services"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/utils"
)

type PlaybackHandler struct {
	Config *configs.Config
	UserDB database.UserInterface
}

func NewPlaybackHandler(config *configs.Config, userDB database.UserInterface) *PlaybackHandler {
	return &PlaybackHandler{
		Config: config,
		UserDB: userDB,
	}
}

func (h *PlaybackHandler) Play(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("device_id")

	var playReq dao.PlaybackRequest
	err := json.NewDecoder(r.Body).Decode(&playReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cookie, _ := utils.LerCookie(r)

	err = services.PlayContext(playReq, deviceID, cookie["session_id"], h.UserDB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, _ := strconv.Atoi(cookie["session_id"])
	user, err := h.UserDB.FindBySessionId(id)
	if user == nil || err != nil {
		utils.DeletarCookie(w)
		http.Redirect(w, r, "/login", 302)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *PlaybackHandler) Pause(w http.ResponseWriter, r *http.Request) {
	cookie, _ := utils.LerCookie(r)
	if cookie["session_id"] == "" {
		http.Redirect(w, r, "/login", 302)
		return
	}

	deviceID := r.URL.Query().Get("device_id")

	err := services.PausePlayback(deviceID, cookie["session_id"], h.UserDB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, _ := strconv.Atoi(cookie["session_id"])
	user, err := h.UserDB.FindBySessionId(id)
	if user == nil || err != nil {
		utils.DeletarCookie(w)
		http.Redirect(w, r, "/login", 302)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *PlaybackHandler) Next(w http.ResponseWriter, r *http.Request) {
	cookie, _ := utils.LerCookie(r)

	if cookie["session_id"] == "" {
		http.Redirect(w, r, "/login", 302)
		return
	}

	deviceID := r.URL.Query().Get("device_id")

	err := services.SkipToNextSong(deviceID, cookie["session_id"], h.UserDB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, _ := strconv.Atoi(cookie["session_id"])
	user, err := h.UserDB.FindBySessionId(id)
	if user == nil || err != nil {
		utils.DeletarCookie(w)
		http.Redirect(w, r, "/login", 302)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *PlaybackHandler) Previous(w http.ResponseWriter, r *http.Request) {
	cookie, _ := utils.LerCookie(r)
	if cookie["session_id"] == "" {
		http.Redirect(w, r, "/login", 302)
		return
	}

	deviceID := r.URL.Query().Get("device_id")

	err := services.SkipToPreviousSong(deviceID, cookie["session_id"], h.UserDB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, _ := strconv.Atoi(cookie["session_id"])
	user, err := h.UserDB.FindBySessionId(id)
	if user == nil || err != nil {
		utils.DeletarCookie(w)
		http.Redirect(w, r, "/login", 302)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *PlaybackHandler) State(w http.ResponseWriter, r *http.Request) {
	cookie, _ := utils.LerCookie(r)

	if cookie["session_id"] == "" {
		http.Redirect(w, r, "/login", 302)
		return
	}

	id, _ := strconv.Atoi(cookie["session_id"])
	user, err := h.UserDB.FindBySessionId(id)
	if user == nil || err != nil {
		utils.DeletarCookie(w)
		http.Redirect(w, r, "/login", 302)
		return
	}

	playBackState, err := services.GetCurrentPlayingSong(user, h.UserDB, h.Config)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(playBackState)
}
