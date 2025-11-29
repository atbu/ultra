pub fn map_through_rotor(signal: u8, wiring: [u8; 26], starting_position: u8) -> u8 {
    let after_entrance_offset = (signal + starting_position) % 26;
    let after_lookup = wiring[after_entrance_offset as usize];
    let after_exit_offset = (after_lookup - starting_position) % 26;

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
        assert_eq!(map_through_rotor(char_to_index('A'), wiring, 1), char_to_index('J'));
    }
}