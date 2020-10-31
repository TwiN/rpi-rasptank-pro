package controller

import (
	"fmt"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	"time"
)

const (
	ArmBus     = 1
	ArmAddress = 0x40

	BaseVerticalServoPin = "1"
)

type Arm struct {
	Driver *i2c.PCA9685Driver
}

func NewArm(rpi *raspi.Adaptor) *Arm {
	arm := &Arm{
		Driver: i2c.NewPCA9685Driver(rpi, i2c.WithBus(ArmBus), i2c.WithAddress(ArmAddress)),
	}
	if err := arm.Driver.SetPWMFreq(50.0); err != nil {
		panic(err)
	}
	return arm
}

func (a *Arm) Move() {
	err := a.Driver.SetPWMFreq(50.0)
	if err != nil {
		fmt.Printf("failed to set PWM freq to 50.0: %s\n", err.Error())
	}
	if err := a.Driver.ServoWrite(BaseVerticalServoPin, 90); err != nil {
		fmt.Println(err)
	}
	time.Sleep(1 * time.Second)
}
