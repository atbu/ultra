# ultra

An Enigma simulator.

Named after the designation given by 🇬🇧 British military intelligence to wartime signals intelligence obtained by
breaking encrypted enemy communications.
Click [here](https://en.wikipedia.org/wiki/Ultra_(cryptography)) for more information.

In its current state, the project simulates an **Enigma I** machine used by the German Armed Forces (Wehrmacht) during
the Second World War.

The primary components of the machine are:
 - Rotors
   - Each rotor has 26 electrical contacts on each face, representing the letters of the alphabet. It contains 26 wires,
   directly linking each letter to another letter, for example an H could be transformed into a Y.
   - Each rotor also has 26 possible starting positions, again representing the letters of the alphabet.
   - To add more complexity, each rotor has a 'ring setting' which offsets the alphabetical values from their electrical
   wiring.
   - On every key press, the rightmost rotor rotates. When a rotor completes a full rotation (26 key presses), the rotor
   to the left will rotate once.
   - A strange concept called double stepping occurs where the middle rotor will rotate
   when the rightmost rotor has completed a full revolution, or if it has completed a full revolution itself, it will
   rotate once more.
   - Enigma I has 3 rotor slots, although the Enigma M4 later used by the Kriegsmarine (German Navy) contained 4.
   - Enigma operators usually had a choice of 5 rotors, numbered *I*, *II*, *III*, *IV* and *V*.
 - Reflector
   - Contains thirteen wires which pair letters together. For example, if B was electrically wired to U, a B entering
   the reflector would leave as a U and vice versa.
   - The way the reflector was designed means, however, that no letter can ever be encrypted to itself - this was a
   severe flaw with Enigma that was exploited by codebreakers.
   - There are three different reflector configurations, labelled *A*, *B* and *C*.
 - Plugboard
   - Allows the operator to easily pair letters together in a similar way to the reflector.
   - For example, if an operator plugged one end of a wire into B and the other end into K, a B entering the plugboard
   would leave as a K.
   - Usually up to 10 pairs would be configured. If a letter does not have a counterpart, it leaves the plugboard as
   itself.

Operation of the machine:
1. The operator presses a key.
2. The signal first goes to the plugboard, which would transform the letter into another letter, if it had a
corresponding letter.
3. The rotors would rotate _before_ any signal passes through them.
   - Rotor stepping is explained above, including _double stepping_.
4. The signal then passes through the rotors from right to left, being transformed at each step, until it reaches the
reflector, at which point it is transformed into its paired letter. It then goes back through the rotors in reverse
order, i.e. left to right.
5. It then reaches the plugboard again, where the signal is transformed into its paired letter yet again (if it has a
paired letter).
6. Finally, the signal goes to the lamps in the centre of the machine. The lamp for the letter corresponding to the
signal will light up.