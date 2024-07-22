package main

import "kibonga/quickbits/internal/models"

type templateData struct {
	Bit  *models.Bit
	Bits []*models.Bit
}
