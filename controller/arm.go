package controller

import (
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
	"time"
)

const (
	armBus     = 1
	armAddress = 0x40
)

type Arm struct {
	Driver *i2c.PCA9685Driver

	BaseHorizontalServo, BaseVerticalServo, ClawServo, ClawVerticalServo, CameraVerticalServo servo
}

func NewArm(rpi *raspi.Adaptor) *Arm {
	return &Arm{
		Driver: i2c.NewPCA9685Driver(rpi, i2c.WithBus(armBus), i2c.WithAddress(armAddress)),
		BaseHorizontalServo: servo{
			Pin:     "0",
			Default: 75,
			Min:     0,
			Max:     105,
		},
		BaseVerticalServo: servo{
			Pin:     "1",
			Default: 150,
			Min:     0,
			Max:     180,
		},
		ClawServo: servo{
			Pin:     "2",
			Default: 85,
			Min:     0,
			Max:     85,
		},
		ClawVerticalServo: servo{
			Pin:     "3",
			Default: 90,
			Min:     0,
			Max:     90,
		},
		CameraVerticalServo: servo{
			Pin:     "4",
			Default: 70,
			Min:     0,
			Max:     100,
		},
	}
}

func (a *Arm) Center() {
	if err := a.Driver.SetPWMFreq(50.0); err != nil {
		log.Printf("failed to set PWM freq to 50.0: %s", err.Error())
	}
	a.BaseHorizontalServo.MoveDefault(a.Driver)
	a.BaseVerticalServo.MoveDefault(a.Driver)
	a.ClawServo.MoveDefault(a.Driver)
	a.ClawVerticalServo.MoveDefault(a.Driver)
	a.CameraVerticalServo.MoveDefault(a.Driver)
}

func (a *Arm) Relax() {
	_ = a.Driver.SetPWMFreq(0.0)
}

func (a *Arm) ClawGrab() {
	if err := a.Driver.SetPWMFreq(50.0); err != nil {
		log.Printf("failed to set PWM freq to 50.0: %s", err.Error())
	}
	a.ClawServo.Move(a.Driver, 0)
}

func (a *Arm) ClawRelease() {
	if err := a.Driver.SetPWMFreq(50.0); err != nil {
		log.Printf("failed to set PWM freq to 50.0: %s", err.Error())
	}
	a.ClawServo.MoveDefault(a.Driver)
}

func (a *Arm) Sweep() {
	if err := a.Driver.SetPWMFreq(50.0); err != nil {
		log.Printf("failed to set PWM freq to 50.0: %s", err.Error())
	}
	for i := 0; i < 90; i++ {
		if err := a.BaseHorizontalServo.Move(a.Driver, i); err != nil {
			log.Println(err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (a *Arm) StraightUp() {
	if err := a.Driver.SetPWMFreq(50.0); err != nil {
		log.Printf("failed to set PWM freq to 50.0: %s", err.Error())
	}
	a.Center()
	a.BaseVerticalServo.MoveMin(a.Driver)
	a.ClawVerticalServo.MoveMin(a.Driver)
}

func (a *Arm) PushUpLeft() {
	a.pushUp(150)
}

func (a *Arm) PushUpRight() {
	a.pushUp(0)
}

func (a *Arm) pushUp(baseHorizontalServoAngle int) {
	if err := a.Driver.SetPWMFreq(50.0); err != nil {
		log.Printf("failed to set PWM freq to 50.0: %s", err.Error())
	}
	a.StraightUp()
	time.Sleep(500 * time.Millisecond)
	a.BaseHorizontalServo.Move(a.Driver, baseHorizontalServoAngle)
	time.Sleep(300 * time.Millisecond)
	for i := 90; i > 0; i -= 3 {
		a.BaseVerticalServo.Move(a.Driver, i)
		time.Sleep(100 * time.Millisecond)
	}
	a.ClawVerticalServo.Move(a.Driver, 50)
	time.Sleep(100 * time.Millisecond)
	a.ClawServo.MoveMin(a.Driver)
	time.Sleep(time.Second)
	a.Center()
	time.Sleep(time.Second)
	a.Relax()
}

func (a *Arm) LookAt(x, y int) {
	if err := a.Driver.SetPWMFreq(50.0); err != nil {
		log.Printf("failed to set PWM freq to 50.0: %s", err.Error())
	}
	a.BaseHorizontalServo.Move(a.Driver, x)
	a.CameraVerticalServo.Move(a.Driver, y)
	time.Sleep(time.Second)
	a.Relax()
}
