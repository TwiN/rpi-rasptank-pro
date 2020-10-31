package main

import (
	"github.com/TwinProduction/rpi-rasptank-pro/controller"
	"github.com/TwinProduction/rpi-rasptank-pro/display"
	"github.com/mcuadros/go-rpi-ws281x"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
	"image"
	"image/color"
	"image/draw"
	"log"
	"time"
)

// 12: right DC motor forward
// 13: right DC motor backward
// 37: left DC motor backward
// 40: left DC motor forward

// 29: LED on HAT board
// 31: LED on HAT board
// 33: LED on HAT board

// 32: LED on LEFT, BACK

func main() {
	rpi := raspi.NewAdaptor()
	screen := display.NewDisplay(rpi)
	engine := controller.NewEngine(rpi)
	//led := gpio.NewLedDriver(rpi, os.Args[1])
	//work := func() {
	//	gobot.Every(3*time.Second, func() {
	//		log.Println("Toggling")
	//		err := led.Toggle()
	//		if err != nil {
	//			log.Println(err)
	//		}
	//	})
	//}
	//pca9685 := i2c.NewPCA9685Driver(rpi)
	ws2812Config := &ws281x.DefaultConfig
	ws2812Config.Pin = 32

	led := gpio.NewLedDriver(rpi, "32")
	work := func() {
		if err := screen.DisplayIP(); err != nil {
			log.Printf("Failed to write on screen: %s", err.Error())
		}
		gobot.Every(3*time.Second, func() {
			c, _ := ws281x.NewCanvas(8, 4, ws2812Config)
			if err := c.Initialize(); err != nil {
				log.Println(err)
			}
			draw.Draw(c, c.Bounds(), image.NewUniform(color.White), image.ZP, draw.Over)
			if err := c.Render(); err != nil {
				log.Println(err)
			}
			time.Sleep(time.Second * 1)

			// don't forget close the canvas, if not you leds may remain on
			c.Close()
			//err := led.Toggle()
			//if err != nil {
			//	log.Println(err)
			//}
			//engine.Right()
		})
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{rpi},
		[]gobot.Device{screen.Driver, engine.LeftMotor, engine.RightMotor, led},
		work,
	)

	robot.Start()
}
