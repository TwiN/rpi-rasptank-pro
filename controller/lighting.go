package controller

import (
	ws281x "github.com/rpi-ws281x/rpi-ws281x-go"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
)

const (
	boardLedAPin = "29"
	boardLedBPin = "31"
	boardLedCPin = "33"
)

type Lighting struct {
	BoardLedA, BoardLedB, BoardLedC *gpio.LedDriver

	ws2812 *ws281x.WS2811
}

func NewLighting(rpi *raspi.Adaptor) *Lighting {
	opt := ws281x.DefaultOptions
	opt.Channels[0].Brightness = 255
	opt.Channels[0].LedCount = 7
	opt.Channels[0].GpioPin = 12
	ws2812, err := ws281x.MakeWS2811(&opt)
	if err != nil {
		log.Printf("[NewLighting] Failed to create lighting: %s", err.Error())
	}
	if err := ws2812.Init(); err != nil {
		log.Println(err)
	}
	return &Lighting{
		BoardLedA: gpio.NewLedDriver(rpi, boardLedAPin),
		BoardLedB: gpio.NewLedDriver(rpi, boardLedBPin),
		BoardLedC: gpio.NewLedDriver(rpi, boardLedCPin),
		ws2812:    ws2812,
	}
}

func (l *Lighting) RedSideLights() {
	l.sideLights(0x0000FF)
}

func (l *Lighting) GreenSideLights() {
	l.sideLights(0xFF0000)
}

func (l *Lighting) BlueSideLights() {
	l.sideLights(0x00FF00)
}

func (l *Lighting) WhiteSideLights() {
	l.sideLights(0xFFFFFF)
}

func (l *Lighting) YellowSideLights() {
	l.sideLights(0xFF00FF)
}

func (l *Lighting) DarkOrangeSideLights() {
	l.sideLights(0x6500FF)
}

func (l *Lighting) OrangeSideLights() {
	l.sideLights(0xA500FF)
}

func (l *Lighting) ClearSideLights() {
	l.sideLights(0x000000)
}

func (l *Lighting) sideLights(value uint32) {
	if len(l.ws2812.Leds(0)) == 0 {
		return
	}
	for i := range l.ws2812.Leds(0) {
		l.ws2812.Leds(0)[i] = value
	}
	if err := l.ws2812.Render(); err != nil {
		log.Println(err)
	}
	_ = l.ws2812.Wait()
}

func (l *Lighting) ToggleBoardLeds() {
	l.BoardLedA.Toggle()
	l.BoardLedB.Toggle()
	l.BoardLedC.Toggle()
}
