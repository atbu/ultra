use serde::Deserialize;
use std::path::PathBuf;
use clap::Parser;
use toml;
use ultra_core::{EnigmaMachine, Plugboard, Reflector, ReflectorConfiguration, Rotor, RotorConfiguration};

#[derive(Deserialize)]
struct MachineConfig {
    rotors: RotorsConfig,
    reflector: ReflectorConfig,
    plugboard: PlugboardConfig
}

#[derive(Deserialize)]
struct RotorsConfig {
    left: RotorConfig,
    middle: RotorConfig,
    right: RotorConfig
}

#[derive(Deserialize)]
struct RotorConfig {
    configuration: String,
    starting_position: char,
    ring_setting: char
}

impl RotorConfig {
    fn validate(&mut self, name: &str) {
        if self.starting_position.is_ascii_alphabetic() {
            if !self.starting_position.is_ascii_uppercase() {
                self.starting_position = self.starting_position.to_ascii_uppercase();
            }
        } else {
            panic!("{} rotor starting position must be a letter.", name);
        }

        if self.ring_setting.is_ascii_alphabetic() {
            if !self.ring_setting.is_ascii_uppercase() {
                self.ring_setting = self.ring_setting.to_ascii_uppercase();
            }
        } else {
            panic!("{} rotor ring setting must be a letter.", name);
        }
    }
}

#[derive(Deserialize)]
struct ReflectorConfig {
    configuration: char
}

#[derive(Deserialize)]
struct PlugboardConfig {
    configuration: String
}

#[derive(Parser)]
struct CLI {
    #[arg(short, long, value_name = "INPUT")]
    input: Option<String>,

    #[arg(short, long, value_name = "CONFIG", default_value = "MachineConfiguration.toml")]
    config: Option<PathBuf>
}

fn main() {
    let cli = CLI::parse();

    let config_path = match cli.config.as_deref() {
        Some(config_path) => config_path,
        None => panic!("Machine configuration file not found.")
    };

    let mut machine = read_machine_configuration(config_path.to_str().unwrap());

    let mut plugboard_pairs: Vec<(char, char)> = Vec::new();

    if !machine.plugboard.configuration.is_empty() {
        if !machine.plugboard.configuration.chars().all(|x| x.is_ascii()) {
            panic!("Plugboard configuration must only contain letters.")
        }

        for i in (0..(machine.plugboard.configuration.len() - 1)).step_by(2) {
            let x = machine.plugboard.configuration.as_bytes()[i].to_ascii_uppercase() as char;
            let y = machine.plugboard.configuration.as_bytes()[i + 1].to_ascii_uppercase() as char;

            plugboard_pairs.push((x, y));
        }
    }

    let (
        left_rotor_config,
        middle_rotor_config,
        right_rotor_config
    ) = validate_rotor_configuration_strings(&machine);

    let reflector_config = validate_reflector_config_string(&machine);

    machine.rotors.left.validate("Left");
    machine.rotors.middle.validate("Middle");
    machine.rotors.right.validate("Right");

    print_machine_configuration(&machine, &plugboard_pairs);

    let mut enigma_machine: EnigmaMachine = EnigmaMachine {
        left_rotor: Rotor::new(left_rotor_config, machine.rotors.left.starting_position, machine.rotors.left.ring_setting),
        middle_rotor: Rotor::new(middle_rotor_config, machine.rotors.middle.starting_position, machine.rotors.middle.ring_setting),
        right_rotor: Rotor::new(right_rotor_config, machine.rotors.right.starting_position, machine.rotors.right.ring_setting),
        reflector: Reflector::new(reflector_config),
        plugboard: Plugboard::new(&machine.plugboard.configuration)
    };

    if let Some(input) = cli.input.as_deref() {
        println!("{}", enigma_machine.process(input));
    }
}

fn validate_rotor_configuration_strings(machine: &MachineConfig) -> (RotorConfiguration, RotorConfiguration, RotorConfiguration) {
    let left = match machine.rotors.left.configuration.as_str() {
        "I" => RotorConfiguration::I,
        "II" => RotorConfiguration::II,
        "III" => RotorConfiguration::III,
        "IV" => RotorConfiguration::IV,
        "V" => RotorConfiguration::V,
        "VI" => RotorConfiguration::VI,
        "VII" => RotorConfiguration::VII,
        "VIII" => RotorConfiguration::VIII,
        _ => panic!("Unrecognised left rotor configuration: {}", machine.rotors.left.configuration)
    };

    let middle = match machine.rotors.middle.configuration.as_str() {
        "I" => RotorConfiguration::I,
        "II" => RotorConfiguration::II,
        "III" => RotorConfiguration::III,
        "IV" => RotorConfiguration::IV,
        "V" => RotorConfiguration::V,
        "VI" => RotorConfiguration::VI,
        "VII" => RotorConfiguration::VII,
        "VIII" => RotorConfiguration::VIII,
        _ => panic!("Unrecognised middle rotor configuration: {}", machine.rotors.middle.configuration)
    };

    let right = match machine.rotors.right.configuration.as_str() {
        "I" => RotorConfiguration::I,
        "II" => RotorConfiguration::II,
        "III" => RotorConfiguration::III,
        "IV" => RotorConfiguration::IV,
        "V" => RotorConfiguration::V,
        "VI" => RotorConfiguration::VI,
        "VII" => RotorConfiguration::VII,
        "VIII" => RotorConfiguration::VIII,
        _ => panic!("Unrecognised right rotor configuration: {}", machine.rotors.right.configuration)
    };

    (left, middle, right)
}

fn validate_reflector_config_string(machine: &MachineConfig) -> ReflectorConfiguration {
    match machine.reflector.configuration {
        'A' => ReflectorConfiguration::A,
        'B' => ReflectorConfiguration::B,
        'C' => ReflectorConfiguration::C,
        _ => panic!("Unrecognised reflector configuration: {}", machine.reflector.configuration)
    }
}

// TODO should probably return Result type
fn read_machine_configuration(path: &str) -> MachineConfig {
    let content = match std::fs::read_to_string(path) {
        Ok(content) => content,
        Err(e) => panic!("Failed to read machine configuration file: {}", e)
    };

    let config = match toml::from_str::<MachineConfig>(&content) {
        Ok(config) => config,
        Err(e) => panic!("Failed to parse machine configuration: {}", e)
    };

    config
}

fn print_machine_configuration(machine: &MachineConfig, plugboard_pairs: &Vec<(char, char)>) {
    println!("Rotor choices/order (Walzenlage):");
    println!("\tLeft:");
    println!("\t\tConfiguration: {}", machine.rotors.left.configuration);
    println!("\t\tStarting position (Grundstellung): {}", machine.rotors.left.starting_position);
    println!("\t\tRing setting (Ringstellung): {}", machine.rotors.left.ring_setting);
    println!("\tMiddle:");
    println!("\t\tConfiguration: {}", machine.rotors.middle.configuration);
    println!("\t\tStarting position (Grundstellung): {}", machine.rotors.middle.starting_position);
    println!("\t\tRing setting (Ringstellung): {}", machine.rotors.middle.ring_setting);
    println!("\tRight:");
    println!("\t\tConfiguration: {}", machine.rotors.right.configuration);
    println!("\t\tStarting position (Grundstellung): {}", machine.rotors.right.starting_position);
    println!("\t\tRing setting (Ringstellung): {}", machine.rotors.right.ring_setting);
    println!("Reflector (Umkehrwalze):");
    println!("\tConfiguration: {}", machine.reflector.configuration);
    println!("Plugboard (Steckerbrett):");

    if !plugboard_pairs.is_empty() {
        for pair in plugboard_pairs {
            println!("\t- {} <-> {}", pair.0, pair.1);
        }
    } else {
        println!("\tNot configured.");
    }

    println!();
}