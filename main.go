package main

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
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
	board := raspi.NewAdaptor()
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
	screen := i2c.NewGroveLcdDriver(board)
	adaFruit := i2c.NewAdafruitMotorHatDriver(rpi)
	work := func() {
		gobot.Every(3*time.Second, func() {
			screen.Write("hello")
			log.Println("o.o")
			var speed int32 = 10 // 255 = full speed!
			if err := adaFruit.SetDCMotorSpeed(2, speed); err != nil {
				log.Println(err)
			}
			if err := adaFruit.RunDCMotor(2, i2c.AdafruitForward); err != nil {
				log.Println(err)
			}
			time.Sleep(100 * time.Millisecond)
			if err := adaFruit.RunDCMotor(2, i2c.AdafruitRelease); err != nil {
				log.Println(err)
			}
			time.Sleep(100 * time.Millisecond)
			if err := adaFruit.RunDCMotor(2, i2c.AdafruitBackward); err != nil {
				log.Println(err)
			}
			time.Sleep(100 * time.Millisecond)
			if err := adaFruit.RunDCMotor(2, i2c.AdafruitRelease); err != nil {
				log.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{rpi, board},
		//[]gobot.Device{led},
		[]gobot.Device{adaFruit, screen},
		work,
	)

	robot.Start()
}
