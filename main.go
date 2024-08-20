package main

import (
	"ecommerce/database"
	"ecommerce/handlers"
	"ecommerce/middleware"
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

	// Aplicar el middleware a las rutas protegidas
	http.HandleFunc("/login", middleware.AuthMiddleware(handlers.LoginHandler()))
	http.HandleFunc("/register", middleware.AuthMiddleware(handlers.RegisterHandler()))
	http.HandleFunc("/", handlers.HomeHandler())
	http.HandleFunc("/logout", handlers.LogoutHandler())

	// Inicio del servidor
	log.Println("Servidor iniciado en :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
