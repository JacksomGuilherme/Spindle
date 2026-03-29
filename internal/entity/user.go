package entity

import "time"

type User struct {
	ID            int
	SpotifyUserId string
	AccessToken   string
	RefreshToken  string
	ExpiresAt     time.Time
}

func NewUser(spotifyUserId, accesToken, refreshToken string, expiresIn time.Time) *User {
	return &User{
		SpotifyUserId: spotifyUserId,
		AccessToken:   accesToken,
		RefreshToken:  refreshToken,
		ExpiresAt:     expiresIn,
	}
}
