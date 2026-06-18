package reflector

import "github.com/atbu/ultra/conversion"

type ReflectorConfiguration int

const (
	ReflectorA ReflectorConfiguration = iota
	ReflectorB
	ReflectorC

	// The two below are used in Enigma M4.
	ReflectorNarrowA
	ReflectorNarrowB
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
	case ReflectorNarrowA:
		wiringString = "ENKQAUYWJICOPBLMDXZVFTHRGS"
	case ReflectorNarrowB:
		wiringString = "RDOBJNTKVEHMLFCWZAXGYIPSUQ"
	}

	return &Reflector{
		conversion.ConvertWiringStringToArray(wiringString),
	}
}
