package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/configs"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/infra/database"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/entity"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/utils"
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

	scope := "user-read-private user-read-email user-read-playback-state user-modify-playback-state user-read-currently-playing app-remote-control playlist-read-private user-read-recently-played user-follow-read"
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

	if state == "" {
		http.Redirect(w, r, "/?error=state_mismatch", http.StatusFound)
		return
	}

	token, err := utils.ExchangeCodeForToken(code, h.Config)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, _ := getSpotifyUser(token.AccessToken)
	h.PairingStore.Authenticate(state, userID)

	expiresAt := time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	user := entity.NewUser(userID, token.AccessToken, token.RefreshToken, expiresAt)

	h.UserDB.Create(user)

	utils.ExecutarTemplate(w, "callback", nil)
}

func (h *SpotifyLoginHandler) Status(w http.ResponseWriter, r *http.Request) {
	pairingID := r.URL.Query().Get("state")

	p, ok := h.PairingStore.Get(pairingID)
	if !ok {
		http.Error(w, "not found", 404)
		return
	}

	if p.Authenticated {
		err := utils.SalvarCookie(w, p.UserID)
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
