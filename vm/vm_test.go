package vm

import (
	"fmt"
	"testing"
)

func TestVM_op8XY4(t *testing.T) {
	var (
		vm VM
		o  opcode
	)

	o = opcode(0)<<4 | opcode(7)
	o = opcode(0)<<8 | opcode(13)

	fmt.Println("O", o)
	fmt.Println("1", uint16(o)&0x0F00)
	fmt.Println("2", uint16(o)&0x00F0)
	fmt.Println(vm.op8XY4(0x07E0))
}
