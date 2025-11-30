/// Converts a character to a zero-based index related to its position in the alphabet, where 'A'
/// equals 0, 'B' equals 1, ..., 'Z' equals 25.
pub fn char_to_index(char: char) -> u8 {
    char as u8 - 65
}

/// Converts a zero-based index into a character, where 0 equals 'A', 1 equals 'B', ..., 25 equals
/// 'Z'.
pub fn index_to_char(index: u8) -> char {
    (index + 65) as char
}

/// Converts a wiring string, e.g. `"EKMFLGDQVZNTOWYHXUSPAIBRCJ"` to an array of its character
/// indexes, like `[4, 10, 12, ...]`.
pub fn wiring_string_to_array(wiring_string: &str) -> [u8; 26] {
    let mut wiring: [u8; 26] = [0; 26];

    for (index, character) in wiring_string.chars().enumerate() {
        wiring[index] = char_to_index(character);
    }

    wiring
}

pub struct EnigmaMachine {
    left_rotor: Rotor,
    middle_rotor: Rotor,
    right_rotor: Rotor
}

impl EnigmaMachine {
    /// Takes the key that was pressed by the user and returns the character that would light up.
    // fn press_key(&self, key: char) -> char {
    //     // rotate rotors
    //     // pass signal through rotors
    // }

    fn rotate_rotors(&mut self) {
        // Checks double stepping functionality of middle rotor.
        let middle_rotor_should_rotate = self.right_rotor.position == self.right_rotor.notch
            || self.right_rotor.position == self.right_rotor.notch + 1;

        let left_rotor_should_rotate = self.middle_rotor.position == self.middle_rotor.notch;

        self.right_rotor.rotate(); // right rotor always rotates

        if middle_rotor_should_rotate {
            self.middle_rotor.rotate();
        }

        if left_rotor_should_rotate {
            self.left_rotor.rotate();
        }
    }
}

struct Rotor {
    wiring: [u8; 26],
    notch: u8,
    position: u8
}

impl Rotor {
    fn new(wiring_string: &str, notch: char, position: char) -> Self {
        Self {
            wiring: wiring_string_to_array(wiring_string),
            notch: char_to_index(notch),
            position: char_to_index(position)
        }
    }

    /// Increment the rotor's position by 1. Use modulo to ensure that the position will wrap
    /// around, i.e. if the rotor rotates on Z, it will wrap around to A rather than being stuck
    /// in limbo at some non-existent 27th letter.
    fn rotate(&mut self) {
        self.position = (self.position + 1) % 26
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    /// Makes checking the rotor configuration easier.
    fn check_rotor_configuration(machine: &EnigmaMachine, left_rotor: char, middle_rotor: char, right_rotor: char) -> bool {
        machine.left_rotor.position == char_to_index(left_rotor)
        && machine.middle_rotor.position == char_to_index(middle_rotor)
        && machine.right_rotor.position == char_to_index(right_rotor)
    }

    #[test]
    fn test_rotor_stepping() {
        let mut machine: EnigmaMachine = EnigmaMachine {
            left_rotor: Rotor::new("BDFHJLCPRTXVZNYEIWGAKMUSQO", 'V', 'A'),
            middle_rotor: Rotor::new("AJDKSIRUXBLHWTMCQGZNPYFVOE", 'E', 'D'),
            right_rotor: Rotor::new("EKMFLGDQVZNTOWYHXUSPAIBRCJ", 'Q', 'O')
        };

        assert!(check_rotor_configuration(&machine, 'A', 'D', 'O'));

        machine.rotate_rotors();
        assert!(check_rotor_configuration(&machine, 'A', 'D', 'P'));

        machine.rotate_rotors();
        assert!(check_rotor_configuration(&machine, 'A', 'D', 'Q'));

        machine.rotate_rotors();
        assert!(check_rotor_configuration(&machine, 'A', 'E', 'R'));

        machine.rotate_rotors();
        assert!(check_rotor_configuration(&machine, 'B', 'F', 'S'));

        machine.rotate_rotors();
        assert!(check_rotor_configuration(&machine, 'B', 'F', 'T'));

        machine.rotate_rotors();
        assert!(check_rotor_configuration(&machine, 'B', 'F', 'U'));
    }

    #[test]
    fn test_char_to_index() {
        assert_eq!(char_to_index('A'), 0);
        assert_eq!(char_to_index('B'), 1);
        assert_eq!(char_to_index('C'), 2);
        assert_eq!(char_to_index('Z'), 25);
    }

    #[test]
    fn test_index_to_char() {
        assert_eq!(index_to_char(0), 'A');
        assert_eq!(index_to_char(1), 'B');
        assert_eq!(index_to_char(2), 'C');
        assert_eq!(index_to_char(25), 'Z');
    }
}