package handlers

import (
	"ecommerce/database"
	"ecommerce/models"
	"html/template"
	"log"
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

// RegisterHandler maneja la lógica de registro.
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

			log.Println(user)
			log.Print(err)

			// Usuario ya existe
			tmpl.ExecuteTemplate(w, "message", &Message{
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

		// Crear el nuevo usuario utilizando la función InsertUser
		err = models.InsertUser(database.DB, username, hashedPassword)
		if err != nil {
			tmpl.Execute(w, &Message{
				Type:    "error",
				Content: "Error al registrar el usuario. Inténtalo de nuevo.",
			})
			return
		}

		// Registro exitoso
		tmpl.Execute(w, &Message{
			Type:    "success",
			Content: "Registro completado exitosamente. Bienvenido, " + template.HTMLEscapeString(username) + "!",
		})
	}
}
