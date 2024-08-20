package models

import "database/sql"

type Cart struct {
	ID        int
	UserID    int
	ProductID int
	Quantity  int
	Product   Product
}

// Obtener todos los artículos en el carrito de un usuario
func GetCartItems(db *sql.DB, userID int) ([]Cart, error) {
	query := `
        SELECT carts.id, carts.user_id, carts.product_id, carts.quantity, products.name, products.price, products.image_url 
        FROM carts
        JOIN products ON carts.product_id = products.id
        WHERE carts.user_id = $1
    `

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cartItems []Cart
	for rows.Next() {
		var cartItem Cart
		var product Product
		err := rows.Scan(&cartItem.ID, &cartItem.UserID, &cartItem.ProductID, &cartItem.Quantity, &product.Name, &product.Price, &product.ImageURL)
		if err != nil {
			return nil, err
		}
		cartItem.Product = product
		cartItems = append(cartItems, cartItem)
	}
	return cartItems, nil
}

// Añadir un producto al carrito
func AddToCart(db *sql.DB, userID int, productID int) error {
	query := `
        INSERT INTO carts (user_id, product_id, quantity)
        VALUES ($1, $2, 1)
        ON CONFLICT (user_id, product_id)
        DO UPDATE SET quantity = carts.quantity + 1
    `
	_, err := db.Exec(query, userID, productID)
	return err
}

// Eliminar un producto del carrito
func RemoveFromCart(db *sql.DB, userID int, productID int) error {
	query := `
        DELETE FROM carts
        WHERE user_id = $1 AND product_id = $2
    `
	_, err := db.Exec(query, userID, productID)
	return err
}

// Vaciar el carrito de un usuario
func ClearCart(db *sql.DB, userID int) error {
	query := `
        DELETE FROM carts
        WHERE user_id = $1
    `
	_, err := db.Exec(query, userID)
	return err
}
