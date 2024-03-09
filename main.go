package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

// Define a struct
type Products struct {
	ID           int
	Name         string
	Quantity     int
	Code_Value   string
	Is_Published bool
	Expiration   string
	Price        float64
}

func Charge_Products() []Products {
	//Read JSON file
	file, err := os.Open("products.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	//Decode JSON file
	var products []Products
	err = json.NewDecoder(file).Decode(&products)
	if err != nil {
		log.Fatal(err)
	}
	return products
}
func main() {
	Products := Charge_Products()
	//Show Products
	fmt.Println(Products)
	//Begin server
	router := chi.NewRouter()
	//Get that shows pong and 200 OK
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}

}
