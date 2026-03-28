package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/configs"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/infra/database"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/dao"
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

	scope := "user-read-private user-read-email"
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

	token, err := exchangeCodeForToken(code, h.Config)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, _ := getSpotifyUser(token.AccessToken)
	h.PairingStore.Authenticate(state, userID)

	user := entity.NewUser(userID, token.RefreshToken, token.ExpiresIn)

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

func exchangeCodeForToken(code string, config *configs.Config) (*dao.TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", "http://127.0.0.1:8080/auth/spotify/callback")

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	req.SetBasicAuth(config.AppClientID, config.AppClientSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var token dao.TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
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
