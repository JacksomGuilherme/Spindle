package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/configs"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/infra/database"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/dao"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/entity"
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/utils"
)

func GetUserPlaylists(limit, offset int, user *entity.User, userDB database.UserInterface, config *configs.Config) *dao.SpotifyAPIResponse {
	url := fmt.Sprintf("https://api.spotify.com/v1/me/playlists?limit=%d&offset=%d", limit, offset)
	req, err := http.NewRequest("GET", url, nil)
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

	return &apiResponse
}

func GetUserFollowedArtists(limit int, after string, user *entity.User, userDB database.UserInterface, config *configs.Config) *dao.SpotifyAPIResponse {
	url := fmt.Sprintf("https://api.spotify.com/v1/me/following?type=artist&limit=%d", limit)

	if after != "" {
		url += "&after=" + after
	}

	req, err := http.NewRequest("GET", url, nil)
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

	var apiArtistsResponse dao.SpotifyArtistsAPIResponse

	json.NewDecoder(resp.Body).Decode(&apiArtistsResponse)

	return &dao.SpotifyAPIResponse{
		Items:   apiArtistsResponse.Artists.Items,
		Cursors: apiArtistsResponse.Artists.Cursors,
	}
}

func GetUserSavedAlbums(limit, offset int, user *entity.User, userDB database.UserInterface, config *configs.Config) *dao.SpotifyAPIResponse {
	url := fmt.Sprintf("https://api.spotify.com/v1/me/albums?limit=%d&offset=%d", limit, offset)
	req, err := http.NewRequest("GET", url, nil)
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

	var apiAlbumResponse dao.SpotifyAlbumAPIResponse

	json.NewDecoder(resp.Body).Decode(&apiAlbumResponse)

	var items []dao.SpotifyItem

	for _, a := range apiAlbumResponse.Albums {
		items = append(items, a.Album)
	}

	return &dao.SpotifyAPIResponse{
		Next:  apiAlbumResponse.Next,
		Items: items,
	}
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

func GetCurrentPlayingSong(user *entity.User, userDB database.UserInterface, config *configs.Config) (*dao.SpotifyPlaybackStateAPIResponse, error) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/currently-playing", nil)
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

	var apiResponse dao.SpotifyPlaybackStateAPIResponse

	json.NewDecoder(resp.Body).Decode(&apiResponse)

	return &apiResponse, nil
}

func PlayContext(playRequest dao.PlaybackRequest, deviceID, sessionID string, userDB database.UserInterface) error {
	bodyBytes, err := json.Marshal(playRequest)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.spotify.com/v1/me/player/play?device_id=%s", deviceID)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	user, err := userDB.FindBySpotifyUserId(sessionID)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+user.AccessToken)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func PausePlayback(deviceID, sessionID string, userDB database.UserInterface) error {
	url := fmt.Sprintf("https://api.spotify.com/v1/me/player/pause?device_id=%s", deviceID)

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}

	user, err := userDB.FindBySpotifyUserId(sessionID)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+user.AccessToken)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func SkipToNextSong(deviceID, sessionID string, userDB database.UserInterface) error {
	url := fmt.Sprintf("https://api.spotify.com/v1/me/player/next?device_id=%s", deviceID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	user, err := userDB.FindBySpotifyUserId(sessionID)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+user.AccessToken)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func SkipToPreviousSong(deviceID, sessionID string, userDB database.UserInterface) error {
	url := fmt.Sprintf("https://api.spotify.com/v1/me/player/previous?device_id=%s", deviceID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	user, err := userDB.FindBySpotifyUserId(sessionID)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+user.AccessToken)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
