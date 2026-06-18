package conversion

import "testing"

func TestConvertCharToIndex(t *testing.T) {
	tests := []struct {
		input          rune
		expectedOutput int
	}{
		{'A', 0},
		{'B', 1},
		{'C', 2},
		{'Z', 25},
	}

	for _, testData := range tests {
		got := ConvertCharToIndex(testData.input)
		if got != testData.expectedOutput {
			t.Errorf(
				"ConvertCharToIndex(%c) failed - got=%d, expected=%d",
				testData.input,
				got,
				testData.expectedOutput,
			)
		}
	}
}

func TestConvertIndexToChar(t *testing.T) {
	tests := []struct {
		input          int
		expectedOutput rune
	}{
		{0, 'A'},
		{1, 'B'},
		{2, 'C'},
		{25, 'Z'},
	}

	for _, testData := range tests {
		got := ConvertIndexToChar(testData.input)
		if got != testData.expectedOutput {
			t.Errorf(
				"ConvertIndexToChar(%d) failed - got=%c, expected=%c",
				testData.input,
				got,
				testData.expectedOutput,
			)
		}
	}
}
