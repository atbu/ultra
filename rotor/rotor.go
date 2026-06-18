package rotor

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

const ASCII_A_VALUE = 65

func convertCharToIndex(char rune) int {
	return int(char + ASCII_A_VALUE)
}

func convertWiringStringToArray(wiringString string) [26]int {
	var wiring [26]int

	for index, char := range wiringString {
		wiring[index] = convertCharToIndex(char)
	}

	return wiring
}

func inverseWiringArray(wiringArray [26]int) [26]int {
	var invertedArray [26]int

	for index, value := range wiringArray {
		invertedArray[value] = index
	}

	return invertedArray
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

	wiring := convertWiringStringToArray(wiringString)

	var notchIndexes []int
	for _, notchCharacter := range notchCharacters {
		notchIndexes = append(notchIndexes, convertCharToIndex(notchCharacter))
	}

	return &Rotor{
		wiring,
		inverseWiringArray(wiring),
		notchIndexes,
		convertCharToIndex(startingPosition),
		convertCharToIndex(ringSetting),
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
