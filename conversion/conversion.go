package conversion

// The ASCIIValueOfA that we use for converting between characters and their indexes, so that we can
// translate the ASCII value of 65 to an index of 0 and vice versa.
const ASCIIValueOfA = 65

// ConvertCharToIndex converts a character/rune like to its zero-indexed alphabetical value, e.g. 'A' = 0, 'Z' = 25.
func ConvertCharToIndex(char rune) int {
	return int(char - ASCIIValueOfA)
}

// ConvertIndexToChar converts a zero-indexed alphabetical value to its corresponding character/rune, e.g. 0 = 'A',
// 25 = 'Z'.
func ConvertIndexToChar(index int) rune {
	return rune(index + 65)
}

// ConvertWiringStringToArray converts a wiring string, e.g. "EKMFLGDQVZNTOWYHXUSPAIBRCJ", to an array of corresponding
// indexes.
//
// The values in this wiring string and therefore the array relate to the other characters that each alphabetical
// character corresponds to, so in the example above, 'A' is translated to 'E', 'B' is translated to 'K' and so on.
//
// Each rotor has its own wiring string and therefore a different way it maps characters to other characters.
func ConvertWiringStringToArray(wiringString string) [26]int {
	var wiring [26]int

	for index, char := range wiringString {
		wiring[index] = ConvertCharToIndex(char)
	}

	return wiring
}

// InverseWiringArray converts an array of character indexes to another similar array, but where the keys and values
// are swapped.
//
// This is used where a value passes through a rotor one way, then later on in the process it passes through the same
// rotor but in the opposite direction. We pre-compute the values and store them separately so that we don't have
// to find inverse values on the fly when we're trying to process something.
func InverseWiringArray(wiringArray [26]int) [26]int {
	var invertedArray [26]int

	for index, value := range wiringArray {
		invertedArray[value] = index
	}

	return invertedArray
}
