package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

const (
	DB_PASSWORD = "1234"
	DB_NAME     = "postgres"
)

func createProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product Product
		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = db.QueryRow("INSERT INTO product(name, price) VALUES($1, $2) RETURNING id", product.Name, product.Price).Scan(&product.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(product)
	}
}

func getProducts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, price FROM product")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var products []Product
		for rows.Next() {
			var product Product
			err := rows.Scan(&product.ID, &product.Name, &product.Price)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			products = append(products, product)
		}
		json.NewEncoder(w).Encode(products)
	}
}

func getProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}
		var product Product
		err = db.QueryRow("SELECT id, name, price FROM product WHERE id = $1", id).Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Product not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		json.NewEncoder(w).Encode(product)
	}
}

func updateProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}
		var product Product
		err = json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		product.ID = id
		_, err = db.Exec("UPDATE product SET name=$1, price=$2 WHERE id=$3", product.Name, product.Price, product.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(product)
	}
}

func deleteProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}
		_, err = db.Exec("DELETE FROM product WHERE id = $1", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func main() {
	connStr := fmt.Sprintf("password=%s dbname=%s sslmode=disable", DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/products", createProduct(db)).Methods("POST")
	r.HandleFunc("/products", getProducts(db)).Methods("GET")
	r.HandleFunc("/products/{id}", getProduct(db)).Methods("GET")
	r.HandleFunc("/products/{id}", updateProduct(db)).Methods("PUT")
	r.HandleFunc("/products/{id}", deleteProduct(db)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
