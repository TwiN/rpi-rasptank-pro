package main

import (
	"github.com/TwinProduction/rpi-rasptank-pro/controller"
	"github.com/TwinProduction/rpi-rasptank-pro/display"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
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

	led := gpio.NewLedDriver(rpi, "32")
	work := func() {
		if err := screen.DisplayIP(); err != nil {
			log.Printf("Failed to write on screen: %s", err.Error())
		}
		gobot.Every(2*time.Second, func() {
			led.Brightness(3)
			err := led.Toggle()
			if err != nil {
				log.Println(err)
			}
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
