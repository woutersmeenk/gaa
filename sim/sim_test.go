package sim

import (
	"gaa/network"
	"gaa/svg"
	"image/color"
	"io/ioutil"
	"math/rand"
	"reflect"
	"testing"
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
	{simState{-1, -1, 0, 0}, []float64{1, 1, 1}},
	{simState{0, 0, 0, 0}, []float64{-1, -1, 1}},
	{simState{imageWidth - 1, imageHeight - 1, 0, 0}, []float64{1, 1, 1}},
	{simState{imageWidth / 2, imageHeight / 2, 0, 0}, []float64{0, 0, 1}},
	{simState{imageWidth, imageHeight, 0, 0}, []float64{-1, -1, 1}},
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
	{[]float64{-1, -1, -1, -1, -1, -1, -1, -1}, networkOutput{-velocityLimit, -velocityLimit, false, color.RGBA{0, 0, 0, 0}, 0}},
	{[]float64{1, 1, 1, 1, 1, 1, 1, 1}, networkOutput{velocityLimit, velocityLimit, true, color.RGBA{255, 255, 255, 255}, maxWidth}},
	{[]float64{0, 0, 0, 0, 0, 0, 0, 0}, networkOutput{0, 0, false, color.RGBA{127, 127, 127, 127}, maxWidth / 2}},
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

func BenchmarkSimulate(b *testing.B) {
	r := rand.New(rand.NewSource(328932186))
	for n := 0; n < b.N; n++ {
		net := network.New(3, 8, 10, r)
		Simulate(net, svg.New(ioutil.Discard))
	}
}
