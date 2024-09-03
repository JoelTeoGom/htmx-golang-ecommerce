# HTMX Golang E-commerce

## Description

This project is a simple e-commerce application developed in Go (Golang) using HTMX for frontend interactivity. The project includes basic features such as product management, shopping cart, user authentication, and invoice generation.

## Features

- **User Authentication**: User registration and login.
- **Product Management**: Listing, creating, updating, and deleting products.
- **Shopping Cart**: Add and remove products from the cart.
- **Invoicing**: Generation and management of invoices.
- **HTMX Interactivity**: Implementing frontend interactivity without writing JavaScript.

## Project Structure

- `main.go`: Main application file.
- `handlers/`: Contains the controllers for handling different application routes (`auth.go`, `cart.go`, `home.go`, `invoice.go`, `product.go`).
- `models/`: Defines the data models (`cart.go`, `invoice.go`, `product.go`, `user.go`).
- `middleware/`: Contains the middleware for the application (`middleware.go`).
- `static/`: Static files like images.
- `database/`: Database connection and management (`db.go`).

## Requirements

- **Golang**: Version 1.16 or higher.
- **Docker**: To run the application in a container.

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/yourusername/htmx-golang-ecommerce.git
