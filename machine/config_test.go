package machine

import (
	"testing"
)

func TestNewFromConfig(t *testing.T) {
	m, err := NewFromConfig(Config{
		Rotors:    [3]string{"I", "II", "III"},
		Positions: "AAA",
		Rings:     "AAA",
		Reflector: "B",
	})
	if err != nil {
		t.Fatalf("NewFromConfig returned an error: %s", err)
	}

	runTest(t, m, "AAAAA", "BDZGO")
}

func TestNewFromConfigIsCaseInsensitive(t *testing.T) {
	m, err := NewFromConfig(Config{
		Rotors:    [3]string{"i", "ii", "iii"},
		Positions: "mck",
		Rings:     "aaa",
		Reflector: "b",
		Plugboard: "",
	})
	if err != nil {
		t.Fatalf("NewFromConfig returned an error: %s", err)
	}

	runTest(t, m, "ENIGMA", "QMJIDO")
}

// With Beta at position A and the thin B reflector, an Enigma M4 behaves exactly like an Enigma M3 with reflector B.
func TestNewFromConfigFourthRotorCompatibility(t *testing.T) {
	m, err := NewFromConfig(Config{
		Rotors:         [3]string{"I", "II", "III"},
		Positions:      "AAA",
		Rings:          "AAA",
		FourthRotor:    "Beta",
		FourthPosition: "A",
		Reflector:      "B-thin",
	})
	if err != nil {
		t.Fatalf("NewFromConfig returned an error: %s", err)
	}

	runTest(t, m, "AAAAA", "BDZGO")
}

func TestNewFromConfigPlugboard(t *testing.T) {
	withPlugs, err := NewFromConfig(Config{
		Rotors:    [3]string{"I", "II", "III"},
		Positions: "AAA",
		Rings:     "AAA",
		Reflector: "B",
		Plugboard: "ab cd",
	})
	if err != nil {
		t.Fatalf("NewFromConfig returned an error: %s", err)
	}

	if withPlugs.Plugboard == nil || len(withPlugs.Plugboard.Configuration) != 2 {
		t.Fatalf("expected a plugboard with two pairs, got %+v", withPlugs.Plugboard)
	}
}

func TestNewFromConfigErrors(t *testing.T) {
	valid := func() Config {
		return Config{
			Rotors:    [3]string{"I", "II", "III"},
			Positions: "AAA",
			Rings:     "AAA",
			Reflector: "B",
		}
	}

	tests := []struct {
		name   string
		modify func(c *Config)
	}{
		{"unknown rotor", func(c *Config) { c.Rotors[0] = "IX" }},
		{"duplicate rotor", func(c *Config) { c.Rotors[2] = "I" }},
		{"short positions", func(c *Config) { c.Positions = "AA" }},
		{"non-letter positions", func(c *Config) { c.Positions = "A1A" }},
		{"short rings", func(c *Config) { c.Rings = "" }},
		{"unknown reflector", func(c *Config) { c.Reflector = "D" }},
		{"unknown fourth rotor", func(c *Config) {
			c.FourthRotor, c.FourthPosition, c.Reflector = "Delta", "A", "B-thin"
		}},
		{"bad fourth position", func(c *Config) {
			c.FourthRotor, c.FourthPosition, c.Reflector = "Beta", "", "B-thin"
		}},
		{"fourth rotor with wide reflector", func(c *Config) { c.FourthRotor, c.FourthPosition = "Beta", "A" }},
		{"thin reflector without fourth rotor", func(c *Config) { c.Reflector = "C-thin" }},
		{"uneven plugboard", func(c *Config) { c.Plugboard = "ABC" }},
		{"non-letter plugboard", func(c *Config) { c.Plugboard = "A1" }},
		{"repeated plugboard letter", func(c *Config) { c.Plugboard = "AB AC" }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := valid()
			test.modify(&c)
			if m, err := NewFromConfig(c); err == nil {
				t.Errorf("expected an error, got machine %+v", m)
			}
		})
	}
}
