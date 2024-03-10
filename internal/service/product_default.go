package product_default

import (
	"proyect/go-web/internal"
)

func NewDeafultProduct(rp internal.Product_Repository) *DeafultProduct {
	return &DeafultProduct{
		repository: rp,
	}
}

type DeafultProduct struct {
	rp internal.Product_Repository
}

func (t *DeafultProduct) Save(product *internal.Product_Struct) (err error) {
	err = t.rp.Save(product)
	return
}
