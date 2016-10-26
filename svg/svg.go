package svg

import (
	"fmt"
	"io"
)

type SVG struct {
	writer io.Writer
}

func New(writer io.Writer, width, height int) SVG {
	fmt.Fprintf(writer, "<svg width=\"%v\" height=\"%v\">\n", width, height)
	return SVG{writer}
}

func (svg SVG) Line(x1, y1, x2, y2, width float64, r, g, b, a uint8) {
	fmt.Fprintf(svg.writer, "<line x1=\"%v\" y1=\"%v\" x2=\"%v\" y2=\"%v\" stroke-width=\"%v\"/>\n",
		x1, y1, x2, y2, width)
}

func (svg SVG) Close() {
	fmt.Fprint(svg.writer, "</svg>\n")
}
