package main

import (
	"ecommerce/database"
	"ecommerce/handlers"
	"log"
	"net/http"
)

func main() {
	// Inicializar la conexión a la base de datos
	database.InitDB()
	defer database.DB.Close()

	// Manejar archivos estáticos
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Configuración de rutas

	// Configuración de rutas con middleware de autenticación
	http.HandleFunc("/", handlers.AuthMiddleware(handlers.HomeHandler()))
	http.HandleFunc("/register", handlers.AuthMiddleware(handlers.RegisterHandler()))
	http.HandleFunc("/login", handlers.AuthMiddleware(handlers.LoginHandler()))
	http.HandleFunc("/logout", handlers.LogoutHandler())

	// Inicio del servidor
	log.Println("Servidor iniciado en :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
