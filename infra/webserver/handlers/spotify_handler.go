package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/JacksomGuilherme/Spindle/configs"
	"github.com/JacksomGuilherme/Spindle/infra/database"
	"github.com/JacksomGuilherme/Spindle/internal/entity"
	"github.com/JacksomGuilherme/Spindle/internal/utils"
)

type SpotifyLoginHandler struct {
	PairingStore *utils.PairingStore
	Config       *configs.Config
	UserDB       database.UserInterface
}

func NewSpotifyLoginHandler(pairingStore *utils.PairingStore, config *configs.Config, userDB database.UserInterface) *SpotifyLoginHandler {
	return &SpotifyLoginHandler{
		PairingStore: pairingStore,
		Config:       config,
		UserDB:       userDB,
	}
}

func (h *SpotifyLoginHandler) Auth(w http.ResponseWriter, r *http.Request) {
	pairingID := r.URL.Query().Get("pairing_id")

	scope := "user-library-read user-read-private user-read-email user-read-playback-state user-modify-playback-state user-read-currently-playing app-remote-control playlist-read-private user-read-recently-played user-follow-read"
	url := fmt.Sprintf("https://accounts.spotify.com/authorize?response_type=code&client_id=%s&scope=%s&redirect_uri=%s:%s/auth/spotify/callback&state=%s",
		h.Config.AppClientID,
		scope,
		h.Config.Dns,
		h.Config.WebServerPort,
		pairingID,
	)

	http.Redirect(w, r, url, 302)
}

func (h *SpotifyLoginHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" {
		utils.ExecutarTemplate(w, "callback.html", map[string]interface{}{
			"Error": "Missing code",
		})
		return
	}

	if state == "" || !h.PairingStore.IsValid(state) {
		utils.ExecutarTemplate(w, "callback.html", map[string]interface{}{
			"Error": "Invalid or expired session",
		})
		return
	}

	token, err := utils.ExchangeCodeForToken(code, h.Config)
	if err != nil {
		utils.ExecutarTemplate(w, "callback.html", map[string]interface{}{
			"Error": "Failed to authenticate with Spotify",
		})
		return
	}

	userID, err := getSpotifyUser(token.AccessToken)
	if err != nil || userID == "" {
		utils.ExecutarTemplate(w, "callback.html", map[string]interface{}{
			"Error": "Spotify user not found",
		})
		return
	}

	h.PairingStore.Authenticate(state, userID)
	h.PairingStore.Consume(state)

	expiresAt := time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	encryptedAccess, _ := utils.Encrypt(token.AccessToken, h.Config.EncryptionKey)
	encryptedRefresh, _ := utils.Encrypt(token.RefreshToken, h.Config.EncryptionKey)

	existingUser, _ := h.UserDB.FindBySpotifyUserId(userID)
	if existingUser == nil {
		user := entity.NewUser(userID, encryptedAccess, encryptedRefresh, expiresAt)
		h.UserDB.Create(user)
	} else {
		existingUser.AccessToken = encryptedAccess
		existingUser.RefreshToken = encryptedRefresh
		existingUser.ExpiresAt = expiresAt
		h.UserDB.Update(existingUser)
	}

	utils.ExecutarTemplate(w, "callback.html", map[string]interface{}{
		"Success": true,
	})
}

func (h *SpotifyLoginHandler) Status(w http.ResponseWriter, r *http.Request) {
	pairingID := r.URL.Query().Get("state")

	p, ok := h.PairingStore.Get(pairingID)
	if !ok {
		http.Error(w, "not found", 404)
		return
	}

	if p.Authenticated {
		user, err := h.UserDB.FindBySpotifyUserId(p.UserID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = utils.SalvarCookie(w, strconv.Itoa(user.ID))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.PairingStore.Delete(pairingID)

		w.Write([]byte("ok"))
		return
	}

	w.Write([]byte("pending"))
}

func getSpotifyUser(accessToken string) (string, error) {
	req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data struct {
		ID string `json:"id"`
	}

	json.NewDecoder(resp.Body).Decode(&data)

	return data.ID, nil
}
