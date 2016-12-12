package canvas

import (
	"fmt"
	"image/color"
	"math"

	"honnef.co/go/js/dom"
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

type canvas struct {
	ctx *dom.CanvasRenderingContext2D
}

type Canvas interface {
	Line(start, end, startDir, endDir Vector, startWidth, endWidth float64, startColor, endColor color.RGBA)
}

func New(ctx *dom.CanvasRenderingContext2D) Canvas {
	return canvas{ctx}
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
	ctx := c.ctx
	ctx.BeginPath()
	ctx.MoveTo(int(startLeft.X), int(startLeft.Y))
	ctx.QuadraticCurveTo(int(leftC.X), int(leftC.Y), int(endLeft.X), int(endLeft.Y))
	ctx.LineTo(int(endRight.X), int(endRight.Y))
	ctx.QuadraticCurveTo(int(rightC.X), int(rightC.Y), int(startRight.X), int(startRight.Y))
	ctx.ClosePath()
	ctx.StrokeStyle = "none"
	ctx.FillStyle = fmt.Sprintf("rgb(%v,%v,%v)", endColor.R, endColor.G, endColor.B)
	ctx.Fill()
}
