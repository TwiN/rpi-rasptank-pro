package input

import (
	"fmt"
	"time"

	"github.com/TwiN/rpi-rasptank-pro/controller"
	"gobot.io/x/gobot/platforms/keyboard"
)

type Keyboard struct {
	Driver *keyboard.Driver
}

func NewKeyboard() *Keyboard {
	return &Keyboard{Driver: keyboard.NewDriver()}
}

func (k *Keyboard) HandleKeyboardEvents(vehicle *controller.Vehicle) {
	k.Driver.On(keyboard.Key, func(data interface{}) {
		key := data.(keyboard.KeyEvent)
		switch key.Key {
		case keyboard.ArrowUp:
			fmt.Println("UP")
			vehicle.Forward()
			time.Sleep(50 * time.Millisecond)
			vehicle.Stop()
		case keyboard.ArrowDown:
			fmt.Println("DOWN")
			vehicle.Backward()
			time.Sleep(50 * time.Millisecond)
			vehicle.Stop()
		case keyboard.ArrowRight:
			fmt.Println("RIGHT")
			vehicle.Right()
			time.Sleep(50 * time.Millisecond)
			vehicle.Stop()
		case keyboard.ArrowLeft:
			fmt.Println("LEFT")
			vehicle.Left()
			time.Sleep(50 * time.Millisecond)
			vehicle.Stop()
		default:
			fmt.Println("UNSUPPORTED KEY:", key, key.Char)
		}
	})
}
