package svg

import (
	"fmt"
	"io"
)

type svg struct {
	writer io.Writer
}

type SVG interface {
	Open(width, height int)
	Line(x1, y1, x2, y2, width float64, r, g, b, a uint8)
	Close()
}

func New(writer io.Writer) SVG {
	return svg{writer}
}

func (s svg) Open(width, height int) {
	fmt.Fprintf(s.writer, "<svg width=\"%v\" height=\"%v\" xmlns=\"http://www.w3.org/2000/svg\">\n", width, height)
}

func (s svg) Line(x1, y1, x2, y2, width float64, r, g, b, a uint8) {
	fmt.Fprintf(s.writer, "<line x1=\"%v\" y1=\"%v\" x2=\"%v\" y2=\"%v\" stroke-width=\"%v\" stroke=\"rgb(%v,%v,%v)\"/>\n",
		x1, y1, x2, y2, width, r, g, b)
}

func (s svg) Close() {
	fmt.Fprint(s.writer, "</svg>\n")
}
