package reflector

import "github.com/atbu/ultra/conversion"

type ReflectorConfiguration int

// The constants below represent the reflector configurations found in Enigma machines.
// Reflectors A, B and C are found in all Enigma machines, whereas the narrow reflectors are found in Enigma M4, and
// are used when the fourth rotor is installed.
const (
	ReflectorA ReflectorConfiguration = iota
	ReflectorB
	ReflectorC
	ReflectorNarrowB
	ReflectorNarrowC
)

// The Reflector structure simply contains the wiring array which describes how each character is mapped to another
// character.
type Reflector struct {
	Wiring [26]int
}

// The New function constructs a new Reflector from its configuration.
func New(configuration ReflectorConfiguration) *Reflector {
	var wiringString string

	// The wiring strings of each reflector are defined here.
	switch configuration {
	case ReflectorA:
		wiringString = "EJMZALYXVBWFCRQUONTSPIKHGD"
	case ReflectorB:
		wiringString = "YRUHQSLDPXNGOKMIEBFZCWVJAT"
	case ReflectorC:
		wiringString = "FVPJIAOYEDRZXWGCTKUQSBNMHL"
	case ReflectorNarrowB:
		wiringString = "ENKQAUYWJICOPBLMDXZVFTHRGS"
	case ReflectorNarrowC:
		wiringString = "RDOBJNTKVEHMLFCWZAXGYIPSUQ"
	}

	// Convert the wiring string to an array and return the constructed Reflector.
	return &Reflector{
		conversion.ConvertWiringStringToArray(wiringString),
	}
}

// The MapSignal function defines how we map a signal through a reflector - we simply look up the other character our
// input signal refers to and output it.
func (r *Reflector) MapSignal(signal int) int {
	return r.Wiring[signal]
}
