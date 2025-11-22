package main

import (
	"fmt"

	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/bmxx80"
)

type BME280 struct {
	dev *bmxx80.Dev
}

func NewBME280() (*BME280, error) {
	bus, err := i2creg.Open("")
	if err != nil {
		return nil, fmt.Errorf("failed to open I2C bus: %w", err)
	}
	// defer bus.Close()

	dev, err := bmxx80.NewI2C(bus, 0x76, &bmxx80.DefaultOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to open device: %w", err)
	}

	return &BME280{dev: dev}, nil
}

// Read temperature and pressure
func (b *BME280) Read() (string, string, error) {
	var env physic.Env
	if err := b.dev.Sense(&env); err != nil {
		return "", "", fmt.Errorf("failed to read sensor: %w", err)
	}

	return env.Temperature.String(), env.Pressure.String(), nil
}
