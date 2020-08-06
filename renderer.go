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
	p.color = colornames.Skyblue
	pp = &p
	return
}

// PixelRenderer is a renderer for the Pixel library
type PixelRenderer struct {
	color color.RGBA
	win   *pixelgl.Window
	g     graphics

	screenMultiplier float64
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
	imd.Color = pixel.RGB(255, 255, 255)
	p.drawSquare(imd, 0, 20)
	p.drawSquare(imd, 2, 21)
	p.drawSquare(imd, 3, 22)

	if p.win.Closed() {
		return errors.ErrIsClosed
	}

	p.win.Update()
	return
}
