package sim

import (
	"gaa/network"
	"io/ioutil"
	"math/rand"
	"testing"
)

func BenchmarkSimulate(b *testing.B) {
	r := rand.New(rand.NewSource(328932186))
	for n := 0; n < b.N; n++ {
		net := network.New(3, 8, 10, r)
		Simulate(net, ioutil.Discard)
	}
}
