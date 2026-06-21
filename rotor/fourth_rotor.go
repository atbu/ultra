package rotor

import "github.com/atbu/ultra/conversion"

type FourthRotorConfiguration int

// The constants below contain the possible fourth rotor configurations.
const (
	FourthRotorBeta FourthRotorConfiguration = iota
	FourthRotorGamma
)

// The FourthRotor structure contains all the components of a fourth rotor. Note that, compared to a standard Enigma
// rotor, a fourth rotor does not contain any ring settings, and also doesn't contain any notches because it doesn't
// rotate.
type FourthRotor struct {
	// Used for mapping a signal through the rotor and is derived from the rotor's wiring string.
	Wiring [26]int
	// Used for mapping a signal through the rotor in reverse (from the left rather than the right) and is derived
	// from the wiring array above which is in turn derived from the rotor's wiring string.
	InverseWiring [26]int
	// The current position the rotor is in.
	Position int
}

// The NewFourthRotor function constructs a new FourthRotor using its configuration and starting position (as a rune).
func NewFourthRotor(configuration FourthRotorConfiguration, startingPosition rune) *FourthRotor {
	var wiringString string

	// Each rotor's wiring string is defined here.
	switch configuration {
	case FourthRotorBeta:
		wiringString = "LEYJVCNIXWPBQMDRTAKZGFUHOS"
	case FourthRotorGamma:
		wiringString = "FSOKANUERHMBTIYCWLQPZXVGJD"
	}

	// Convert the wiring string to a wiring array so we can use the data from it.
	wiring := conversion.ConvertWiringStringToArray(wiringString)

	// Return the constructed FourthRotor.
	return &FourthRotor{
		wiring,
		// Inverse the wiring array and store it for future operations.
		conversion.InverseWiringArray(wiring),
		conversion.ConvertCharToIndex(startingPosition),
	}
}

// The MapSignal function defines how we map a signal through a fourth rotor.
func (fr *FourthRotor) MapSignal(signal int, inverse bool) int {
	// Apply the difference caused by the rotor's current position to the signal, and account for overflows.
	contactIn := (signal + fr.Position) % 26

	// Map the signal using the Wiring or InverseWiring slices, depending on which way we're passing through the rotor.
	var contactOut int
	if inverse {
		contactOut = fr.InverseWiring[contactIn]
	} else {
		contactOut = fr.Wiring[contactIn]
	}

	// Again, apply the difference caused by the rotor's current position to the signal, and account for overflows.
	signalOut := ((contactOut + 26) - fr.Position) % 26

	return signalOut
}
