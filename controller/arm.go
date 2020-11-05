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

	DefaultBaseHorizontalPosition = 75
	DefaultBaseVerticalPosition   = 150
	DefaultClawPosition           = 85
	DefaultClawVerticalPosition   = 90
	DefaultCameraVerticalPosition = 70
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
	if err := a.Driver.ServoWrite(BaseHorizontalServoPin, DefaultBaseHorizontalPosition); err != nil {
		fmt.Println(err)
	}
	if err := a.Driver.ServoWrite(BaseVerticalServoPin, DefaultBaseVerticalPosition); err != nil {
		fmt.Println(err)
	}
	if err := a.Driver.ServoWrite(ClawServoPin, DefaultClawPosition); err != nil {
		fmt.Println(err)
	}
	if err := a.Driver.ServoWrite(ClawVerticalServoPin, DefaultClawVerticalPosition); err != nil {
		fmt.Println(err)
	}
	if err := a.Driver.ServoWrite(CameraVerticalServoPin, DefaultCameraVerticalPosition); err != nil {
		fmt.Println(err)
	}
}

func (a *Arm) Relax() {
	_ = a.Driver.SetPWMFreq(0.0)
}

func (a *Arm) ClawGrab() {
	if err := a.Driver.SetPWMFreq(50.0); err != nil {
		log.Printf("failed to set PWM freq to 50.0: %s", err.Error())
	}
	if err := a.Driver.ServoWrite(ClawServoPin, byte(0)); err != nil {
		fmt.Println(err)
	}
}

func (a *Arm) ClawRelease() {
	if err := a.Driver.SetPWMFreq(50.0); err != nil {
		log.Printf("failed to set PWM freq to 50.0: %s", err.Error())
	}
	if err := a.Driver.ServoWrite(ClawServoPin, byte(DefaultClawPosition)); err != nil {
		fmt.Println(err)
	}
}

func (a *Arm) Sweep() {
	if err := a.Driver.SetPWMFreq(50.0); err != nil {
		log.Printf("failed to set PWM freq to 50.0: %s", err.Error())
	}
	for i := 0; i < 90; i++ {
		if err := a.Driver.ServoWrite(BaseHorizontalServoPin, byte(i)); err != nil {
			fmt.Println(err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (a *Arm) StraightUp() {
	if err := a.Driver.SetPWMFreq(50.0); err != nil {
		log.Printf("failed to set PWM freq to 50.0: %s", err.Error())
	}
	a.Center()
	if err := a.Driver.ServoWrite(BaseVerticalServoPin, 180); err != nil {
		fmt.Println(err)
	}
	if err := a.Driver.ServoWrite(ClawVerticalServoPin, 0); err != nil {
		fmt.Println(err)
	}
	if err := a.Driver.ServoWrite(CameraVerticalServoPin, 0); err != nil {
		fmt.Println(err)
	}
}

func (a *Arm) PushUpLeft() {
	a.pushUp(150)
}

func (a *Arm) PushUpRight() {
	a.pushUp(0)
}

func (a *Arm) pushUp(baseHorizontalServoAngle byte) {
	if err := a.Driver.SetPWMFreq(50.0); err != nil {
		log.Printf("failed to set PWM freq to 50.0: %s", err.Error())
	}
	a.StraightUp()
	time.Sleep(500 * time.Millisecond)
	if err := a.Driver.ServoWrite(BaseHorizontalServoPin, baseHorizontalServoAngle); err != nil {
		fmt.Println(err)
	}
	time.Sleep(300 * time.Millisecond)
	for i := 90; i > 0; i -= 3 {
		if err := a.Driver.ServoWrite(BaseVerticalServoPin, byte(i)); err != nil {
			fmt.Println(err)
		}
		time.Sleep(100 * time.Millisecond)
	}
	if err := a.Driver.ServoWrite(ClawVerticalServoPin, 50); err != nil {
		fmt.Println(err)
	}
	time.Sleep(100 * time.Millisecond)
	if err := a.Driver.ServoWrite(ClawServoPin, 0); err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Second)
	a.Center()
	time.Sleep(time.Second)
	a.Relax()
}

func (a *Arm) LookAt(x, y int) {
	if x > DefaultBaseHorizontalPosition+30 {
		x = DefaultBaseHorizontalPosition + 30
	} else if x < DefaultBaseHorizontalPosition-30 {
		x = DefaultBaseHorizontalPosition - 30
	}
	if y > DefaultCameraVerticalPosition+30 {
		y = DefaultCameraVerticalPosition + 30
	} else if y < DefaultCameraVerticalPosition-30 {
		y = DefaultCameraVerticalPosition - 30
	}
	if err := a.Driver.SetPWMFreq(50.0); err != nil {
		log.Printf("failed to set PWM freq to 50.0: %s", err.Error())
	}
	if err := a.Driver.ServoWrite(BaseHorizontalServoPin, byte(x)); err != nil {
		fmt.Println(err)
	}
	if err := a.Driver.ServoWrite(CameraVerticalServoPin, byte(y)); err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Second)
	a.Relax()
}
