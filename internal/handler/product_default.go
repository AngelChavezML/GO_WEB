package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"proyect/go-web/internal"
	"proyect/go-web/platform/web"
	"time"
)

// Add the missing import for the internal package

func NewDeafultProduct(task map[int]internal.Product_Struct, lastID int) *DeafultProduct {
	//Default Values
	defaultProducts := make(map[int]internal.Product_Struct)
	defaultLastID := 0
	if task != nil {
		defaultProducts = task
	}
	if lastID != 0 {
		defaultLastID = lastID
	}
	return &NewDeafultProduct{
		products: defaultProducts,
		lastID:   defaultLastID,
	}
}
func Pong(w http.ResponseWriter, r *http.Request) {
	web.ResponseJSON(w, http.StatusOK, map[string]string{"message": "pong"})
}

type ProductsJSON struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	Code_Value   string  `json:"code_value"`
	Is_Published bool    `json:"is_published"`
	Expiration   string  `json:"expiration"`
	Price        float64 `json:"price"`
}
type ProductsRequestBody struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	Code_Value   string  `json:"code_value"`
	Is_Published bool    `json:"is_published"`
	Expiration   string  `json:"expiration"`
	Price        float64 `json:"price"`
}
type DeafultProduct struct {
	product map[int]internal.Product_Struct
	lastID  int
}

func (d *DeafultProduct) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//Request
		// read into bytes
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			web.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "error reading request body"})
			return
		}
		// parse to map (dynamic)

		bodyMap := map[string]interface{}{}
		if err := json.Unmarshal(bytes, &bodyMap); err != nil {
			web.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "error parsing request body"})
			return
		}
		// Validate
		if _, ok := bodyMap["name"]; !ok {
			web.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "name is required"})
			return
		}
		if _, ok := bodyMap["quantity"]; !ok {
			web.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "quantity is required"})
			return
		}
		if _, ok := bodyMap["code_value"]; !ok {
			web.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "code_value is required"})
			return
		}
		if _, ok := bodyMap["expiration"]; !ok {
			web.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "expiration is required"})
			return
		}
		if _, ok := bodyMap["price"]; !ok {
			web.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "price is required"})
			return
		}
		// Check is code_value is unique
		for _, p := range d.product {
			if p.Code_Value == bodyMap["code_value"].(string) {
				web.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "code_value must be unique"})
				return
			}
		}
		// Check if expiration date is in format XX/XX/XXXX
		if _, err := time.Parse("01/02/2006", bodyMap["expiration"].(string)); err != nil {
			web.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "expiration must be in format MM/DD/YYYY"})
			return
		}
		// Parse json to struct (static)
		var product ProductsRequestBody
		if err := json.Unmarshal(bytes, &product); err != nil {
			web.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "error parsing request body"})
			return
		}
		d.lastID++
		products := internal.Product_Struct{
			ID:           d.lastID,
			Name:         product.Name,
			Quantity:     product.Quantity,
			Code_Value:   product.Code_Value,
			Is_Published: product.Is_Published,
			Expiration:   product.Expiration,
			Price:        product.Price,
		}
		//Store
		d.product[d.lastID] = products
		data := ProductsJSON{
			ID:           products.ID,
			Name:         products.Name,
			Quantity:     products.Quantity,
			Code_Value:   products.Code_Value,
			Is_Published: products.Is_Published,
			Expiration:   products.Expiration,
			Price:        products.Price,
		}
		web.ResponseJSON(w, http.StatusCreated, map[string]any{
			"Message": "Product created successfully",
			"Product": data,
		})
	}
}
