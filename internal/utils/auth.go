package utils

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/configs"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/infra/database"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/dao"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/entity"
)

func ExchangeCodeForToken(code string, config *configs.Config) (*dao.TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", "https://192.168.15.2:8080/auth/spotify/callback")

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

func GetValidAccessToken(user *entity.User, userDB database.UserInterface, config *configs.Config) (string, error) {
	if time.Now().Before(user.ExpiresAt) {
		return user.AccessToken, nil
	}

	token, err := refreshAccessToken(user.RefreshToken, config)
	if err != nil {
		return "", err
	}

	user.AccessToken = token.AccessToken
	user.ExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	userDB.Update(user)

	return user.AccessToken, nil
}

func refreshAccessToken(refreshToken string, config *configs.Config) (*dao.TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	req.SetBasicAuth(config.AppClientID, config.AppClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	var token dao.TokenResponse
	json.NewDecoder(resp.Body).Decode(&token)

	return &token, nil
}
