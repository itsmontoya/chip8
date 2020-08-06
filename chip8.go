package main

import (
	"context"

	"github.com/itsmontoya/chip8/vm"
)

// New will return a new instance of Chip8
func New(screenMultiplier float64) *Chip8 {
	var c Chip8
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.screenMultiplier = screenMultiplier
	c.errC = make(chan error, 2)
	return &c
}

// Chip8 manages the various Chip8 components
type Chip8 struct {
	ctx    context.Context
	cancel func()

	screenMultiplier float64

	errC chan error
}

func (c *Chip8) run() {
	var (
		vm  vm.VM
		p   *PixelRenderer
		err error
	)

	if err = vm.Load("./roms/Chip8 Picture.ch8"); err != nil {
		// Error encountered while loading file, return
		c.errC <- err
		return
	}

	// Initialize a new instance of Pixel
	if p, err = newPixel(c.screenMultiplier); err != nil {
		// Error encountered while initializing pixel, return
		c.errC <- err
		return
	}

	// Initialize VM
	vm.Initialize(p)

	// Run the VM and pass the returning value to the error channel
	c.errC <- vm.Run(c.ctx)
}
