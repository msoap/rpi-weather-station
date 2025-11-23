package main

import (
	"fmt"
	"math/rand"
)

type BME280Fake struct{}

func NewBME280Fake() (*BME280Fake, error) {
	return &BME280Fake{}, nil
}

func (b *BME280Fake) Read() (string, string, error) {
	temp := 19.0 + rand.Float64()*5.0
	pres := 100.0 + rand.Float64()*10.0
	return fmt.Sprintf("%.2f°C", temp), fmt.Sprintf("%.3f kPa", pres), nil
}
