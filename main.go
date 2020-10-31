package main

import (
	"fmt"
	"github.com/TwinProduction/rpi-rasptank-pro/display"
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

	adaFruit := i2c.NewAdafruitMotorHatDriver(rpi, i2c.WithBus(1))
	work := func() {
		err := display.DrawString(screen, fmt.Sprintf("%s", GetLocalIP()))
		if err != nil {
			log.Printf("Failed to write on display: %s", err.Error())
		}
		gobot.Every(3*time.Second, func() {
			//pca9685.SetPWMFreq(50)
			//pca9685.SetPWM()
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
		[]gobot.Connection{rpi},
		[]gobot.Device{screen},
		work,
	)

	robot.Start()
}
