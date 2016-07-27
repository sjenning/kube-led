package minnowboard

import (
	"fmt"

	"github.com/davecheney/gpio"
)

var (
	LED_D2 = 360
)

type Pin struct {
	pin gpio.Pin
}

func NewPin(pin int) (*Pin, error) {
	gpioPin, err := gpio.OpenPin(pin, gpio.ModeOutput)
	if err != nil {
		return nil, fmt.Errorf("Error opening pin: %s\n", err)
	}
	return &Pin{pin: gpioPin}, nil
}

func (p *Pin) On() {
	p.pin.Clear()
}

func (p *Pin) Off() {
	p.pin.Set()
}

func (p *Pin) Close() {
	p.pin.Clear()
	p.pin.Close()
}
