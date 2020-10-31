package sensor

import (
	"fmt"
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

// NewUltrasonicSensor creates a new UltrasonicSensor
func NewUltrasonicSensor() *UltrasonicSensor {
	// Try to open  , panic, but if not, when we at least get some measure of guarantee
	// that later on when this needs to opened, rpio won't panic
	if err := rpio.Open(); err != nil {
		panic(fmt.Errorf("unable to open rpio: %s", err.Error()))
	}
	rpio.Close()
	return &UltrasonicSensor{
		triggerPin: rpio.Pin(TriggerPin),
		echoPin:    rpio.Pin(EchoPin),
	}
}

// MeasureDistance measures the distance by sending a high powered ultrasonic sound wave, waiting for its return
// on the echo pin, and using the time taken to calculate the distance.
func (us *UltrasonicSensor) MeasureDistance() float32 {
	_ = rpio.Open()
	defer rpio.Close()
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
	for i := 0; i <= Limit && us.echoPin.Read() != rpio.High; i++ {
		// If we hit the Limit, then the reading is no longer valid
		// Therefore, we're better off returning 0 as a way to say
		// that there's either a problem, or something is bouncing
		// the sound wave elsewhere
		if i == Limit {
			return 0
		}
	}
	start = time.Now()
	for i := 0; i <= Limit && us.echoPin.Read() != rpio.Low; i++ {
		// If we hit the limit, then the reading is no longer valid
		// Therefore, we're better off returning 0 as a way to say
		// that there's either a problem, or something is bouncing
		// the sound wave elsewhere
		if i == Limit {
			return 0
		}
	}
	end = time.Now()
	return (float32(end.UnixNano()-start.UnixNano()) * (SpeedOfSoundInCentimetersPerSecond / 2)) / float32(time.Second)
}

// MeasureDistanceReliably calls MeasureDistance thrice and returns the lowest distance measured.
// This allows a "safer" distance to be measured.
func (us *UltrasonicSensor) MeasureDistanceReliably() float32 {
	var measure, lowestMeasuredDistance float32
	for i := 0; i < 3; i++ {
		measure = us.MeasureDistance()
		fmt.Printf("measure %d: %f\n", i, measure)
		if lowestMeasuredDistance == 0 || lowestMeasuredDistance > measure {
			lowestMeasuredDistance = measure
		}
	}
	return lowestMeasuredDistance
}
