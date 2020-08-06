package main

import "fmt"

type opcode uint16

func (o opcode) toHex() hex {
	return hex(fmt.Sprintf("%04X", o))
}
