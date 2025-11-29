pub struct Rotor {
    wiring: [u8; 26],
    position: u8,
    notch: u8
}

impl Rotor {
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

pub fn press_key(signal: u8, rotors: &mut [&mut Rotor; 3]) -> u8 {
    // There has to be a better way of doing this
    if rotors[0].rotate() {
        if rotors[1].rotate() {
            rotors[2].rotate();
        }
    }

    map_through_rotor(
        map_through_rotor(
            map_through_rotor(
                signal,
                rotors[0]
            ),
            rotors[1]
        ),
        rotors[2],
    )
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
        let wiring_string: &str = "EKMFLGDQVZNTOWYHXUSPAIBRCJ";
        let mut wiring: [u8; 26] = [0; 26];
        for (index, character) in wiring_string.chars().enumerate() {
            wiring[index] = char_to_index(character);
        }

        let mut rotor_a: Rotor = Rotor {
            wiring,
            position: char_to_index('Q'),
            notch: char_to_index('Q')
        };

        let mut rotor_b: Rotor = Rotor {
            wiring,
            position: char_to_index('A'),
            notch: char_to_index('Q')
        };

        let mut rotor_c: Rotor = Rotor {
            wiring,
            position: char_to_index('A'),
            notch: char_to_index('Q')
        };

        let mut rotors: [&mut Rotor; 3] = [&mut rotor_a, &mut rotor_b, &mut rotor_c];

        assert_eq!(rotors[0].position, char_to_index('Q'));
        assert_eq!(rotors[1].position, char_to_index('A'));
        assert_eq!(rotors[2].position, char_to_index('A'));
        press_key(char_to_index('A'), &mut rotors);
        assert_eq!(rotors[0].position, char_to_index('R'));
        assert_eq!(rotors[1].position, char_to_index('B'));
        assert_eq!(rotors[2].position, char_to_index('A'));
    }
}