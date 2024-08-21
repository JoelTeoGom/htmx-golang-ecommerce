package models

import (
	"database/sql"
	"time"
)

// Invoice representa una factura
type Invoice struct {
	ID        int
	UserID    int
	Total     float64
	CreatedAt time.Time
	Items     []InvoiceItem
}

// InvoiceItem representa un producto dentro de una factura
type InvoiceItem struct {
	ID        int
	InvoiceID int
	ProductID int
	Quantity  int
	Price     float64
	Product   Product
}

// InvoiceData estructura que encapsula la factura y los elementos de la factura.
type InvoiceData struct {
	Invoice Invoice
	Items   []InvoiceItem
}

// Crear una nueva factura
func CreateInvoice(db *sql.DB, userID int, total float64) (int, error) {
	query := `
        INSERT INTO invoices (user_id, total)
        VALUES ($1, $2)
        RETURNING id
    `
	var invoiceID int
	err := db.QueryRow(query, userID, total).Scan(&invoiceID)
	if err != nil {
		return 0, err
	}
	return invoiceID, nil
}

// Añadir un producto a la factura
func AddInvoiceItem(db *sql.DB, invoiceID int, productID int, quantity int, price float64) error {
	query := `
        INSERT INTO invoice_items (invoice_id, product_id, quantity, price)
        VALUES ($1, $2, $3, $4)
    `
	_, err := db.Exec(query, invoiceID, productID, quantity, price)
	return err
}

// Obtener todas las facturas de un usuario
func GetInvoicesByUserID(db *sql.DB, userID int) ([]Invoice, error) {
	query := `
        SELECT id, user_id, total, created_at
        FROM invoices
        WHERE user_id = $1
    `
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []Invoice
	for rows.Next() {
		var invoice Invoice
		err := rows.Scan(&invoice.ID, &invoice.UserID, &invoice.Total, &invoice.CreatedAt)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}
	return invoices, nil
}

// Obtener todos los productos de una factura
func GetInvoiceItems(db *sql.DB, invoiceID int) ([]InvoiceItem, error) {
	query := `
        SELECT invoice_items.id, invoice_items.invoice_id, invoice_items.product_id, invoice_items.quantity, invoice_items.price, products.name, products.image_url
        FROM invoice_items
        JOIN products ON invoice_items.product_id = products.id
        WHERE invoice_items.invoice_id = $1
    `
	rows, err := db.Query(query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []InvoiceItem
	for rows.Next() {
		var item InvoiceItem
		err := rows.Scan(&item.ID, &item.InvoiceID, &item.ProductID, &item.Quantity, &item.Price, &item.Product.Name, &item.Product.ImageURL)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// GetInvoiceByID obtiene una factura por su ID
func GetInvoiceByID(db *sql.DB, invoiceID int) (*Invoice, error) {
	query := `
		SELECT invoices.id, invoices.user_id, invoices.total, invoices.created_at
		FROM invoices
		WHERE invoices.id = $1
	`

	row := db.QueryRow(query, invoiceID)

	var invoice Invoice
	err := row.Scan(&invoice.ID, &invoice.UserID, &invoice.Total, &invoice.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No existe la factura
		}
		return nil, err
	}

	// Obtener los productos de la factura
	itemsQuery := `
		SELECT invoice_items.id, invoice_items.product_id, invoice_items.quantity, invoice_items.price, products.name, products.image_url
		FROM invoice_items
		JOIN products ON invoice_items.product_id = products.id
		WHERE invoice_items.invoice_id = $1
	`

	rows, err := db.Query(itemsQuery, invoice.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []InvoiceItem
	for rows.Next() {
		var item InvoiceItem
		err := rows.Scan(&item.ID, &item.ProductID, &item.Quantity, &item.Price, &item.Product.Name, &item.Product.ImageURL)
		if err != nil {
			return nil, err
		}
		item.InvoiceID = invoice.ID
		items = append(items, item)
	}

	invoice.Items = items
	return &invoice, nil
}
