package canvas

import (
	"fmt"
	"io"
)

type canvas struct {
	writer io.Writer
}

type Canvas interface {
	Open(width, height int)
	Line(x1, y1, x2, y2, width float64, r, g, b, a uint8)
	Close()
}

func New(writer io.Writer) Canvas {
	return canvas{writer}
}

func (s canvas) Open(width, height int) {
	fmt.Fprintf(s.writer, "<svg width=\"%v\" height=\"%v\" xmlns=\"http://www.w3.org/2000/svg\">\n", width, height)
}

func (s canvas) Line(x1, y1, x2, y2, width float64, r, g, b, a uint8) {
	fmt.Fprintf(s.writer, "<line x1=\"%v\" y1=\"%v\" x2=\"%v\" y2=\"%v\" stroke-width=\"%v\" stroke=\"rgb(%v,%v,%v)\"/>\n",
		x1, y1, x2, y2, width, r, g, b)
}

func (s canvas) Close() {
	fmt.Fprint(s.writer, "</svg>\n")
}
