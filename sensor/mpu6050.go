package sensor

import (
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

	calibratedAccelerometerOffsetX, calibratedAccelerometerOffsetY, calibratedAccelerometerOffsetZ int16
	calibratedGyroscopeOffsetX, calibratedGyroscopeOffsetY, calibratedGyroscopeOffsetZ             int16
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
	m.calibratedAccelerometerOffsetX = m.Driver.Accelerometer.X
	m.calibratedAccelerometerOffsetY = m.Driver.Accelerometer.Y
	m.calibratedAccelerometerOffsetZ = m.Driver.Accelerometer.Z
	m.calibratedGyroscopeOffsetX = m.Driver.Gyroscope.X
	m.calibratedGyroscopeOffsetY = m.Driver.Gyroscope.Y
	m.calibratedGyroscopeOffsetZ = m.Driver.Gyroscope.Z
	log.Printf("[Calibrate] ax=%d; ay=%d; az=%d; gx=%d; gy=%d; gz=%d", m.Driver.Accelerometer.X, m.Driver.Accelerometer.Y, m.Driver.Accelerometer.Z, m.Driver.Gyroscope.X, m.Driver.Gyroscope.Y, m.Driver.Gyroscope.Z)
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
	// Normalize data
	m.Driver.Accelerometer.X = (m.Driver.Accelerometer.X - m.calibratedAccelerometerOffsetX) / 1000
	m.Driver.Accelerometer.Y = (m.Driver.Accelerometer.Y - m.calibratedAccelerometerOffsetY) / 1000
	m.Driver.Accelerometer.Z = (m.Driver.Accelerometer.Z - m.calibratedAccelerometerOffsetZ) / 1000
	m.Driver.Gyroscope.X -= m.calibratedGyroscopeOffsetX
	m.Driver.Gyroscope.Y -= m.calibratedGyroscopeOffsetY
	m.Driver.Gyroscope.Z -= m.calibratedGyroscopeOffsetZ

	// pitch is used to figure out whether the bot fell forward or backward
	pitch := -(math.Atan2(float64(m.Driver.Accelerometer.X), math.Sqrt(float64(m.Driver.Accelerometer.Y*m.Driver.Accelerometer.Y+m.Driver.Accelerometer.Z*m.Driver.Accelerometer.Z))) * 180.0) / math.Pi
	// roll is used to figure out whether the bot fell on its right or left side
	// A negative roll means that the bot fell on its right side
	// A positive roll means that the bot fell on its left side
	roll := (math.Atan2(float64(m.Driver.Accelerometer.Y), float64(m.Driver.Accelerometer.Z)) * 180.0) / math.Pi

	if pitch != 0 || roll != 0 {
		log.Printf("[FallDetected] pitch=%.0f; roll=%.0f", pitch, roll)
	}

	if pitch > FallDetectionThreshold || pitch < -FallDetectionThreshold || roll > FallDetectionThreshold || roll < -FallDetectionThreshold {
		if roll < 0 {
			log.Println("[FallDetected] fell on the right side")
		} else {
			log.Println("[FallDetected] fell on the left side")
		}
		return true, roll < 0
	}
	return false, false
}
