package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"periph.io/x/host/v3"
)

const (
	updateDefaultDelay = time.Second * 10
)

type BME280Reader interface {
	Read() (string, string, error)
}

func main() {
	updateDelay := flag.Duration("update-delay", updateDefaultDelay, "Update delay")
	useFakeBME := flag.Bool("fake-bme", false, "Use fake BME280 sensor for testing")
	useTermScreen := flag.Bool("term-screen", false, "Use terminal screen for testing")
	flag.Parse()

	var (
		disp   dispDrawler
		err    error
		exitCh chan struct{}
	)
	if *useTermScreen {
		logFile, err := os.OpenFile("term.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		log.SetOutput(logFile)

		disp, exitCh, err = NewTermScreen()
		if err != nil {
			log.Fatalf("Failed to initialize terminal screen: %v", err)
		}
	} else {
		if _, err := host.Init(); err != nil {
			log.Fatal(err)
		}

		dispSH1106, err := NewSH1106()
		if err != nil {
			log.Fatalf("Failed to initialize SH1106: %v", err)
		}
		dispSH1106.Init()
		disp = dispSH1106
	}

	scr := NewScreen(disp)
	defer scr.Finish()
	scr.Clear()

	width, height := disp.Size()

	scr.Rect(0, 0, width, height, 1)
	scr.HLine(0, 2, width, 1)
	scr.HLine(0, 4, width, 1)

	var bme BME280Reader
	if *useFakeBME {
		bme, err = NewBME280Fake()
	} else {
		bme, err = NewBME280()
	}
	if err != nil {
		log.Fatalf("Failed to initialize BME280: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	shiftX, shiftY, i := 0, 0, 0 // for shift printed text to avoid burn-in oled screen
MAINLOOP:
	for {
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

		scr.FillRect(1, 5, width-2, height-6, 0) // Clear area

		scr.DrawTextCentered(64+shiftX, 15+shiftY, "Temp: "+temp)
		scr.DrawTextCentered(64+shiftX, 35+shiftY, "Pres: "+pres)

		scr.Update()
		i++

		select {
		case <-sigChan:
			break MAINLOOP
		case <-exitCh:
			break MAINLOOP
		case <-time.Tick(*updateDelay):
		}
	}

	log.Println("Program terminated")
}
