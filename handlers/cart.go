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
		log.Println(productID)
		err = models.RemoveFromCart(database.DB, userID, productID)
		if err != nil {
			http.Error(w, "No se pudo eliminar el producto del carrito", http.StatusInternalServerError)
			return
		}
		w.Header().Set("HX-Trigger-After-Swap", "updateCart")
	}
}

//CARTHOME

// AddToCartHomeHandler agrega un producto al carrito y devuelve la vista del carrito para el home
func AddToCartHomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetUserIDFromCookie(r)
		if err != nil || userID == 0 {
			http.Error(w, "No se pudo obtener el usuario de la sesión", http.StatusUnauthorized)
			return
		}

		productID, _ := strconv.Atoi(r.FormValue("id"))
		err = models.AddToCart(database.DB, userID, productID)
		if err != nil {
			http.Error(w, "No se pudo agregar el producto al carrito", http.StatusInternalServerError)
			return
		}

		//creara un evento updateCart
		w.Header().Set("HX-Trigger-After-Swap", "updateCart")
		// Ahora renderizamos el carrito para el home
		RenderCartHome(w, r, userID)
	}
}

// RenderCartHome renderiza el carrito simplificado para el home
func RenderCartHome(w http.ResponseWriter, r *http.Request, userID int) {
	cartItems, err := models.GetCartItems(database.DB, userID)
	if err != nil {
		http.Error(w, "Error al obtener el carrito", http.StatusInternalServerError)
		return
	}

	var total float64
	for _, item := range cartItems {
		total += item.Product.Price * float64(item.Quantity)
	}

	data := struct {
		CartItems []models.Cart
		Total     float64
	}{
		CartItems: cartItems,
		Total:     total,
	}

	tmpl := template.Must(template.ParseFiles("templates/cartHome.html"))
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error al renderizar la plantilla del carrito", http.StatusInternalServerError)
	}
}

// CheckoutCartHome maneja el checkout desde el carrito y limpia el carrito del usuario
func CheckoutCartHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetUserIDFromCookie(r)
		if err != nil || userID == 0 {
			http.Error(w, "No se pudo obtener el usuario de la sesión", http.StatusUnauthorized)
			return
		}

		// Obtener los items del carrito antes de limpiar
		cartItems, err := models.GetCartItems(database.DB, userID)
		if err != nil {
			http.Error(w, "Error al obtener el carrito", http.StatusInternalServerError)
			return
		}

		// Calcular el total
		var total float64
		for _, item := range cartItems {
			total += item.Product.Price * float64(item.Quantity)
		}

		// Crear una nueva factura
		invoiceID, err := models.CreateInvoice(database.DB, userID, total)
		if err != nil {
			http.Error(w, "Error al crear la factura", http.StatusInternalServerError)
			return
		}

		// Añadir los productos a la factura
		for _, item := range cartItems {
			err := models.AddInvoiceItem(database.DB, invoiceID, item.ProductID, item.Quantity, item.Product.Price)
			if err != nil {
				http.Error(w, "Error al añadir productos a la factura", http.StatusInternalServerError)
				return
			}
		}

		// Limpiar el carrito
		err = models.ClearCart(database.DB, userID)
		if err != nil {
			http.Error(w, "Error al limpiar el carrito", http.StatusInternalServerError)
			return
		}

		// Obtener la factura recién creada
		invoice, err := models.GetInvoiceByID(database.DB, invoiceID)
		if err != nil {
			http.Error(w, "Error al obtener la factura", http.StatusInternalServerError)
			return
		}
		inv := *invoice

		invoiceItems, err := models.GetInvoiceItems(database.DB, invoiceID)
		if err != nil {
			http.Error(w, "Error al obtener los items de la factura", http.StatusInternalServerError)
			return
		}

		log.Println(invoice)
		invoiceData := models.InvoiceData{
			Invoice: inv,          // Suponiendo que tienes un objeto `invoice` del tipo `models.Invoice`
			Items:   invoiceItems, // Suponiendo que tienes una lista de `invoiceItems` del tipo `[]models.InvoiceItem`
		}

		// Determinar la plantilla a utilizar basado en la URL solicitada
		var tmpl *template.Template
		if r.URL.Path == "/cart/checkoutHome" {
			w.Header().Set("HX-Trigger-After-Swap", "updateCart")
			tmpl = template.Must(template.ParseFiles("templates/invoice.html"))
		} else {
			log.Println("sdfahjdshjfdsahjkfdshjkadfhsfhdjks")
			tmpl = template.Must(template.ParseFiles("templates/invoice_cart.html"))

		}

		// Renderizar la plantilla seleccionada
		err = tmpl.Execute(w, invoiceData)
		if err != nil {
			http.Error(w, "Error al renderizar la plantilla de la factura", http.StatusInternalServerError)
		}
	}
}

// CartCountHandler maneja la solicitud para obtener el número de ítems en el carrito
func CartCountHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtener el ID del usuario desde la cookie JWT
		userID, err := GetUserIDFromCookie(r)
		if err != nil || userID == 0 {
			// Si no hay usuario autenticado, devolver 0 como conteo del carrito
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("0"))
			return
		}

		// Obtener el número de ítems en el carrito desde la base de datos
		itemCount, err := models.GetCartItemCount(database.DB, userID)
		if err != nil {
			http.Error(w, "Error al obtener el conteo de ítems en el carrito", http.StatusInternalServerError)
			return
		}

		// Devolver el número de ítems como texto simple
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.Itoa(itemCount)))
	}
}
