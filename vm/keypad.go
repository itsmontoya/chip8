package vm

// Keypad represents a keypad state
type Keypad [16]byte

// Set will set the state for
func (k *Keypad) Set(index int, state bool) {
	k[index] = boolToByte(state)
}
