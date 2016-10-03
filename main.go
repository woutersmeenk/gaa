package main

import (
	//"fmt"
	"gaa/network"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
	"math/rand"
)

const (
	NetworkInputs  = 3
	NetworkOutputs = 8
	HiddenNeurons  = 10
	ImageHeight    = 500
	ImageWidth     = 500
	Steps          = 5000
)

type networkInput struct {
	x, y float64
	bias float64
}

func (input *networkInput) TransInputs() []float64 {
	return []float64{
		(float64(input.x)/ImageHeight*2 - 1),
		(float64(input.y)/ImageWidth*2 - 1),
		input.bias,
	}
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

func performAction(output *networkOutput, state *simState, gc *draw2dimg.GraphicContext) {
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
	if state.x < 0 {
		state.x = ImageWidth + state.x
	}
	if state.y < 0 {
		state.y = ImageHeight + state.y
	}
	if state.y > ImageHeight {
		state.y = state.y - ImageHeight
	}

	if state.x > ImageHeight {
		state.x = state.x - ImageHeight
	}
	gc.SetStrokeColor(output.color)
	gc.SetLineWidth(output.width)
	gc.MoveTo(oldX, oldY)
	gc.LineTo(state.x, state.y)
	gc.Stroke()
}

func main() {
	dest := image.NewRGBA(image.Rect(0, 0, ImageWidth, ImageHeight))
	gc := draw2dimg.NewGraphicContext(dest)
	r := rand.New(rand.NewSource(5135264217))
	net := network.New(NetworkInputs, NetworkOutputs, HiddenNeurons, r)
	input := &networkInput{x: 0, y: 0, bias: 1}
	state := &simState{x: 0, y: 0, dx: 0, dy: 0}
	for t := 0; t < Steps; t++ {
		var output = &networkOutput{}
		net.Eval(input, output)
		performAction(output, state, gc)
		//fmt.Printf("state: %v output: %v\n", *state, *output)
		input.x, input.y = state.x, state.y
	}
	fn := "output.png"
	err := draw2dimg.SaveToPngFile(fn, dest)
	if err != nil {
		panic(err)
	}

}
