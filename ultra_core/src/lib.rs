/// Converts a character to a zero-based index related to its position in the alphabet, where 'A'
/// equals 0, 'B' equals 1, ..., 'Z' equals 25.
fn char_to_index(char: char) -> u8 {
    char as u8 - 65
}

/// Converts a zero-based index into a character, where 0 equals 'A', 1 equals 'B', ..., 25 equals
/// 'Z'.
fn index_to_char(index: u8) -> char {
    (index + 65) as char
}

/// Converts a wiring string, e.g. `"EKMFLGDQVZNTOWYHXUSPAIBRCJ"` to an array of its character
/// indexes, like `[4, 10, 12, ...]`.
fn wiring_string_to_array(wiring_string: &str) -> [u8; 26] {
    let mut wiring: [u8; 26] = [0; 26];

    for (index, character) in wiring_string.chars().enumerate() {
        wiring[index] = char_to_index(character);
    }

    wiring
}

/// Converts a wiring string, e.g. `"EKMFLGDQVZNTOWYHXUSPAIBRCJ"` to an array of its character
/// indexes, like `[4, 10, 12, ...]`.
fn inverse_wiring_array(wiring_array: [u8; 26]) -> [u8; 26] {
    let mut inverted: [u8; 26] = [0; 26];

    for (original_index, &original_value) in wiring_array.iter().enumerate() {
        inverted[original_value as usize] = original_index as u8;
    }

    inverted
}

/// Represents an Enigma machine (currently only an Enigma I although other models might become
/// supported in the future).
pub struct EnigmaMachine {
    /// The rotor in the left-hand slot of the machine.
    pub left_rotor: Rotor,
    /// The rotor in the middle slot of the machine.
    pub middle_rotor: Rotor,
    /// The rotor in the right-hand slot of the machine.
    pub right_rotor: Rotor,
    /// The machine's reflector.
    pub reflector: Reflector,
    /// The machine's plugboard configuration (if configured).
    pub plugboard: Option<Plugboard>
}

impl EnigmaMachine {
    /// Takes the key that was pressed by the user and returns the character that would light up.
    fn press_key(&mut self, key: char) -> char {
        let mut signal: u8 = char_to_index(key);

        signal = self.map_through_plugboard(signal);

        self.rotate_rotors();

        signal = self.right_rotor.map_signal(signal, false);
        signal = self.middle_rotor.map_signal(signal, false);
        signal = self.left_rotor.map_signal(signal, false);

        signal = self.reflector.map_signal(signal);

        signal = self.left_rotor.map_signal(signal, true);
        signal = self.middle_rotor.map_signal(signal, true);
        signal = self.right_rotor.map_signal(signal, true);

        signal = self.map_through_plugboard(signal);

        index_to_char(signal)
    }

    /// Rotates the rotors to the next state.
    fn rotate_rotors(&mut self) {
        let middle_in_notch = self.middle_rotor.position == self.middle_rotor.notch;
        let right_in_notch = self.right_rotor.position == self.right_rotor.notch;

        self.right_rotor.rotate();

        if middle_in_notch {
            self.left_rotor.rotate();
        }

        if middle_in_notch || right_in_notch {
            self.middle_rotor.rotate();
        }
    }

    /// Map a single character through the plugboard to its matching character (if any).
    fn map_through_plugboard(&self, signal: u8) -> u8 {
        // If we have a plugboard, map the signal through it.
        // If not, let the signal pass through as it is.
        match &self.plugboard {
            Some(plugboard) => plugboard.map_signal(signal),
            None => signal
        }
    }

    /// Set the rotors to the positions specified by parameters.
    #[allow(dead_code)] // Used in tests but shows as dead code, so override compiler warnings.
    fn update_rotor_positions(&mut self, left_rotor: char, middle_rotor: char, right_rotor: char) {
        self.left_rotor.position = char_to_index(left_rotor);
        self.middle_rotor.position = char_to_index(middle_rotor);
        self.right_rotor.position = char_to_index(right_rotor);
    }

    /// Take a message as a string, process each character individually using this Enigma machine
    /// and return the output.
    pub fn process(&mut self, message: &str) -> String {
        // We know that the output won't be any longer/shorter than the input as letters are
        // translated 1:1.
        let mut output = String::with_capacity(message.len());

        if message.is_empty() {
            panic!("Plaintext cannot be empty.");
        }

        // Check if all characters are uppercase letters. If not, panic.
        if !message.chars().all(|x| x.is_ascii_uppercase()) {
            panic!("Plaintext must consist only of uppercase letters.");
        }

        for character in message.chars() {
            output.push(self.press_key(character));
        }

        output
    }
}

/// Represents a single Rotor in an Enigma machine.
pub struct Rotor {
    /// This rotor's mappings - i.e. if this wiring array was `[3, 7, 12, 8, 5]` then this would
    /// indicate that the character `A` maps through this rotor to the character `D`, that the
    /// character `B` maps to the character `H` and so on.
    wiring: [u8; 26],
    /// This rotor's inverse mappings, i.e. the wiring array with its indexes and values swapped.
    inverse_wiring: [u8; 26],
    /// This rotor's notch position, i.e. the position at which it will cause the next rotor to
    /// rotate.
    notch: u8,
    /// This rotor's current position.
    position: u8,
    /// This rotor's current ring setting, i.e. its wiring offset.
    ring_setting: u8
}

/// Represents the five standard rotor configurations of an Enigma I.
pub enum RotorConfiguration {
    I,
    II,
    III,
    IV,
    V
}

impl Rotor {
    /// Constructs a Rotor. Planning to add custom rotor configurations rather than having to use
    /// the enum to use one of the five standard variants.
    pub fn new(rotor_configuration: RotorConfiguration, starting_position: char, ring_setting: char) -> Self {
        // Map the RotorConfiguration enum value to the actual wiring string and notch position.
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
            // Store the reversed variant of the wiring array so it can be used in processing
            // without having to calculate on the fly each time.
            inverse_wiring: inverse_wiring_array(wiring),
            notch: char_to_index(notch),
            position: char_to_index(starting_position),
            ring_setting: char_to_index(ring_setting)
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
    /// If the `inverse` Boolean parameter is true, the inverse wiring array is used (for return
    /// runs through the rotor set). If it is not, the standard wiring array is used.
    fn map_signal(&self, signal: u8, inverse: bool) -> u8 {
        let delta = ((self.position + 26) - self.ring_setting) % 26;
        let contact_in = (signal + delta) % 26;
        let contact_out = match inverse {
            false => self.wiring[contact_in as usize],
            true => self.inverse_wiring[contact_in as usize]
        };
        let signal_out = ((contact_out + 26) - delta) % 26;

        signal_out
    }
}

/// Represents a reflector in an Enigma machine.
pub struct Reflector {
    /// The mapping of the reflector, representing what each character leaves the reflector as.
    wiring: [u8; 26]
}

/// Represents the three standard reflector configurations of an Enigma I.
pub enum ReflectorConfiguration {
    A,
    B,
    C
}

impl Reflector {
    /// Constructs a Reflector. Planning to add custom reflector configurations rather than having
    /// to use the enum to use one of the three standard variants.
    pub fn new(reflector_configuration: ReflectorConfiguration) -> Self {
        // Map the ReflectorConfiguration enum value to the actual wiring string.
        let wiring_string = match reflector_configuration {
            ReflectorConfiguration::A => "EJMZALYXVBWFCRQUONTSPIKHGD",
            ReflectorConfiguration::B => "YRUHQSLDPXNGOKMIEBFZCWVJAT",
            ReflectorConfiguration::C => "FVPJIAOYEDRZXWGCTKUQSBNMHL"
        };

        Self {
            wiring: wiring_string_to_array(wiring_string)
        }
    }

    /// Maps a signal through the reflector, returning the index of the character as it leaves
    /// the reflector.
    fn map_signal(&self, signal: u8) -> u8 {
        self.wiring[signal as usize]
    }
}

/// Represents the plugboard of an Enigma machine. Contains a `Vector` of `(u8, u8)` tuples.
/// The two `u8` values in the tuple correspond to the pairing of letters on a plugboard.
/// A tuple of `(3, 8)` means that an `D` would leave the plugboard as `I` and vice versa.
pub struct Plugboard {
    configuration: Vec<(u8, u8)>
}

impl Plugboard {
    /// This constructor returns an `Option<Self>` because there is the potential that if the string
    /// is empty, we will consider the machine to have no plugboard.
    /// The plugboard configuration is passed as a string like "ABCDEFGH" where A will be paired
    /// with B, C will be paired with D, E will be paired with F and G will be paired with H.
    pub fn new(plugboard_string: &str) -> Option<Self> {
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

        // Convert the plugboard string to a `Vec<(u8, u8)>` as explained in the documentation
        // of the Plugboard struct.
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
    fn test_case_1_sanity_check() {
        let mut machine: EnigmaMachine = EnigmaMachine {
            left_rotor: Rotor::new(RotorConfiguration::I, 'A', 'A'),
            middle_rotor: Rotor::new(RotorConfiguration::II, 'A', 'A'),
            right_rotor: Rotor::new(RotorConfiguration::III, 'A', 'A'),
            reflector: Reflector::new(ReflectorConfiguration::B),
            plugboard: None
        };

        assert_eq!(machine.process("AAAAA"), "BDZGO");
    }

    #[test]
    fn test_case_2_reciprocity() {
        let mut machine: EnigmaMachine = EnigmaMachine {
            reflector: Reflector::new(ReflectorConfiguration::B),
            left_rotor: Rotor::new(RotorConfiguration::I, 'M', 'A'),
            middle_rotor: Rotor::new(RotorConfiguration::II, 'C', 'A'),
            right_rotor: Rotor::new(RotorConfiguration::III, 'K', 'A'),
            plugboard: None
        };

        assert_eq!(machine.process("ENIGMA"), "QMJIDO");

        machine.update_rotor_positions('M', 'C', 'K');

        assert_eq!(machine.process("QMJIDO"), "ENIGMA");
    }

    #[test]
    fn test_case_3_normal_turnover() {
        let mut machine: EnigmaMachine = EnigmaMachine {
            reflector: Reflector::new(ReflectorConfiguration::B),
            left_rotor: Rotor::new(RotorConfiguration::I, 'K', 'A'),
            middle_rotor: Rotor::new(RotorConfiguration::II, 'D', 'A'),
            right_rotor: Rotor::new(RotorConfiguration::III, 'O', 'A'),
            plugboard: None
        };

        assert_eq!(machine.process("AAAAA"), "JWZBJ");

        machine.update_rotor_positions('K', 'D', 'U');

        assert_eq!(machine.process("AAAAA"), "YWDVQ");
    }

    #[test]
    fn test_case_4_double_step() {
        let mut machine: EnigmaMachine = EnigmaMachine {
            reflector: Reflector::new(ReflectorConfiguration::B),
            left_rotor: Rotor::new(RotorConfiguration::I, 'A', 'A'),
            middle_rotor: Rotor::new(RotorConfiguration::II, 'D', 'A'),
            right_rotor: Rotor::new(RotorConfiguration::III, 'U', 'A'),
            plugboard: None
        };

        assert_eq!(machine.process("AAAAA"), "EQIBM");
    }

    #[test]
    fn test_case_5_plugboard() {
        let mut machine: EnigmaMachine = EnigmaMachine {
            reflector: Reflector::new(ReflectorConfiguration::B),
            left_rotor: Rotor::new(RotorConfiguration::I, 'Z', 'A'),
            middle_rotor: Rotor::new(RotorConfiguration::II, 'Z', 'A'),
            right_rotor: Rotor::new(RotorConfiguration::III, 'Z', 'A'),
            plugboard: Plugboard::new("ABCDEFGH")
        };

        assert_eq!(machine.process("AAAAA"), "UZYRQ");
    }

    #[test]
    fn test_case_6_full_integration() {
        let mut machine: EnigmaMachine = EnigmaMachine {
            reflector: Reflector::new(ReflectorConfiguration::B),
            left_rotor: Rotor::new(RotorConfiguration::II, 'A', 'A'),
            middle_rotor: Rotor::new(RotorConfiguration::IV, 'B', 'A'),
            right_rotor: Rotor::new(RotorConfiguration::V, 'L', 'A'),
            plugboard: Plugboard::new("BQCRDIEJKWMTOSPXUZGH")
        };

        assert_eq!(
            machine.process("EVERYTHINGISGOINGEXTREMELYWELL"),
            "LLSDWFYUVEVDHBJVTWWECZNWYXLCNX"
        );
    }

    #[test]
    fn ring_test_case_1() {
        let mut machine: EnigmaMachine = EnigmaMachine {
            reflector: Reflector::new(ReflectorConfiguration::B),
            left_rotor: Rotor::new(RotorConfiguration::I, 'A', 'B'),
            middle_rotor: Rotor::new(RotorConfiguration::II, 'A', 'B'),
            right_rotor: Rotor::new(RotorConfiguration::III, 'A', 'B'),
            plugboard: None
        };

        assert_eq!(machine.process("AAAAA"), "EWTYX");
    }

    #[test]
    fn ring_test_case_2() {
        let mut machine: EnigmaMachine = EnigmaMachine {
            reflector: Reflector::new(ReflectorConfiguration::B),
            left_rotor: Rotor::new(RotorConfiguration::I, 'K', 'A'),
            middle_rotor: Rotor::new(RotorConfiguration::II, 'D', 'A'),
            right_rotor: Rotor::new(RotorConfiguration::III, 'T', 'B'),
            plugboard: None
        };

        assert_eq!(machine.process("AAAA"), "JTIN");
    }

    #[test]
    fn ring_test_case_3() {
        let mut machine: EnigmaMachine = EnigmaMachine {
            reflector: Reflector::new(ReflectorConfiguration::B),
            left_rotor: Rotor::new(RotorConfiguration::I, 'G', 'R'),
            middle_rotor: Rotor::new(RotorConfiguration::II, 'U', 'T'),
            right_rotor: Rotor::new(RotorConfiguration::III, 'M', 'M'),
            plugboard: Plugboard::new("AKSORILP")
        };

        assert_eq!(machine.process("HELLOWORLD"), "CDKSEVMKXJ");
    }

    #[test]
    fn random_test_case_1() {
        let mut machine = EnigmaMachine {
            left_rotor: Rotor::new(RotorConfiguration::II, 'F', 'D'),
            middle_rotor: Rotor::new(RotorConfiguration::II, 'P', 'W'),
            right_rotor: Rotor::new(RotorConfiguration::II, 'K', 'L'),
            reflector: Reflector::new(ReflectorConfiguration::B),
            plugboard: Plugboard::new("JWYLFKREVPXTHOBCMQZG")
        };

        assert_eq!(machine.process("UISDNUUINSNASIAASNUUIDIADIADDDNNNS"), "LNLCJPIFILXIKZPROFOZATVGWZUWOFBFVB");
    }

    #[test]
    fn random_test_case_2() {
        let mut machine = EnigmaMachine {
            left_rotor: Rotor::new(RotorConfiguration::III, 'K', 'B'),
            middle_rotor: Rotor::new(RotorConfiguration::III, 'E', 'T'),
            right_rotor: Rotor::new(RotorConfiguration::V, 'C', 'H'),
            reflector: Reflector::new(ReflectorConfiguration::B),
            plugboard: Plugboard::new("ACPQUHYFWRMJOSKTDIVG")
        };

        assert_eq!(machine.process("EEBXZZEBZXNXLBLBBZNLZBLNLNBBNBLLNBXNEZLB"), "BPNAFCWSDBGAFDIQPKGHXNFMXIGIKLXPKTPORWOX");
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