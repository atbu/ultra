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
    reflector: Reflector,
    plugboard: Option<Plugboard>
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

        // If we have a plugboard, map the signal through it.
        // If not, let the signal pass through as it is.
        let signal = match &self.plugboard {
            Some(plugboard) => plugboard.map_signal(signal),
            None => signal
        };

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

enum RotorConfiguration {
    I,
    II,
    III,
    IV,
    V
}

impl Rotor {
    fn new(rotor_configuration: RotorConfiguration) -> Self {
        let (wiring_string, notch) = match rotor_configuration {
            RotorConfiguration::I => ("EKMFLGDQVZNTOWYHXUSPAIBRCJ", 'Q'),
            RotorConfiguration::II => ("AJDKSIRUXBLHWTMCQGZNPYFVOE", 'E'),
            RotorConfiguration::III => ("BDFHJLCPRTXVZNYEIWGAKMUSQO", 'V'),
            RotorConfiguration::IV => ("ESOVPZJAYQUIRHXLNFTGKDCMWB", 'J'),
            RotorConfiguration::V => ("VZBRGITYUPSDNHLXAWMJQOFECK", 'Z')
        };

        let wiring: [u8; 26] = wiring_string_to_array(wiring_string);

        Self {
            wiring,
            inverse_wiring: inverse_wiring_array(wiring),
            notch: char_to_index(notch),
            position: char_to_index('A')
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

enum ReflectorConfiguration {
    A,
    B,
    C
}

struct Reflector {
    wiring: [u8; 26]
}

impl Reflector {
    fn new(reflector_configuration: ReflectorConfiguration) -> Self {
        let wiring_string = match reflector_configuration {
            ReflectorConfiguration::A => "EJMZALYXVBWFCRQUONTSPIKHGD",
            ReflectorConfiguration::B => "YRUHQSLDPXNGOKMIEBFZCWVJAT",
            ReflectorConfiguration::C => "FVPJIAOYEDRZXWGCTKUQSBNMHL"
        };

        Self {
            wiring: wiring_string_to_array(wiring_string)
        }
    }

    fn map_signal(&self, signal: u8) -> u8 {
        self.wiring[signal as usize]
    }
}

struct Plugboard {
    configuration: Vec<(u8, u8)>
}

impl Plugboard {
    /// This constructor returns an Option<Self> because there is the potential that if the string
    /// is empty, or if
    fn new(plugboard_string: &str) -> Option<Self> {
        // If the plugboard string is empty, we can return None as we have no plugboard.
        if plugboard_string.is_empty() {
            return None
        }

        // Make sure that plugboard string has an even number of characters, we don't want a plug
        // that is disconnected on one end.
        if plugboard_string.len() % 2 != 0 {
            let invalid_connection = plugboard_string.as_bytes()[plugboard_string.len() - 1] as char;
            panic!("Plugboard configuration has an invalid connection: {} is not connected to any other character.", invalid_connection);
        }

        let mut plugboard_pairs: Vec<(u8, u8)> = Vec::new();
        for i in (0..(plugboard_string.len() - 1)).step_by(2) {
            let x = char_to_index(plugboard_string.as_bytes()[i] as char);
            let y = char_to_index(plugboard_string.as_bytes()[i + 1] as char);

            plugboard_pairs.push((x, y));
        }

        Some(
            Self {
                configuration: plugboard_pairs
            }
        )
    }

    /// Maps a signal through the plugboard.
    /// Iterate through all plugboard pairs. If there is a pair containing the signal,
    /// return the value on the other end of the pair.
    /// If there is no plugboard pair containing this signal, just return the signal itself
    /// as it is not paired with any other signal.
    fn map_signal(&self, signal: u8) -> u8 {
        for i in &self.configuration {
            if i.0 == signal {
                return i.1
            }

            if i.1 == signal {
                return i.0
            }
        }

        signal
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_rotor_stepping() {
        let mut machine: EnigmaMachine = EnigmaMachine {
            left_rotor: Rotor::new(RotorConfiguration::I),
            middle_rotor: Rotor::new(RotorConfiguration::II),
            right_rotor: Rotor::new(RotorConfiguration::III),
            reflector: Reflector::new(ReflectorConfiguration::B),
            plugboard: Plugboard::new("BL")
        };

        assert_eq!(machine.press_key('A'), 'L');
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