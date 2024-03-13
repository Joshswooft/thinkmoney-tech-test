package sku

import (
	"errors"
	"unicode"
)

var (
	ErrNoSpecialCharacters = errors.New("a SKU must not contain any special characters")
)

// a SKU known as Stock Keeping Unit is a unique identifier for a product
// You should initialize SKUs using the New() func
type SKU struct {
	value rune
}

// creates and validates a sku from a given rune
func New(value rune) (SKU, error) {
	upperSku := unicode.ToUpper(value)
	if err := validate(upperSku); err != nil {
		return SKU{}, err
	}
	return SKU{value: upperSku}, nil
}

// allows us to use log formatting
func (s SKU) String() string {
	return string(s.value)
}

func (s SKU) Value() rune {
	return s.value
}

func validate(r rune) error {
	if !unicode.IsLetter(r) {
		return ErrNoSpecialCharacters
	}
	return nil
}
