package handlers

import (
	"ecommerce/database"
	"ecommerce/models"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Message struct {
	Type    string
	Content string
}

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func generateToken(userID int, username string) (string, error) {
	// Crear el token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userID,
		"username": username,
		"exp":      time.Now().Add(1 * time.Hour).Unix(),
	})

	// Firmar el token con la clave secreta
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func setTokenCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   os.Getenv("NODE_ENV") == "production",
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(1 * time.Hour),
	}

	http.SetCookie(w, cookie)
}

func RegisterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl := template.Must(template.ParseFiles("templates/register.html"))
			tmpl.Execute(w, nil)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		// Verificar si el usuario ya existe
		user, err := models.GetUserByUsername(database.DB, username)
		tmpl := template.Must(template.ParseFiles("templates/message.html"))

		if err == nil && user != nil {
			// Usuario ya existe
			tmpl.Execute(w, &Message{
				Type:    "error",
				Content: "El nombre de usuario ya está registrado",
			})
			return
		}

		// Hashear la contraseña
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			tmpl.Execute(w, &Message{
				Type:    "error",
				Content: "Error al crear el usuario. Inténtalo de nuevo.",
			})
			return
		}

		// Crear el nuevo usuario
		err = models.InsertUser(database.DB, username, hashedPassword)
		if err != nil {
			tmpl.Execute(w, &Message{
				Type:    "error",
				Content: "Error al registrar el usuario. Inténtalo de nuevo.",
			})
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			tmpl := template.Must(template.ParseFiles("templates/login.html"))
			tmpl.Execute(w, nil)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		tmpl := template.Must(template.ParseFiles("templates/message.html"))

		user, err := models.GetUserByUsername(database.DB, username)
		if err != nil || user == nil {
			tmpl.Execute(w, &Message{
				Type:    "error",
				Content: "Nombre de usuario incorrecto",
			})
			return
		}

		// Verificar la contraseña con bcrypt
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			tmpl.Execute(w, &Message{
				Type:    "error",
				Content: "Contraseña incorrecta",
			})
			return
		}

		// Generar el token JWT
		token, err := generateToken(user.ID, user.Username)
		if err != nil {
			tmpl.Execute(w, &Message{
				Type:    "error",
				Content: "Error al generar la sesión. Inténtalo de nuevo.",
			})
			return
		}

		// Configurar la cookie con el token
		setTokenCookie(w, token)

		// Inicio de sesión exitoso
		tmpl.Execute(w, &Message{
			Type:    "success",
			Content: "Inicio de sesión exitoso. Bienvenido, " + template.HTMLEscapeString(username) + "!",
		})
	}
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verificar si la cookie de token existe
		cookie, err := r.Cookie("token")
		if err != nil {
			// Si no hay cookie, redirigir a login
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Verificar el token JWT
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err == nil && token.Valid {
			// Si el token es válido, redirigir a la página principal
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		// Si el token es inválido, redirigir a login
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{
			Name:     "token",
			Value:    "",
			HttpOnly: true,
			Secure:   os.Getenv("NODE_ENV") == "production",
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
			Expires:  time.Now(),
		}

		http.SetCookie(w, cookie)

		// Redirigir al login
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}
