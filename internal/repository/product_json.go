package repository

import (
	"encoding/json"
	"log"
	"os"
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
