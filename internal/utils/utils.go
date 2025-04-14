package utils

import "avito/internal/models"

var allowedCities = map[models.PvzCity]bool{
	models.Kazan:           true,
	models.Moscow:          true,
	models.SaintPetersburg: true,
}

var allowedProductTypes = map[models.ProductType]bool{
	models.Shoes:       true,
	models.Clothes:     true,
	models.Electronics: true,
}

var skipPaths = map[string]bool{
	"/dummyLogin": true,
	"/register":   true,
	"/login":      true,
}

func IsCityAllowed(city models.PvzCity) bool {
	return allowedCities[city]
}

func IsProductTypeAllowed(productType models.ProductType) bool {
	return allowedProductTypes[productType]
}

func IsPathSkipped(path string) bool {
	return skipPaths[path]
}
