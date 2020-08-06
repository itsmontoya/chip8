package main

import (
	"context"
	"flag"
	"os"

	"github.com/Hatch1fy/errors"
	"github.com/faiface/pixel/pixelgl"
	"github.com/hatchify/closer"
	"github.com/hatchify/scribe"
)

var (
	out         = scribe.New("Chip8")
	close       = closer.New()
	ctx, cancel = context.WithCancel(context.Background())
)

func main() {
	var screenMultiplier float64
	flag.Float64Var(&screenMultiplier, "screenMultiplier", 8, "How many true pixels represent each single Chip8 pixel.")
	flag.Parse()

	c := New(screenMultiplier)
	go func() {
		err := close.Wait()
		c.cancel()
		c.errC <- err
	}()

	pixelgl.Run(c.run)
	err := <-c.errC
	exit(err)
}

func onError(err error, fn func(error) bool) {
	if err == nil {
		return
	}

	fn(err)
}

func exit(err error) {
	switch err {
	case nil, errors.ErrIsClosed:
		os.Exit(0)

	default:
		out.Errorf("error encountered: %v", err)
		os.Exit(1)
	}
}
