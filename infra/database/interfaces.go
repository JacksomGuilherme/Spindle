package database

import "github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/entity"

type UserInterface interface {
	Create(user *entity.User) error
	FindBySpotifyUserId(spotifyUserId string) (*entity.User, error)
}
