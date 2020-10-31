package main

import (
	"fmt"
	"github.com/TwinProduction/rpi-rasptank-pro/controller"
	"github.com/TwinProduction/rpi-rasptank-pro/display"
	"github.com/TwinProduction/rpi-rasptank-pro/sensor"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
	"math"
	"time"
)

// 12: left DC motor backward
// 13: left DC motor forward
// 37: right DC motor forward
// 40: right DC motor backward

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
		var lastDistanceFromObstacle float32
		var stuckCounter int
		gobot.Every(1*time.Second, func() {
			distanceFromObstacle := ultrasonicSensor.MeasureDistanceReliably()
			log.Printf("distance from obstacle: %f", distanceFromObstacle)
			if math.Round(float64(lastDistanceFromObstacle)) == math.Round(float64(distanceFromObstacle)) {
				fmt.Println("doesn't look like you moved much..")
				stuckCounter++
				log.Println("going right")
				if stuckCounter%2 == 0 {
					vehicle.Right()
				} else {
					vehicle.Backward()
				}
				msToSleep := stuckCounter * 50
				if msToSleep > 1000 {
					msToSleep = 1000
				}
				time.Sleep(time.Duration(msToSleep) * time.Millisecond)
				vehicle.Stop()
			} else {
				stuckCounter = 0
				if distanceFromObstacle == 0 {
					// If the distance was 0, there's probably something blocking the sensor, so we'll just turn
					log.Println("going right")
					vehicle.Right()
					time.Sleep(100 * time.Millisecond)
					vehicle.Stop()
				} else if distanceFromObstacle < 3 {
					log.Println("going backward")
					vehicle.Backward()
					time.Sleep(250 * time.Millisecond)
					vehicle.Stop()
				} else if distanceFromObstacle < 35 {
					log.Println("going left")
					vehicle.Left()
					time.Sleep(100 * time.Millisecond)
					vehicle.Stop()
				} else {
					log.Println("going forward")
					vehicle.Forward()
					time.Sleep(500 * time.Millisecond)
					vehicle.Stop()
				}
			}

			lastDistanceFromObstacle = distanceFromObstacle
		})
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{rpi},
		[]gobot.Device{screen.Driver, vehicle.LeftMotor, vehicle.RightMotor, led},
		work,
	)

	robot.Start()
}
