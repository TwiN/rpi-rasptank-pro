package controller

import (
	"errors"
	"fmt"
	"log"

	"gobot.io/x/gobot/drivers/i2c"
)

type servo struct {
	Pin     string
	Default int
	Min     int
	Max     int

	current    int
	checkPoint int
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

func (s *servo) MoveMax(driver *i2c.PCA9685Driver) error {
	return s.move(driver, s.Max)
}

func (s *servo) MoveMin(driver *i2c.PCA9685Driver) error {
	return s.move(driver, s.Min)
}

func (s *servo) Checkpoint(checkpoint bool) *servo {
	if checkpoint {
		s.checkPoint = s.current
	}
	return s
}

func (s *servo) MoveToCheckpoint(driver *i2c.PCA9685Driver) error {
	if s.checkPoint == -1 {
		return errors.New("checkpoint not set")
	}
	return s.move(driver, s.checkPoint)
}

func (s *servo) MoveToCheckpointAndClear(driver *i2c.PCA9685Driver) error {
	err := s.MoveToCheckpoint(driver)
	s.checkPoint = -1
	return err
}

func (s *servo) move(driver *i2c.PCA9685Driver, value int) error {
	s.current = value
	fmt.Printf("pin=%s moved to value=%d\n", s.Pin, value)
	err := driver.ServoWrite(s.Pin, byte(value))
	if err != nil {
		log.Println(err)
	}
	//time.Sleep(10*time.Millisecond)
	return err
}
