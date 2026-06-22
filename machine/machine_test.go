package machine

import (
	"testing"

	"github.com/atbu/ultra/conversion"
	"github.com/atbu/ultra/plugboard"
	"github.com/atbu/ultra/reflector"
	"github.com/atbu/ultra/rotor"
)

func (m *Machine) updateRotorPositions(left rune, middle rune, right rune) {
	m.LeftRotor.Position = conversion.ConvertCharToIndex(left)
	m.MiddleRotor.Position = conversion.ConvertCharToIndex(middle)
	m.RightRotor.Position = conversion.ConvertCharToIndex(right)
}

func createPlugboard(t *testing.T, plugboardString string) *plugboard.Plugboard {
	pb, err := plugboard.New(plugboardString)
	if err != nil {
		t.Errorf("failed to create new plugboard with string %s: %s", plugboardString, err)
	}

	return pb
}

func runTest(t *testing.T, machine *Machine, input string, expectedOutput string) {
	actualOutput := machine.Process(input)
	if actualOutput != expectedOutput {
		t.Errorf("machine.Process(\"%s\") failed - got=%s, expected=%s", input, actualOutput, expectedOutput)
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

	runTest(t, machine, "AAAAA", "BDZGO")
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

	runTest(t, machine, "ENIGMA", "QMJIDO")

	machine.updateRotorPositions('M', 'C', 'K')

	runTest(t, machine, "QMJIDO", "ENIGMA")
}

func TestNormalTurnover(t *testing.T) {
	machine := &Machine{
		rotor.NewRotor(rotor.RotorI, 'K', 'A'),
		rotor.NewRotor(rotor.RotorII, 'D', 'A'),
		rotor.NewRotor(rotor.RotorIII, 'O', 'A'),
		nil,
		reflector.New(reflector.ReflectorB),
		nil,
	}

	runTest(t, machine, "AAAAA", "JWZBJ")

	machine.updateRotorPositions('K', 'D', 'U')

	runTest(t, machine, "AAAAA", "YWDVQ")
}

func TestDoubleStep(t *testing.T) {
	machine := &Machine{
		rotor.NewRotor(rotor.RotorI, 'A', 'A'),
		rotor.NewRotor(rotor.RotorII, 'D', 'A'),
		rotor.NewRotor(rotor.RotorIII, 'U', 'A'),
		nil,
		reflector.New(reflector.ReflectorB),
		nil,
	}

	runTest(t, machine, "AAAAA", "EQIBM")
}

func TestPlugboard(t *testing.T) {
	machine := &Machine{
		rotor.NewRotor(rotor.RotorI, 'Z', 'A'),
		rotor.NewRotor(rotor.RotorII, 'Z', 'A'),
		rotor.NewRotor(rotor.RotorIII, 'Z', 'A'),
		nil,
		reflector.New(reflector.ReflectorB),
		createPlugboard(t, "ABCDEFGH"),
	}

	runTest(t, machine, "AAAAA", "UZYRQ")
}

func TestFullIntegration(t *testing.T) {
	machine := &Machine{
		rotor.NewRotor(rotor.RotorII, 'A', 'A'),
		rotor.NewRotor(rotor.RotorIV, 'B', 'A'),
		rotor.NewRotor(rotor.RotorV, 'L', 'A'),
		nil,
		reflector.New(reflector.ReflectorB),
		createPlugboard(t, "BQCRDIEJKWMTOSPXUZGH"),
	}

	runTest(t, machine, "EVERYTHINGISGOINGEXTREMELYWELL", "LLSDWFYUVEVDHBJVTWWECZNWYXLCNX")
}

func TestRingCase1(t *testing.T) {
	machine := &Machine{
		rotor.NewRotor(rotor.RotorI, 'A', 'B'),
		rotor.NewRotor(rotor.RotorII, 'A', 'B'),
		rotor.NewRotor(rotor.RotorIII, 'A', 'B'),
		nil,
		reflector.New(reflector.ReflectorB),
		nil,
	}

	runTest(t, machine, "AAAAA", "EWTYX")
}

func TestRingCase2(t *testing.T) {
	machine := &Machine{
		rotor.NewRotor(rotor.RotorI, 'K', 'A'),
		rotor.NewRotor(rotor.RotorII, 'D', 'A'),
		rotor.NewRotor(rotor.RotorIII, 'T', 'B'),
		nil,
		reflector.New(reflector.ReflectorB),
		nil,
	}

	runTest(t, machine, "AAAA", "JTIN")
}

func TestRingCase3(t *testing.T) {
	machine := &Machine{
		rotor.NewRotor(rotor.RotorI, 'G', 'R'),
		rotor.NewRotor(rotor.RotorII, 'U', 'T'),
		rotor.NewRotor(rotor.RotorIII, 'M', 'M'),
		nil,
		reflector.New(reflector.ReflectorB),
		createPlugboard(t, "AKSORILP"),
	}

	runTest(t, machine, "HELLOWORLD", "CDKSEVMKXJ")
}

func TestProcessEmptyString(t *testing.T) {
	// note: this machine configuration wouldn't be possible in real life because an operator would only
	// have one of each rotor (as far as I know)
	machine := &Machine{
		rotor.NewRotor(rotor.RotorII, 'F', 'D'),
		rotor.NewRotor(rotor.RotorII, 'P', 'W'),
		rotor.NewRotor(rotor.RotorII, 'K', 'L'),
		nil,
		reflector.New(reflector.ReflectorB),
		createPlugboard(t, "JWYLFKREVPXTHOBCMQZG"),
	}

	runTest(t, machine, "", "")
}

func TestRandomCase1(t *testing.T) {
	// note: this machine configuration wouldn't be possible in real life because an operator would only
	// have one of each rotor (as far as I know)
	machine := &Machine{
		rotor.NewRotor(rotor.RotorII, 'F', 'D'),
		rotor.NewRotor(rotor.RotorII, 'P', 'W'),
		rotor.NewRotor(rotor.RotorII, 'K', 'L'),
		nil,
		reflector.New(reflector.ReflectorB),
		createPlugboard(t, "JWYLFKREVPXTHOBCMQZG"),
	}

	runTest(t, machine, "UISDNUUINSNASIAASNUUIDIADIADDDNNNS", "LNLCJPIFILXIKZPROFOZATVGWZUWOFBFVB")
}

func TestRandomCase2(t *testing.T) {
	// note: this machine configuration wouldn't be possible in real life because an operator would only
	// have one of each rotor (as far as I know)
	machine := &Machine{
		rotor.NewRotor(rotor.RotorIII, 'K', 'B'),
		rotor.NewRotor(rotor.RotorIII, 'E', 'T'),
		rotor.NewRotor(rotor.RotorV, 'C', 'H'),
		nil,
		reflector.New(reflector.ReflectorB),
		createPlugboard(t, "ACPQUHYFWRMJOSKTDIVG"),
	}

	runTest(t, machine, "EEBXZZEBZXNXLBLBBZNLZBLNLNBBNBLLNBXNEZLB", "BPNAFCWSDBGAFDIQPKGHXNFMXIGIKLXPKTPORWOX")
}

func assertRotorPosition(t *testing.T, rotor *rotor.Rotor, position int) {
	if rotor.Position != position {
		t.Errorf("rotor position assertion failed - got=%d, expected=%d", rotor.Position, position)
	}
}

func TestMultipleRotorTurnovers(t *testing.T) {
	machine := &Machine{
		rotor.NewRotor(rotor.RotorII, 'K', 'B'),
		rotor.NewRotor(rotor.RotorIII, 'E', 'T'),
		rotor.NewRotor(rotor.RotorVI, 'Z', 'H'),
		nil,
		reflector.New(reflector.ReflectorB),
		nil,
	}

	assertRotorPosition(t, machine.RightRotor, 25) // Z, not yet turned over
	assertRotorPosition(t, machine.MiddleRotor, 4) // E, not yet turned over

	// Right rotor is currently in Z position so should turn over on next rotation.
	machine.PressKey('A')

	assertRotorPosition(t, machine.RightRotor, 0)  // A, has turned over
	assertRotorPosition(t, machine.MiddleRotor, 5) // F, has turned over because right rotor was at notch

	machine.updateRotorPositions('K', 'F', 'M')

	assertRotorPosition(t, machine.RightRotor, 12) // M, not yet turned over
	assertRotorPosition(t, machine.MiddleRotor, 5) // F, not yet turned over

	machine.PressKey('A')

	assertRotorPosition(t, machine.RightRotor, 13) // N, has turned over
	assertRotorPosition(t, machine.MiddleRotor, 6) // F, has turned over because right rotor was at notch
}

// Tests that the original UKW-B, and narrow UKW-B reflector and Beta fourth rotor are equivalent.
func TestMachineEquivalency(t *testing.T) {
	enigmaI := &Machine{
		rotor.NewRotor(rotor.RotorI, 'A', 'A'),
		rotor.NewRotor(rotor.RotorII, 'A', 'A'),
		rotor.NewRotor(rotor.RotorIII, 'A', 'A'),
		nil,
		reflector.New(reflector.ReflectorB),
		nil,
	}

	enigmaM4 := &Machine{
		rotor.NewRotor(rotor.RotorI, 'A', 'A'),
		rotor.NewRotor(rotor.RotorII, 'A', 'A'),
		rotor.NewRotor(rotor.RotorIII, 'A', 'A'),
		rotor.NewFourthRotor(rotor.FourthRotorBeta, 'A'),
		reflector.New(reflector.ReflectorNarrowB),
		nil,
	}

	enigmaIOutput := enigmaI.Process("ABCDEFG")
	enigmaM4Output := enigmaM4.Process("ABCDEFG")
	if enigmaIOutput != enigmaM4Output {
		t.Errorf(
			"Enigma I and Enigma M4 fourth rotor/reflector equivalency test failed - I=%s, M4=%s",
			enigmaIOutput,
			enigmaM4Output,
		)
	}
}
