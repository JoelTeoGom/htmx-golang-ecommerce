package handlers

import (
	"ecommerce/database"
	"ecommerce/models"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
)

// Función para obtener el ID del usuario desde la cookie JWT
func GetUserIDFromCookie(r *http.Request) (int, error) {
	// Obtener la cookie
	cookie, err := r.Cookie("token")
	if err != nil {
		return 0, err
	}

	// Parsear el token
	tokenStr := cookie.Value
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	// Verificar si el token es válido
	if err != nil || !token.Valid {
		return 0, err
	}

	// Obtener el ID del usuario desde los claims del token
	userID, ok := (*claims)["id"].(float64)
	if !ok {
		return 0, err
	}

	return int(userID), nil
}

func CartHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetUserIDFromCookie(r)
		if err != nil || userID == 0 {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		cartItems, err := models.GetCartItems(database.DB, userID)
		if err != nil {
			http.Error(w, "Error al obtener el carrito", http.StatusInternalServerError)
			return
		}

		log.Println(cartItems)

		user, err := models.GetUserById(database.DB, userID)
		if err != nil {
			http.Error(w, "Error al obtener el user", http.StatusInternalServerError)
			return
		}

		username := user.Username
		var total float64
		for _, item := range cartItems {
			total += item.Product.Price * float64(item.Quantity)
		}

		data := struct {
			IsAuthenticated bool
			Username        string
			CartItems       []models.Cart
			Total           float64
		}{
			Username:        username,
			CartItems:       cartItems,
			Total:           total,
			IsAuthenticated: true,
		}

		tmpl := template.Must(template.ParseFiles("templates/cart.html"))
		tmpl.Execute(w, data)
	}
}

func AddToCartHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetUserIDFromCookie(r)
		if err != nil || userID == 0 {
			http.Error(w, "No se pudo obtener el usuario de la sesión", http.StatusUnauthorized)
			return
		}

		productID, _ := strconv.Atoi(r.FormValue("id"))

		log.Println(productID)
		err = models.AddToCart(database.DB, userID, productID)

		if err != nil {
			http.Error(w, "No se pudo agregar el producto al carrito", http.StatusInternalServerError)
			return
		}

		CartHandler().ServeHTTP(w, r) // Refrescar el carrito después de añadir el producto
	}
}

func RemoveFromCartHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetUserIDFromCookie(r)
		if err != nil || userID == 0 {
			http.Error(w, "No se pudo obtener el usuario de la sesión", http.StatusUnauthorized)
			return
		}

		productID, _ := strconv.Atoi(r.FormValue("id"))

		err = models.RemoveFromCart(database.DB, userID, productID)
		if err != nil {
			http.Error(w, "No se pudo eliminar el producto del carrito", http.StatusInternalServerError)
			return
		}

	}
}
