package handlers

import (
	"html/template"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

// HomeHandler maneja la página principal
func HomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var username string
		var isAuthenticated bool

		// Verificar si hay una cookie válida
		cookie, err := r.Cookie("token")
		if err == nil {
			// Si la cookie existe, intentar parsear el token
			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

			if err == nil && token.Valid {
				// Obtener los claims del token (como el nombre de usuario)
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					username = claims["username"].(string)
					isAuthenticated = true
				}
			}
		}

		// Cargar la plantilla
		tmpl := template.Must(template.ParseFiles("templates/home.html"))
		// Pasar los datos a la plantilla
		data := struct {
			IsAuthenticated bool
			Username        string
		}{
			IsAuthenticated: isAuthenticated,
			Username:        username,
		}

		tmpl.Execute(w, data)
	}
}
