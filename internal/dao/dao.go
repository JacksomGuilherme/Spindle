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

type SpotifyPlaybackItem struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	Artists    []SpotifyItem `json:"artists"`
	ContextURI string        `json:"uri"`
}

type SpotifyAPIResponse struct {
	Items   []SpotifyItem `json:"items"`
	Next    string        `json:"next"`
	Cursors SpotifyCursor `json:"cursors"`
}

type SpotifyCursor struct {
	After  string `json:"after"`
	Before string `json:"before"`
}

type SpotifyArtist struct {
	Cursors SpotifyCursor `json:"cursors"`
	Items   []SpotifyItem `json:"items"`
}

type SpotifyAlbum struct {
	Album SpotifyItem `json:"album"`
}

type SpotifyContext struct {
	URI string `json:"uri"`
}

type SpotifyArtistsAPIResponse struct {
	Artists SpotifyArtist `json:"artists"`
}

type SpotifyAlbumAPIResponse struct {
	Next   string         `json:"next"`
	Albums []SpotifyAlbum `json:"items"`
}

type SpotifyPlaybackStateAPIResponse struct {
	Context    SpotifyContext      `json:"context"`
	Item       SpotifyPlaybackItem `json:"item"`
	IsPlaying  bool                `json:"is_playing"`
	ProgressMs int                 `json:"progress_ms"`
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
	ContextURI string      `json:"context_uri"`
	Offset     interface{} `json:"offset,omitempty"`
	PositionMs string      `json:"position_ms"`
}
