package main

// Not used, it's doesn't work :(

import (
	"log"

	"github.com/mdp/monochromeoled"
	"golang.org/x/exp/io/i2c"
)

var oled *monochromeoled.OLED

func init() {
	var err error
	oled, err = monochromeoled.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, 0x3c, 128, 32)
	if err != nil {
		log.Fatal(err)
	}
}

// OLEDShow shows something
func OLEDShow() {
	oled.SetPixel(1, 10, 0x1)
	oled.SetPixel(2, 10, 0x1)
	oled.SetPixel(3, 10, 0x1)
	oled.SetPixel(4, 10, 0x1)
	oled.SetPixel(5, 10, 0x1)
	oled.SetPixel(7, 10, 0x1)
	oled.SetPixel(8, 10, 0x1)
	oled.SetPixel(9, 10, 0x1)
	oled.SetPixel(10, 10, 0x1)
	oled.SetPixel(11, 10, 0x1)
	oled.SetPixel(12, 10, 0x1)
	oled.SetPixel(13, 10, 0x1)
}

// OLEDClose closes the LCD
func OLEDClose() {
	oled.Close()
}
