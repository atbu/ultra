package reflector

import (
	"testing"

	"github.com/atbu/ultra/conversion"
)

func TestConstructor(t *testing.T) {
	tests := []struct {
		input        ReflectorConfiguration
		wiringString string
	}{
		{ReflectorA, "EJMZALYXVBWFCRQUONTSPIKHGD"},
		{ReflectorB, "YRUHQSLDPXNGOKMIEBFZCWVJAT"},
		{ReflectorC, "FVPJIAOYEDRZXWGCTKUQSBNMHL"},
		{ReflectorNarrowB, "ENKQAUYWJICOPBLMDXZVFTHRGS"},
		{ReflectorNarrowC, "RDOBJNTKVEHMLFCWZAXGYIPSUQ"},
	}

	for _, testData := range tests {
		got := New(testData.input)
		wantWiring := conversion.ConvertWiringStringToArray(testData.wiringString)
		if got.Wiring != wantWiring {
			t.Errorf("New() in reflector package returned unexpected wiring - got=%v, expected=%v", got.Wiring, wantWiring)
		}
	}
}

func TestMapSignal(t *testing.T) {
	tests := []struct {
		configuration  ReflectorConfiguration
		input          rune
		expectedOutput rune
	}{
		{ReflectorA, 'A', 'E'},
		{ReflectorNarrowB, 'Z', 'S'},
		{ReflectorC, 'C', 'P'},
		{ReflectorB, 'B', 'R'},
		{ReflectorNarrowC, 'M', 'L'},
		{ReflectorA, 'Y', 'G'},
		{ReflectorC, 'H', 'Y'},
	}

	for _, testData := range tests {
		reflector := New(testData.configuration)
		got := reflector.MapSignal(conversion.ConvertCharToIndex(testData.input))
		expectedOutput := conversion.ConvertCharToIndex(testData.expectedOutput)
		if got != expectedOutput {
			t.Errorf("MapSignal() in reflector package returned unexpected output - got=%v, expected=%v", got, expectedOutput)
		}
	}
}
