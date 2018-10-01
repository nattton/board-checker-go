package models

import (
	"fmt"
	"strconv"
)

func RoundFloat(amount float64) float64 {
	roundAmount := fmt.Sprintf("%.2f", amount)
	amount, _ = strconv.ParseFloat(roundAmount, 64)
	return amount
}
