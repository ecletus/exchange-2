package exchange_test

import "github.com/moisespsena-go/aorm"

type Product struct {
	aorm.Model
	Code       string
	Name       string
	Price      float64
	Tag        *string
	Category   Category
	CategoryID *uint
}

type Category struct {
	aorm.Model
	Name string
}
