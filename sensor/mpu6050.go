package sensor

import (
	"fmt"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
	"math"
)

const (
	mpu6050Bus     = 1
	mpu6050Address = 0x68

	FallDetectionThreshold = 110
)

type MPU6050GyroscopeAccelerometerTemperatureSensor struct {
	Driver *i2c.MPU6050Driver

	calibratedAccelerometerX, calibratedAccelerometerY, calibratedAccelerometerZ int16
	calibratedGyroscopeX, calibratedGyroscopeY, calibratedGyroscopeZ             int16
}

func NewMPU6050GyroscopeAccelerometerTemperatureSensor(rpi *raspi.Adaptor) *MPU6050GyroscopeAccelerometerTemperatureSensor {
	return &MPU6050GyroscopeAccelerometerTemperatureSensor{
		Driver: i2c.NewMPU6050Driver(rpi, i2c.WithBus(mpu6050Bus), i2c.WithAddress(mpu6050Address)),
	}
}

func (m *MPU6050GyroscopeAccelerometerTemperatureSensor) GetTemperature() (int16, error) {
	if err := m.Driver.GetData(); err != nil {
		return 0, err
	}
	return m.Driver.Temperature, nil
}

func (m *MPU6050GyroscopeAccelerometerTemperatureSensor) Calibrate() error {
	if err := m.Driver.GetData(); err != nil {
		return err
	}
	if m.Driver.Accelerometer.X == 0 && m.Driver.Accelerometer.Y == 0 && m.Driver.Accelerometer.Z == 0 {
		// data not there yet, let's try again
		return m.Calibrate()
	}
	m.calibratedAccelerometerX = m.Driver.Accelerometer.X
	m.calibratedAccelerometerY = m.Driver.Accelerometer.Y
	m.calibratedAccelerometerZ = m.Driver.Accelerometer.Z
	m.calibratedGyroscopeX = m.Driver.Gyroscope.X
	m.calibratedGyroscopeY = m.Driver.Gyroscope.Y
	m.calibratedGyroscopeZ = m.Driver.Gyroscope.Z
	fmt.Printf("[calibrated] ax=%d; ay=%d; az=%d; gx=%d; gy=%d; gz=%d\n", m.Driver.Accelerometer.X, m.Driver.Accelerometer.Y, m.Driver.Accelerometer.Z, m.Driver.Gyroscope.X, m.Driver.Gyroscope.Y, m.Driver.Gyroscope.Z)
	return nil
}

// FallDetected detects if a fall was detected as well as whether the fall was likely on the right side
func (m *MPU6050GyroscopeAccelerometerTemperatureSensor) FallDetected() (bool, bool) {
	if err := m.Driver.GetData(); err != nil {
		return false, false
	}
	if m.Driver.Accelerometer.X == 0 && m.Driver.Accelerometer.Y == 0 && m.Driver.Accelerometer.Z == 0 {
		// data not there yet, let's try again
		return m.FallDetected()
	}
	//fmt.Printf("[before] ax=%d; ay=%d; az=%d; gx=%d; gy=%d; gz=%d\n", m.Driver.Accelerometer.X, m.Driver.Accelerometer.Y, m.Driver.Accelerometer.Z, m.Driver.Gyroscope.X, m.Driver.Gyroscope.Y, m.Driver.Gyroscope.Z)
	// Normalize data
	m.Driver.Accelerometer.X = (m.Driver.Accelerometer.X - m.calibratedAccelerometerX) / 1000
	m.Driver.Accelerometer.Y = (m.Driver.Accelerometer.Y - m.calibratedAccelerometerY) / 1000
	m.Driver.Accelerometer.Z = (m.Driver.Accelerometer.Z - m.calibratedAccelerometerZ) / 1000
	m.Driver.Gyroscope.X -= m.calibratedGyroscopeX
	m.Driver.Gyroscope.Y -= m.calibratedGyroscopeY
	m.Driver.Gyroscope.Z -= m.calibratedGyroscopeZ
	//fmt.Printf("[after] ax=%d; ay=%d; az=%d; gx=%d; gy=%d; gz=%d\n", m.Driver.Accelerometer.X, m.Driver.Accelerometer.Y, m.Driver.Accelerometer.Z, m.Driver.Gyroscope.X, m.Driver.Gyroscope.Y, m.Driver.Gyroscope.Z)

	// pitch is used to figure out whether the bot fell forward or backward
	pitch := -(math.Atan2(float64(m.Driver.Accelerometer.X), math.Sqrt(float64(m.Driver.Accelerometer.Y*m.Driver.Accelerometer.Y+m.Driver.Accelerometer.Z*m.Driver.Accelerometer.Z))) * 180.0) / math.Pi
	// roll is used to figure out whether the bot fell on its right or left side
	// A negative roll means that the bot fell on its right side
	// A positive roll means that the bot fell on its left side
	roll := (math.Atan2(float64(m.Driver.Accelerometer.Y), float64(m.Driver.Accelerometer.Z)) * 180.0) / math.Pi

	fmt.Printf("pitch=%.0f; roll=%.0f\n", pitch, roll)

	if pitch > FallDetectionThreshold || pitch < -FallDetectionThreshold || roll > FallDetectionThreshold || roll < -FallDetectionThreshold {
		if roll < 0 {
			log.Println("fell on the right side")
		} else {
			log.Println("fell on the left side")
		}
		return true, roll < 0
	}
	return false, false
}
