package plugboard

import (
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

func TestPlugboardConstructor(t *testing.T) {
	tests := []string{
		"ABCD",
		"RKWDPT",
		"QPALZMWOXNEICBRUVFGH",
	}

	for _, input := range tests {
		got, err := New(input)
		if err != nil {
			t.Errorf("New() in plugboard package returned unexpected error - got=%v, expected=%v", err, nil)
		}

		// Each pair in the plugboard is derived from two consecutive characters in the input string, so we can
		// derive the expected pairs from the input rather than hand-listing them.
		wantPairs := len(input) / 2
		if len(got.Configuration) != wantPairs {
			t.Errorf("New() in plugboard package returned unexpected number of pairs - got=%v, expected=%v", len(got.Configuration), wantPairs)
			continue
		}

		for i, pair := range got.Configuration {
			wantA := conversion.ConvertCharToIndex(rune(input[i*2]))
			wantB := conversion.ConvertCharToIndex(rune(input[i*2+1]))
			if pair.a != wantA || pair.b != wantB {
				t.Errorf("New() in plugboard package returned unexpected pair at index %d - got=%v, expected=%v", i, pair, PlugboardPair{wantA, wantB})
			}
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
