package handlers

import (
	"ecommerce/database"
	"ecommerce/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// ProductListHandler maneja la lista de productos
func ProductListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := models.GetAllProducts(database.DB)
		if err != nil {
			http.Error(w, "Error al obtener productos", http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.ParseFiles("templates/products.html"))
		tmpl.ExecuteTemplate(w, "products", products)
	}
}

// ProductSearchHandler maneja la búsqueda de productos
func ProductSearchHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("search")
		products, err := models.SearchProducts(database.DB, query)
		if err != nil {
			http.Error(w, "Error al buscar productos", http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.ParseFiles("templates/products.html"))
		tmpl.ExecuteTemplate(w, "products", products)
	}
}

func ProductFilterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productType := r.URL.Query().Get("type")
		log.Println("Filtrando productos por tipo:", productType)

		products, err := models.GetProductsByType(database.DB, productType)
		if err != nil {
			http.Error(w, "Error al buscar productos", http.StatusInternalServerError)
			return
		}

		log.Println(products)

		tmpl := template.Must(template.ParseFiles("templates/products.html"))
		tmpl.ExecuteTemplate(w, "products", products)
	}
}

// ProductDetailHandler maneja la página de detalles de un producto
func ProductDetailHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtener el ID del producto desde la URL
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			http.Error(w, "ID de producto no válido", http.StatusBadRequest)
			return
		}

		product, err := models.GetProductByID(database.DB, id)
		if err != nil {
			http.Error(w, "Error al obtener producto", http.StatusInternalServerError)
			return
		}
		if product == nil {
			http.NotFound(w, r)
			return
		}

		tmpl := template.Must(template.ParseFiles("templates/product_detail.html"))
		tmpl.Execute(w, product)
	}
}
