package main

import (
	"fmt"

	"github.com/piotrjaromin/gpio"
)

type Direction uint

const (
	Input Direction = iota
	Output
)

type State uint

const (
	High State = iota
	Low
)

type PinConfig struct {
	direction Direction
	state     State
}

type PinState struct {
	config PinConfig
	pin    gpio.Pin
}

var pins map[uint]PinState

func initGpio(pinNum uint, direction Direction) error {
	var newPin PinState
	var err error

	pinfo, keyExists := pins[pinNum]

	if keyExists {
		if pinfo.config.direction == direction {
			return nil
		}
		pinfo.pin.Close()
	}

	switch direction {
	case Output:
		newPin.pin, err = gpio.NewOutput(pinNum, false)
		newPin.config.state = Low
	case Input:
		newPin.pin, err = gpio.NewInput(pinNum)
	}

	if err != nil {
		return err
	}

	newPin.config.direction = direction

	pins[pinNum] = newPin

	return nil
}

func setOutput(pinNum uint, state State) error {
	var err error

	pinfo, keyExists := pins[pinNum]

	if !keyExists {
		return fmt.Errorf("setOutput: pin %v has not been set", pinNum)
	}

	if pinfo.config.direction == Input {
		return fmt.Errorf("setOutput: pin %v is configured as input", pinNum)
	}

	switch state {
	case High:
		err = pinfo.pin.High()
	case Low:
		err = pinfo.pin.Low()
	}

	return err
}

func readInput(pinNum uint) (State, error) {
	pinfo, keyExists := pins[pinNum]

	if !keyExists {
		return Low, fmt.Errorf("readInput: pin %v has not been set", pinNum)
	}

	if pinfo.config.direction == Output {
		return Low, fmt.Errorf("readInput: pin %v is configured as output", pinNum)
	}

	val, err := pinfo.pin.Read()

	sval := State(val)

	return sval, err
}

func main() {
	initGpio(0, Output)
}
