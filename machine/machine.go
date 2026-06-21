package machine

import (
	"strings"

	"github.com/atbu/ultra/conversion"
	"github.com/atbu/ultra/plugboard"
	"github.com/atbu/ultra/reflector"
	"github.com/atbu/ultra/rotor"
)

// The Machine structure defines the components of an Enigma machine.
type Machine struct {
	LeftRotor   *rotor.Rotor
	MiddleRotor *rotor.Rotor
	RightRotor  *rotor.Rotor
	FourthRotor *rotor.FourthRotor // optional, a 'nil' value indicates that there is no fourth rotor present
	Reflector   *reflector.Reflector
	Plugboard   *plugboard.Plugboard // optional, a 'nil' value indicates that no plugs are connected to the plugboard
}

// The PressKey function simulates pressing a single key on an Enigma machine. It takes the key to be pressed as a rune
// parameter, and returns the output of the character which would light up on the lampboard, also as a rune.
func (m *Machine) PressKey(key rune) rune {
	// We first convert the signal in rune form to a character index so we can more easily manipulate it.
	signal := conversion.ConvertCharToIndex(key)

	// In the first stage of cryptography, we map through the plugboard first. See the `plugboard` package for more
	// details on how this works.
	signal = m.MapThroughPlugboard(signal)

	// Pressing a key causes the rotors to rotate, so we do this accordingly.
	m.RotateRotors()

	// We map through the rotors from right to left, and we use an `inverse` value of false, since we are not passing
	// through the machine in reverse.
	signal = m.RightRotor.MapSignal(signal, false)
	signal = m.MiddleRotor.MapSignal(signal, false)
	signal = m.LeftRotor.MapSignal(signal, false)
	signal = m.MapThroughFourthRotor(signal, false)

	// We then map through the reflector, which simply maps characters to other characters using a basic wiring string.
	signal = m.Reflector.MapSignal(signal)

	// We then pass through the rotors in reverse, from left to right. Since we are passing through in reverse, we set
	// `inverse` to true so that we're using the inverse values from the wiring string.
	signal = m.MapThroughFourthRotor(signal, true)
	signal = m.LeftRotor.MapSignal(signal, true)
	signal = m.MiddleRotor.MapSignal(signal, true)
	signal = m.RightRotor.MapSignal(signal, true)

	// We then go back through the plugboard again.
	signal = m.MapThroughPlugboard(signal)

	// Then convert the character index back to a rune and output it.
	return conversion.ConvertIndexToChar(signal)
}

func (m *Machine) RotateRotors() {
	middleInNotch := m.MiddleRotor.IsInNotch()
	rightInNotch := m.RightRotor.IsInNotch()

	// It's worth noting that the fourth rotor in Enigma M4 machines doesn't rotate so that's why we don't have any
	// logic for that here.

	// The right rotor always rotates.
	m.RightRotor.Rotate()

	// Every rotor has a 'notch' value which, when reached, triggers the next rotor to rotate. Therefore, if the middle
	// rotor is at its notch value, we rotate the left hand rotor.
	if middleInNotch {
		m.LeftRotor.Rotate()
	}

	// We check if the right rotor at its notch value, if so we rotate the middle rotor.
	// Also, a quirk of Enigma is 'double stepping', where the middle rotor will also rotate itself if it is at its
	// notch position, so we account for that here.
	if middleInNotch || rightInNotch {
		m.MiddleRotor.Rotate()
	}
}

// The MapThroughPlugboard function maps a signal as a character index, through the plugboard, outputting it
// as another character index.
func (m *Machine) MapThroughPlugboard(signal int) int {
	// If we don't have a plugboard, just return the signal unchanged.
	if m.Plugboard == nil {
		return signal
	}

	// If we do have a plugboard then call its own mapping function and return its value.
	return m.Plugboard.MapSignal(signal)
}

// The MapThroughFourthRotor function maps a signal as a character index, through the fourth rotor, outputting it as
// another character index. It also takes a Boolean value called inverse as a parameter.
// If `inverse` is true, then we are passing through the fourth rotor in reverse (from the left rather than from the
// right), so we should use its inverse wiring array.
func (m *Machine) MapThroughFourthRotor(signal int, inverse bool) int {
	// If we don't have a fourth rotor, just return the signal unchanged.
	if m.FourthRotor == nil {
		return signal
	}

	// If we do have a fourth rotor then call its own mapping function and return its value.
	return m.FourthRotor.MapSignal(signal, inverse)
}

// The Process function takes a message as a string, processes each character individually by calling the PressKey
// function, and returns an output string.
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
