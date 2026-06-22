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

func TestConvertBetweenIndexAndChar(t *testing.T) {
	tests := []rune{
		'A',
		'B',
		'C',
		'X',
		'Y',
		'Z',
	}

	for _, testData := range tests {
		got := ConvertIndexToChar(ConvertCharToIndex(testData))
		if got != testData {
			t.Errorf(
				"ConvertIndexToChar and ConvertCharToIndex reciprocity test failed - got=%c, expected=%c",
				got,
				testData,
			)
		}
	}
}

func TestConvertWiringStringToArray(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput []int
	}{
		{"ABCDEFGHIJKLMNOPQRSTUVWXYZ", []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25}},
		{"DCBAFEIHGMLKONJRQPVUTSZYXW", []int{3, 2, 1, 0, 5, 4, 8, 7, 6, 12, 11, 10, 14, 13, 9, 17, 16, 15, 21, 20, 19, 18, 25, 24, 23, 22}},
		{"ZYXWVUTSRQPONMLKJIHGFEDCBA", []int{25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}},
	}

	for _, testData := range tests {
		got := ConvertWiringStringToArray(testData.input)
		if got != [26]int(testData.expectedOutput) {
			t.Errorf("ConvertWiringStringToArray(%s) failed,\ngot=%d,\nexpected=%d", testData.input, got, testData.expectedOutput)
		}
	}
}

func TestInverseWiringArray(t *testing.T) {
	tests := []struct {
		input          []int
		expectedOutput []int
	}{
		{[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25}, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25}},
		{[]int{19, 20, 21, 22, 23, 24, 25, 15, 16, 17, 18, 10, 11, 12, 13, 14, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3}, []int{22, 23, 24, 25, 16, 17, 18, 19, 20, 21, 11, 12, 13, 14, 15, 7, 8, 9, 10, 0, 1, 2, 3, 4, 5, 6}},
		{[]int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 25, 24, 23, 22, 21}, []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 25, 24, 23, 22, 21}},
	}

	for _, testData := range tests {
		got := InverseWiringArray([26]int(testData.input))
		if got != [26]int(testData.expectedOutput) {
			t.Errorf("InverseWiringArray(%d) failed,\ngot=%d,\nexpected=%d", testData.input, got, testData.expectedOutput)
		}
	}
}
