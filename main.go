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

	// Rutas para productos
	http.HandleFunc("/products", handlers.ProductListHandler())
	http.HandleFunc("/products/search", handlers.ProductSearchHandler())
	http.HandleFunc("/product", handlers.ProductDetailHandler())
	http.HandleFunc("/product/create", handlers.ProductCreateHandler())

	// Rutas para el carrito (protegidas por autenticación)
	http.HandleFunc("/cart", middleware.AuthMiddleware(handlers.CartHandler()))
	http.HandleFunc("/cart/add", middleware.AuthMiddleware(handlers.AddToCartHandler()))
	http.HandleFunc("/cart/remove", middleware.AuthMiddleware(handlers.RemoveFromCartHandler()))

	// Inicio del servidor
	log.Println("Servidor iniciado en :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
