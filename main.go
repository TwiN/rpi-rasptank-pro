package main

import (
	"fmt"
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

func main() {
	rpi := raspi.NewAdaptor()
	screen := display.CreateDriver(rpi)
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

	motor := gpio.NewMotorDriver(rpi, "40")
	work := func() {
		err := display.DrawString(screen, fmt.Sprintf("%s", GetLocalIP()))
		if err != nil {
			log.Printf("Failed to write on display: %s", err.Error())
		}
		gobot.Every(3*time.Second, func() {
			err := motor.On()
			if err != nil {
				log.Println(err)
			}
			time.Sleep(time.Second)
			motor.Off()
		})
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{rpi},
		[]gobot.Device{screen, motor},
		work,
	)

	robot.Start()
}
