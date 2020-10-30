package main

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
	"time"
)

// 12: right DC motor

func main() {
	adaptor := raspi.NewAdaptor()
	led := gpio.NewLedDriver(adaptor, "7")
	work := func() {
		gobot.Every(3*time.Second, func() {
			log.Println("Toggling")
			err := led.Toggle()
			if err != nil {
				log.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{adaptor},
		[]gobot.Device{led},
		work,
	)

	robot.Start()
}
