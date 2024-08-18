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
	http.HandleFunc("/", handlers.LoginPage)
	http.HandleFunc("/login", handlers.LoginHandler())
	http.HandleFunc("/register", handlers.RegisterHandler())

	// Inicio del servidor
	log.Println("Servidor iniciado en :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
