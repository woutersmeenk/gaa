package sim

import (
	"fmt"
	"gaa/network"
	"gaa/svg"
	"image/color"
	"io"
	"math"
)

const (
	imageHeight = 500
	imageWidth  = 500
	steps       = 1000
)

type networkInput struct {
	x, y float64
	bias float64
}

func (input *networkInput) TransInputs() []float64 {
	result := []float64{
		(input.x/imageHeight*2 - 1),
		(input.y/imageWidth*2 - 1),
		input.bias,
	}
	return result
}

type networkOutput struct {
	dx, dy       float64    // vector representing direction of velocity or acceleration in pixels per timestep (squared)
	acceleration bool       // When true above vector is interperted as acceleration otherwise as velocity
	color        color.RGBA // Color of line
	width        float64    // width of line in pixels
}

func (output *networkOutput) TransOutputs(outputSlice []float64) {
	output.dx = float64(outputSlice[0])
	output.dy = float64(outputSlice[1])
	output.acceleration = outputSlice[2] > 0
	r := uint8(((outputSlice[3] + 1) / 2) * 256)
	g := uint8(((outputSlice[4] + 1) / 2) * 256)
	b := uint8(((outputSlice[5] + 1) / 2) * 256)
	a := uint8(((outputSlice[6] + 1) / 2) * 256)
	output.color = color.RGBA{r, g, b, a}
	output.width = ((outputSlice[7] + 1) / 2) * 10
}

type simState struct {
	x, y   float64
	dx, dy float64
}

func performAction(output *networkOutput, state *simState, svg svg.SVG) {
	oldX, oldY := state.x, state.y
	if output.acceleration {
		state.dx += output.dx
		state.dy += output.dy
	} else {
		state.dx = output.dx
		state.dy = output.dy
	}
	state.x += state.dx
	state.y += state.dy
	x1, y1 := modPos(oldX, imageWidth), modPos(oldY, imageHeight)
	x2, y2 := modPos(state.x, imageWidth), modPos(state.y, imageHeight)
	width := output.width
	r, g, b, a := output.color.R, output.color.G, output.color.B, output.color.A
	svg.Line(x1, y1, x2, y2, width, r, g, b, a)
}

func modPos(n float64, d float64) float64 {
	result := math.Mod(n, d)
	if result < 0 {
		result += d
	}
	return result
}

func Simulate(net *network.Network, writer io.Writer) {
	svg := svg.New(writer, imageWidth, imageHeight)
	initialX := imageWidth / 2.0
	initialY := imageHeight / 2.0
	state := &simState{x: initialX, y: initialY, dx: 0, dy: 0}
	input := &networkInput{x: initialX, y: initialY, bias: 1}
	for t := 0; t < steps; t++ {
		var output = &networkOutput{}
		net.Activate(input, output)
		performAction(output, state, svg)
		input.x, input.y = state.x, state.y
	}
	svg.Close()
}
