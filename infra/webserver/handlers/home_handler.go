package handlers

import (
	"net/http"

	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/internal/utils"
)

func Home(w http.ResponseWriter, r *http.Request) {
	cookie, _ := utils.LerCookie(r)

	if cookie["session_id"] == "" {
		http.Redirect(w, r, "/login", 302)
		return
	}

	utils.ExecutarTemplate(w, "home", nil)
}
