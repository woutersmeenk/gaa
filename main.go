package main

import (
	"bufio"
	"fmt"
	"gaa/network"
	"gaa/sim"
	"gaa/svg"
	"math/rand"
	"os"
)

const (
	hiddenNeurons = 10
	seed          = 64877563487
)

func main() {
	for i := 0; i < 20; i++ {
		r := rand.New(rand.NewSource(seed + int64(i)))
		net := network.New(sim.NetworkInputs, sim.NetworkOutputs, hiddenNeurons, r)
		fn := fmt.Sprintf("output/output-%v.svg", i)
		f, err := os.Create(fn)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		writer := bufio.NewWriter(f)
		defer writer.Flush()
		sim.Simulate(net, svg.New(writer))
	}
}
