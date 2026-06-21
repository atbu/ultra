package plugboard

import (
	"errors"
	"strings"

	"github.com/atbu/ultra/conversion"
)

var ErrUnevenNumberOfCharacters = errors.New("plugboard configuration invalid due to an uneven number of characters in plugboard string")

// The PlugboardPair structure represents a pair of characters that are linked in both ways. In real Enigma machines,
// a plugboard exists, usually with ten wires with plugs on either end, where one end is plugged into one letter, and
// another plug is plugged into another letter.
//
// For example, if 'A' is connected to 'G' in this way, if 'A' is mapped
// through the plugboard, we will get 'G' out, and likewise if 'G' is mapped through the plugboard, we will get 'A'.
//
// If a letter is not connected to any other, it passes through the plugboard unchanged.
type PlugboardPair struct {
	a, b int
}

// The Plugboard structure contains a slice of all the plugboard pairs currently connected.
type Plugboard struct {
	Configuration []PlugboardPair
}

// The New function constructs a new plugboard from a string which contains all the pairs, e.g. "AGRJFILO" would
// represent a plugboard configuration where A and G map to each other, R and J map to each other, and so on.
func New(plugboardString string) (*Plugboard, error) {
	// If the plugboardString is empty, we can return nil as we have no plugboard.
	if len(strings.TrimSpace(plugboardString)) == 0 {
		return nil, nil
	}

	// The plugboard string must contain pairs - if we have an uneven number of characters then a letter is plugged in
	// but not connected to any other letter which shouldn't be allowed, so we return an error.
	if len(plugboardString)%2 != 0 {
		return nil, ErrUnevenNumberOfCharacters
	}

	// Iterate over the pairs of characters in the plugboard string and add them to the configuration slice.
	var configuration []PlugboardPair
	for i := range len(plugboardString) - 1 {
		// Ignore every other character as we're looking for pairs.
		if i%2 != 0 {
			continue
		}

		x := conversion.ConvertCharToIndex(rune(plugboardString[i]))
		y := conversion.ConvertCharToIndex(rune(plugboardString[i+1]))

		configuration = append(configuration, PlugboardPair{x, y})
	}

	// Return a constructed Plugboard configured as above.
	return &Plugboard{configuration}, nil
}

// The MapSignal function defines how a signal passes through a plugboard. We receive a signal as a character index, and
// return a character index as output.
func (p *Plugboard) MapSignal(signal int) int {
	// We check every pair in the plugboard configuration to see if the character we're looking for is contained in it.
	// If it is, we return the other character in the pair.
	for _, pair := range p.Configuration {
		if pair.a == signal {
			return pair.b
		}

		if pair.b == signal {
			return pair.a
		}
	}

	// If we don't find a pair with our character in it, just return the signal unchanged because this represents
	// a character with no corresponding plug in the plugboard.
	return signal
}
