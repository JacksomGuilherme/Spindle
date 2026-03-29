package dao

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type PlaylistImagesResponse struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type PlaylistsResponse struct {
	ID         string                   `json:"id"`
	Name       string                   `json:"name"`
	ImageURL   []PlaylistImagesResponse `json:"images"`
	ContextURI string
}

type PlaylistsAPIResponse struct {
	Items []PlaylistsResponse `json:"items"`
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
