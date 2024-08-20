package handlers

import (
	"ecommerce/database"
	"ecommerce/models"
	"html/template"
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

// ProductCreateHandler maneja la creación de un nuevo producto (ejemplo básico)
func ProductCreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl := template.Must(template.ParseFiles("templates/product_create.html"))
			tmpl.Execute(w, nil)
			return
		}

		// Obtener datos del formulario
		name := r.FormValue("name")
		description := r.FormValue("description")
		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			http.Error(w, "Precio no válido", http.StatusBadRequest)
			return
		}
		imageURL := r.FormValue("image_url")

		// Insertar producto en la base de datos
		product := models.Product{
			Name:        name,
			Description: description,
			Price:       price,
			ImageURL:    imageURL,
		}

		_, err = models.InsertProduct(database.DB, product)
		if err != nil {
			http.Error(w, "Error al crear el producto", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/products", http.StatusSeeOther)
	}
}
