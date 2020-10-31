package sensor

import (
	"github.com/stianeikeland/go-rpio"
	"time"
)

const (
	EchoPin    = 8
	TriggerPin = 11

	// SpeedOfSoundInCentimetersPerSecond is the speed of sound in centimeters per second
	SpeedOfSoundInCentimetersPerSecond = 34300

	// Limit is the maximum amount of iterations to wait for a response on the echo pin
	// This is to make sure that in case the sound wave is never received, the function
	// won't hang indefinitely
	Limit = 100000
)

// UltrasonicSensor is a sensor for HC-SR04
type UltrasonicSensor struct {
	triggerPin rpio.Pin
	echoPin    rpio.Pin
}

func NewUltrasonicSensor() *UltrasonicSensor {
	return &UltrasonicSensor{
		triggerPin: rpio.Pin(TriggerPin),
		echoPin:    rpio.Pin(EchoPin),
	}
}

func (us *UltrasonicSensor) MeasureDistance() float32 {
	// Set echo pin as INPUT and trigger pin as OUTPUT
	us.echoPin.Input()
	us.triggerPin.Output()
	// Clear trigger pin
	us.triggerPin.Low()
	time.Sleep(5 * time.Microsecond)
	// Transmit HIGH output from trigger pin for 10μs
	us.triggerPin.High()
	time.Sleep(10 * time.Microsecond)
	us.triggerPin.Low()

	var start, end time.Time
	for i := 0; i < Limit && us.echoPin.Read() != rpio.High; i++ {
	}
	start = time.Now()
	for i := 0; i < Limit && us.echoPin.Read() != rpio.Low; i++ {
	}
	end = time.Now()
	return (float32(end.UnixNano()-start.UnixNano()) * (SpeedOfSoundInCentimetersPerSecond / 2)) / float32(time.Second)
}
