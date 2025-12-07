# rpi-weather-station

RPI a simple weather station

## Hardware

- Raspberry Pi 1
- BMP280 sensor
- SH1106 OLED display, 128x64 screen, 4 wire SPI

## Build for Raspberry Pi 1*

```sh
make build-rpi
```

## BMP280 connections

| BMP280 Pin | RPI Pin      |
|------------|--------------|
| VCC        | 3.3V         |
| GND        | GND          |
| SCL        | SCL (GPIO 3) |
| SDA        | SDA (GPIO 2) |

## SH1106 connections

| SH1106 Pin | RPI Pin        |
|------------|----------------|
| VCC        | 3.3V           |
| GND        | GND            |
| DIN        | MOSI (GPIO 10) |
| CLK        | SCLK (GPIO 11) |
| CS         | CE0 (GPIO 8)   |
| DC         | GPIO 25        |
| RST        | GPIO 24        |

## Reallife screenshot

![scr](https://github.com/user-attachments/assets/e1cad8f1-8bcc-4b48-887e-5ea758caa637)
