package database

import (
	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (u *UserRepository) Create(user *entity.User) error {
	return u.DB.Create(user).Error
}

func (u *UserRepository) Update(user *entity.User) error {
	return u.DB.Save(user).Error
}

func (u *UserRepository) Delete(sessionID int) error {
	user, err := u.FindBySessionId(sessionID)
	if err != nil {
		return err
	}
	return u.DB.Delete(user).Error
}

func (u *UserRepository) FindBySpotifyUserId(spotifyUserId string) (*entity.User, error) {
	var user entity.User
	if err := u.DB.Where("spotify_user_id = ?", spotifyUserId).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) FindBySessionId(sessionID int) (*entity.User, error) {
	var user entity.User
	if err := u.DB.Where("id = ?", sessionID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
