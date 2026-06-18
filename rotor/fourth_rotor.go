package rotor

import "github.com/atbu/ultra/conversion"

type FourthRotorConfiguration int

const (
	FourthRotorBeta FourthRotorConfiguration = iota
	FourthRotorGamma
)

type FourthRotor struct {
	Wiring        [26]int
	InverseWiring [26]int
	Position      int
}

func NewFourthRotor(configuration FourthRotorConfiguration, startingPosition rune) *FourthRotor {
	var wiringString string

	switch configuration {
	case FourthRotorBeta:
		wiringString = "LEYJVCNIXWPBQMDRTAKZGFUHOS"
	case FourthRotorGamma:
		wiringString = "FSOKANUERHMBTIYCWLQPZXVGJD"
	}

	wiring := conversion.ConvertWiringStringToArray(wiringString)

	return &FourthRotor{
		wiring,
		conversion.InverseWiringArray(wiring),
		conversion.ConvertCharToIndex(startingPosition),
	}
}

func (fr *FourthRotor) MapSignal(signal int, inverse bool) int {
	contactIn := (signal + fr.Position) % 26

	var contactOut int
	if inverse {
		contactOut = fr.InverseWiring[contactIn]
	} else {
		contactOut = fr.Wiring[contactIn]
	}

	signalOut := ((contactOut + 26) - fr.Position) % 26

	return signalOut
}
