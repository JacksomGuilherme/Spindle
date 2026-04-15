package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JacksomGuilherme/Spindle/configs"
	"github.com/JacksomGuilherme/Spindle/infra/database"
	"github.com/JacksomGuilherme/Spindle/internal/dao"
	"github.com/JacksomGuilherme/Spindle/internal/entity"
	"github.com/JacksomGuilherme/Spindle/internal/services"
	"github.com/JacksomGuilherme/Spindle/internal/utils"
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

	id, _ := strconv.Atoi(cookie["session_id"])
	user, err := h.UserDB.FindBySessionId(id)
	if user == nil || err != nil {
		utils.DeletarCookie(w)
		http.Redirect(w, r, "/login", 302)
		return
	}

	pageParam := r.URL.Query().Get("page")

	page := 0
	if pageParam != "" {
		p, _ := strconv.Atoi(pageParam)
		page = p
	}

	limit := 12
	offset := page * limit

	content := getContent("playlists", limit, offset, "", user, h.UserDB, h.Config)
	if content == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userDevices, err := services.GetUserDevices(user, h.UserDB, h.Config)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var activeDevice *dao.Device

	for _, device := range userDevices {
		if device.IsActive {
			activeDevice = &device
			break
		}
	}

	if activeDevice == nil && len(userDevices) > 0 {
		activeDevice = &userDevices[0]
	}

	utils.ExecutarTemplate(w, "home.html", map[string]interface{}{
		"Content":         content.Items,
		"Page":            page,
		"DisplayPage":     page + 1,
		"NextCursor":      content.Cursors.After,
		"HasNext":         content.Next != "",
		"HasPrevious":     page > 0,
		"ActiveDevice":    activeDevice,
		"DeviceConnected": activeDevice != nil,
		"Tab":             "playlists",
	})
}

func (h *HomeHandler) TabContent(w http.ResponseWriter, r *http.Request) {
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

	tab := r.URL.Query().Get("tab")
	if tab == "" {
		tab = "playlists"
	}

	pageParam := r.URL.Query().Get("page")
	after := r.URL.Query().Get("after")

	page := 0
	if pageParam != "" {
		p, _ := strconv.Atoi(pageParam)
		page = p
	}

	limit := 12
	offset := page * limit

	content := getContent(tab, limit, offset, after, user, h.UserDB, h.Config)
	if content == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	utils.ExecutarTemplate(w, "tab-content", map[string]interface{}{
		"Content":     content.Items,
		"Page":        page,
		"DisplayPage": page + 1,
		"NextCursor":  content.Cursors.After,
		"HasNext":     content.Next != "" || content.Cursors.After != "",
		"HasPrevious": page > 0 || tab == "artists",
		"Tab":         tab,
	})
}

func (h *HomeHandler) ActiveDevice(w http.ResponseWriter, r *http.Request) {
	cookie, _ := utils.LerCookie(r)

	id, _ := strconv.Atoi(cookie["session_id"])
	user, err := h.UserDB.FindBySessionId(id)
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

	var activeDevice *dao.Device

	for _, device := range userDevices {
		if device.IsActive {
			activeDevice = &device
			break
		}
	}

	if activeDevice == nil && len(userDevices) > 0 {
		activeDevice = &userDevices[0]
	}

	json.NewEncoder(w).Encode(activeDevice)
}

func getContent(tab string, limit, offset int, after string, user *entity.User, userDB database.UserInterface, config *configs.Config) *dao.SpotifyAPIResponse {
	switch tab {
	case "playlists":
		return services.GetUserPlaylists(limit, offset, user, userDB, config)
	case "artists":
		return services.GetUserFollowedArtists(limit, after, user, userDB, config)
	case "albums":
		return services.GetUserSavedAlbums(limit, offset, user, userDB, config)
	}
	return nil
}
