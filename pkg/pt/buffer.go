package pt

import (
	"image"
)

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

type PixelBuffer struct {
	width  int
	height int
	buff   []Pixel
}

func NewPxlBufferAR(height int, aspect float64) *PixelBuffer {
	width := float64(height) * aspect
	return NewPxlBuffer(int(width), height)
}

func NewPxlBuffer(width, height int) *PixelBuffer {
	return &PixelBuffer{
		width:  width,
		height: height,
		buff:   make([]Pixel, width*height),
	}
}

func (b *PixelBuffer) addSample(x, y int, c Color) {
	b.buff[y*b.width+x].addSample(c)
}

func (b *PixelBuffer) h() int {
	return b.height
}

func (b *PixelBuffer) w() int {
	return b.width
}

func (b *PixelBuffer) ToImage() image.Image {
	topLeft := image.Point{0, 0}
	bottomRight := image.Point{b.width, b.height}
	img := image.NewRGBA(image.Rectangle{topLeft, bottomRight})
	for i, px := range b.buff {
		x := i % b.width
		y := b.height - (i / b.width)
		img.Set(x, y, px.color.goColor())
	}
	return img
}

type FrameBuffer struct {
	width  int
	height int
	buff   []Color
}

func NewFrameBufferAR(height int, aspect float64) *FrameBuffer {
	width := float64(height) * aspect
	return NewFrameBuffer(int(width), height)
}

func NewFrameBuffer(width, height int) *FrameBuffer {
	return &FrameBuffer{
		width:  width,
		height: height,
		buff:   make([]Color, width*height),
	}
}

func (b *FrameBuffer) addSample(x, y int, c Color) {
	b.buff[y*b.width+x] = c
}

func (b *FrameBuffer) w() int {
	return b.width
}

func (b *FrameBuffer) h() int {
	return b.height
}

func (b *FrameBuffer) ToImage() image.Image {
	topLeft := image.Point{0, 0}
	bottomRight := image.Point{b.width, b.height}
	img := image.NewRGBA(image.Rectangle{topLeft, bottomRight})
	for i, color := range b.buff {
		x := i % b.width
		y := b.height - (i / b.width)
		img.Set(x, y, color.goColor())
	}
	return img
}
