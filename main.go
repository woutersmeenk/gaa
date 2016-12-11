package main

import (
	"gaa/canvas"
	"gaa/network"
	"gaa/sim"
	"math/rand"
	"time"

	"honnef.co/go/js/dom"
)

const (
	hiddenNeurons = 10
)

func main() {
	seed := time.Now().UnixNano()
	doc := dom.GetWindow().Document()
	body := doc.GetElementByID("body")
	for i := 0; i < 16; i++ {
		ticker := time.NewTicker(30 * time.Millisecond)
		extraSeed := int64(i)
		r := rand.New(rand.NewSource(seed + extraSeed))
		net := network.New(sim.NetworkInputs, sim.NetworkOutputs, hiddenNeurons, r)
		cv := canvas.New(doc, sim.ImageHeight, sim.ImageWidth)
		htmlCv := cv.GetHTMLCanvas()
		htmlCv.SetAttribute("style", "width: 300; height: 300; border: 1px solid #dcdcdc; margin: 2px")
		body.AppendChild(htmlCv)
		go func() {
			sim.Simulate(net, cv, ticker.C)
		}()
	}
}
