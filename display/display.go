package display

import (
	"github.com/pbnjay/pixfont"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	"image"
	"image/color"
)

const (
	Width  = 128
	Height = 64

	Bus     = 1
	Address = 0x3c
)

func CreateDriver(rpi *raspi.Adaptor) *i2c.SSD1306Driver {
	return i2c.NewSSD1306Driver(rpi, i2c.WithBus(Bus), i2c.WithAddress(Address))
}

func DrawString(driver *i2c.SSD1306Driver, text string) error {
	rectangle := image.Rect(0, 0, Width, Height)
	img := image.NewRGBA(rectangle)
	pixfont.DrawString(img, 10, 10, text, color.White)
	flipped := image.NewRGBA(rectangle)
	for j := 0; j < img.Bounds().Dy(); j++ {
		for i := 0; i < img.Bounds().Dx(); i++ {
			flipped.Set(Width-i, Height-j, img.At(i, j))
		}
	}
	driver.Clear()
	if err := driver.ShowImage(flipped); err != nil {
		return err
	}
	return driver.Display()
}
