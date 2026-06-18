package rotor

import "github.com/atbu/ultra/conversion"

type RotorConfiguration int

// Represents the possible rotors supported by the simulation.
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

type Rotor struct {
	Wiring        [26]int
	InverseWiring [26]int
	Notches       []int
	Position      int
	RingSetting   int
}

func NewRotor(configuration RotorConfiguration, startingPosition rune, ringSetting rune) *Rotor {
	var wiringString string
	var notchCharacters []rune

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

	wiring := conversion.ConvertWiringStringToArray(wiringString)

	var notchIndexes []int
	for _, notchCharacter := range notchCharacters {
		notchIndexes = append(notchIndexes, conversion.ConvertCharToIndex(notchCharacter))
	}

	return &Rotor{
		wiring,
		conversion.InverseWiringArray(wiring),
		notchIndexes,
		conversion.ConvertCharToIndex(startingPosition),
		conversion.ConvertCharToIndex(ringSetting),
	}
}

func (r *Rotor) Rotate() {
	r.Position = (r.Position + 1) % 26
}

func (r *Rotor) MapSignal(signal int, inverse bool) int {
	delta := ((r.Position + 26) - r.RingSetting) % 26
	contactIn := (signal + delta) % 26

	var contactOut int
	if inverse {
		contactOut = r.InverseWiring[contactIn]
	} else {
		contactOut = r.Wiring[contactIn]
	}

	signalOut := ((contactOut + 26) - delta) % 26

	return signalOut
}
