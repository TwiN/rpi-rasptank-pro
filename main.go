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
	//ultrasonicSensor := sensor.NewUltrasonicSensor()
	servo := gpio.NewServoDriver(rpi, "1")

	led := gpio.NewLedDriver(rpi, "32")
	work := func() {
		if err := screen.DisplayIP(); err != nil {
			log.Printf("Failed to write on screen: %s", err.Error())
		}
		servo.Center()
		//gobot.Every(1*time.Second, func() {
		//})
		//var lastDistanceFromObstacle float32
		//var lastDirection string
		//var stuckCounter int
		//var wentBackAndForth bool
		for {
			//distanceFromObstacle := ultrasonicSensor.MeasureDistanceReliably()
			//log.Printf("distance from obstacle: %f", distanceFromObstacle)
			//if wentBackAndForth || math.Round(float64(lastDistanceFromObstacle)) == math.Round(float64(distanceFromObstacle)) {
			//	wentBackAndForth = false
			//	fmt.Println("doesn't look like you moved much..")
			//	stuckCounter++
			//	if stuckCounter%2 == 0 {
			//		log.Println("going right")
			//		vehicle.Right()
			//	} else {
			//		log.Println("going backward")
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
			//		log.Println("going right")
			//		vehicle.Right()
			//		time.Sleep(300 * time.Millisecond)
			//		vehicle.Stop()
			//	} else if distanceFromObstacle < 15 {
			//		log.Println("going backward")
			//		vehicle.Backward()
			//		time.Sleep(500 * time.Millisecond)
			//		vehicle.Stop()
			//	} else {
			//		log.Println("going forward")
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
		[]gobot.Device{screen.Driver, vehicle.LeftMotor, vehicle.RightMotor, led, servo},
		work,
	)

	robot.Start()
}
