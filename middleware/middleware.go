package middleware

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Leer la cookie "token"
		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			// Si no hay cookie o está vacía, continuar con la solicitud
			next.ServeHTTP(w, r)
			return
		}

		// Verificar y analizar el token JWT
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			// Si el token es inválido, continuar con la solicitud
			next.ServeHTTP(w, r)
			return
		}

		// Solo redirigir si el usuario intenta acceder a login o register y ya está autenticado
		if r.URL.Path == "/login" || r.URL.Path == "/register" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Si ya está en la página principal (/), simplemente continuar
		next.ServeHTTP(w, r)
	}
}
