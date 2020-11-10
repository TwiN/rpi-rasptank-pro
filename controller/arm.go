package controller

import (
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
	"sync"
	"time"
)

const (
	armBus     = 1
	armAddress = 0x40
)

type Arm struct {
	sync.Mutex
	Driver *i2c.PCA9685Driver

	BaseHorizontalServo, BaseVerticalServo, ClawServo, ClawVerticalServo, CameraVerticalServo *servo
}

func NewArm(rpi *raspi.Adaptor) *Arm {
	return &Arm{
		Driver: i2c.NewPCA9685Driver(rpi, i2c.WithBus(armBus), i2c.WithAddress(armAddress)),
		BaseHorizontalServo: &servo{
			Pin:     "0",
			Default: 75,
			Min:     0,
			Max:     170,
		},
		BaseVerticalServo: &servo{
			Pin:     "1",
			Default: 150,
			Min:     0,
			Max:     180,
		},
		ClawServo: &servo{
			Pin:     "2",
			Default: 85,
			Min:     0,
			Max:     85,
		},
		ClawVerticalServo: &servo{
			Pin:     "3",
			Default: 70,
			Min:     0,
			Max:     120,
		},
		CameraVerticalServo: &servo{
			Pin:     "4",
			Default: 70,
			Min:     0,
			Max:     100,
		},
	}
}

func (a *Arm) MoveToDefaultPosition() {
	a.Lock()
	defer a.Unlock()
	a.WakeUp()
	a.moveToDefaultPosition()
	time.Sleep(500 * time.Millisecond)
	a.Relax()
}

func (a *Arm) moveToDefaultPosition() {
	a.BaseHorizontalServo.MoveDefault(a.Driver)
	a.BaseVerticalServo.MoveDefault(a.Driver)
	a.ClawServo.MoveDefault(a.Driver)
	a.ClawVerticalServo.MoveDefault(a.Driver)
	a.CameraVerticalServo.MoveDefault(a.Driver)
}

func (a *Arm) Relax() {
	_ = a.Driver.SetPWMFreq(0.0)
}

func (a *Arm) WakeUp() {
	if err := a.Driver.SetPWMFreq(50.0); err != nil {
		log.Printf("failed to set PWM freq to 50.0: %s", err.Error())
	}
}

func (a *Arm) ClawGrab() {
	a.Lock()
	defer a.Unlock()
	a.WakeUp()
	a.ClawServo.Move(a.Driver, 0)
}

func (a *Arm) ClawRelease() {
	a.Lock()
	defer a.Unlock()
	a.WakeUp()
	a.ClawServo.MoveDefault(a.Driver)
	time.Sleep(100 * time.Millisecond)
	a.Relax()
}

func (a *Arm) ClawExtendGrab(checkpointBeforeMove bool) {
	a.Lock()
	defer a.Unlock()
	a.WakeUp()
	a.ClawVerticalServo.Checkpoint(checkpointBeforeMove).Move(a.Driver, 0)
	a.BaseVerticalServo.Checkpoint(checkpointBeforeMove).Move(a.Driver, 50)
	time.Sleep(400 * time.Millisecond)
	a.ClawServo.Checkpoint(checkpointBeforeMove).MoveMin(a.Driver)
}

func (a *Arm) ClawExtendRelease(returnToCheckpoint bool) {
	a.Lock()
	defer a.Unlock()
	a.WakeUp()
	if returnToCheckpoint {
		a.ClawVerticalServo.MoveToCheckpointAndClear(a.Driver)
		a.BaseVerticalServo.MoveToCheckpointAndClear(a.Driver)
		a.ClawServo.MoveToCheckpointAndClear(a.Driver)
	} else {
		a.ClawVerticalServo.MoveDefault(a.Driver)
		a.BaseVerticalServo.MoveDefault(a.Driver)
		a.ClawServo.MoveDefault(a.Driver)
	}
	time.Sleep(200 * time.Millisecond)
	a.Relax()
}

func (a *Arm) Sweep() {
	a.Lock()
	defer a.Unlock()
	a.WakeUp()
	for i := a.BaseHorizontalServo.Min; i < a.BaseHorizontalServo.Max; i += 3 {
		a.BaseHorizontalServo.Move(a.Driver, i)
		time.Sleep(30 * time.Millisecond)
	}
	a.moveToDefaultPosition()
	time.Sleep(time.Second)
	a.Relax()
}

func (a *Arm) LookAt(x, y int) {
	a.Lock()
	defer a.Unlock()
	a.WakeUp()
	a.BaseHorizontalServo.Move(a.Driver, x)
	a.CameraVerticalServo.Move(a.Driver, y)
	time.Sleep(300 * time.Millisecond)
	a.Relax()
}

func (a *Arm) MoveBaseHorizontal(degrees int) {
	a.Lock()
	defer a.Unlock()
	a.WakeUp()
	a.BaseHorizontalServo.Move(a.Driver, a.BaseHorizontalServo.Default+degrees)
	a.Relax()
}

func (a *Arm) MoveBaseVertical(degrees int) {
	a.Lock()
	defer a.Unlock()
	a.WakeUp()
	a.BaseVerticalServo.Move(a.Driver, a.BaseVerticalServo.Default+-degrees)
}

func (a *Arm) MoveClawVertical(degrees int) {
	a.Lock()
	defer a.Unlock()
	a.WakeUp()
	a.ClawVerticalServo.Move(a.Driver, a.ClawVerticalServo.Default+degrees)
}

func (a *Arm) MoveClaw(degrees int) {
	a.Lock()
	defer a.Unlock()
	a.WakeUp()
	a.ClawServo.Move(a.Driver, a.ClawServo.Default+degrees)
}

func (a *Arm) PushUpLeft() {
	a.Lock()
	defer a.Unlock()
	a.pushUp(150)
}

func (a *Arm) PushUpRight() {
	a.Lock()
	defer a.Unlock()
	a.pushUp(0)
}

func (a *Arm) pushUp(baseHorizontalServoAngle int) {
	a.WakeUp()
	a.moveToDefaultPosition()
	time.Sleep(100 * time.Millisecond)
	a.extendUp()
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
	a.moveToDefaultPosition()
	time.Sleep(time.Second)
	a.Relax()
}

// extendUp extends the arm in a straight position
func (a *Arm) extendUp() {
	a.BaseVerticalServo.MoveMax(a.Driver)
	a.ClawVerticalServo.MoveMin(a.Driver)
	a.CameraVerticalServo.MoveMin(a.Driver)
}
