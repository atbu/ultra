package machine

import (
	"errors"
	"fmt"
	"strings"

	"github.com/atbu/ultra/plugboard"
	"github.com/atbu/ultra/reflector"
	"github.com/atbu/ultra/rotor"
)

// The Config structure describes a machine's settings using the names an operator would use, e.g. rotor "IV" or
// reflector "B". It is designed to be easy to fill in from user input (such as the web UI) and turned into a Machine
// with the NewFromConfig function.
type Config struct {
	// The rotors fitted to the left, middle and right slots, e.g. ["I", "II", "III"].
	Rotors [3]string `json:"rotors"`
	// The starting positions of the left, middle and right rotors, e.g. "AAA".
	Positions string `json:"positions"`
	// The ring settings of the left, middle and right rotors, e.g. "AAA".
	Rings string `json:"rings"`
	// The optional Enigma M4 fourth rotor, "Beta" or "Gamma". An empty string means no fourth rotor is fitted.
	FourthRotor string `json:"fourthRotor"`
	// The starting position of the fourth rotor, e.g. "A". Ignored when there is no fourth rotor.
	FourthPosition string `json:"fourthPosition"`
	// The reflector, "A", "B" or "C", or "B-thin" / "C-thin" for an Enigma M4 with a fourth rotor fitted.
	Reflector string `json:"reflector"`
	// The plugboard pairs, e.g. "AB CD EF". Spaces are ignored and an empty string means no plugs are connected.
	Plugboard string `json:"plugboard"`
}

var rotorNames = map[string]rotor.RotorConfiguration{
	"I":    rotor.RotorI,
	"II":   rotor.RotorII,
	"III":  rotor.RotorIII,
	"IV":   rotor.RotorIV,
	"V":    rotor.RotorV,
	"VI":   rotor.RotorVI,
	"VII":  rotor.RotorVII,
	"VIII": rotor.RotorVIII,
}

var fourthRotorNames = map[string]rotor.FourthRotorConfiguration{
	"BETA":  rotor.FourthRotorBeta,
	"GAMMA": rotor.FourthRotorGamma,
}

var reflectorNames = map[string]reflector.ReflectorConfiguration{
	"A":      reflector.ReflectorA,
	"B":      reflector.ReflectorB,
	"C":      reflector.ReflectorC,
	"B-THIN": reflector.ReflectorNarrowB,
	"C-THIN": reflector.ReflectorNarrowC,
}

var ErrDuplicateRotor = errors.New("the same rotor cannot be fitted to more than one slot")

// The NewFromConfig function validates a Config and constructs the Machine it describes.
func NewFromConfig(c Config) (*Machine, error) {
	var rotors [3]*rotor.Rotor
	positions := strings.ToUpper(c.Positions)
	rings := strings.ToUpper(c.Rings)

	if !isLetters(positions, 3) {
		return nil, fmt.Errorf("rotor positions %q must be exactly three letters", c.Positions)
	}
	if !isLetters(rings, 3) {
		return nil, fmt.Errorf("ring settings %q must be exactly three letters", c.Rings)
	}

	seen := map[string]bool{}
	for i, name := range c.Rotors {
		name = strings.ToUpper(name)
		configuration, ok := rotorNames[name]
		if !ok {
			return nil, fmt.Errorf("unknown rotor %q", c.Rotors[i])
		}
		if seen[name] {
			return nil, ErrDuplicateRotor
		}
		seen[name] = true
		rotors[i] = rotor.NewRotor(configuration, rune(positions[i]), rune(rings[i]))
	}

	reflectorName := strings.ToUpper(c.Reflector)
	reflectorConfiguration, ok := reflectorNames[reflectorName]
	if !ok {
		return nil, fmt.Errorf("unknown reflector %q", c.Reflector)
	}
	thinReflector := strings.HasSuffix(reflectorName, "-THIN")

	// On a real Enigma M4, the fourth rotor only fits alongside a thin reflector, and a thin reflector needs the
	// fourth rotor to fill the space, so we require one whenever the other is present.
	var fourthRotor *rotor.FourthRotor
	if c.FourthRotor != "" {
		configuration, ok := fourthRotorNames[strings.ToUpper(c.FourthRotor)]
		if !ok {
			return nil, fmt.Errorf("unknown fourth rotor %q", c.FourthRotor)
		}
		position := strings.ToUpper(c.FourthPosition)
		if !isLetters(position, 1) {
			return nil, fmt.Errorf("fourth rotor position %q must be a single letter", c.FourthPosition)
		}
		if !thinReflector {
			return nil, errors.New("a fourth rotor must be used with a thin reflector (B-thin or C-thin)")
		}
		fourthRotor = rotor.NewFourthRotor(configuration, rune(position[0]))
	} else if thinReflector {
		return nil, errors.New("a thin reflector must be used with a fourth rotor (Beta or Gamma)")
	}

	pb, err := parsePlugboard(c.Plugboard)
	if err != nil {
		return nil, err
	}

	return &Machine{
		LeftRotor:   rotors[0],
		MiddleRotor: rotors[1],
		RightRotor:  rotors[2],
		FourthRotor: fourthRotor,
		Reflector:   reflector.New(reflectorConfiguration),
		Plugboard:   pb,
	}, nil
}

// The parsePlugboard function checks a plugboard string, e.g. "AB CD", only contains letters which are each used once,
// before constructing a Plugboard from it.
func parsePlugboard(s string) (*plugboard.Plugboard, error) {
	pairs := strings.ToUpper(strings.Join(strings.Fields(s), ""))

	if !isLetters(pairs, len(pairs)) {
		return nil, fmt.Errorf("plugboard %q must only contain letters", s)
	}

	seen := map[rune]bool{}
	for _, letter := range pairs {
		if seen[letter] {
			return nil, fmt.Errorf("plugboard letter %c is used more than once", letter)
		}
		seen[letter] = true
	}

	return plugboard.New(pairs)
}

// The isLetters function reports whether s is exactly length characters long and only contains the letters A to Z.
func isLetters(s string, length int) bool {
	if len(s) != length {
		return false
	}

	for _, char := range s {
		if char < 'A' || char > 'Z' {
			return false
		}
	}

	return true
}
