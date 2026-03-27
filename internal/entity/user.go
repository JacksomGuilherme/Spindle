package entity

type User struct {
	ID            int
	SpotifyUserId string
	RefreshToken  string
	ExpiresIn     int
}

func NewUser(spotifyUserId, refreshToken string, expiresIn int) *User {
	return &User{
		SpotifyUserId: spotifyUserId,
		RefreshToken:  refreshToken,
		ExpiresIn:     expiresIn,
	}
}
