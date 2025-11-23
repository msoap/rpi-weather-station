package main

import (
	"fmt"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3/rpi"
)

const (
	dispW = 128
	dispH = 64
)

// SH1106 driver
type SH1106 struct {
	spi     spi.Conn
	dc      gpio.PinOut
	reset   gpio.PinOut
	spiPort spi.PortCloser
	buf     [dispW * dispH / 8]byte
}

// Create new SH1106 driver
func NewSH1106() (*SH1106, error) {
	// SPI0.0 (CE0)
	spiPort, err := spireg.Open("/dev/spidev0.0")
	if err != nil {
		return nil, fmt.Errorf("failed to open SPI: %w", err)
	}

	spiDev, err := spiPort.Connect(8*physic.MegaHertz, spi.Mode0, 8)
	if err != nil {
		return nil, fmt.Errorf("failed to connect SPI: %w", err)
	}

	// GPIO DC и RESET
	dc := rpi.P1_22    // GPIO25
	reset := rpi.P1_18 // GPIO24
	if err := dc.Out(gpio.Low); err != nil {
		return nil, fmt.Errorf("failed to set DC pin: %w", err)
	}
	if err := reset.Out(gpio.High); err != nil {
		return nil, fmt.Errorf("failed to set RESET pin: %w", err)
	}

	return &SH1106{
		spi:     spiDev,
		dc:      dc,
		reset:   reset,
		spiPort: spiPort,
	}, nil
}

// Send command
func (d *SH1106) cmd(c byte) {
	d.dc.Out(gpio.Low)
	d.spi.Tx([]byte{c}, nil)
}

// Send data
func (d *SH1106) data(b []byte) {
	d.dc.Out(gpio.High)
	d.spi.Tx(b, nil)
}

// Initialize display
func (d *SH1106) Init() {
	// Hardware reset
	d.reset.Out(gpio.Low)
	time.Sleep(10 * time.Millisecond)
	d.reset.Out(gpio.High)

	// --- SH1106 INIT SEQUENCE ---
	d.cmd(0xAE) // display off
	d.cmd(0xD5)
	d.cmd(0x80) // clock
	d.cmd(0xA8)
	d.cmd(0x3F) // multiplex 64
	d.cmd(0xD3)
	d.cmd(0x00) // display offset
	d.cmd(0x40) // start line 0

	d.cmd(0xAD)
	d.cmd(0x8B) // charge pump enable
	d.cmd(0xA1) // segment remap
	d.cmd(0xC8) // COM scan direction
	d.cmd(0xDA)
	d.cmd(0x12) // COM pins
	d.cmd(0x81)
	d.cmd(0x7F) // contrast
	d.cmd(0xA4) // display follows RAM
	d.cmd(0xA6) // normal display
	d.cmd(0xAF) // display ON

	d.Clear()
	d.Update()
}

// Clear buffer
func (d *SH1106) Clear() {
	for i := range d.buf {
		d.buf[i] = 0
	}
}

// Set pixel
func (d *SH1106) SetPixel(x, y int, on bool) {
	if x < 0 || x >= dispW || y < 0 || y >= dispH {
		return
	}
	page := y / 8
	index := page*dispW + x
	mask := byte(1 << (y % 8))

	if on {
		d.buf[index] |= mask
	} else {
		d.buf[index] &^= mask
	}
}

// Get pixel
func (d *SH1106) GetPixel(x, y int) bool {
	if x < 0 || x >= dispW || y < 0 || y >= dispH {
		return false
	}
	page := y / 8
	index := page*dispW + x
	mask := byte(1 << (y % 8))
	return d.buf[index]&mask != 0
}

func (d *SH1106) Update() {
	for page := 0; page < dispH/8; page++ {
		d.cmd(0xB0 | byte(page)) // set page
		d.cmd(0x00)              // lower column = 0
		d.cmd(0x10)              // higher column = 0

		pageData := make([]byte, 132)
		copy(pageData[2:], d.buf[page*dispW:(page+1)*dispW])
		d.data(pageData)
	}
}

func (d *SH1106) Size() (width, height int) {
	return dispW, dispH
}

func (d *SH1106) Finish() {
	d.Clear()
	d.Update()
	d.spiPort.Close()
}
