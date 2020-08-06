package main

import "context"

// New will return a new instance of Chip8
func New(screenMultiplier float64) *Chip8 {
	var c Chip8
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.screenMultiplier = screenMultiplier
	c.errC = make(chan error, 1)
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
		p   *PixelRenderer
		vm  VM
		err error
	)

	// Initialize a new instance of Pixel
	if p, err = newPixel(c.screenMultiplier); err != nil {
		// Error encountered while initializing pixel, return
		c.errC <- err
		return
	}

	// Initialize VM
	vm.Initialize(p)

	// Run the VM and pass the returning value to the error channel
	c.errC <- vm.Run(ctx)
}
