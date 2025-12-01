import random

possible_rotor_configurations = ["RotorConfiguration::I", "RotorConfiguration::II", "RotorConfiguration::III", "RotorConfiguration::IV", "RotorConfiguration::V"]
possible_reflector_configurations = ["ReflectorConfiguration::A", "ReflectorConfiguration::B", "ReflectorConfiguration::C"]
alphabet = list("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

rotor_configurations = []
for i in range(3):
    configuration = random.choice(possible_rotor_configurations)
    starting_position = random.choice(alphabet)
    ring_setting = random.choice(alphabet)
    
    rotor_configurations.append(f"Rotor::new({configuration}, '{starting_position}', '{ring_setting}')")

reflector_configuration = f"Reflector::new({random.choice(possible_reflector_configurations)})"

plugboard_configuration = ""
for i in range(20):
    x = random.choice(alphabet)
    alphabet.remove(x)
    plugboard_configuration += x
plugboard_configuration = f"Plugboard::new(\"{plugboard_configuration}\")"

plaintext = ""
for i in range(random.randint(10,50)):
    plaintext += random.choice(alphabet)

four_spaces = "    "

print("#[test]")
print("fn test_case() {")
print(f"{four_spaces}let mut machine = EnigmaMachine {{")
print(f"{four_spaces*2}left_rotor: {rotor_configurations[0]},")
print(f"{four_spaces*2}middle_rotor: {rotor_configurations[1]},")
print(f"{four_spaces*2}right_rotor: {rotor_configurations[2]},")
print(f"{four_spaces*2}reflector: {reflector_configuration},")
print(f"{four_spaces*2}plugboard: {plugboard_configuration}")
print(f"{four_spaces}}};")
print()
print(f"{four_spaces}assert_eq!(machine.process(\"{plaintext}\"), OUTPUT_HERE);")
print("}")
