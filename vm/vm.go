package vm

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"time"
)

var (
	// ErrRendererNotSet is returned when a VMs Renderer has not been set before calling VM.Run
	ErrRendererNotSet = errors.New("cannot run, renderer not set")
)

const (
	errInvalidOpcodeFmt = "invalid opcode, %s is not supported"
)

const (
	cyclesPerSecond  = 1
	durationPerCycle = time.Second / cyclesPerSecond
)

const (
	errOpcodeNotImplementedFmt = "opcode %s not implemented"
)

// VM emulates a chip8 instance
// The system's memory map
// 0x000-0x1FF - Chip 8 interpreter (contains font set in emu)
// 0x050-0x0A0 - Used for the built in 4x5 pixel font set (0-F)
// 0x200-0xFFF - Program ROM and work RAM
type VM struct {
	memory    memory
	registers [16]byte
	stack     [16]uint16

	programCounter uint16
	indexRegister  uint16
	stackPointer   uint16
	currentOpcode  opcode

	Graphics Graphics
	keypad   Keypad

	// Flags
	needsDraw bool

	// Timers
	delayTimer byte
	soundTimer byte

	// Renderer
	r Renderer
}

// Initialize will initialize the VM
func (v *VM) Initialize(r Renderer) {
	// Clear memory and counters
	v.programCounter = 0x200
	v.indexRegister = 0
	v.currentOpcode = 0

	// Set renderer
	v.r = r

	// Copy fontset bytes to memory starting at 0x50
	copy(v.memory[0x50:], fontset[:])
}

// Load will load a game into the Virtual Machine
func (v *VM) Load(filename string) (err error) {
	var bs []byte
	// Read provided program file
	if bs, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	// Copy program bytes to memory starting at 0x200
	copy(v.memory[0x200:], bs)
	return
}

// Cycle will emulate a chip8 cycle
func (v *VM) Cycle() (needsDraw bool, err error) {
	// Fetch Opcode
	var o opcode
	if o, err = v.fetchOpcode(); err != nil {
		return
	}

	fmt.Println("Opcode!", o.toHex())

	// Execute Opcode
	if err = v.executeOpcode(o); err != nil {
		return
	}

	// Update timers
	v.updateTimers()
	return
}

// SetKeys will set the currently pressed keys
func (v *VM) SetKeys() {
	v.keypad = v.r.GetKeypad()
}

// Run will run the VM until the context expires
func (v *VM) Run(ctx context.Context) (err error) {
	if v.r == nil {
		err = ErrRendererNotSet
		return
	}

	var needsDraw bool
	tkr := time.NewTicker(durationPerCycle)
	for range tkr.C {
		if isDone(ctx) {
			// Context is finished, return
			return
		}

		if needsDraw, err = v.Cycle(); needsDraw {

		}

		if err = v.r.Draw(v.Graphics); err != nil {
			return
		}

		v.SetKeys()
	}

	return
}

func (v *VM) fetchOpcode() (o opcode, err error) {
	// Get first byte from program counter
	firstByte := v.memory[v.programCounter]
	// Get second byte from program counter
	secondByte := v.memory[v.programCounter+1]
	// Combine bytes to become an opcode
	o = opcode(firstByte)<<8 | opcode(secondByte)
	return
}

func (v *VM) execute0x0000(o opcode) (err error) {
	switch o & 0x000F {
	case 0x0000:
		return v.op00E0(o)
	case 0x000E:
		return v.op00EE(o)

	default:
		return v.op0NNN(o)
	}
}

func (v *VM) execute0x1000(o opcode) (err error) {
	return v.op1NNN(o)
}

func (v *VM) execute0x2000(o opcode) (err error) {
	return v.op2NNN(o)
}

func (v *VM) execute0x3000(o opcode) (err error) {
	return v.op3XNN(o)
}

func (v *VM) execute0x4000(o opcode) (err error) {
	return v.op4XNN(o)
}

func (v *VM) execute0x5000(o opcode) (err error) {
	return v.op5XY0(o)
}

func (v *VM) execute0x6000(o opcode) (err error) {
	return v.op6XNN(o)
}

func (v *VM) execute0x7000(o opcode) (err error) {
	return v.op7XNN(o)
}

func (v *VM) execute0x8000(o opcode) (err error) {
	switch o & 0x000F {
	case 0x0000:
		return v.op8XY0(o)
	case 0x0001:
		return v.op8XY0(o)
	case 0x0002:
		return v.op8XY0(o)
	case 0x0003:
		return v.op8XY0(o)
	case 0x0004:
		return v.op8XY0(o)
	case 0x0005:
		return v.op8XY0(o)
	case 0x0006:
		return v.op8XY0(o)
	case 0x0007:
		return v.op8XY0(o)
	case 0x000E:
		return v.op8XY0(o)

	default:
		return fmt.Errorf(errInvalidOpcodeFmt, o.toHex())
	}
}

func (v *VM) execute0x9000(o opcode) (err error) {
	return v.op9XY0(o)
}

func (v *VM) execute0xA000(o opcode) (err error) {
	return v.opANNN(o)
}

func (v *VM) execute0xB000(o opcode) (err error) {
	return v.opBNNN(o)
}

func (v *VM) execute0xC000(o opcode) (err error) {
	return v.opCXNN(o)
}

func (v *VM) execute0xD000(o opcode) (err error) {
	return v.opDXYN(o)
}

func (v *VM) execute0xE000(o opcode) (err error) {
	switch o & 0x000F {
	case 0x000E:
		return v.opEX9E(o)
	case 0x0001:
		return v.opEXA1(o)

	default:
		return fmt.Errorf(errInvalidOpcodeFmt, o.toHex())
	}
}

func (v *VM) execute0xF000(o opcode) (err error) {
	switch o & 0x00FF {
	case 0x0007:
		return v.opFX07(o)
	case 0x000A:
		return v.opFX0A(o)
	case 0x0015:
		return v.opFX15(o)
	case 0x0018:
		return v.opFX18(o)
	case 0x001E:
		return v.opFX1E(o)
	case 0x0029:
		return v.opFX29(o)
	case 0x0033:
		return v.opFX33(o)
	case 0x0055:
		return v.opFX55(o)
	case 0x0065:
		return v.opFX65(o)

	default:
		return fmt.Errorf(errInvalidOpcodeFmt, o.toHex())
	}
}

func (v *VM) executeOpcode(o opcode) (err error) {
	switch o & 0xF000 {
	case 0x0000:
		return v.execute0x0000(o)
	case 0x1000:
		return v.execute0x1000(o)
	case 0x2000:
		return v.execute0x2000(o)
	case 0x3000:
		return v.execute0x3000(o)
	case 0x4000:
		return v.execute0x4000(o)
	case 0x5000:
		return v.execute0x5000(o)
	case 0x6000:
		return v.execute0x6000(o)
	case 0x7000:
		return v.execute0x7000(o)
	case 0x8000:
		return v.execute0x8000(o)
	case 0x9000:
		return v.execute0x9000(o)
	case 0xA000:
		return v.execute0xA000(o)
	case 0xB000:
		return v.execute0xB000(o)
	case 0xC000:
		return v.execute0xC000(o)
	case 0xD000:
		return v.execute0xD000(o)
	case 0xE000:
		return v.execute0xE000(o)
	case 0xF000:
		return v.execute0xF000(o)

	default:
		return fmt.Errorf(errInvalidOpcodeFmt, o.toHex())
	}
}

// Calls machine code routine (RCA 1802 for COSMAC VIP) at address NNN. Not necessary for most ROMs.
func (v *VM) op0NNN(o opcode) (err error) {
	return fmt.Errorf(errOpcodeNotImplementedFmt, "0NNN")
}

// Clears the screen.
func (v *VM) op00E0(o opcode) (err error) {
	v.Graphics.clear()
	v.programCounter += 2
	return
}

// Returns from a subroutine.
func (v *VM) op00EE(o opcode) (err error) {
	return
}

// Jumps to address NNN.
func (v *VM) op1NNN(o opcode) (err error) {
	return
}

// Calls subroutine at NNN.
func (v *VM) op2NNN(o opcode) (err error) {
	// Set current program counter to the stack
	v.stack[v.stackPointer] = v.programCounter
	// Increment stack pointer
	v.stackPointer++
	// Point program counter to NNN
	v.programCounter = uint16(o) & 0x0FFF
	return
}

// Skips the next instruction if VX equals NN. (Usually the next instruction is a jump to skip a code block)
func (v *VM) op3XNN(o opcode) (err error) {
	return
}

// Skips the next instruction if VX doesn't equal NN. (Usually the next instruction is a jump to skip a code block)
func (v *VM) op4XNN(o opcode) (err error) {
	return
}

// Skips the next instruction if VX equals VY. (Usually the next instruction is a jump to skip a code block)
func (v *VM) op5XY0(o opcode) (err error) {
	return
}

// Sets VX to NN.
func (v *VM) op6XNN(o opcode) (err error) {
	// Set VX to NN
	v.registers[(o&0x0F00)>>8] = v.registers[(o&0x00FF)>>8]

	// Increment program counter by 2
	v.programCounter += 2
	return
}

// Adds NN to VX. (Carry flag is not changed)
func (v *VM) op7XNN(o opcode) (err error) {
	return
}

// Sets VX to the value of VY.
func (v *VM) op8XY0(o opcode) (err error) {
	return
}

// Sets VX to VX or VY. (Bitwise OR operation)
func (v *VM) op8XY1(o opcode) (err error) {
	return
}

// Sets VX to VX and VY. (Bitwise AND operation)
func (v *VM) op8XY2(o opcode) (err error) {
	return
}

// Sets VX to VX xor VY.
func (v *VM) op8XY3(o opcode) (err error) {
	return
}

// Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.
func (v *VM) op8XY4(o opcode) (err error) {
	// Get the carry state from the register
	if v.registers[(o&0x00F0)>>4] > (0xFF - v.registers[(o&0x0F00)>>8]) {
		// Carry, set VF to 1
		v.registers[0xF] = 1
	} else {
		// No carry, set VF to 0
		v.registers[0xF] = 0
	}

	// Add VX to VY
	v.registers[(o&0x0F00)>>8] += v.registers[(o&0x00F0)>>4]

	// Increment program counter by 2
	v.programCounter += 2
	return
}

// VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
func (v *VM) op8XY5(o opcode) (err error) {
	return
}

// Stores the least significant bit of VX in VF and then shifts VX to the right by 1.[b]
func (v *VM) op8XY6(o opcode) (err error) {
	return
}

// Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
func (v *VM) op8XY7(o opcode) (err error) {
	return
}

// Stores the most significant bit of VX in VF and then shifts VX to the left by 1.[b]
func (v *VM) op8XYE(o opcode) (err error) {
	return
}

// Skips the next instruction if VX doesn't equal VY. (Usually the next instruction is a jump to skip a code block)
func (v *VM) op9XY0(o opcode) (err error) {
	return
}

// Sets I to the address NNN.
func (v *VM) opANNN(o opcode) (err error) {
	v.indexRegister = uint16(o & 0x0FFF)
	v.programCounter += 2
	return
}

// Jumps to the address NNN plus V0.
func (v *VM) opBNNN(o opcode) (err error) {
	return
}

// Sets VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN.
func (v *VM) opCXNN(o opcode) (err error) {
	return
}

// Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N pixels. Each row of 8 pixels is read as bit-coded starting from memory location I; I value doesn’t change after the execution of this instruction. As described above, VF is set to 1 if any screen pixels are flipped from set to unset when the sprite is drawn, and to 0 if that doesn’t happen
func (v *VM) opDXYN(o opcode) (err error) {
	//	var pixel byte
	//	x := v.registers[(o&0x0F00)>>8]
	//	y := v.registers[(o&0x00F0)>>4]
	//	height := uint16(o) & 0x000F

	v.registers[0xF] = 0

	/*

		TODO: Work on this a bit

		for yLine := uint16(0); yLine < height; yLine++ {
			pixel = v.memory[v.indexRegister+yLine]
			for xLine := uint16(0); xLine < 8; xLine++ {
				if pixel&(0x80>>xLine) != 0 {
					if v.Graphics[(x+xLine+((y+yLine)*64))] == 1 {
						v.memory[0xF] = 1
					}

					v.Graphics[x+xLine+((y+yLine)*64)] ^= 1
				}

			}
		}

	*/

	// Set needs draw flag to true
	v.needsDraw = true

	// Increment program counter by two
	v.programCounter += 2
	return
}

// Skips the next instruction if the key stored in VX is pressed. (Usually the next instruction is a jump to skip a code block)
func (v *VM) opEX9E(o opcode) (err error) {
	return
}

// Skips the next instruction if the key stored in VX isn't pressed. (Usually the next instruction is a jump to skip a code block)
func (v *VM) opEXA1(o opcode) (err error) {
	return
}

// Sets VX to the value of the delay timer.
func (v *VM) opFX07(o opcode) (err error) {
	return
}

// A key press is awaited, and then stored in VX. (Blocking Operation. All instruction halted until next key event)
func (v *VM) opFX0A(o opcode) (err error) {
	return
}

// Sets the delay timer to VX.
func (v *VM) opFX15(o opcode) (err error) {
	return
}

// Sets the sound timer to VX.
func (v *VM) opFX18(o opcode) (err error) {
	return
}

// Adds VX to I. VF is not affected.[c]
func (v *VM) opFX1E(o opcode) (err error) {
	return
}

// Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font.
func (v *VM) opFX29(o opcode) (err error) {
	return
}

//  Stores the binary-coded decimal representation of VX, with the most significant of three digits at the address in I, the middle digit at I plus 1, and the least significant digit at I plus 2. (In other words, take the decimal representation of VX, place the hundreds digit in memory at location in I, the tens digit at location I+1, and the ones digit at location I+2.)
func (v *VM) opFX33(o opcode) (err error) {
	// Solution credit to TJA (http://www.multigesture.net/wp-content/uploads/mirror/goldroad/chip8.shtml)
	v.memory[v.indexRegister] = v.registers[(o&0x0F00)>>8] / 100
	v.memory[v.indexRegister+1] = (v.registers[(o&0x0F00)>>8] / 10) % 10
	v.memory[v.indexRegister+2] = (v.registers[(o&0x0F00)>>8] % 100) % 10

	// Increment program counter by 2
	v.programCounter += 2
	return
}

// Stores V0 to VX (including VX) in memory starting at address I. The offset from I is increased by 1 for each value written, but I itself is left unmodified.[d]
func (v *VM) opFX55(o opcode) (err error) {
	return
}

// Fills V0 to VX (including VX) with values from memory starting at address I. The offset from I is increased by 1 for each value written, but I itself is left unmodified.[d]
func (v *VM) opFX65(o opcode) (err error) {
	return
}

func (v *VM) updateTimers() {

}
