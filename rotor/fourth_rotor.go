package rotor

type FourthRotorConfiguration int

const (
	Beta FourthRotorConfiguration = iota
	Gamma
)

type FourthRotor struct {
	Wiring        [26]int
	InverseWiring [26]int
	Position      int
}

func NewFourthRotor(configuration FourthRotorConfiguration, startingPosition rune) *FourthRotor {
	var wiringString string

	switch configuration {
	case Beta:
		wiringString = "LEYJVCNIXWPBQMDRTAKZGFUHOS"
	case Gamma:
		wiringString = "FSOKANUERHMBTIYCWLQPZXVGJD"
	}

	wiring := convertWiringStringToArray(wiringString)

	return &FourthRotor{
		wiring,
		inverseWiringArray(wiring),
		convertCharToIndex(startingPosition),
	}
}
