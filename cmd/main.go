package main

import (
	"encoding/json"
	"net/http"
	"proyect/go-web/internal/repository"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func main() {
	Products := repository.Charge_Products()
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
	//Post that creates a new product
	router.Post("/products", func(w http.ResponseWriter, r *http.Request) {
		//Id is not necesary but if is passed it should be unique
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		var newProduct Products_Struct
		err := json.NewDecoder(r.Body).Decode(&newProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Invalid product"))
			return
		}
		//Any field can be empty except ID and Is_Published
		if !bool(newProduct.Is_Published) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Fields Missing"))
			return
		}
		//Check if ID is unique
		for _, product := range Products {
			if product.ID == newProduct.ID {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("400 - ID already exists"))
				return
			}
		}
		//Check if Code_Value is unique
		for _, product := range Products {
			if product.Code_Value == newProduct.Code_Value {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("400 - Code_Value already exists"))
				return
			}
		}
		//Field expiration must be in format XX/XX/XXXX
		//Check if expiration is in format XX/XX/XXXX
		if len(newProduct.Expiration) != 10 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Expiration must be in format XX/XX/XXXX"))
			return
		}
		if newProduct.Expiration[2] != '/' || newProduct.Expiration[5] != '/' {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Expiration must be in format XX/XX/XXXX"))
			return
		}
		//Types of fields must be equal to the struct
		//Check if ID is int
		if newProduct.ID != int(newProduct.ID) && newProduct.ID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Invalid type of ID"))
			return
		}
		//Check if Quantity is int
		if newProduct.Quantity != int(newProduct.Quantity) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Invalid type of Quantity"))
			return
		}
		//Check if Price is float64
		if newProduct.Price != float64(newProduct.Price) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Invalid type of Price"))
			return
		}
		//Check if Is_Published is bool
		if newProduct.Is_Published != bool(newProduct.Is_Published) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Invalid type of Is_Published"))
			return
		}
		//Check if Name is string
		if newProduct.Name != string(newProduct.Name) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Invalid type of Name"))
			return
		}
		//Check if Code_Value is string
		if newProduct.Code_Value != string(newProduct.Code_Value) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Invalid type of Code_Value"))
			return
		}
		//Check if Expiration is string
		if newProduct.Expiration != string(newProduct.Expiration) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Invalid type of Expiration"))
			return
		}
		//If ID is 0, assign a new ID
		if newProduct.ID == 0 {
			newProduct.ID = len(Products) + 1
		}
		//Add new product to slice
		Products = append(Products, newProduct)
		//Show 201 - Product created and the new product
		w.Write([]byte("201 - Product created"))
		json.NewEncoder(w).Encode(newProduct)

	})
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
