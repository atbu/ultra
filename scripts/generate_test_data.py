import random

possible_rotor_configurations = ["rotor.RotorI", "rotor.RotorII", "rotor.RotorIII", "rotor.RotorIV", "rotor.RotorV", "rotor.RotorVI", "rotor.RotorVII", "rotor.RotorVIII"]
possible_fourth_rotor_configurations = ["rotor.FourthRotorBeta", "rotor.FourthRotorGamma"]
possible_reflector_configurations = ["reflector.ReflectorA", "reflector.ReflectorB", "reflector.ReflectorC", "reflector.ReflectorNarrowB", "reflector.ReflectorNarrowC"]
alphabet = list("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

rotor_configurations = []
for i in range(3):
    configuration = random.choice(possible_rotor_configurations)
    starting_position = random.choice(alphabet)
    ring_setting = random.choice(alphabet)
    
    rotor_configurations.append(f"rotor.NewRotor({configuration}, '{starting_position}', '{ring_setting}')")

fourth_rotor = random.choice(possible_fourth_rotor_configurations)
fourth_rotor_starting_position = random.choice(alphabet)
fourth_rotor_configuration = f"rotor.NewFourthRotor({fourth_rotor}, '{fourth_rotor_starting_position}')"

reflector_configuration = f"reflector.New({random.choice(possible_reflector_configurations)})"

plugboard_configuration = ""
for i in range(20):
    x = random.choice(alphabet)
    alphabet.remove(x)
    plugboard_configuration += x
plugboard_configuration = f"createPlugboard(t, \"{plugboard_configuration}\")"

plaintext = ""
for i in range(random.randint(10,50)):-
    plaintext += random.choice(alphabet)

print("func Test(t *testing.T) {")
print(f"	machine := &Machine{{")
print(f"		{rotor_configurations[0]},")
print(f"		{rotor_configurations[1]},")
print(f"		{rotor_configurations[2]},")
print(f"		{fourth_rotor_configuration},")
print(f"		{reflector_configuration},")
print(f"		{plugboard_configuration},")
print(f"	}}")
print()
print(f"	runTest(t, machine, \"{plaintext}\", \"OUTPUT_HERE\")")
print("}")
