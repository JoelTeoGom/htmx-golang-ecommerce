package handlers

import (
	"ecommerce/database"
	"ecommerce/models"
	"html/template"
	"net/http"
)

// Función para manejar el checkout y generar una factura
func CheckoutAndCreateInvoice() http.HandlerFunc {
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

		// Crear la factura
		invoiceID, err := models.CreateInvoice(database.DB, userID, total)
		if err != nil {
			http.Error(w, "Error al crear la factura", http.StatusInternalServerError)
			return
		}

		// Añadir los productos a la factura
		for _, item := range cartItems {
			err = models.AddInvoiceItem(database.DB, invoiceID, item.ProductID, item.Quantity, item.Product.Price)
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

		// Obtener la factura recién creada y los productos asociados
		invoice, err := models.GetInvoicesByUserID(database.DB, userID)
		if err != nil {
			http.Error(w, "Error al obtener la factura", http.StatusInternalServerError)
			return
		}

		items, err := models.GetInvoiceItems(database.DB, invoiceID)
		if err != nil {
			http.Error(w, "Error al obtener los productos de la factura", http.StatusInternalServerError)
			return
		}

		// Renderizar la plantilla de la factura
		data := struct {
			Invoice models.Invoice
			Items   []models.InvoiceItem
		}{
			Invoice: invoice[0],
			Items:   items,
		}

		tmpl := template.Must(template.ParseFiles("templates/invoice.html"))
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Error al renderizar la plantilla de la factura", http.StatusInternalServerError)
		}
	}
}
