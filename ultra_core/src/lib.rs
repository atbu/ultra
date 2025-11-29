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
    fn new(wiring_string: &str, position: char, notch: char) -> Self {
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

pub fn rotate_rotors(mut rotors: RotorSet) -> RotorSet {
    let mut should_rotate = true;
    for rotor in &mut rotors {
        if should_rotate {
            rotor.rotate();
        }

        if rotor.position == rotor.notch {
            should_rotate = true;
        } else {
            should_rotate = false;
        }
    }

    rotors
}

pub fn press_key(signal: char, rotors: RotorSet) -> (u8, RotorSet) {
    let rotors = rotate_rotors(rotors);
    let signal: u8 = char_to_index(signal);

    let mut signal: u8 = signal;
    for rotor in &rotors {
        signal = map_through_rotor(signal, rotor);
    }

    (signal, rotors)
}

pub fn map_through_rotor(signal: u8, rotor: &Rotor) -> u8 {
    let after_entrance_offset = (signal + rotor.position) % 26;
    let after_lookup = rotor.wiring[after_entrance_offset as usize];
    let after_exit_offset = (after_lookup - rotor.position) % 26;

    after_exit_offset
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

        let rotor_i = Rotor::new(ROTOR_I_WIRING, START_POSITION, ROTOR_I_TURNOVER);

        // rotor at index 0 is rightmost, ascending order from right to left
        let rotors: RotorSet = [rotor_i];

        let (signal, rotors) = press_key('A', rotors);
        assert_eq!(signal, char_to_index('K'));
        assert_eq!(rotors[0].position, 2);

        let (signal, rotors) = press_key('A', rotors);
        assert_eq!(signal, char_to_index('C'));
        assert_eq!(rotors[0].position, 3);

        let (signal, rotors) = press_key('T', rotors);
        assert_eq!(signal, char_to_index('N'));
        assert_eq!(rotors[0].position, 4);
    }
}