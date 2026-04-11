package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/configs"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/infra/database"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/entity"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/utils"
	"github.com/skip2/go-qrcode"
)

type LoginHandler struct {
	PairingStore *utils.PairingStore
	Config       *configs.Config
	UserDB       database.UserInterface
}

func NewLoginHandler(pairingStore *utils.PairingStore, config *configs.Config, userDB database.UserInterface) *LoginHandler {
	return &LoginHandler{
		PairingStore: pairingStore,
		Config:       config,
		UserDB:       userDB,
	}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	pairingID := entity.NewID().String()

	h.PairingStore.Create(pairingID)
	url := fmt.Sprintf("%s:%s/auth/spotify/login?pairing_id=%s", h.Config.Dns, h.Config.WebServerPort, pairingID)

	utils.ExecutarTemplate(w, "login.html", map[string]interface{}{
		"PairingID": pairingID,
		"AuthUrl":   url,
	})
}

func (h *LoginHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, _ := utils.LerCookie(r)
	sessionID := cookie["session_id"]
	utils.DeletarCookie(w)

	id, _ := strconv.Atoi(sessionID)
	h.UserDB.Delete(id)
	http.Redirect(w, r, "/login", 302)
}

func (h *LoginHandler) GetQRCode(w http.ResponseWriter, r *http.Request) {
	pairingID := r.URL.Query().Get("pairing_id")

	url := fmt.Sprintf("%s:%s/auth/spotify/login?pairing_id=%s", h.Config.Dns, h.Config.WebServerPort, pairingID)

	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(png)
}
