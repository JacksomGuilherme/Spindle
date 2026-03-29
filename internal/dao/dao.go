package dao

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type SpotifyItemImages struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type SpotifyItem struct {
	ID         string              `json:"id"`
	Name       string              `json:"name"`
	Images     []SpotifyItemImages `json:"images"`
	ContextURI string              `json:"uri"`
}

type SpotifyAPIResponse struct {
	Items []SpotifyItem `json:"items"`
}

type SpotifyArtist struct {
	Items []SpotifyItem `json:"items"`
}

type SpotifyAlbum struct {
	Album SpotifyItem `json:"album"`
}

type SpotifyArtistsAPIResponse struct {
	Artists SpotifyArtist `json:"artists"`
}

type SpotifyAlbumAPIResponse struct {
	Albums []SpotifyAlbum `json:"items"`
}

type Device struct {
	ID       string `json:"id"`
	IsActive bool   `json:"is_active"`
	Name     string `json:"name"`
}

type DevicesAPIResponse struct {
	Devices []Device `json:"devices"`
}

type PlaybackRequest struct {
	ContextURI string `json:"context_uri"`
}
