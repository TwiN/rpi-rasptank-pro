package main

import (
	"fmt"
	"github.com/TwinProduction/rpi-rasptank-pro/controller"
	"github.com/TwinProduction/rpi-rasptank-pro/display"
	"github.com/rpi-ws281x/rpi-ws281x-go"
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
	hw := ws2811.HwDetect()
	fmt.Printf("Hardware Type    : %d\n", hw.Type)
	fmt.Printf("Hardware Version : 0x%08X\n", hw.Version)
	fmt.Printf("Periph base      : 0x%08X\n", hw.PeriphBase)
	fmt.Printf("Video core base  : 0x%08X\n", hw.VideocoreBase)
	fmt.Printf("Description      : %v\n", hw.Desc)

	ws2812Options := &ws2811.DefaultOptions
	ws2812Options.Channels[0].GpioPin = 32
	ws2812, err := ws2811.MakeWS2811(ws2812Options)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ws2812.Leds(1))
	if err := ws2812.Render(); err != nil {
		fmt.Println(err)
	}

	led := gpio.NewLedDriver(rpi, "32")
	work := func() {
		if err := screen.DisplayIP(); err != nil {
			log.Printf("Failed to write on screen: %s", err.Error())
		}
		gobot.Every(3*time.Second, func() {
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
