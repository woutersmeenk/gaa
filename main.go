package main

import (
	"math/rand"
	"strings"
	"time"

	"github.com/woutersmeenk/gaa/canvas"
	"github.com/woutersmeenk/gaa/network"
	"github.com/woutersmeenk/gaa/sim"

	"strconv"

	"fmt"

	"honnef.co/go/js/dom"
)

const (
	hiddenNeurons = 10
)

var windows [16]*window

type window struct {
	ticker *time.Ticker
	net    network.Network
	cv     canvas.Canvas
}

func newWindow(id int, params parameters, ctx *dom.CanvasRenderingContext2D) window {
	var net network.Network
	seed := params.seed
	for _, step := range params.steps {
		seed += int64(step)
		rnd := rand.New(rand.NewSource(seed))
		if net == nil {
			net = network.New(sim.NetworkInputs, sim.NetworkOutputs, hiddenNeurons, rnd)
		} else {
			net = net.Mutate(rnd)
		}
	}
	cv := canvas.New(ctx)
	return window{nil, net, cv}
}

func (w *window) start() {
	w.ticker = time.NewTicker(30 * time.Millisecond)
	go func() {
		sim.Simulate(w.net, w.cv, w.ticker.C)
	}()
}

func stopAllWindows() {
	for _, win := range windows {
		win.ticker.Stop()
	}
}

type parameters struct {
	seed  int64
	steps []int
}

func decodeQueryString() (result parameters, err error) {
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
				intStep, err := strconv.ParseUint(string(step), 16, 4)
				if err != nil {
					return result, err
				}
				result.steps = append(result.steps, int(intStep))
			}
		}
	}
	return result, nil
}

func (p parameters) encodeQueryString() string {
	var steps []byte
	for _, step := range p.steps {
		steps = strconv.AppendInt(steps, int64(step), 16)
	}
	seed := strconv.FormatInt(p.seed, 16)
	return fmt.Sprintf("seed=%v&steps=%v", seed, string(steps))
}

func (p parameters) addStep(step int) (result parameters) {
	result.seed = p.seed
	result.steps = make([]int, len(p.steps))
	copy(result.steps, p.steps)
	result.steps = append(result.steps, step)
	return result
}

func main() {
	params, err := decodeQueryString()
	if err != nil {
		panic(err)
	}
	if params.seed == 0 {
		params.seed = time.Now().UnixNano()
		dom.GetWindow().Location().Search = params.encodeQueryString()
	}
	doc := dom.GetWindow().Document()
	stopLink := doc.GetElementByID("stop")
	stopLink.AddEventListener("click", false, func(e dom.Event) {
		stopAllWindows()
	})
	windowsEl := doc.GetElementByID("windows")
	lastStep := -1
	if len(params.steps) > 0 {
		lastStep = params.steps[len(params.steps)-1]
	}
	for windowID := 0; windowID < 16; windowID++ {
		htmlCanvas := doc.CreateElement("canvas").(*dom.HTMLCanvasElement)
		link := doc.CreateElement("a").(*dom.HTMLAnchorElement)
		link.AppendChild(htmlCanvas)
		windowsEl.AppendChild(link)

		htmlCanvas.Height = sim.ImageHeight
		htmlCanvas.Width = sim.ImageWidth
		htmlCanvas.SetClass("window-selected")

		newParams := params
		if lastStep != windowID {
			newParams = params.addStep(windowID)
			htmlCanvas.SetClass("window-not-selected")
		}

		link.Href = "index.html?" + newParams.encodeQueryString()

		win := newWindow(windowID, newParams, htmlCanvas.GetContext2d())
		windows[windowID] = &win
		windows[windowID].start()
	}
}
