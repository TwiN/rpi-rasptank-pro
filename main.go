package main

import (
	"github.com/TwinProduction/rpi-rasptank-pro/controller"
	"github.com/TwinProduction/rpi-rasptank-pro/display"
	"github.com/TwinProduction/rpi-rasptank-pro/sensor"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
	"time"
)

// 12: right DC motor backward
// 13: right DC motor forward
// 37: left DC motor forward
// 40: left DC motor backward

// 29: LED on HAT board
// 31: LED on HAT board
// 33: LED on HAT board

// 32: LED on LEFT, BACK

// 11: Ultrasonic Trigger
// 8:  Ultrasonic Echo

func main() {
	rpi := raspi.NewAdaptor()
	screen := display.NewDisplay(rpi)
	vehicle := controller.NewVehicle(rpi)

	led := gpio.NewLedDriver(rpi, "32")
	ultrasonicSensor := sensor.NewUltrasonicSensor()
	work := func() {
		if err := screen.DisplayIP(); err != nil {
			log.Printf("Failed to write on screen: %s", err.Error())
		}
		gobot.Every(1*time.Second, func() {
			distanceFromObstacle := ultrasonicSensor.MeasureDistanceReliably()
			log.Printf("distance from obstacle: %f", distanceFromObstacle)
			if distanceFromObstacle < 10 {
				log.Println("going backward")
				vehicle.Backward()
				time.Sleep(250 * time.Millisecond)
				vehicle.Stop()
			} else if distanceFromObstacle < 35 {
				log.Println("going right")
				vehicle.Right()
				time.Sleep(100 * time.Millisecond)
				vehicle.Stop()
			} else {
				log.Println("going forward")
				vehicle.Forward()
				time.Sleep(500 * time.Millisecond)
				vehicle.Stop()
			}
		})
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{rpi},
		[]gobot.Device{screen.Driver, vehicle.LeftMotor, vehicle.RightMotor, led},
		work,
	)

	robot.Start()
}
