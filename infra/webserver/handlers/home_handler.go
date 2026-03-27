package handlers

import (
	"net/http"

	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/utils"
)

func Home(w http.ResponseWriter, r *http.Request) {
	utils.ExecutarTemplate(w, "home", nil)
}
