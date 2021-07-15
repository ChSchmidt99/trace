package pt

// TODO: Make Buffer interface
type Buffer struct {
	Width  int
	Height int
	Buff   []Color
}

func NewBufferAspect(height int, aspect float64) *Buffer {
	width := float64(height) * aspect
	return NewBuffer(int(width), height)
}

func NewBuffer(width, height int) *Buffer {
	return &Buffer{
		Width:  width,
		Height: height,
		Buff:   make([]Color, width*height),
	}
}

func (b *Buffer) setPixel(x, y int, c Color) {
	b.Buff[y*b.Width+x] = c
}
