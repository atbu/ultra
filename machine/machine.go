package machine

import (
	"strings"

	"github.com/atbu/ultra/conversion"
	"github.com/atbu/ultra/plugboard"
	"github.com/atbu/ultra/reflector"
	"github.com/atbu/ultra/rotor"
)

type Machine struct {
	LeftRotor   rotor.Rotor
	MiddleRotor rotor.Rotor
	RightRotor  rotor.Rotor
	FourthRotor *rotor.FourthRotor
	Reflector   reflector.Reflector
	Plugboard   *plugboard.Plugboard
}

func (m *Machine) PressKey(key rune) rune {
	signal := conversion.ConvertCharToIndex(key)

	signal = m.MapThroughPlugboard(signal)

	m.RotateRotors()

	signal = m.RightRotor.MapSignal(signal, false)
	signal = m.MiddleRotor.MapSignal(signal, false)
	signal = m.LeftRotor.MapSignal(signal, false)
	signal = m.MapThroughFourthRotor(signal, false)

	signal = m.Reflector.MapSignal(signal)

	signal = m.MapThroughFourthRotor(signal, true)
	signal = m.LeftRotor.MapSignal(signal, true)
	signal = m.MiddleRotor.MapSignal(signal, true)
	signal = m.RightRotor.MapSignal(signal, true)

	signal = m.MapThroughPlugboard(signal)

	return conversion.ConvertIndexToChar(signal)
}

func (m *Machine) RotateRotors() {
	m.RightRotor.Rotate()

	if m.MiddleRotor.IsInNotch() {
		m.LeftRotor.Rotate()
	}

	if m.MiddleRotor.IsInNotch() || m.RightRotor.IsInNotch() {
		m.MiddleRotor.Rotate()
	}
}

func (m *Machine) MapThroughPlugboard(signal int) int {
	if m.Plugboard == nil {
		return signal
	}

	return m.Plugboard.MapSignal(signal)
}

func (m *Machine) MapThroughFourthRotor(signal int, inverse bool) int {
	if m.FourthRotor == nil {
		return signal
	}

	return m.FourthRotor.MapSignal(signal, inverse)
}

func (m *Machine) Process(message string) string {
	message = strings.ToUpper(strings.TrimSpace(message))

	if len(message) == 0 {
		return ""
	}

	var output string
	for _, character := range message {
		output = output + string(m.PressKey(character))
	}

	return output
}
