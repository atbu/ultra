package machine

import (
	"testing"

	"github.com/atbu/ultra/conversion"
	"github.com/atbu/ultra/reflector"
	"github.com/atbu/ultra/rotor"
)

func (m *Machine) updateRotorPositions(left rune, middle rune, right rune) {
	m.LeftRotor.Position = conversion.ConvertCharToIndex(left)
	m.MiddleRotor.Position = conversion.ConvertCharToIndex(middle)
	m.RightRotor.Position = conversion.ConvertCharToIndex(right)
}

func runTest(t *testing.T, machine *Machine, input string, expectedOutput string) {
	actualOutput := machine.Process(input)
	if actualOutput != expectedOutput {
		t.Errorf("machine.Process(%s) failed - got=%s, expected=%s", input, actualOutput, expectedOutput)
	}
}

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
	expectedOutput := "BDZGO"
	runTest(t, machine, input, expectedOutput)
}

func TestReciprocity(t *testing.T) {
	machine := &Machine{
		rotor.NewRotor(rotor.RotorI, 'M', 'A'),
		rotor.NewRotor(rotor.RotorII, 'C', 'A'),
		rotor.NewRotor(rotor.RotorIII, 'K', 'A'),
		nil,
		reflector.New(reflector.ReflectorB),
		nil,
	}

	input := "ENIGMA"
	expectedOutput := "QMJIDO"
	runTest(t, machine, input, expectedOutput)

	machine.updateRotorPositions('M', 'C', 'K')

	input = "AAAAA"
	expectedOutput = "YWDVQ"
	runTest(t, machine, input, expectedOutput)
}
