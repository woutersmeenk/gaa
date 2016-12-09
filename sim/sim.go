package sim

import (
	"gaa/canvas"
	"gaa/network"
	"image/color"
	"math"
	"time"
)

const (
	NetworkInputs  = 3
	NetworkOutputs = 8
	imageHeight    = 500
	imageWidth     = 500
	steps          = 1000
	velocityLimit  = 5
	maxWidth       = 30
)

type simState struct {
	pos, dir canvas.Vector
	width    float64
	color    color.RGBA
}

func (state *simState) TransInputs() []float64 {
	pos := wrappedLoc(state.pos)
	result := []float64{
		(pos.X/imageWidth*2 - 1),
		(pos.Y/imageHeight*2 - 1),
		1,
	}
	return result
}

type networkOutput struct {
	dir          canvas.Vector // vector representing direction of velocity or acceleration in pixels per timestep (squared)
	acceleration bool          // When true above vector is interperted as acceleration otherwise as velocity
	color        color.RGBA    // Color of line
	width        float64       // width of line in pixels
}

func (output *networkOutput) TransOutputs(outputSlice []float64) {
	output.dir = canvas.Vector{}
	output.dir.X = outputSlice[0] * velocityLimit
	output.dir.Y = outputSlice[1] * velocityLimit
	output.acceleration = outputSlice[2] > 0
	r := uint8(((outputSlice[3] + 1) / 2) * 255)
	g := uint8(((outputSlice[4] + 1) / 2) * 255)
	b := uint8(((outputSlice[5] + 1) / 2) * 255)
	a := uint8(((outputSlice[6] + 1) / 2) * 255)
	output.color = color.RGBA{r, g, b, a}
	output.width = ((outputSlice[7] + 1) / 2) * maxWidth
}

func performAction(output *networkOutput, state *simState, c canvas.Canvas) {
	oldPos := state.pos
	oldDir := state.dir
	var acc canvas.Vector
	if output.acceleration {
		acc = output.dir
		state.dir = state.dir.Add(acc)
	} else {
		acc = output.dir.Sub(state.dir)
		state.dir = output.dir
	}
	state.dir.X = applyLimit(state.dir.X, velocityLimit)
	state.dir.Y = applyLimit(state.dir.Y, velocityLimit)

	// calculate new position y = y0 + v0*t + 1/2*a*t^2
	state.pos = state.pos.Add(state.dir).Add(acc.Mul(0.5))
	posDelta := state.pos.Sub(oldPos)

	from := wrappedLoc(oldPos)
	to := wrappedLoc(state.pos)
	// We draw the line in two pieces (if needed) to prevent drawing accros the image border
	// forward from the current location and backward from the next location
	c.Line(from, from.Add(posDelta), oldDir, state.dir, state.width, output.width, state.color, output.color)
	if from.Add(posDelta) != state.pos {
		c.Line(to.Sub(posDelta), to, oldDir, state.dir, state.width, output.width, state.color, output.color)
	}
	state.width = output.width
	state.color = output.color
}

func applyLimit(val, limit float64) float64 {
	if val > limit {
		return limit
	}
	if val < -limit {
		return limit
	}
	return val
}

func wrappedLoc(pos canvas.Vector) canvas.Vector {
	return canvas.Vector{modPos(pos.X, imageWidth), modPos(pos.Y, imageHeight)}
}

func modPos(n float64, d float64) float64 {
	result := math.Mod(n, d)
	if result < 0 {
		result += d
	}
	return result
}

func Simulate(net network.Network, c canvas.Canvas) {
	c.Open(imageWidth, imageHeight)
	defer c.Close()
	initialPos := canvas.Vector{imageWidth / 2.0, imageHeight / 2.0}
	state := &simState{pos: initialPos,
		dir:   canvas.Vector{1, 1},
		width: 10,
		color: color.RGBA{0, 0, 0, 0}}
	for t := 0; t < steps; t++ {
		if t%10 == 0 {
			time.Sleep(0 * time.Second)
		}
		var output = &networkOutput{}
		net.Activate(state, output)
		performAction(output, state, c)
	}
}
