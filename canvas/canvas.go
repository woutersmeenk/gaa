package canvas

import (
	"fmt"
	"image/color"
	"io"
	"math"
)

type Vector struct {
	X, Y float64
}

func (v Vector) Add(v2 Vector) Vector {
	return Vector{v.X + v2.X, v.Y + v2.Y}
}

func (v Vector) Sub(v2 Vector) Vector {
	return Vector{v.X - v2.X, v.Y - v2.Y}
}

func (v Vector) Mul(scalar float64) Vector {
	return Vector{v.X * scalar, v.Y * scalar}
}

func (v Vector) Mag() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vector) Norm() Vector {
	mag := v.Mag()
	return Vector{v.X / mag, v.Y / mag}
}

func (v Vector) PerpendicularClock() Vector {
	return Vector{-v.Y, v.X}
}

func (v Vector) PerpendicularCounterClock() Vector {
	return Vector{v.Y, -v.X}
}

func (v Vector) AsPath() string {
	return fmt.Sprintf("%v %v", v.X, v.Y)
}

type canvas struct {
	writer io.Writer
}

type Canvas interface {
	Open(width, height int)
	Line(start, end, startDir, endDir Vector, startWidth, endWidth float64, startColor, endColor color.RGBA)
	Close()
}

func New(writer io.Writer) Canvas {
	return canvas{writer}
}

func (c canvas) Open(width, height int) {
	fmt.Fprintf(c.writer, "<svg width=\"%v\" height=\"%v\" xmlns=\"http://www.w3.org/2000/svg\">\n", width, height)
}

func (c canvas) Line(start, end, startDir, endDir Vector, startWidth, endWidth float64, startColor, endColor color.RGBA) {
	// Calculate the shape of the line. We calculate the four points of the trapazoid that describes it

	startDirNorm := startDir.Norm()
	startLeft := start.Add(startDirNorm.PerpendicularCounterClock().Mul(startWidth / 2))
	startRight := start.Add(startDirNorm.PerpendicularClock().Mul(startWidth / 2))

	leftC := startLeft.Add(startDir.Mul(0.5))

	endDirNorm := endDir.Norm()
	endLeft := end.Add(endDirNorm.PerpendicularCounterClock().Mul(endWidth / 2))
	endRight := end.Add(endDirNorm.PerpendicularClock().Mul(endWidth / 2))

	rightC := endRight.Sub(startDir.Mul(0.5))

	fmt.Fprintf(c.writer, `
<!-- %v %v %v %v %v %v -->
<path d="M %v 
		 Q %v %v
		 L %v
		 Q %v %v
		 Z" 
	  stroke="none" fill="rgb(%v,%v,%v)"/>`,
		startWidth, endWidth, start, end, startDirNorm, endDirNorm,
		startLeft.AsPath(),
		leftC.AsPath(), endLeft.AsPath(),
		endRight.AsPath(),
		rightC.AsPath(), startRight.AsPath(),
		endColor.R, endColor.G, endColor.B)
}

func (c canvas) Close() {
	fmt.Fprint(c.writer, "</svg>\n")
}
