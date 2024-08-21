package models

import (
	"database/sql"
	"log"
)

// Product representa un producto electrónico
type Product struct {
	ID          int
	Name        string
	Type        string // Nuevo campo para el tipo de producto
	Description string
	Price       float64
	ImageURL    string
}

// SearchProducts busca productos por nombre o descripción
func SearchProducts(db *sql.DB, query string) ([]Product, error) {
	query = "%" + query + "%"
	rows, err := db.Query("SELECT id, name, type, description, price, image_url FROM products WHERE name ILIKE $1", query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Type, &product.Description, &product.Price, &product.ImageURL); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

// GetAllProducts obtiene todos los productos de la base de datos
func GetAllProducts(db *sql.DB) ([]Product, error) {
	query := "SELECT id, name, type, description, price, image_url FROM products"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Type, &product.Description, &product.Price, &product.ImageURL)
		if err != nil {
			log.Println("Error al escanear producto:", err)
			continue
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// GetProductByID obtiene un producto por su ID
func GetProductByID(db *sql.DB, id int) (*Product, error) {
	query := "SELECT id, name, type, description, price, image_url FROM products WHERE id = $1"
	row := db.QueryRow(query, id)

	var product Product
	err := row.Scan(&product.ID, &product.Name, &product.Type, &product.Description, &product.Price, &product.ImageURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Producto no encontrado
		}
		return nil, err
	}

	return &product, nil
}

// InsertProduct inserta un nuevo producto en la base de datos
func InsertProduct(db *sql.DB, product Product) (int, error) {
	query := "INSERT INTO products (name, type, description, price, image_url) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	var id int
	err := db.QueryRow(query, product.Name, product.Type, product.Description, product.Price, product.ImageURL).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetProductsByType(db *sql.DB, productType string) ([]Product, error) {
	query := "SELECT id, name, type, description, price, image_url FROM products WHERE type = $1"
	rows, err := db.Query(query, productType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Type, &product.Description, &product.Price, &product.ImageURL)
		if err != nil {
			log.Println("Error al escanear producto:", err)
			continue
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
