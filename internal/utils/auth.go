package utils

import (
	"encoding/json"
	"fmt"
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
	data.Set("redirect_uri", fmt.Sprintf("%s:%s/auth/spotify/callback", config.Dns, config.WebServerPort))

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
	accessToken, err := Decrypt(user.AccessToken, config.EncryptionKey)
	if err != nil {
		return "", err
	}

	if time.Now().Before(user.ExpiresAt) {
		return accessToken, nil
	}

	refreshToken, err := Decrypt(user.RefreshToken, config.EncryptionKey)
	if err != nil {
		return "", err
	}

	token, err := refreshAccessToken(refreshToken, config)
	if err != nil {
		return "", err
	}

	encryptedAccess, _ := Encrypt(token.AccessToken, config.EncryptionKey)
	user.AccessToken = encryptedAccess
	user.ExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	if token.RefreshToken != "" {
		encryptedRefresh, _ := Encrypt(token.RefreshToken, config.EncryptionKey)
		user.RefreshToken = encryptedRefresh
	}

	userDB.Update(user)
	return token.AccessToken, nil
}

func refreshAccessToken(refreshToken string, config *configs.Config) (*dao.TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(config.AppClientID, config.AppClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var token dao.TokenResponse
	json.NewDecoder(resp.Body).Decode(&token)

	return &token, nil
}
