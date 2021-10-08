package display

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/TwiN/rpi-rasptank-pro/util"
	"github.com/pbnjay/pixfont"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

const (
	Width  = 128
	Height = 64

	bus     = 1
	address = 0x3c
)

type Display struct {
	Driver *i2c.SSD1306Driver
}

func NewDisplay(rpi *raspi.Adaptor) *Display {
	return &Display{
		Driver: i2c.NewSSD1306Driver(rpi, i2c.WithBus(bus), i2c.WithAddress(address)),
	}
}

func (d *Display) DisplayIP() error {
	return d.DrawString(util.GetLocalIP())
}

func (d *Display) DisplayIPAndText(text string) error {
	return d.DrawString(fmt.Sprintf("%s\n%s", util.GetLocalIP(), text))
}

func (d *Display) DrawString(text string) error {
	rectangle := image.Rect(0, 0, Width, Height)
	img := image.NewRGBA(rectangle)
	lines := strings.Split(text, "\n")
	for number, line := range lines {
		pixfont.DrawString(img, 0, 10+(number*10), line, color.White)
	}
	flipped := image.NewRGBA(rectangle)
	for j := 0; j < img.Bounds().Dy(); j++ {
		for i := 0; i < img.Bounds().Dx(); i++ {
			flipped.Set(Width-i, Height-j, img.At(i, j))
		}
	}
	d.Driver.Clear()
	if err := d.Driver.ShowImage(flipped); err != nil {
		return err
	}
	return d.Driver.Display()
}
