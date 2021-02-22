package controller

import (
	"log"

	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

const (
	LeftMotorForwardPin   = "13"
	LeftMotorBackwardPin  = "12"
	RightMotorForwardPin  = "37"
	RightMotorBackwardPin = "40"

	DirectionForward  = "forward"
	DirectionBackward = "backward"
	DirectionLeft     = "left"
	DirectionRight    = "right"
	DirectionNone     = "none"
)

type Vehicle struct {
	LeftMotor, RightMotor *gpio.MotorDriver

	LastDirection string
}

func NewVehicle(rpi *raspi.Adaptor) *Vehicle {
	vehicle := &Vehicle{
		LeftMotor:  gpio.NewMotorDriver(rpi, LeftMotorForwardPin),
		RightMotor: gpio.NewMotorDriver(rpi, RightMotorForwardPin),
	}
	vehicle.LeftMotor.SetName("left-motor")
	vehicle.LeftMotor.ForwardPin = LeftMotorForwardPin
	vehicle.LeftMotor.BackwardPin = LeftMotorBackwardPin
	vehicle.RightMotor.SetName("right-motor")
	vehicle.RightMotor.ForwardPin = RightMotorForwardPin
	vehicle.RightMotor.BackwardPin = RightMotorBackwardPin
	return vehicle
}

func (v *Vehicle) Forward() {
	v.LastDirection = DirectionForward
	v.move(DirectionForward, DirectionForward)
}

func (v *Vehicle) Backward() {
	v.LastDirection = DirectionBackward
	v.move(DirectionBackward, DirectionBackward)
}

func (v *Vehicle) Left() {
	v.LastDirection = DirectionLeft
	v.move(DirectionForward, DirectionBackward)
}

func (v *Vehicle) Right() {
	v.LastDirection = DirectionRight
	v.move(DirectionBackward, DirectionForward)
}

func (v *Vehicle) Stop() {
	//v.move(DirectionNone, DirectionNone)
	v.LeftMotor.Off()
	v.RightMotor.Off()
}

func (v *Vehicle) move(leftMotorDirection, rightMotorDirection string) {
	if err := v.LeftMotor.Direction(leftMotorDirection); err != nil {
		log.Printf("Stopping because failed to send direction to left motor: %s", err.Error())
		v.Stop()
	}
	if err := v.RightMotor.Direction(rightMotorDirection); err != nil {
		log.Printf("Stopping because failed to send direction to right motor: %s", err.Error())
		v.Stop()
	}
}
