package plugboard

import (
	"errors"
	"strings"

	"github.com/atbu/ultra/conversion"
)

var ErrUnevenNumberOfCharacters = errors.New("plugboard configuration invalid due to an uneven number of characters in plugboard string")

type PlugboardPair struct {
	a, b int
}

type Plugboard struct {
	Configuration []PlugboardPair
}

func New(plugboardString string) (*Plugboard, error) {
	// If the plugboardString is empty, we can return nil as we have no plugboard.
	if len(strings.TrimSpace(plugboardString)) == 0 {
		return nil, nil
	}

	if len(plugboardString)%2 != 0 {
		return nil, ErrUnevenNumberOfCharacters
	}

	var configuration []PlugboardPair
	for i := range len(plugboardString) - 1 {
		// ignore every other value as we're looking for pairs
		if i%2 != 0 {
			continue
		}

		x := conversion.ConvertCharToIndex(rune(plugboardString[i]))
		y := conversion.ConvertCharToIndex(rune(plugboardString[i+1]))

		configuration = append(configuration, PlugboardPair{x, y})
	}

	return &Plugboard{configuration}, nil
}

func (p *Plugboard) MapSignal(signal int) int {
	for _, pair := range p.Configuration {
		if pair.a == signal {
			return pair.b
		}

		if pair.b == signal {
			return pair.a
		}
	}

	return signal
}
