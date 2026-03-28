package utils

import (
	"net/http"
	"time"

	"github.com/JacksomGuilherme/Kindle-Spotify-Controller/configs"
	"github.com/gorilla/securecookie"
)

var s *securecookie.SecureCookie

func ConfigurarCookies(config *configs.Config) {
	s = securecookie.New([]byte(config.HashKey), []byte(config.BlockKey))
}

func SalvarCookie(w http.ResponseWriter, userID string) error {
	dados := map[string]string{
		"session_id": userID,
	}

	dadosCodificados, erro := s.Encode("dados", dados)
	if erro != nil {
		return erro
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "dados",
		Value: dadosCodificados,
		Path:  "/",
	})

	return nil
}

func LerCookie(r *http.Request) (map[string]string, error) {
	cookie, erro := r.Cookie("dados")
	if erro != nil {
		return nil, erro
	}

	valores := make(map[string]string)
	if erro = s.Decode("dados", cookie.Value, &valores); erro != nil {
		return nil, erro
	}

	return valores, nil

}

func DeletarCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "dados",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})
}
