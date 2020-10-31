package controller

import (
	"fmt"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
	"time"
)

const (
	ArmBus     = 1
	ArmAddress = 0x40

	BaseHorizontalServoPin = "0"
	BaseVerticalServoPin   = "1"
	ClawServoPin           = "2"
	ClawVerticalServoPin   = "3"
	CameraVerticalServoPin = "4"
)

type Arm struct {
	Driver *i2c.PCA9685Driver
}

func NewArm(rpi *raspi.Adaptor) *Arm {
	return &Arm{
		Driver: i2c.NewPCA9685Driver(rpi, i2c.WithBus(ArmBus), i2c.WithAddress(ArmAddress)),
	}
}

func (a *Arm) Center() {
	if err := a.Driver.SetPWMFreq(50.0); err != nil {
		log.Printf("failed to set PWM freq to 50.0: %s", err.Error())
	}
	if err := a.Driver.ServoWrite(BaseHorizontalServoPin, 70); err != nil {
		fmt.Println(err)
	}
	if err := a.Driver.ServoWrite(BaseVerticalServoPin, 90); err != nil {
		fmt.Println(err)
	}
	if err := a.Driver.ServoWrite(ClawServoPin, 90); err != nil {
		fmt.Println(err)
	}
	if err := a.Driver.ServoWrite(ClawVerticalServoPin, 90); err != nil {
		fmt.Println(err)
	}
	if err := a.Driver.ServoWrite(CameraVerticalServoPin, 90); err != nil {
		fmt.Println(err)
	}
}

func (a *Arm) Sweep() {
	if err := a.Driver.SetPWMFreq(50.0); err != nil {
		log.Printf("failed to set PWM freq to 50.0: %s", err.Error())
	}
	for i := 10; i < 130; i++ {
		if err := a.Driver.ServoWrite(BaseHorizontalServoPin, byte(i)); err != nil {
			fmt.Println(err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
