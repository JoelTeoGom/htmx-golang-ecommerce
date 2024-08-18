package models

import (
	"database/sql"
	"log"
)

type User struct {
	ID       int
	Username string
	Password string
}

// GetUserByUsername busca un usuario en la base de datos por su nombre de usuario.
func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	query := "SELECT id, username, password FROM users WHERE username=$1"
	user := &User{}
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Usuario no encontrado
		}
		log.Println("Error al consultar la base de datos:", err)
		return nil, err
	}
	return user, nil
}