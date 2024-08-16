package main

import (
	"ecommerce/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var products = []models.Product{
	{ID: 1, Name: "Laptop", Description: "Una laptop potente para trabajo y juegos.", Price: 999.99, ImageURL: "/static/laptop.jpg"},
	{ID: 2, Name: "Smartphone", Description: "Un smartphone con cámara de alta resolución.", Price: 499.99, ImageURL: "/static/smartphone.jpg"},
	{ID: 3, Name: "Auriculares", Description: "Auriculares con cancelación de ruido.", Price: 199.99, ImageURL: "/static/headphones.jpg"},
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("search")
	filteredProducts := filterProducts(searchQuery)

	var tmplPath string
	if r.Header.Get("HX-Request") != "" {
		tmplPath = filepath.Join("templates", "product_list.html")
	} else {
		tmplPath = filepath.Join("templates", "index.html")
	}

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Error al cargar la plantilla", http.StatusInternalServerError)
		return
	}

	data := struct {
		Products []models.Product
	}{
		Products: filteredProducts,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error al renderizar la plantilla", http.StatusInternalServerError)
	}
}

func filterProducts(query string) []models.Product {
	if query == "" {
		return products
	}

	var filtered []models.Product
	for _, product := range products {
		if strings.Contains(strings.ToLower(product.Name), strings.ToLower(query)) {
			filtered = append(filtered, product)
		}
	}
	return filtered
}
