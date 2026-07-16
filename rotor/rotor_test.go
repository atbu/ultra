package rotor

import (
	"reflect"
	"testing"

	"github.com/atbu/ultra/conversion"
)

func TestNewRotor(t *testing.T) {
	tests := []struct {
		configuration    RotorConfiguration
		wiringString     string
		notches          []rune
		startingPosition rune
		ringSetting      rune
	}{
		{RotorI, "EKMFLGDQVZNTOWYHXUSPAIBRCJ", []rune{'Q'}, 'A', 'A'},
		{RotorII, "AJDKSIRUXBLHWTMCQGZNPYFVOE", []rune{'E'}, 'A', 'A'},
		{RotorIII, "BDFHJLCPRTXVZNYEIWGAKMUSQO", []rune{'V'}, 'M', 'F'},
		{RotorIV, "ESOVPZJAYQUIRHXLNFTGKDCMWB", []rune{'J'}, 'Z', 'B'},
		{RotorV, "VZBRGITYUPSDNHLXAWMJQOFECK", []rune{'Z'}, 'A', 'A'},
		{RotorVI, "JPGVOUMFYQBENHZRDKASXLICTW", []rune{'Z', 'M'}, 'K', 'C'},
		{RotorVII, "NZJHGRCXMYSWBOUFAIVLPEKQDT", []rune{'Z', 'M'}, 'A', 'A'},
		{RotorVIII, "FKQHTLXOCBJSPDZRAMEWNIUYGV", []rune{'Z', 'M'}, 'Y', 'X'},
	}

	for _, testData := range tests {
		got := NewRotor(testData.configuration, testData.startingPosition, testData.ringSetting)

		wantWiring := conversion.ConvertWiringStringToArray(testData.wiringString)
		if got.Wiring != wantWiring {
			t.Errorf("NewRotor() in rotor package returned unexpected wiring - got=%v, expected=%v", got.Wiring, wantWiring)
		}

		wantInverseWiring := conversion.InverseWiringArray(wantWiring)
		if got.InverseWiring != wantInverseWiring {
			t.Errorf("NewRotor() in rotor package returned unexpected inverse wiring - got=%v, expected=%v", got.InverseWiring, wantInverseWiring)
		}

		var wantNotches []int
		for _, notch := range testData.notches {
			wantNotches = append(wantNotches, conversion.ConvertCharToIndex(notch))
		}
		if !reflect.DeepEqual(got.Notches, wantNotches) {
			t.Errorf("NewRotor() in rotor package returned unexpected notches - got=%v, expected=%v", got.Notches, wantNotches)
		}

		wantPosition := conversion.ConvertCharToIndex(testData.startingPosition)
		if got.Position != wantPosition {
			t.Errorf("NewRotor() in rotor package returned unexpected position - got=%v, expected=%v", got.Position, wantPosition)
		}

		wantRingSetting := conversion.ConvertCharToIndex(testData.ringSetting)
		if got.RingSetting != wantRingSetting {
			t.Errorf("NewRotor() in rotor package returned unexpected ring setting - got=%v, expected=%v", got.RingSetting, wantRingSetting)
		}
	}
}

func TestRotate(t *testing.T) {
	tests := []struct {
		startingPosition rune
		expectedPosition int
	}{
		// A normal step simply increments the position by one.
		{'A', conversion.ConvertCharToIndex('B')},
		{'M', conversion.ConvertCharToIndex('N')},
		// Rotating from 'Z' (index 25) must wrap around back to 'A' (index 0).
		{'Z', conversion.ConvertCharToIndex('A')},
	}

	for _, testData := range tests {
		rotor := NewRotor(RotorI, testData.startingPosition, 'A')
		rotor.Rotate()
		if rotor.Position != testData.expectedPosition {
			t.Errorf("Rotate() in rotor package returned unexpected position - got=%v, expected=%v", rotor.Position, testData.expectedPosition)
		}
	}
}

func TestMapSignal(t *testing.T) {
	tests := []struct {
		configuration    RotorConfiguration
		startingPosition rune
		ringSetting      rune
		input            rune
		inverse          bool
		expectedOutput   rune
	}{
		// With position 'A' and ring setting 'A' the delta is zero, so the signal maps straight through the wiring.
		{RotorI, 'A', 'A', 'A', false, 'E'},
		{RotorII, 'A', 'A', 'A', false, 'A'},
		// Inverse mapping with a zero delta uses the inverse wiring. RotorI maps 'A'->'E', so 'E' maps back to 'A'.
		{RotorI, 'A', 'A', 'E', true, 'A'},
		// A non-zero position offsets the signal (Enigma stepping math).
		{RotorI, 'B', 'A', 'A', false, 'J'},
		// A non-zero ring setting offsets the signal in the opposite direction.
		{RotorI, 'A', 'B', 'A', false, 'K'},
	}

	for _, testData := range tests {
		rotor := NewRotor(testData.configuration, testData.startingPosition, testData.ringSetting)
		got := rotor.MapSignal(conversion.ConvertCharToIndex(testData.input), testData.inverse)
		expectedOutput := conversion.ConvertCharToIndex(testData.expectedOutput)
		if got != expectedOutput {
			t.Errorf("MapSignal() in rotor package returned unexpected output - got=%v, expected=%v", got, expectedOutput)
		}
	}
}

// TestMapSignalRoundTrip verifies that mapping a signal forward and then back through a rotor returns the original
// signal for every starting signal, across a range of positions and ring settings. This exercises the position /
// ring-setting delta wrap-around math without needing hand-computed expected values.
func TestMapSignalRoundTrip(t *testing.T) {
	positions := []rune{'A', 'F', 'N', 'Z'}
	ringSettings := []rune{'A', 'C', 'T', 'Z'}

	for _, position := range positions {
		for _, ringSetting := range ringSettings {
			rotor := NewRotor(RotorIII, position, ringSetting)
			for signal := 0; signal < 26; signal++ {
				forward := rotor.MapSignal(signal, false)
				got := rotor.MapSignal(forward, true)
				if got != signal {
					t.Errorf("MapSignal() round-trip in rotor package failed for position=%c ringSetting=%c signal=%d - got=%v, expected=%v", position, ringSetting, signal, got, signal)
				}
			}
		}
	}
}

func TestIsInNotch(t *testing.T) {
	tests := []struct {
		configuration    RotorConfiguration
		startingPosition rune
		expectedOutput   bool
	}{
		// RotorI has a single notch at 'Q'.
		{RotorI, 'Q', true},
		{RotorI, 'A', false},
		// RotorVI has two notches, at 'Z' and 'M'. Both must be recognised.
		{RotorVI, 'Z', true},
		{RotorVI, 'M', true},
		{RotorVI, 'A', false},
	}

	for _, testData := range tests {
		rotor := NewRotor(testData.configuration, testData.startingPosition, 'A')
		got := rotor.IsInNotch()
		if got != testData.expectedOutput {
			t.Errorf("IsInNotch() in rotor package returned unexpected output - got=%v, expected=%v", got, testData.expectedOutput)
		}
	}
}
