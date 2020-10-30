package main

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
	"time"
)

var (
	// Min pulse length out of 4096
	servoMin = 150
	// Max pulse length out of 4096
	servoMax = 700
	// Limiting the max this servo can rotate (in deg)
	maxDegree = 180
	// Number of degrees to increase per call
	degIncrease = 10
	yawDeg      = 90
)

func degree2pulse(deg int) int32 {
	pulse := servoMin
	pulse += ((servoMax - servoMin) / maxDegree) * deg
	return int32(pulse)
}

func adafruitServoMotorRunner(a *i2c.AdafruitMotorHatDriver) (err error) {
	log.Printf("Servo Motor Run Loop...\n")

	var channel byte = 1
	deg := 90

	// Do not need to set this every run loop
	freq := 60.0
	if err = a.SetServoMotorFreq(freq); err != nil {
		log.Printf("%s", err.Error())
		return
	}
	// start in the middle of the 180-deg range
	pulse := degree2pulse(deg)
	if err = a.SetServoMotorPulse(channel, 0, pulse); err != nil {
		log.Printf(err.Error())
		return
	}
	// INCR
	pulse = degree2pulse(deg + degIncrease)
	if err = a.SetServoMotorPulse(channel, 0, pulse); err != nil {
		log.Printf(err.Error())
		return
	}
	time.Sleep(2000 * time.Millisecond)
	// DECR
	pulse = degree2pulse(deg - degIncrease)
	if err = a.SetServoMotorPulse(channel, 0, pulse); err != nil {
		log.Printf(err.Error())
		return
	}
	return
}

func main() {
	adaptor := raspi.NewAdaptor()
	hatDriver := i2c.NewAdafruitMotorHatDriver(adaptor)
	work := func() {
		gobot.Every(5*time.Second, func() {
			err := adafruitServoMotorRunner(hatDriver)
			panic(err)
		})
	}

	robot := gobot.NewRobot("adaFruitBot",
		[]gobot.Connection{adaptor},
		[]gobot.Device{hatDriver},
		work,
	)

	robot.Start()
}
