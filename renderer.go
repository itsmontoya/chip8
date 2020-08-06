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
	GetKeypad() Keypad
}

func newPixel(screenMultiplier float64) (pp *PixelRenderer, err error) {
	var p PixelRenderer
	cfg := makeConfig("Pixel Rocks!", screenMultiplier)
	if p.win, err = pixelgl.NewWindow(cfg); err != nil {
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
	imd *imdraw.IMDraw
	g   graphics

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

	// Bottom left corner
	p.imd.Push(pixel.V(x+0, y+0))
	// Bottom right corner
	p.imd.Push(pixel.V(x+p.screenMultiplier, y+0))
	// Top right corner
	p.imd.Push(pixel.V(x+p.screenMultiplier, y+p.screenMultiplier))
	// Top left corner
	p.imd.Push(pixel.V(x+0, y+p.screenMultiplier))

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

func (p *PixelRenderer) getKeyValue(key pixelgl.Button) byte {
	val := p.win.Pressed(key)
	return boolToByte(val)
}

// Draw will draw to the screen
func (p *PixelRenderer) Draw(g graphics) (err error) {
	// Draw the pixels which changed in the new graphics state
	g.forEachDelta(p.g, p.drawPixel)

	if p.win.Closed() {
		// Window has been closed, return
		return errors.ErrIsClosed
	}

	// Update window (swap buffers)
	p.win.Update()
	return
}

// GetKeypad will get the current keypad
func (p *PixelRenderer) GetKeypad() (k Keypad) {
	// 1234
	k[0] = p.getKeyValue(pixelgl.Key1)
	k[1] = p.getKeyValue(pixelgl.Key2)
	k[2] = p.getKeyValue(pixelgl.Key3)
	k[3] = p.getKeyValue(pixelgl.Key4)

	// QUER
	k[4] = p.getKeyValue(pixelgl.KeyQ)
	k[5] = p.getKeyValue(pixelgl.KeyW)
	k[6] = p.getKeyValue(pixelgl.KeyE)
	k[7] = p.getKeyValue(pixelgl.KeyR)

	// ASDF
	k[8] = p.getKeyValue(pixelgl.KeyA)
	k[9] = p.getKeyValue(pixelgl.KeyS)
	k[10] = p.getKeyValue(pixelgl.KeyD)
	k[11] = p.getKeyValue(pixelgl.KeyF)

	// ZXCV
	k[12] = p.getKeyValue(pixelgl.KeyZ)
	k[13] = p.getKeyValue(pixelgl.KeyX)
	k[14] = p.getKeyValue(pixelgl.KeyC)
	k[15] = p.getKeyValue(pixelgl.KeyV)
	return
}
