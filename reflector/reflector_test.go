package reflector

import (
	"reflect"
	"testing"

	"github.com/atbu/ultra/conversion"
)

// todo to structure better

func TestConstructor(t *testing.T) {
	tests := []struct {
		input          ReflectorConfiguration
		expectedOutput *Reflector
	}{
		{ReflectorA, &Reflector{[26]int{
			4,
			9,
			12,
			25,
			0,
			11,
			24,
			23,
			21,
			1,
			22,
			5,
			2,
			17,
			16,
			20,
			14,
			13,
			19,
			18,
			15,
			8,
			10,
			7,
			6,
			3,
		}}},
		{ReflectorB, &Reflector{[26]int{
			24,
			17,
			20,
			7,
			16,
			18,
			11,
			3,
			15,
			23,
			13,
			6,
			14,
			10,
			12,
			8,
			4,
			1,
			5,
			25,
			2,
			22,
			21,
			9,
			0,
			19,
		}}},
		{ReflectorC, &Reflector{[26]int{
			5,
			21,
			15,
			9,
			8,
			0,
			14,
			24,
			4,
			3,
			17,
			25,
			23,
			22,
			6,
			2,
			19,
			10,
			20,
			16,
			18,
			1,
			13,
			12,
			7,
			11,
		}}},
		{ReflectorNarrowB, &Reflector{[26]int{
			4,
			13,
			10,
			16,
			0,
			20,
			24,
			22,
			9,
			8,
			2,
			14,
			15,
			1,
			11,
			12,
			3,
			23,
			25,
			21,
			5,
			19,
			7,
			17,
			6,
			18,
		}}},
		{ReflectorNarrowC, &Reflector{[26]int{
			17,
			3,
			14,
			1,
			9,
			13,
			19,
			10,
			21,
			4,
			7,
			12,
			11,
			5,
			2,
			22,
			25,
			0,
			23,
			6,
			24,
			8,
			15,
			18,
			20,
			16,
		}}},
	}

	for _, testData := range tests {
		got := New(testData.input)
		if !reflect.DeepEqual(got, testData.expectedOutput) {
			t.Errorf("New() in reflector package returned unexpected output - got=%v, expected=%v", got, testData.expectedOutput)
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
