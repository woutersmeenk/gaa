package sim

import (
	"image/color"
	"reflect"
	"testing"

	"github.com/woutersmeenk/gaa/canvas"
)

var modPosTests = []struct {
	n, d, res float64
}{
	{-2, 5, 3},
	{-1, 5, 4},
	{0, 5, 0},
	{4, 5, 4},
	{5, 5, 0},
	{6, 5, 1},
}

func TestModPos(t *testing.T) {
	for _, tt := range modPosTests {
		res := modPos(tt.n, tt.d)
		if res != tt.res {
			t.Errorf("modPos(%v, %v) = %v want %v", tt.n, tt.d, res, tt.res)
		}
	}
}

var transInputTests = []struct {
	state simState
	res   []float64
}{
	{simState{canvas.Vector{-1, -1}, canvas.Vector{}, 0, color.RGBA{}}, []float64{0.996, 0.996, 1}},
	{simState{canvas.Vector{0, 0}, canvas.Vector{}, 0, color.RGBA{}}, []float64{-1, -1, 1}},
	{simState{canvas.Vector{ImageWidth - 1, ImageHeight - 1}, canvas.Vector{}, 0, color.RGBA{}}, []float64{0.996, 0.996, 1}},
	{simState{canvas.Vector{ImageWidth / 2, ImageHeight / 2}, canvas.Vector{}, 0, color.RGBA{}}, []float64{0, 0, 1}},
	{simState{canvas.Vector{ImageWidth, ImageHeight}, canvas.Vector{}, 0, color.RGBA{}}, []float64{-1, -1, 1}},
}

func TestTransInputs(t *testing.T) {
	for _, tt := range transInputTests {
		res := tt.state.TransInputs()
		if !reflect.DeepEqual(res, tt.res) {
			t.Errorf("TransInputs(%v) = %v want %v", tt.state, res, tt.res)
		}
	}
}

var transOutputTests = []struct {
	output []float64
	res    networkOutput
}{
	{[]float64{-1, -1, -1, -1, -1, -1, -1, -1}, networkOutput{canvas.Vector{-velocityLimit, -velocityLimit}, false, color.RGBA{0, 0, 0, 0}, 0}},
	{[]float64{1, 1, 1, 1, 1, 1, 1, 1}, networkOutput{canvas.Vector{velocityLimit, velocityLimit}, true, color.RGBA{255, 255, 255, 255}, maxWidth}},
	{[]float64{0, 0, 0, 0, 0, 0, 0, 0}, networkOutput{canvas.Vector{0, 0}, false, color.RGBA{127, 127, 127, 127}, maxWidth / 2}},
}

func TestTransOutput(t *testing.T) {
	for _, tt := range transOutputTests {
		res := networkOutput{}
		res.TransOutputs(tt.output)
		if !reflect.DeepEqual(res, tt.res) {
			t.Errorf("TransOutputs(%v) = %v want %v", tt.output, res, tt.res)
		}
	}
}
