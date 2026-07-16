package rotor

import (
	"testing"

	"github.com/atbu/ultra/conversion"
)

func TestNewFourthRotor(t *testing.T) {
	tests := []struct {
		configuration    FourthRotorConfiguration
		wiringString     string
		startingPosition rune
	}{
		{FourthRotorBeta, "LEYJVCNIXWPBQMDRTAKZGFUHOS", 'A'},
		{FourthRotorGamma, "FSOKANUERHMBTIYCWLQPZXVGJD", 'K'},
	}

	for _, testData := range tests {
		got := NewFourthRotor(testData.configuration, testData.startingPosition)

		wantWiring := conversion.ConvertWiringStringToArray(testData.wiringString)
		if got.Wiring != wantWiring {
			t.Errorf("NewFourthRotor() in rotor package returned unexpected wiring - got=%v, expected=%v", got.Wiring, wantWiring)
		}

		wantInverseWiring := conversion.InverseWiringArray(wantWiring)
		if got.InverseWiring != wantInverseWiring {
			t.Errorf("NewFourthRotor() in rotor package returned unexpected inverse wiring - got=%v, expected=%v", got.InverseWiring, wantInverseWiring)
		}

		wantPosition := conversion.ConvertCharToIndex(testData.startingPosition)
		if got.Position != wantPosition {
			t.Errorf("NewFourthRotor() in rotor package returned unexpected position - got=%v, expected=%v", got.Position, wantPosition)
		}
	}
}

func TestFourthRotorMapSignal(t *testing.T) {
	tests := []struct {
		configuration    FourthRotorConfiguration
		startingPosition rune
		input            rune
		inverse          bool
		expectedOutput   rune
	}{
		// With position 'A' the delta is zero, so the signal maps straight through the wiring.
		{FourthRotorBeta, 'A', 'A', false, 'L'},
		{FourthRotorGamma, 'A', 'A', false, 'F'},
		// Inverse mapping with a zero delta uses the inverse wiring. Beta maps 'A'->'L', so 'L' maps back to 'A'.
		{FourthRotorBeta, 'A', 'L', true, 'A'},
	}

	for _, testData := range tests {
		rotor := NewFourthRotor(testData.configuration, testData.startingPosition)
		got := rotor.MapSignal(conversion.ConvertCharToIndex(testData.input), testData.inverse)
		expectedOutput := conversion.ConvertCharToIndex(testData.expectedOutput)
		if got != expectedOutput {
			t.Errorf("MapSignal() in rotor package returned unexpected output - got=%v, expected=%v", got, expectedOutput)
		}
	}
}

// TestFourthRotorMapSignalRoundTrip verifies that mapping a signal forward and then back through a fourth rotor
// returns the original signal, across a range of positions. This exercises the position delta wrap-around math
// without needing hand-computed expected values.
func TestFourthRotorMapSignalRoundTrip(t *testing.T) {
	positions := []rune{'A', 'F', 'N', 'Z'}

	for _, position := range positions {
		rotor := NewFourthRotor(FourthRotorGamma, position)
		for signal := 0; signal < 26; signal++ {
			forward := rotor.MapSignal(signal, false)
			got := rotor.MapSignal(forward, true)
			if got != signal {
				t.Errorf("MapSignal() round-trip in rotor package failed for position=%c signal=%d - got=%v, expected=%v", position, signal, got, signal)
			}
		}
	}
}
