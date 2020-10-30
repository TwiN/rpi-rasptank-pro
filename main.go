package main

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	adaptor := raspi.NewAdaptor()
	gobot.NewRobot()
	driver := i2c.NewPCA9685Driver(adaptor)
	err := driver.Start()
	if err != nil {
		panic(err)
	}
}
