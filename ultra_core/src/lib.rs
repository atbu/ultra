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

pub fn inverse_wiring_array(wiring_array: [u8; 26]) -> [u8; 26] {
    let mut inverted: [u8; 26] = [0; 26];

    for (original_index, &original_value) in wiring_array.iter().enumerate() {
        inverted[original_value as usize] = original_index as u8;
    }

    inverted
}

pub struct EnigmaMachine {
    left_rotor: Rotor,
    middle_rotor: Rotor,
    right_rotor: Rotor,
    reflector: Reflector
}

impl EnigmaMachine {
    /// Takes the key that was pressed by the user and returns the character that would light up.
    fn press_key(&mut self, key: char) -> char {
        self.rotate_rotors();

        let mut signal: u8 = char_to_index(key);
        signal = self.right_rotor.map_signal(signal);
        signal = self.middle_rotor.map_signal(signal);
        signal = self.left_rotor.map_signal(signal);
        signal = self.reflector.map_signal(signal);
        signal = self.left_rotor.map_signal_inverse(signal);
        signal = self.middle_rotor.map_signal_inverse(signal);
        signal = self.right_rotor.map_signal_inverse(signal);

        println!("Entered as {}, leaving as {}", key, index_to_char(signal));

        index_to_char(signal)
    }

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
    inverse_wiring: [u8; 26],
    notch: u8,
    position: u8
}

impl Rotor {
    fn new(wiring_string: &str, notch: char, position: char) -> Self {
        let wiring: [u8; 26] = wiring_string_to_array(wiring_string);

        Self {
            wiring,
            inverse_wiring: inverse_wiring_array(wiring),
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

    /// Maps a signal through a single Rotor, taking into account the rotation of the rotor.
    /// https://en.wikipedia.org/wiki/Enigma_rotor_details#Rotor_offset
    fn map_signal(&self, signal: u8) -> u8 {
        let after_shift_in = (signal + self.position) % 26;
        let mapped_value = self.wiring[after_shift_in as usize];
        let after_shift_out = (mapped_value + 26 - self.position) % 26;

        after_shift_out
    }

    fn map_signal_inverse(&self, signal: u8) -> u8 {
        let after_shift_in = (signal + self.position) % 26;
        let mapped_value = self.inverse_wiring[after_shift_in as usize];
        let after_shift_out = (mapped_value + 26 - self.position) % 26;

        after_shift_out
    }
}

struct Reflector {
    wiring: [u8; 26]
}

impl Reflector {
    fn new(wiring_string: &str) -> Self {
        Self {
            wiring: wiring_string_to_array(wiring_string)
        }
    }

    fn map_signal(&self, signal: u8) -> u8 {
        self.wiring[signal as usize]
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    // /// Makes checking the rotor configuration easier.
    // fn check_rotor_configuration(machine: &EnigmaMachine, left_rotor: char, middle_rotor: char, right_rotor: char) -> bool {
    //     machine.left_rotor.position == char_to_index(left_rotor)
    //     && machine.middle_rotor.position == char_to_index(middle_rotor)
    //     && machine.right_rotor.position == char_to_index(right_rotor)
    // }

    #[test]
    fn test_rotor_stepping() {
        let mut machine: EnigmaMachine = EnigmaMachine {
            left_rotor: Rotor::new("EKMFLGDQVZNTOWYHXUSPAIBRCJ", 'Q', 'A'),
            middle_rotor: Rotor::new("AJDKSIRUXBLHWTMCQGZNPYFVOE", 'E', 'A'),
            right_rotor: Rotor::new("BDFHJLCPRTXVZNYEIWGAKMUSQO", 'V', 'A'),
            reflector: Reflector::new("YRUHQSLDPXNGOKMIEBFZCWVJAT")
        };

        assert_eq!(machine.press_key('A'), 'B');
        assert_eq!(machine.press_key('A'), 'D');
        assert_eq!(machine.press_key('A'), 'Z');
        assert_eq!(machine.press_key('A'), 'G');
        assert_eq!(machine.press_key('A'), 'O');
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