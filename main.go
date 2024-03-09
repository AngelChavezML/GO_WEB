package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Define a struct
type Products_Struct struct {
	ID           int
	Name         string
	Quantity     int
	Code_Value   string
	Is_Published bool
	Expiration   string
	Price        float64
}

func Charge_Products() []Products_Struct {
	//Read JSON file
	file, err := os.Open("products.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	//Decode JSON file
	var products []Products_Struct
	err = json.NewDecoder(file).Decode(&products)
	if err != nil {
		log.Fatal(err)
	}
	return products
}
func main() {
	Products := Charge_Products()
	//Show Products
	//fmt.Println(Products)
	//Begin server
	router := chi.NewRouter()
	//Get that shows pong and 200 OK
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})
	//Get that shows all products
	router.Get("/products", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Products)
	})
	//Get that shows a product by ID
	router.Get("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		id := chi.URLParam(r, "id")
		productID, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Invalid product ID"))
			return
		}
		for _, product := range Products {
			if product.ID == productID {
				json.NewEncoder(w).Encode(product)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Product not found"))
	})
	//Get that seach products that their price is higher than a given price
	router.Get("/products/search/{priceGt}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		priceGt := chi.URLParam(r, "priceGt")
		priceGtFloat, err := strconv.ParseFloat(priceGt, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Invalid price"))
			return
		}
		var productsFound []Products_Struct
		for _, product := range Products {
			if product.Price > priceGtFloat {
				productsFound = append(productsFound, product)
			}
		}
		json.NewEncoder(w).Encode(productsFound)
	})
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
