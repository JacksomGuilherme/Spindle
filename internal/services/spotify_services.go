package services

import (
	"encoding/json"
	"net/http"

	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/configs"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/infra/database"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/dao"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/entity"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/utils"
)

func GetUserPlaylists(user *entity.User, userDB database.UserInterface, config *configs.Config) []dao.SpotifyItem {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/playlists", nil)
	if err != nil {
		return nil
	}
	accessToken, err := utils.GetValidAccessToken(user, userDB, config)
	if err != nil {
		return nil
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var apiResponse dao.SpotifyAPIResponse

	json.NewDecoder(resp.Body).Decode(&apiResponse)

	return apiResponse.Items
}

func GetUserFollowedArtists(user *entity.User, userDB database.UserInterface, config *configs.Config) []dao.SpotifyItem {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/following?type=artist", nil)
	if err != nil {
		return nil
	}
	accessToken, err := utils.GetValidAccessToken(user, userDB, config)
	if err != nil {
		return nil
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var apiResponse dao.SpotifyArtistsAPIResponse

	json.NewDecoder(resp.Body).Decode(&apiResponse)

	return apiResponse.Artists.Items
}

func GetUserSavedAlbums(user *entity.User, userDB database.UserInterface, config *configs.Config) []dao.SpotifyItem {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/albums", nil)
	if err != nil {
		return nil
	}
	accessToken, err := utils.GetValidAccessToken(user, userDB, config)
	if err != nil {
		return nil
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var apiResponse dao.SpotifyAlbumAPIResponse

	json.NewDecoder(resp.Body).Decode(&apiResponse)

	var items []dao.SpotifyItem

	for _, a := range apiResponse.Albums {
		items = append(items, a.Album)
	}

	return items
}

func GetUserDevices(user *entity.User, userDB database.UserInterface, config *configs.Config) ([]dao.Device, error) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/devices", nil)
	if err != nil {
		return nil, err
	}

	accessToken, err := utils.GetValidAccessToken(user, userDB, config)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResponse dao.DevicesAPIResponse

	json.NewDecoder(resp.Body).Decode(&apiResponse)

	return apiResponse.Devices, nil
}
