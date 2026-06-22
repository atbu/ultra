package plugboard

import (
	"reflect"
	"testing"

	"github.com/atbu/ultra/conversion"
)

func TestEmptyPlugboardString(t *testing.T) {
	tests := []string{
		"",
		" ",
		"\n",
		"\t",
		"\n\t ",
		"\r\n\t ",
	}

	for _, testData := range tests {
		got, err := New(testData)
		if got != nil {
			t.Errorf("New() in plugboard package returned unexpected output - got=%v, expected=%v", got, nil)
		}
		if err != nil {
			t.Errorf("New() in plugboard package returned unexpected error - got=%v, expected=%v", err, nil)
		}
	}
}

func TestUnevenPlugboardString(t *testing.T) {
	tests := []string{
		"ABC",
		"ABCDE",
		"ZJFIWDL",
		"R",
		"DPWERMS",
	}

	for _, testData := range tests {
		got, err := New(testData)
		if got != nil {
			t.Errorf("New() in plugboard package returned unexpected output - got=%v, expected=%v", got, nil)
		}
		if err == nil {
			t.Errorf("New() in plugboard package returned no error when we expected it to - got=%v, expected=%v", err, ErrUnevenNumberOfCharacters)
		}
	}
}

func generatePlugboardPair(a rune, b rune) PlugboardPair {
	return PlugboardPair{int(a - 65), int(b - 65)}
}

func TestPlugboardConstructor(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput *Plugboard
	}{
		{"ABCD", &Plugboard{[]PlugboardPair{
			generatePlugboardPair('A', 'B'),
			generatePlugboardPair('C', 'D'),
		}}},
		{"RKWDPT", &Plugboard{[]PlugboardPair{
			generatePlugboardPair('R', 'K'),
			generatePlugboardPair('W', 'D'),
			generatePlugboardPair('P', 'T'),
		}}},
		{"QPALZMWOXNEICBRUVFGH", &Plugboard{[]PlugboardPair{
			generatePlugboardPair('Q', 'P'),
			generatePlugboardPair('A', 'L'),
			generatePlugboardPair('Z', 'M'),
			generatePlugboardPair('W', 'O'),
			generatePlugboardPair('X', 'N'),
			generatePlugboardPair('E', 'I'),
			generatePlugboardPair('C', 'B'),
			generatePlugboardPair('R', 'U'),
			generatePlugboardPair('V', 'F'),
			generatePlugboardPair('G', 'H'),
		}}},
	}

	for _, testData := range tests {
		got, err := New(testData.input)
		if !reflect.DeepEqual(got, testData.expectedOutput) {
			t.Errorf("New() in plugboard package returned unexpected output - got=%v, expected=%v", got, testData.expectedOutput)
		}
		if err != nil {
			t.Errorf("New() in plugboard package returned unexpected error - got=%v, expected=%v", err, nil)
		}
	}
}

func TestMapSignal(t *testing.T) {
	tests := []struct {
		plugboardConfiguration string
		inputSignal            rune
		expectedOutput         rune
	}{
		{"ABCDEFGH", 'A', 'B'},
		{"CDOFRPIT", 'D', 'C'},
		{"PLEFRD", 'Z', 'Z'},
	}

	for _, testData := range tests {
		plugboard, err := New(testData.plugboardConfiguration)
		if err != nil {
			t.Errorf("Failed to generate plugboard for MapSignal test - err=%s", err)
		}

		got := plugboard.MapSignal(conversion.ConvertCharToIndex(testData.inputSignal))
		expectedOutput := conversion.ConvertCharToIndex(testData.expectedOutput)
		if got != expectedOutput {
			t.Errorf("MapSignal() in plugboard package returned unexpected output - got=%v, expected=%v", got, expectedOutput)
		}
	}
}
