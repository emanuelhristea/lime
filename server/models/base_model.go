package models

import (
	"errors"

	"github.com/emanuelhristea/lime/config"
)

var (
	// ErrKeyNotFound is a ...
	ErrKeyNotFound = errors.New("key not found")

	// ErrLicenseNotFound is a ...
	ErrLicenseNotFound = errors.New("license not found")

	// ErrTariffNotFound is a ...
	ErrTariffNotFound = errors.New("tariff not found")

	// ErrCustomerNotFound is a ...
	ErrCustomerNotFound = errors.New("customer not found")
)

func init() {
	config.DB.AutoMigrate(&Tariff{}, &Customer{}, &Subscription{}, &License{})
}
