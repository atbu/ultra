const NUMBER_OF_ROTORS: usize = 1;

const ROTOR_I_WIRING: &str = "EKMFLGDQVZNTOWYHXUSPAIBRCJ";
const ROTOR_I_TURNOVER: char = 'Q';

// const ROTOR_II_WIRING: &str = "AJDKSIRUXBLHWTMCQGZNPYFVOE";
// const ROTOR_II_TURNOVER: char = 'E';

// const ROTOR_III_WIRING: &str = "BDFHJLCPRTXVZNYEIWGAKMUSQO";
// const ROTOR_III_TURNOVER: char = 'V';

pub struct Rotor {
    wiring: [u8; 26],
    position: u8,
    notch: u8
}

impl Rotor {
    fn new(wiring_string: &str, notch: char, position: char) -> Self {
        let mut wiring: [u8; 26] = [0; 26];

        for (index, character) in wiring_string.chars().enumerate() {
            wiring[index] = char_to_index(character);
        }

        Self {
            wiring,
            position: char_to_index(position),
            notch: char_to_index(notch)
        }
    }

    fn rotate(&mut self) -> bool {
        let mut rotate_neighbour: bool = false;

        // We are at the notch, so tell the neighbouring rotor to rotate too
        if self.position == self.notch {
            rotate_neighbour = true;
        }

        self.position = (self.position + 1) % 26;

        rotate_neighbour
    }
}

type RotorSet = [Rotor; NUMBER_OF_ROTORS];

pub struct EnigmaMachine {
    rotor_set: RotorSet
}

impl EnigmaMachine {
    fn rotate_rotors(&mut self) {
        let mut should_rotate = true;
        for rotor in &mut self.rotor_set {
            if should_rotate {
                rotor.rotate();
            }

            if rotor.position == rotor.notch {
                should_rotate = true;
            } else {
                should_rotate = false;
            }
        }
    }

    fn map_through_rotor(signal: u8, rotor: &Rotor) -> u8 {
        let after_entrance_offset = (signal + rotor.position) % 26;
        let after_lookup = rotor.wiring[after_entrance_offset as usize];
        let after_exit_offset = (after_lookup - rotor.position) % 26;

        after_exit_offset
    }

    pub fn press_key(&mut self, signal: char) -> u8 {
        self.rotate_rotors();
        let signal: u8 = char_to_index(signal);

        let mut signal: u8 = signal;
        for rotor in &self.rotor_set {
            signal = Self::map_through_rotor(signal, rotor);
        }

        signal
    }
}

/// Converts a character to an integer index, where 'A' equals 0, 'B' equals 1, ..., 'Z' equals 25.
pub fn char_to_index(char: char) -> u8 {
    char as u8 - 65
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_char_to_index() {
        assert_eq!(char_to_index('A'), 0);
        assert_eq!(char_to_index('B'), 1);
        assert_eq!(char_to_index('C'), 2);
        assert_eq!(char_to_index('Z'), 25);
    }

    #[test]
    fn simple_transformation() {
        const START_POSITION: char = 'B';

        let mut machine: EnigmaMachine = EnigmaMachine {
            rotor_set: [
                Rotor::new(ROTOR_I_WIRING, ROTOR_I_TURNOVER, START_POSITION)
            ]
        };

        let signal = machine.press_key('A');
        assert_eq!(signal, char_to_index('K'));
        assert_eq!(machine.rotor_set[0].position, 2);

        let signal = machine.press_key('A');
        assert_eq!(signal, char_to_index('C'));
        assert_eq!(machine.rotor_set[0].position, 3);

        let signal = machine.press_key('T');
        assert_eq!(signal, char_to_index('N'));
        assert_eq!(machine.rotor_set[0].position, 4);
    }
}