package controller

import (
	"gobot.io/x/gobot/drivers/i2c"
	"log"
)

type servo struct {
	Pin     string
	Default int
	Min     int
	Max     int
}

func (s *servo) Move(driver *i2c.PCA9685Driver, value int) error {
	if value > s.Max {
		value = s.Max
	} else if value < s.Min {
		value = s.Min
	}
	return s.move(driver, value)
}

func (s *servo) MoveDefault(driver *i2c.PCA9685Driver) error {
	return s.move(driver, s.Default)
}

func (s *servo) MoveMin(driver *i2c.PCA9685Driver) error {
	return s.move(driver, s.Min)
}

func (s *servo) move(driver *i2c.PCA9685Driver, value int) error {
	err := driver.ServoWrite(s.Pin, byte(value))
	if err != nil {
		log.Println(err)
	}
	return err
}
