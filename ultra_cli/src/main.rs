use std::env;
use ultra_core::{
    EnigmaMachine,
    Rotor,
    RotorConfiguration,
    Reflector,
    ReflectorConfiguration,
    Plugboard
};

fn main() {
    let args: Vec<String> = env::args().collect();

    // args[0] is the path to the binary.
    // args[1] is the first argument after the binary itself.
    let plaintext: &str = &args[1];

    let mut machine: EnigmaMachine = EnigmaMachine {
        left_rotor: Rotor::new(RotorConfiguration::I, 'A', 'A'),
        middle_rotor: Rotor::new(RotorConfiguration::II, 'A', 'A'),
        right_rotor: Rotor::new(RotorConfiguration::III, 'A', 'A'),
        reflector: Reflector::new(ReflectorConfiguration::B),
        plugboard: Plugboard::new("ABCD")
    };

    println!("{}", machine.process(plaintext));
}
