package main

import (
	"context"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func isDone(ctx context.Context) (done bool) {
	select {
	case <-ctx.Done():
		return true

	default:
		return false
	}
}

func makeConfig(title string, screenMulitplier float64) (cfg pixelgl.WindowConfig) {
	cfg.Title = title
	cfg.Bounds = pixel.R(0, 0, 64*screenMulitplier, 32*screenMulitplier)
	cfg.VSync = true
	return
}

func getXY(i int) (x, y float64) {
	row := i / 64
	x = float64(i - (row * 64))
	y = float64(row)
	return
}

func boolToByte(val bool) byte {
	if val {
		return 1
	}

	return 0
}
