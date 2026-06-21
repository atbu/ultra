package rotor

import (
	"slices"

	"github.com/atbu/ultra/conversion"
)

type RotorConfiguration int

// These constants represent the rotors supported by the simulation.
const (
	RotorI RotorConfiguration = iota
	RotorII
	RotorIII
	RotorIV
	RotorV
	RotorVI
	RotorVII
	RotorVIII
)

// The Rotor structure contains all the components of an Enigma rotor.
type Rotor struct {
	// Used for mapping a signal through the rotor and is derived from the rotor's wiring string.
	Wiring [26]int
	// Used for mapping a signal through the rotor in reverse (from the left rather than the right) and is derived
	// from the wiring array above which is in turn derived from the rotor's wiring string.
	InverseWiring [26]int
	// The notches of the rotor, which when reached, cause the rotor to the left to rotate (except when 'double
	// stepping' occurs, see the RotateRotors function in the machine package for more information).
	Notches []int
	// The current position the rotor is in, which would be visible by a real-life user as the letter visible through
	// the window above the rotor.
	Position int
	// Used to internally shift the position of the external alphabet ring to the internal wiring/contacts to add an
	// offset.
	RingSetting int
}

// The NewRotor function constructs a new Rotor using its configuration, starting position (as a rune) and ring setting
// (also as a rune).
func NewRotor(configuration RotorConfiguration, startingPosition rune, ringSetting rune) *Rotor {
	var wiringString string
	var notchCharacters []rune

	// Each rotor's wiring string and notch position(s) (some rotors have multiple, as seen below) are defined here.
	switch configuration {
	case RotorI:
		wiringString = "EKMFLGDQVZNTOWYHXUSPAIBRCJ"
		notchCharacters = []rune{'Q'}
	case RotorII:
		wiringString = "AJDKSIRUXBLHWTMCQGZNPYFVOE"
		notchCharacters = []rune{'E'}
	case RotorIII:
		wiringString = "BDFHJLCPRTXVZNYEIWGAKMUSQO"
		notchCharacters = []rune{'V'}
	case RotorIV:
		wiringString = "ESOVPZJAYQUIRHXLNFTGKDCMWB"
		notchCharacters = []rune{'J'}
	case RotorV:
		wiringString = "VZBRGITYUPSDNHLXAWMJQOFECK"
		notchCharacters = []rune{'Z'}
	case RotorVI:
		wiringString = "JPGVOUMFYQBENHZRDKASXLICTW"
		notchCharacters = []rune{'Z', 'M'}
	case RotorVII:
		wiringString = "NZJHGRCXMYSWBOUFAIVLPEKQDT"
		notchCharacters = []rune{'Z', 'M'}
	case RotorVIII:
		wiringString = "FKQHTLXOCBJSPDZRAMEWNIUYGV"
		notchCharacters = []rune{'Z', 'M'}
	}

	// Convert the wiring string to a wiring array so we can use the data from it.
	wiring := conversion.ConvertWiringStringToArray(wiringString)

	// Convert the notch characters to notch indexes (so we can refer to them in the same way we refer to character
	// indexes so we can carry out operations with them).
	var notchIndexes []int
	for _, notchCharacter := range notchCharacters {
		notchIndexes = append(notchIndexes, conversion.ConvertCharToIndex(notchCharacter))
	}

	// Construct a new rotor using the configuration above.
	return &Rotor{
		wiring,
		// Inverse the wiring array and store it for future operations.
		conversion.InverseWiringArray(wiring),
		notchIndexes,
		conversion.ConvertCharToIndex(startingPosition),
		conversion.ConvertCharToIndex(ringSetting),
	}
}

// The Rotate function defines how we rotate a rotor. We add one to its position, then get the remainder from 26, to
// account for any overflow, for example 'Z' is represented by the index 25. If we rotate the rotor, following
// the logic below, we will wrap around to index 0, which corresponds to the character 'A'. If we didn't use the modulo
// operator, we'd get the index 26 which doesn't correspond to a character.
func (r *Rotor) Rotate() {
	r.Position = (r.Position + 1) % 26
}

// The MapSignal function defines how we map a signal through a rotor.
func (r *Rotor) MapSignal(signal int, inverse bool) int {
	// Calculate the difference that the ring setting and current position of the rotor will apply to the input signal.
	// We use (... + 26) and (... % 26) here in a similar way as in the Rotate function - to prevent overflows /
	// underflows so that we wrap around and always have an index which corresponds to a valid character.
	delta := ((r.Position + 26) - r.RingSetting) % 26

	// Apply the difference to the input signal and again account for any overflows / underflows.
	contactIn := (signal + delta) % 26

	// Map the signal using the Wiring or InverseWiring slices, depending on which way we're passing through the rotor.
	var contactOut int
	if inverse {
		contactOut = r.InverseWiring[contactIn]
	} else {
		contactOut = r.Wiring[contactIn]
	}

	// Again, apply the difference caused by the rotor's ring setting / position, accounting for overflows / underflows.
	signalOut := ((contactOut + 26) - delta) % 26

	return signalOut
}

// The IsInNotch function determines whether a rotor its at any of its notch positions or not.
func (r *Rotor) IsInNotch() bool {
	return slices.Contains(r.Notches, r.Position)
}
