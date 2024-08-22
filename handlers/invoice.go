package handlers

import (
	"ecommerce/database"
	"ecommerce/models"
	"html/template"
	"net/http"
)

// Función para manejar el checkout y generar una factura
func GetAllInvoiceByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the user ID from the session
		userID, err := GetUserIDFromCookie(r)
		if err != nil || userID == 0 {
			http.Error(w, "Failed to retrieve user from session", http.StatusUnauthorized)
			return
		}

		// Fetch all invoices for the user from the database
		invoices, err := models.GetInvoicesByUserID(database.DB, userID)
		if err != nil {
			http.Error(w, "Failed to retrieve invoices", http.StatusInternalServerError)
			return
		}

		// Create a slice to hold the invoice data structures to pass to the template
		var allInvoicesData []struct {
			Invoice models.Invoice
			Items   []models.InvoiceItem
		}

		// Loop through each invoice and get its items
		for _, invoice := range invoices {
			items, err := models.GetInvoiceItems(database.DB, invoice.ID)
			if err != nil {
				http.Error(w, "Failed to retrieve invoice items", http.StatusInternalServerError)
				return
			}

			// Append the invoice and its items to the slice
			allInvoicesData = append(allInvoicesData, struct {
				Invoice models.Invoice
				Items   []models.InvoiceItem
			}{
				Invoice: invoice,
				Items:   items,
			})
		}

		// Render the template with all invoices data
		tmpl := template.Must(template.ParseFiles("templates/invoices.html"))
		err = tmpl.Execute(w, allInvoicesData)
		if err != nil {
			http.Error(w, "Error rendering invoice template", http.StatusInternalServerError)
		}
	}
}
