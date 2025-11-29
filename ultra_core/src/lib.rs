const NUMBER_OF_ROTORS: usize = 3;

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
    let mut should_rotate = false;
    for rotor in &mut rotors {
        println!("{should_rotate}");

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
    let signal_index = char_to_index(signal);

    let signal = map_through_rotor(
        map_through_rotor(
            map_through_rotor(
                signal_index,
                &rotors[0]
            ),
            &rotors[1]
        ),
        &rotors[2],
    );

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
        let wiring_string = "EKMFLGDQVZNTOWYHXUSPAIBRCJ";

        let rotor_a = Rotor::new(wiring_string, 'Q', 'Q');
        let rotor_b = Rotor::new(wiring_string, 'A', 'Q');
        let rotor_c = Rotor::new(wiring_string, 'A', 'Q');

        // rotor at index 0 is rightmost, ascending order from right to left
        let rotors: RotorSet = [rotor_a, rotor_b, rotor_c];
        let (_, rotors) = press_key('A', rotors);

        assert_eq!(rotors[0].position, char_to_index('Q'));
        assert_eq!(rotors[1].position, char_to_index('A') + 1);
        assert_eq!(rotors[2].position, char_to_index('A'));
    }
}