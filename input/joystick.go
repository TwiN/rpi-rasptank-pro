package input

import (
	"fmt"
	"github.com/TwinProduction/rpi-rasptank-pro/controller"
	"gobot.io/x/gobot/platforms/joystick"
)

type Joystick struct {
	Driver *joystick.Driver
}

func NewJoystick(adaptor *joystick.Adaptor) *Joystick {
	return &Joystick{
		Driver: joystick.NewDriver(adaptor, joystick.Xbox360),
	}
}

func (j *Joystick) Handle(vehicle *controller.Vehicle, arm *controller.Arm) {
	j.Driver.On(joystick.UpPress, func(data interface{}) {
		vehicle.Forward()
	})
	j.Driver.On(joystick.UpRelease, func(data interface{}) {
		vehicle.Stop()
	})
	j.Driver.On(joystick.DownPress, func(data interface{}) {
		vehicle.Backward()
	})
	j.Driver.On(joystick.DownRelease, func(data interface{}) {
		vehicle.Stop()
	})
	j.Driver.On(joystick.LeftPress, func(data interface{}) {
		vehicle.Left()
	})
	j.Driver.On(joystick.LeftRelease, func(data interface{}) {
		vehicle.Stop()
	})
	j.Driver.On(joystick.RightPress, func(data interface{}) {
		vehicle.Right()
	})
	j.Driver.On(joystick.RightRelease, func(data interface{}) {
		vehicle.Stop()
	})

	j.Driver.On(joystick.LeftX, func(data interface{}) {
		arm.MoveBaseHorizontal(int(data.(int16) / 500))
	})
	j.Driver.On(joystick.LeftY, func(data interface{}) {
		arm.MoveBaseVertical(int(data.(int16) / 250))
	})
	j.Driver.On(joystick.RightX, func(data interface{}) {
		fmt.Println("right_x:", data)

		arm.MoveClawVertical(int(data.(int16) / 250))
	})
	//j.Driver.On(joystick.RightY, func(data interface{}) {
	//	fmt.Println("right_y:", data)
	//	arm.MoveClaw(int(data.(int16) / 250))
	//})
	j.Driver.On(joystick.APress, func(data interface{}) {
		arm.ClawGrab()
	})
	j.Driver.On(joystick.ARelease, func(data interface{}) {
		arm.ClawRelease()
	})
}
