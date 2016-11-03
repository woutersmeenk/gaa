package sim

import (
	"gaa/network"
	"gaa/svg"
	"image/color"
	"math"
)

const (
	NetworkInputs  = 3
	NetworkOutputs = 8
	imageHeight    = 500
	imageWidth     = 500
	steps          = 1000
)

type simState struct {
	x, y   float64
	bias   float64 // always 1
	dx, dy float64
}

func (state *simState) TransInputs() []float64 {
	x, y := wrappedLoc(state.x, state.y)
	result := []float64{
		(x/imageWidth*2 - 1),
		(y/imageHeight*2 - 1),
		state.bias,
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
	output.dx = outputSlice[0] * 5
	output.dy = outputSlice[1] * 5
	output.acceleration = outputSlice[2] > 0
	r := uint8(((outputSlice[3] + 1) / 2) * 255)
	g := uint8(((outputSlice[4] + 1) / 2) * 255)
	b := uint8(((outputSlice[5] + 1) / 2) * 255)
	a := uint8(((outputSlice[6] + 1) / 2) * 255)
	output.color = color.RGBA{r, g, b, a}
	output.width = ((outputSlice[7] + 1) / 2) * 10
}

func performAction(output *networkOutput, state *simState, svg svg.SVG) {
	oldX, oldY := state.x, state.y
	if output.acceleration {
		state.dx += output.dx
		state.dy += output.dy
		if state.dy > 5 {
			state.dy = 5
		}
		if state.dy < -5 {
			state.dy = -5
		}
		if state.dx > 5 {
			state.dx = 5
		}
		if state.dx < -5 {
			state.dx = -5
		}
	} else {
		state.dx = output.dx
		state.dy = output.dy
	}
	state.x += state.dx
	state.y += state.dy
	x1, y1 := wrappedLoc(oldX, oldY)
	x2, y2 := wrappedLoc(state.x, state.y)
	width := output.width
	r, g, b, a := output.color.R, output.color.G, output.color.B, output.color.A
	svg.Line(x1, y1, x2, y2, width, r, g, b, a)

}

func wrappedLoc(x, y float64) (float64, float64) {
	return modPos(x, imageWidth), modPos(y, imageHeight)
}

func modPos(n float64, d float64) float64 {
	result := math.Mod(n, d)
	if result < 0 {
		result += d
	}
	return result
}

func Simulate(net network.Network, s svg.SVG) {
	s.Open(imageWidth, imageHeight)
	defer s.Close()
	initialX := imageWidth / 2.0
	initialY := imageHeight / 2.0
	state := &simState{x: initialX, y: initialY, bias: 1, dx: 0, dy: 0}
	for t := 0; t < steps; t++ {
		var output = &networkOutput{}
		net.Activate(state, output)
		performAction(output, state, s)
	}
}
