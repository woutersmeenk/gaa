package main

import (
	"bytes"
	"fmt"
	"gaa/canvas"
	"gaa/network"
	"gaa/sim"
	"math/rand"

	"honnef.co/go/js/dom"
)

const (
	hiddenNeurons = 10
	seed          = 64877563487
)

func main() {
	doc := dom.GetWindow().Document()
	body := doc.GetElementByID("body")
	f := doc.GetElementByID("foo")
	println("test")
	println(fmt.Sprintf("f: %v, body: %v", f, body))
	f.AddEventListener("click", false, func(event dom.Event) {
		go func() {
			f.SetInnerHTML("Generating...")
			for i := 0; i < 20; i++ {
				extraSeed := int64(i)
				r := rand.New(rand.NewSource(seed + extraSeed))
				net := network.New(sim.NetworkInputs, sim.NetworkOutputs, hiddenNeurons, r)
				var writer bytes.Buffer
				sim.Simulate(net, canvas.New(&writer))
				div := doc.CreateElement("div")
				div.SetInnerHTML(writer.String())
				div.SetAttribute("style", "width: 100; height: 100")
				body.AppendChild(div)
			}
		}()
	})
}
