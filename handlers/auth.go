package handlers

import (
	"ecommerce/database"
	"ecommerce/models"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Message struct {
	Type    string
	Content string
}

// LoginPage muestra la página de inicio de sesión.
func LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(w, nil)
}

// LoginHandler maneja la lógica de autenticación.
func LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		tmpl := template.Must(template.ParseFiles("templates/message.html"))

		user, err := models.GetUserByUsername(database.DB, username)
		if err != nil || user == nil {
			tmpl.ExecuteTemplate(w, "message", &Message{
				Type:    "error",
				Content: "Nombre de usuario incorrecto",
			})
			return
		}

		// Verificar la contraseña con bcrypt
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			tmpl.ExecuteTemplate(w, "message", &Message{
				Type:    "error",
				Content: "Contraseña incorrecta",
			})
			return
		}

		tmpl.ExecuteTemplate(w, "message", &Message{
			Type:    "success",
			Content: "Inicio de sesión exitoso. Bienvenido, " + template.HTMLEscapeString(username) + "!",
		})
	}
}
