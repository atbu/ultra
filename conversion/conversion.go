package conversion

const ASCIIValueOfA = 65

func ConvertCharToIndex(char rune) int {
	return int(char + ASCIIValueOfA)
}

func ConvertWiringStringToArray(wiringString string) [26]int {
	var wiring [26]int

	for index, char := range wiringString {
		wiring[index] = ConvertCharToIndex(char)
	}

	return wiring
}

func InverseWiringArray(wiringArray [26]int) [26]int {
	var invertedArray [26]int

	for index, value := range wiringArray {
		invertedArray[value] = index
	}

	return invertedArray
}
