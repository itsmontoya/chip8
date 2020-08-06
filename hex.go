package main

const (
	op0NNN hex = "0NNN"
	op00E0 hex = "00E0"
	op00EE hex = "00EE"
	op1NNN hex = "1NNN"
	op2NNN hex = "2NNN"
	op3XNN hex = "3XNN"
	op4XNN hex = "4XNN"
	op5XY0 hex = "5XY0"
	op6XNN hex = "6XNN"
	op7XNN hex = "7XNN"
	op8XY0 hex = "8XY0"
	op8XY1 hex = "8XY1"
	op8XY2 hex = "8XY2"
	op8XY3 hex = "8XY3"
	op8XY4 hex = "8XY4"
	op8XY5 hex = "8XY5"
	op8XY6 hex = "8XY6"
	op8XY7 hex = "8XY7"
	op8XYE hex = "8XYE"
	op9XY0 hex = "9XY0"
	opANNN hex = "ANNN"
	opBNNN hex = "BNNN"
	opCXNN hex = "CXNN"
	opDXYN hex = "DXYN"
	opEX9E hex = "EX9E"
	opEXA1 hex = "EXA1"
	opFX07 hex = "FX07"
	opFX0A hex = "FX0A"
	opFX15 hex = "FX15"
	opFX18 hex = "FX18"
	opFX1E hex = "FX1E"
	opFX29 hex = "FX29"
	opFX33 hex = "FX33"
	opFX55 hex = "FX55"
	opFX65 hex = "FX65"
)

type hex string

func (h hex) isMatch(in hex) (isMatch bool) {
	for i := 0; i < 4; i++ {
		b := h[i]
		switch b {
		case 'X', 'Y', 'N':
			continue

		default:
			if in[i] != b {
				return false
			}
		}
	}

	return true
}
