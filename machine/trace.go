package machine

import "github.com/atbu/ultra/conversion"

// The Component type names a part of the machine that a signal passes through.
type Component string

const (
	ComponentPlugboard   Component = "plugboard"
	ComponentRightRotor  Component = "right"
	ComponentMiddleRotor Component = "middle"
	ComponentLeftRotor   Component = "left"
	ComponentFourthRotor Component = "fourth"
	ComponentReflector   Component = "reflector"
)

// The TraceStep structure records the letter entering and leaving one component during a key press.
type TraceStep struct {
	Component Component
	In        rune
	Out       rune
}

// The PressKeyWithTrace function behaves exactly like PressKey, but also returns every step the signal took through
// the machine, in order. This lets a user interface draw the path of the signal from the key to the lamp.
//
// The plugboard is always included in the trace (with the letter unchanged if no plug is connected), but the fourth
// rotor only appears if one is fitted.
func (m *Machine) PressKeyWithTrace(key rune) (rune, []TraceStep) {
	var trace []TraceStep
	signal := conversion.ConvertCharToIndex(key)

	step := func(component Component, mapSignal func(int) int) {
		out := mapSignal(signal)
		trace = append(trace, TraceStep{
			component, conversion.ConvertIndexToChar(signal), conversion.ConvertIndexToChar(out),
		})
		signal = out
	}

	forward := func(r interface{ MapSignal(int, bool) int }) func(int) int {
		return func(s int) int { return r.MapSignal(s, false) }
	}
	backward := func(r interface{ MapSignal(int, bool) int }) func(int) int {
		return func(s int) int { return r.MapSignal(s, true) }
	}

	step(ComponentPlugboard, m.MapThroughPlugboard)

	// As in PressKey, the rotors step before the signal passes through them.
	m.RotateRotors()

	step(ComponentRightRotor, forward(m.RightRotor))
	step(ComponentMiddleRotor, forward(m.MiddleRotor))
	step(ComponentLeftRotor, forward(m.LeftRotor))
	if m.FourthRotor != nil {
		step(ComponentFourthRotor, forward(m.FourthRotor))
	}

	step(ComponentReflector, m.Reflector.MapSignal)

	if m.FourthRotor != nil {
		step(ComponentFourthRotor, backward(m.FourthRotor))
	}
	step(ComponentLeftRotor, backward(m.LeftRotor))
	step(ComponentMiddleRotor, backward(m.MiddleRotor))
	step(ComponentRightRotor, backward(m.RightRotor))

	step(ComponentPlugboard, m.MapThroughPlugboard)

	return conversion.ConvertIndexToChar(signal), trace
}
