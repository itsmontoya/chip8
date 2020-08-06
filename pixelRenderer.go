package main

import (
	"image/color"

	"github.com/Hatch1fy/errors"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/itsmontoya/chip8/vm"
	"golang.org/x/image/colornames"
)

func newPixel(screenMultiplier float64) (pp *PixelRenderer, err error) {
	var p PixelRenderer
	p.cfg = makeConfig("Pixel Rocks!", screenMultiplier)

	// Initialize a new Pixel window
	if p.win, err = pixelgl.NewWindow(p.cfg); err != nil {
		return
	}

	p.imd = imdraw.New(nil)
	p.screenMultiplier = screenMultiplier
	p.clearColor = colornames.Skyblue
	p.offColor = color.RGBA{255, 255, 255, 0}
	p.onColor = color.RGBA{255, 255, 255, 255}

	// Set reference to PixelRenderer
	pp = &p
	return
}

// PixelRenderer is a renderer for the Pixel library
type PixelRenderer struct {
	win *pixelgl.Window
	cfg pixelgl.WindowConfig
	imd *imdraw.IMDraw
	g   vm.Graphics

	screenMultiplier float64

	clearColor color.RGBA
	offColor   color.RGBA
	onColor    color.RGBA
}

func (p *PixelRenderer) drawPixel(i int, val byte) {
	p.setColor(val)
	x, y := getXY(i)
	p.drawSquare(x, y)
	p.g[i] = val
}

func (p *PixelRenderer) drawSquare(x, y float64) {
	// Multiply X value by screen multiplier
	x *= p.screenMultiplier
	// Multiply Y value by screen multiplier
	y *= p.screenMultiplier
	// Inverse Y
	y = p.cfg.Bounds.H() - y

	// Bottom left corner
	p.imd.Push(pixel.V(x+0, y+0))
	// Bottom right corner
	p.imd.Push(pixel.V(x+p.screenMultiplier, y+0))
	// Top right corner
	p.imd.Push(pixel.V(x+p.screenMultiplier, y-p.screenMultiplier))
	// Top left corner
	p.imd.Push(pixel.V(x+0, y-p.screenMultiplier))

	// Complete shape
	p.imd.Polygon(0)

	// Draw shape to window buffers
	p.imd.Draw(p.win)
}

func (p *PixelRenderer) setColor(val byte) {
	if val == 0 {
		// Value is unset, use "off" color
		p.imd.Color = p.offColor
		return
	}

	// Value is set, use "on" color
	p.imd.Color = p.onColor
}

// Draw will draw to the screen
func (p *PixelRenderer) Draw(g vm.Graphics) (err error) {
	// Draw the pixels which changed in the new Graphics state
	g.ForEachDelta(p.g, p.drawPixel)

	if p.win.Closed() {
		// Window has been closed, return
		return errors.ErrIsClosed
	}

	// Update window (swap buffers)
	p.win.Update()
	return
}

// GetKeypad will get the current keypad
func (p *PixelRenderer) GetKeypad() (k vm.Keypad) {
	// 1234
	k.Set(0, p.win.Pressed(pixelgl.Key1))
	k.Set(1, p.win.Pressed(pixelgl.Key2))
	k.Set(2, p.win.Pressed(pixelgl.Key3))
	k.Set(3, p.win.Pressed(pixelgl.Key4))

	// QUER
	k.Set(4, p.win.Pressed(pixelgl.KeyQ))
	k.Set(5, p.win.Pressed(pixelgl.KeyW))
	k.Set(6, p.win.Pressed(pixelgl.KeyE))
	k.Set(7, p.win.Pressed(pixelgl.KeyR))

	// ASDF
	k.Set(8, p.win.Pressed(pixelgl.KeyA))
	k.Set(9, p.win.Pressed(pixelgl.KeyS))
	k.Set(10, p.win.Pressed(pixelgl.KeyD))
	k.Set(11, p.win.Pressed(pixelgl.KeyF))

	// ZXCV
	k.Set(12, p.win.Pressed(pixelgl.KeyZ))
	k.Set(13, p.win.Pressed(pixelgl.KeyX))
	k.Set(14, p.win.Pressed(pixelgl.KeyC))
	k.Set(15, p.win.Pressed(pixelgl.KeyV))
	return
}
