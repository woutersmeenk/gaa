package main

import (
	"bufio"
	"fmt"
	"gaa/network"
	"gaa/sim"
	"math/rand"
	"os"
)

const (
	networkInputs  = 3
	networkOutputs = 8
	hiddenNeurons  = 10
	seed           = 5135264217
)

func main() {
	for i := 0; i < 20; i++ {
		r := rand.New(rand.NewSource(seed + int64(i)))
		net := network.New(networkInputs, networkOutputs, hiddenNeurons, r)
		fn := fmt.Sprintf("output/output-%v.svg", i)
		f, err := os.Create(fn)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		writer := bufio.NewWriter(f)
		defer writer.Flush()
		sim.Simulate(net, writer)
	}
}
