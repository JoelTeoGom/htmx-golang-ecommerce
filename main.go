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
	http.HandleFunc("/login", middleware.LogAuthMiddleware(handlers.LoginHandler()))
	http.HandleFunc("/register", middleware.LogAuthMiddleware(handlers.RegisterHandler()))
	http.HandleFunc("/", handlers.HomeHandler())
	http.HandleFunc("/logout", handlers.LogoutHandler())

	// Rutas para productos
	http.HandleFunc("/products", handlers.ProductListHandler())
	http.HandleFunc("/products/search", handlers.ProductSearchHandler())
	http.HandleFunc("/product", handlers.ProductDetailHandler())
	//http.HandleFunc("/product/create", handlers.ProductCreateHandler())
	http.HandleFunc("/products/filter", handlers.ProductFilterHandler())

	// Rutas para el carrito (protegidas por autenticación)
	http.HandleFunc("/cart", middleware.AuthMiddleware(handlers.CartHandler()))
	http.HandleFunc("/cart/add", middleware.AuthMiddleware(handlers.AddToCartHandler()))
	http.HandleFunc("/cart/remove", middleware.AuthMiddleware(handlers.RemoveFromCartHandler()))

	http.HandleFunc("/cart/checkout", middleware.AuthMiddleware(handlers.CheckoutCartHome()))
	http.HandleFunc("/cart/checkoutHome", middleware.AuthMiddleware(handlers.CheckoutCartHome()))

	http.HandleFunc("/cart/addHome", middleware.AuthMiddleware(handlers.AddToCartHomeHandler()))

	http.HandleFunc("/invoices", middleware.AuthMiddleware(handlers.GetAllInvoiceByUser()))

	http.HandleFunc("/api/cart-count", middleware.AuthMiddleware(handlers.CartCountHandler()))

	// Inicio del servidor
	log.Println("Servidor iniciado en :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
