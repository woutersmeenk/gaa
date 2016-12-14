package main

import (
	"gaa/canvas"
	"gaa/network"
	"gaa/sim"
	"math/rand"
	"strings"
	"time"

	"strconv"

	"fmt"

	"honnef.co/go/js/dom"
)

const (
	hiddenNeurons = 10
)

type window struct {
	rnd    *rand.Rand
	ticker *time.Ticker
	net    network.Network
	cv     canvas.Canvas
}

func newWindow(seed int64, ctx *dom.CanvasRenderingContext2D) window {
	rnd := rand.New(rand.NewSource(seed))
	net := network.New(sim.NetworkInputs, sim.NetworkOutputs, hiddenNeurons, rnd)
	cv := canvas.New(ctx)
	return window{rnd, nil, net, cv}
}

func (w window) start() {
	w.ticker = time.NewTicker(30 * time.Millisecond)
	go func() {
		sim.Simulate(w.net, w.cv, w.ticker.C)
	}()
}

func (w window) stop() {
	w.ticker.Stop()
}

var windows [16]window

type params struct {
	seed  int64
	steps []int
}

func decodeQueryString() (result params, err error) {
	qs := dom.GetWindow().Location().Search[1:]
	for _, keyValue := range strings.Split(qs, "&") {
		kva := strings.Split(keyValue, "=")
		if len(kva) != 2 {
			continue
		}
		key, value := kva[0], kva[1]
		if key == "seed" {
			result.seed, err = strconv.ParseInt(value, 16, 64)
			if err != nil {
				return result, err
			}
		}
		if key == "steps" {
			for _, step := range value {
				intStep, err := strconv.ParseInt(string(step), 16, 4)
				if err != nil {
					return result, err
				}
				result.steps = append(result.steps, int(intStep))
			}
		}
	}
	return result, nil
}

func encodeQueryString(p params) string {
	var steps []byte
	for _, step := range p.steps {
		steps = strconv.AppendInt(steps, int64(step), 16)
	}
	seed := strconv.FormatInt(p.seed, 16)
	return fmt.Sprintf("seed=%v&steps=%v", seed, string(steps))
}

func addStep(p params, step int) (result params) {
	result.seed = p.seed
	result.steps = make([]int, len(p.steps))
	copy(result.steps, p.steps)
	print("====")
	println(result.steps)
	result.steps = append(result.steps, step)
	println(result.steps)
	return result
}

func main() {
	p, err := decodeQueryString()
	if err != nil {
		panic(err)
	}
	if p.seed == 0 {
		p.seed = time.Now().UnixNano()
	}
	doc := dom.GetWindow().Document()
	body := doc.GetElementByID("body")
	for i := 0; i < 16; i++ {
		extraSeed := int64(i)
		htmlCanvas := doc.CreateElement("canvas").(*dom.HTMLCanvasElement)
		htmlCanvas.Height = sim.ImageHeight
		htmlCanvas.Width = sim.ImageWidth
		htmlCanvas.SetAttribute("style", "width: 300; height: 300; border: 1px solid #dcdcdc; margin: 2px")
		println(p.steps)
		println(addStep(p, i).steps)

		qs := encodeQueryString(addStep(p, i))
		link := doc.CreateElement("a").(*dom.HTMLAnchorElement)
		link.Href = "index.html?" + qs
		link.AppendChild(htmlCanvas)
		body.AppendChild(link)
		windows[i] = newWindow(p.seed+extraSeed, htmlCanvas.GetContext2d())
		windows[i].start()
	}
}
