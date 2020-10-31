package controller

import (
	"fmt"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
)

const (
	LeftMotorForward   = "13"
	LeftMotorBackward  = "12"
	RightMotorForward  = "37"
	RightMotorBackward = "40"

	DirectionForward  = "forward"
	DirectionBackward = "backward"
	DirectionNone     = "none"
)

type Vehicle struct {
	LeftMotor, RightMotor *gpio.MotorDriver
}

func NewVehicle(rpi *raspi.Adaptor) *Vehicle {
	vehicle := &Vehicle{
		LeftMotor:  gpio.NewMotorDriver(rpi, LeftMotorForward),
		RightMotor: gpio.NewMotorDriver(rpi, RightMotorForward),
	}
	vehicle.LeftMotor.SetName("left-motor")
	vehicle.LeftMotor.ForwardPin = LeftMotorForward
	vehicle.LeftMotor.BackwardPin = LeftMotorBackward
	vehicle.RightMotor.SetName("right-motor")
	vehicle.RightMotor.ForwardPin = RightMotorForward
	vehicle.RightMotor.BackwardPin = RightMotorBackward
	return vehicle
}

func (v *Vehicle) Forward() {
	v.move(DirectionForward, DirectionForward)
}

func (v *Vehicle) Backward() {
	v.move(DirectionBackward, DirectionBackward)
}

func (v *Vehicle) Left() {
	v.move(DirectionForward, DirectionBackward)
}

func (v *Vehicle) Right() {
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
	fmt.Printf("current speed: %d and %d\n", v.LeftMotor.CurrentSpeed, v.RightMotor.CurrentSpeed)
}
