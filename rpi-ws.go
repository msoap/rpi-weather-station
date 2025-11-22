package main

import (
	"log"
	"strings"
	"time"

	"periph.io/x/host/v3"
)

func main() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	disp, err := NewSH1106()
	if err != nil {
		log.Fatalf("Failed to initialize SH1106: %v", err)
	}
	disp.Init()

	scr := NewScreen(disp)
	scr.Clear()

	width, height := disp.Size()

	scr.Rectangle(0, 0, width, height, true)
	scr.DrawHorizontalLine(0, 2, width, true)
	scr.DrawHorizontalLine(0, 4, width, true)

	bme, err := NewBME280()
	if err != nil {
		log.Fatalf("Failed to initialize BME280: %v", err)
	}

	shiftX, shiftY, i := 0, 0, 0 // for shift printed text to avoid burn-in oled screen
	for range time.Tick(time.Second * 10) {
		switch i % 4 {
		case 0:
			shiftX, shiftY = 0, 0
		case 1:
			shiftX, shiftY = 0, 1
		case 2:
			shiftX, shiftY = 1, 0
		case 3:
			shiftX, shiftY = 0, 1
		}

		temp, pres, err := bme.Read()
		if err != nil {
			log.Fatalf("Failed to read BME280: %v", err)
		}

		temp = strings.ReplaceAll(temp, "°", " ")
		log.Printf("Temperature: %s, Pressure: %s", temp, pres)

		scr.Box(1, 5, 120, 55, false) // Clear area
		scr.DrawTextCentered(64+shiftX, 15+shiftY, "Temp: "+temp)
		scr.DrawTextCentered(64+shiftX, 35+shiftY, "Pres: "+pres)

		disp.Update()
		i++
	}

	log.Println("OLED initialized, test pixels drawn.")
}
