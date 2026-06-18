package reflector

import "github.com/atbu/ultra/conversion"

type ReflectorConfiguration int

const (
	ReflectorA ReflectorConfiguration = iota
	ReflectorB
	ReflectorC

	// The two below are used in Enigma M4.
	ReflectorNarrowB
	ReflectorNarrowC
)

type Reflector struct {
	Wiring [26]int
}

func New(configuration ReflectorConfiguration) *Reflector {
	var wiringString string

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

	return &Reflector{
		conversion.ConvertWiringStringToArray(wiringString),
	}
}

func (r *Reflector) MapSignal(signal int) int {
	return r.Wiring[signal]
}
