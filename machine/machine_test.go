package machine

import (
	"testing"

	"github.com/atbu/ultra/reflector"
	"github.com/atbu/ultra/rotor"
)

func TestSanityCheck(t *testing.T) {
	machine := &Machine{
		rotor.NewRotor(rotor.RotorI, 'A', 'A'),
		rotor.NewRotor(rotor.RotorII, 'A', 'A'),
		rotor.NewRotor(rotor.RotorIII, 'A', 'A'),
		nil,
		reflector.New(reflector.ReflectorB),
		nil,
	}

	input := "AAAAA"
	actualOutput := machine.Process(input)
	expectedOutput := "BDZGO"
	if actualOutput != expectedOutput {
		t.Errorf("machine.Process(%s) failed - got=%s, expected=%s", input, actualOutput, expectedOutput)
	}
}
