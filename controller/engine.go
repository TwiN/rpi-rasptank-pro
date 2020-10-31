package controller

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

type Engine struct {
	LeftMotor, RightMotor *gpio.MotorDriver
}

func NewEngine(rpi *raspi.Adaptor) *Engine {
	engine := &Engine{
		LeftMotor:  gpio.NewMotorDriver(rpi, LeftMotorForward),
		RightMotor: gpio.NewMotorDriver(rpi, RightMotorForward),
	}
	engine.LeftMotor.SetName("left-motor")
	engine.LeftMotor.ForwardPin = LeftMotorForward
	engine.LeftMotor.BackwardPin = LeftMotorBackward
	engine.RightMotor.SetName("right-motor")
	engine.RightMotor.ForwardPin = RightMotorForward
	engine.RightMotor.BackwardPin = RightMotorBackward
	return engine
}

func (e *Engine) Forward() {
	e.move(DirectionForward, DirectionForward)
}

func (e *Engine) Backward() {
	e.move(DirectionBackward, DirectionBackward)
}

func (e *Engine) Left() {
	e.move(DirectionForward, DirectionBackward)
}

func (e *Engine) Right() {
	e.move(DirectionBackward, DirectionForward)
}

func (e *Engine) Stop() {
	//e.move(DirectionNone, DirectionNone)
	e.LeftMotor.Off()
	e.RightMotor.Off()
}

func (e *Engine) move(leftMotorDirection, rightMotorDirection string) {
	if err := e.LeftMotor.Direction(leftMotorDirection); err != nil {
		log.Printf("Stopping because failed to send direction to left motor: %s", err.Error())
		e.Stop()
	}
	if err := e.RightMotor.Direction(rightMotorDirection); err != nil {
		log.Printf("Stopping because failed to send direction to right motor: %s", err.Error())
		e.Stop()
	}
}
