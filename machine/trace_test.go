package machine

import (
	"testing"
)

// PressKeyWithTrace must give exactly the same results as PressKey, across rotor turnovers and with every optional
// component present or absent.
func TestPressKeyWithTraceMatchesPressKey(t *testing.T) {
	configs := []Config{
		{Rotors: [3]string{"I", "II", "III"}, Positions: "AAA", Rings: "AAA", Reflector: "B"},
		{Rotors: [3]string{"IV", "V", "II"}, Positions: "QEV", Rings: "BUL", Reflector: "C", Plugboard: "AV BS CG DL FU HZ IN KM OW RX"},
		{Rotors: [3]string{"VIII", "VI", "VII"}, Positions: "ZMZ", Rings: "MZA", Reflector: "A"},
		{Rotors: [3]string{"II", "IV", "I"}, Positions: "CDF", Rings: "EVA", FourthRotor: "Gamma", FourthPosition: "Q", Reflector: "C-thin", Plugboard: "QW ER TZ"},
	}

	// A long message so every rotor turns over at least once, including double stepping.
	message := ""
	for range 30 {
		message += "THEQUICKBROWNFOXJUMPSOVERTHELAZYDOG"
	}

	for _, config := range configs {
		plain, err := NewFromConfig(config)
		if err != nil {
			t.Fatalf("NewFromConfig(%+v) returned an error: %s", config, err)
		}
		traced, _ := NewFromConfig(config)

		for i, key := range message {
			expected := plain.PressKey(key)
			actual, trace := traced.PressKeyWithTrace(key)

			if actual != expected {
				t.Fatalf("config %+v, key %d (%c): PressKeyWithTrace=%c, PressKey=%c", config, i, key, actual, expected)
			}

			checkTrace(t, trace, key, actual, config.FourthRotor != "")
		}

		if plain.LeftRotor.Position != traced.LeftRotor.Position ||
			plain.MiddleRotor.Position != traced.MiddleRotor.Position ||
			plain.RightRotor.Position != traced.RightRotor.Position {
			t.Errorf("config %+v: rotor positions differ after processing", config)
		}
	}
}

func checkTrace(t *testing.T, trace []TraceStep, key rune, output rune, hasFourthRotor bool) {
	t.Helper()

	expected := []Component{
		ComponentPlugboard, ComponentRightRotor, ComponentMiddleRotor, ComponentLeftRotor,
		ComponentReflector,
		ComponentLeftRotor, ComponentMiddleRotor, ComponentRightRotor, ComponentPlugboard,
	}
	if hasFourthRotor {
		expected = []Component{
			ComponentPlugboard, ComponentRightRotor, ComponentMiddleRotor, ComponentLeftRotor, ComponentFourthRotor,
			ComponentReflector,
			ComponentFourthRotor, ComponentLeftRotor, ComponentMiddleRotor, ComponentRightRotor, ComponentPlugboard,
		}
	}

	if len(trace) != len(expected) {
		t.Fatalf("trace has %d steps, expected %d: %+v", len(trace), len(expected), trace)
	}

	for i, step := range trace {
		if step.Component != expected[i] {
			t.Errorf("trace step %d is %s, expected %s", i, step.Component, expected[i])
		}
		if i > 0 && step.In != trace[i-1].Out {
			t.Errorf("trace step %d starts at %c but the previous step ended at %c", i, step.In, trace[i-1].Out)
		}
	}

	if trace[0].In != key {
		t.Errorf("trace starts at %c, expected the key %c", trace[0].In, key)
	}
	if trace[len(trace)-1].Out != output {
		t.Errorf("trace ends at %c, expected the output %c", trace[len(trace)-1].Out, output)
	}
}
