package motor

import (
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

const (
	LeftMotorForward   = "40"
	LeftMotorBackward  = "37"
	RightMotorForward  = "12"
	RightMotorBackward = "13"

	DirectionForward  = "forward"
	DirectionBackward = "backward"
)

type Motors struct {
	LeftMotor, RightMotor *gpio.MotorDriver
}

func NewMotors(rpi *raspi.Adaptor) *Motors {
	motors := &Motors{
		LeftMotor:  gpio.NewMotorDriver(rpi, LeftMotorForward),
		RightMotor: gpio.NewMotorDriver(rpi, RightMotorForward),
	}
	motors.LeftMotor.SetName("left-motor")
	motors.LeftMotor.ForwardPin = LeftMotorForward
	motors.LeftMotor.BackwardPin = LeftMotorBackward
	motors.RightMotor.SetName("right-motor")
	motors.RightMotor.ForwardPin = RightMotorForward
	motors.RightMotor.BackwardPin = RightMotorBackward
	return motors
}

func (m *Motors) Forward() {
	m.LeftMotor.Direction(DirectionForward)
	m.RightMotor.Direction(DirectionForward)
}

func (m *Motors) Backward() {
	m.LeftMotor.Direction(DirectionBackward)
	m.RightMotor.Direction(DirectionBackward)
}

func (m *Motors) Left() {
	m.LeftMotor.Direction(DirectionBackward)
	m.RightMotor.Direction(DirectionForward)
}

func (m *Motors) Right() {
	m.LeftMotor.Direction(DirectionBackward)
	m.RightMotor.Direction(DirectionForward)
}

func (m *Motors) Stop() {
	m.LeftMotor.Direction(DirectionBackward)
	m.RightMotor.Direction(DirectionForward)
}
