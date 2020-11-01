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
	arm := controller.NewArm(rpi)
	ultrasonicSensor := sensor.NewUltrasonicSensor()

	led := gpio.NewLedDriver(rpi, "32")
	work := func() {
		if err := screen.DisplayIP(); err != nil {
			log.Printf("Failed to write on screen: %s", err.Error())
		}

		arm.Center()
		time.Sleep(time.Second)
		arm.PushUpRight()
		//for {
		//	arm.Grab()
		//	time.Sleep(time.Second)
		//}
		//time.Sleep(time.Second)
		//arm.Sweep()
		//var lastDistanceFromObstacle float32
		//var lastDirection string
		//var stuckCounter int
		//var wentBackAndForth bool
		for {
			distanceFromObstacle := ultrasonicSensor.MeasureDistanceReliably()
			for retries := 0; retries < 3 && distanceFromObstacle == sensor.InvalidMeasurement; retries++ {
				distanceFromObstacle = ultrasonicSensor.MeasureDistanceReliably()
				log.Println("Trying to measure distance again")
			}
			if err := screen.DisplayIPAndText(fmt.Sprintf("\nus: %.0fcm", distanceFromObstacle)); err != nil {
				log.Printf("Failed to write on screen: %s", err.Error())
			}
			//if distanceFromObstacle == sensor.InvalidMeasurement {
			//	log.Println("couldn't measure distance, turning right")
			//	vehicle.Right()
			//	time.Sleep(100 * time.Millisecond)
			//	vehicle.Stop()
			//	continue
			//}
			//log.Printf("distance from obstacle: %f", distanceFromObstacle)
			//if wentBackAndForth || math.Round(float64(lastDistanceFromObstacle)) == math.Round(float64(distanceFromObstacle)) {
			//	wentBackAndForth = false
			//	fmt.Println("doesn't look like you moved much..")
			//	stuckCounter++
			//	if stuckCounter%2 == 0 {
			//		log.Println("turning right")
			//		vehicle.Right()
			//	} else {
			//		log.Println("moving backward")
			//		vehicle.Backward()
			//	}
			//	msToSleep := stuckCounter * 200
			//	if msToSleep > 1000 {
			//		msToSleep = 1000
			//	}
			//	time.Sleep(time.Duration(msToSleep) * time.Millisecond)
			//	vehicle.Stop()
			//} else {
			//	stuckCounter = 0
			//	if distanceFromObstacle == 0 {
			//		// If the distance was 0, there's probably something blocking the sensor, so we'll just turn
			//		log.Println("turning right")
			//		vehicle.Right()
			//		time.Sleep(300 * time.Millisecond)
			//		vehicle.Stop()
			//	} else if distanceFromObstacle < 15 {
			//		log.Println("moving backward")
			//		vehicle.Backward()
			//		time.Sleep(500 * time.Millisecond)
			//		vehicle.Stop()
			//	} else {
			//		log.Println("moving forward")
			//		vehicle.Forward()
			//		time.Sleep(500 * time.Millisecond)
			//		vehicle.Stop()
			//	}
			//}
			//wentBackAndForth = (lastDirection == controller.DirectionLeft && vehicle.LastDirection == controller.DirectionRight) ||
			//	(lastDirection == controller.DirectionRight && vehicle.LastDirection == controller.DirectionLeft) ||
			//	(lastDirection == controller.DirectionForward && vehicle.LastDirection == controller.DirectionBackward) ||
			//	(lastDirection == controller.DirectionBackward && vehicle.LastDirection == controller.DirectionForward)
			//lastDirection = vehicle.LastDirection
			//lastDistanceFromObstacle = distanceFromObstacle
			time.Sleep(100 * time.Millisecond)
		}
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{rpi},
		[]gobot.Device{screen.Driver, vehicle.LeftMotor, vehicle.RightMotor, led, arm.Driver},
		work,
	)

	robot.Start()
}
