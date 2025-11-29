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

pub fn rotate_rotors(mut rotors: [Rotor; 3]) -> [Rotor; 3] {
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

pub fn press_key(signal: u8, rotors: [Rotor; 3]) -> (u8, [Rotor; 3]) {
    let rotors = rotate_rotors(rotors);

    let signal = map_through_rotor(
        map_through_rotor(
            map_through_rotor(
                signal,
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
        let wiring_string: &str = "EKMFLGDQVZNTOWYHXUSPAIBRCJ";
        let mut wiring: [u8; 26] = [0; 26];
        for (index, character) in wiring_string.chars().enumerate() {
            wiring[index] = char_to_index(character);
        }

        let q = char_to_index('Q');
        let a = char_to_index('A');

        let rotor_a: Rotor = Rotor {
            wiring,
            position: q,
            notch: q
        };

        let rotor_b: Rotor = Rotor {
            wiring,
            position: a,
            notch: q
        };

        let rotor_c: Rotor = Rotor {
            wiring,
            position: a,
            notch: q
        };

        // rotor at index 0 is rightmost, ascending order from right to left
        let rotors: [Rotor; 3] = [rotor_a, rotor_b, rotor_c];
        let (_, rotors) = press_key(char_to_index('A'), rotors);

        assert_eq!(rotors[0].position, q);
        assert_eq!(rotors[1].position, a + 1);
        assert_eq!(rotors[2].position, a);
    }
}