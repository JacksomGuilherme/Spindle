package database

import "github.com/JacksomGuilherme/Spindle/internal/entity"

type UserInterface interface {
	Create(user *entity.User) error
	Update(user *entity.User) error
	Delete(sessionID int) error
	FindBySpotifyUserId(spotifyUserId string) (*entity.User, error)
	FindBySessionId(sessionID int) (*entity.User, error)
}
