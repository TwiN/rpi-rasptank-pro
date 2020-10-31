package motor

import (
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
)

const (
	LeftMotorForward   = "40"
	LeftMotorBackward  = "37"
	RightMotorForward  = "12"
	RightMotorBackward = "13"

	DirectionForward  = "forward"
	DirectionBackward = "backward"
	DirectionNone     = "none"
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
	m.move(DirectionForward, DirectionForward)
}

func (m *Motors) Backward() {
	m.move(DirectionBackward, DirectionBackward)
}

func (m *Motors) Left() {
	m.move(DirectionForward, DirectionBackward)
}

func (m *Motors) Right() {
	m.move(DirectionBackward, DirectionForward)
}

func (m *Motors) Stop() {
	//m.move(DirectionNone, DirectionNone)
	m.LeftMotor.Off()
	m.RightMotor.Off()
}

func (m *Motors) move(leftMotorDirection, rightMotorDirection string) {
	if err := m.LeftMotor.Direction(leftMotorDirection); err != nil {
		log.Printf("Stopping because failed to send direction to left motor: %s", err.Error())
		m.Stop()
	}
	if err := m.RightMotor.Direction(rightMotorDirection); err != nil {
		log.Printf("Stopping because failed to send direction to right motor: %s", err.Error())
		m.Stop()
	}
}
