package network

import (
	"math/rand"
)

type Network struct {
	numInputs, numOutputs, numHiddenNeurons int
	// First the inputs then the outputs then the hidden neurons
	neuronValues  []float64
	neuronWeights [][]float64
}

type InputTranslator interface {
	TransInputs() []float64
}

type OutputTranslator interface {
	TransOutputs(outputs []float64)
}

func New(numInputs, numOutputs, numHiddenNeurons int, r *rand.Rand) *Network {
	numNeurons := numInputs + numOutputs + numHiddenNeurons
	neuronValues := make([]float64, numNeurons)
	neuronWeights := make([][]float64, numNeurons)
	// Skip the input neurons they dont have outgoing weights
	for from := numInputs; from < numNeurons; from++ {
		neuronWeights[from] = make([]float64, numNeurons)
		for to := 0; to < numNeurons; to++ {
			neuronWeights[from][to] = r.Float64()*2 - 1
		}
	}
	return &Network{numInputs, numOutputs, numHiddenNeurons, neuronValues, neuronWeights}
}

func (net *Network) Eval(input InputTranslator, output OutputTranslator) {
	inputs := input.TransInputs()
	numNeurons := len(net.neuronValues)
	newNeuronValues := make([]float64, numNeurons)
	for i := 0; i < net.numInputs; i++ {
		net.neuronValues[i] = inputs[i]
		newNeuronValues[i] = inputs[i]
	}
	for from := net.numInputs; from < numNeurons; from++ {
		for to := 0; to < numNeurons; to++ {
			newNeuron := net.neuronWeights[from][to] * net.neuronValues[to]
			newNeuronValues[from] += newNeuron
		}
		newNeuronValues[from] /= float64(numNeurons)
	}
	net.neuronValues = newNeuronValues
	output.TransOutputs(net.neuronValues[net.numInputs : net.numInputs+net.numOutputs])
}
