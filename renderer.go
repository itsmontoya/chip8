package main

import (
	"image/color"

	"github.com/Hatch1fy/errors"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// Renderer will render graphics output
type Renderer interface {
	Draw(graphics) error
}

func newPixel(screenMultiplier float64) (pp *PixelRenderer, err error) {
	var p PixelRenderer
	cfg := makeConfig("Pixel Rocks!", screenMultiplier)
	if p.win, err = pixelgl.NewWindow(cfg); err != nil {
		return
	}

	p.screenMultiplier = screenMultiplier
	p.clearColor = colornames.Skyblue
	p.offColor = color.RGBA{255, 255, 255, 0}
	p.onColor = color.RGBA{255, 255, 255, 255}

	// Set some debug squares for the renderer
	p.g[0] = 1
	p.g[63] = 1
	p.g[len(p.g)-64] = 1
	p.g[len(p.g)-1] = 1

	// Set reference to PixelRenderer
	pp = &p
	return
}

// PixelRenderer is a renderer for the Pixel library
type PixelRenderer struct {
	win *pixelgl.Window
	g   graphics

	screenMultiplier float64

	clearColor color.RGBA
	offColor   color.RGBA
	onColor    color.RGBA
}

func (p *PixelRenderer) drawSquare(imd *imdraw.IMDraw, x, y float64) {
	x *= p.screenMultiplier
	y *= p.screenMultiplier
	imd.Push(pixel.V(x+0, y+0))
	imd.Push(pixel.V(x+p.screenMultiplier, y+0))
	imd.Push(pixel.V(x+p.screenMultiplier, y+p.screenMultiplier))
	imd.Push(pixel.V(x+0, y+p.screenMultiplier))
	imd.Polygon(0)
	imd.Draw(p.win)
}

// Draw will draw to the screen
func (p *PixelRenderer) Draw(g graphics) (err error) {
	imd := imdraw.New(nil)
	for i, val := range p.g {
		if val == 0 {
			imd.Color = p.offColor
		} else {
			imd.Color = p.onColor
		}

		row := i / 64
		cell := float64(i - (row * 64))
		p.drawSquare(imd, cell, float64(row))
	}

	if p.win.Closed() {
		return errors.ErrIsClosed
	}

	p.win.Update()
	return
}
