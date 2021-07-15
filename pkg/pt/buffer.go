package pt

import (
	"image/color"
	"math"
)

// TODO: Check if struct would be faster than interface
type Buffer interface {
	addSample(x, y int, c Color)
	w() int
	h() int
}

type Pixel struct {
	samples int
	color   Color
}

func (px *Pixel) addSample(c Color) {
	px.samples++
	if px.samples == 1 {
		px.color = c
		return
	}
	px.color = px.color.Add(c.Sub(px.color).Div(float64(px.samples)))
}

// Convert to image color and also Gamma correct for gamma=2.0
func (px Pixel) GoColor() color.Color {
	r := math.Sqrt(px.color.X)
	g := math.Sqrt(px.color.Y)
	b := math.Sqrt(px.color.Z)
	return color.RGBA{R: uint8(Clamp(r, 0.0, 1.0) * 255), G: uint8(Clamp(g, 0.0, 1.0) * 255), B: uint8(Clamp(b, 0.0, 1.0) * 255), A: 255}
}

type PixelBuffer struct {
	Width  int
	Height int
	Buff   []Pixel
}

func NewBufferAspect(height int, aspect float64) *PixelBuffer {
	width := float64(height) * aspect
	return NewBuffer(int(width), height)
}

func NewBuffer(width, height int) *PixelBuffer {
	return &PixelBuffer{
		Width:  width,
		Height: height,
		Buff:   make([]Pixel, width*height),
	}
}

func (b *PixelBuffer) addSample(x, y int, c Color) {
	b.Buff[y*b.Width+x].addSample(c)
}

func (b *PixelBuffer) h() int {
	return b.Height
}

func (b *PixelBuffer) w() int {
	return b.Width
}
