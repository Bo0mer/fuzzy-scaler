package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/Bo0mer/fuzzy-scaler/fuzzy"
)

func main() {
	log.Println("Starting...")
	empty := fuzzy.NewTriangularFunc(0.0, 0.30)
	almostEmpty := fuzzy.NewTriangularFunc(0.20, 0.50)
	halfEmpty := fuzzy.NewTriangularFunc(0.40, 0.70)
	almostFull := fuzzy.NewTriangularFunc(0.60, 0.90)
	full := fuzzy.NewTriangularFunc(0.80, 1.00)

	diskSelector := func(m fuzzy.Metric) float64 {
		return m.Disk
	}
	diskFuzz := fuzzy.NewFuzzificator(diskSelector, empty, almostEmpty, halfEmpty, almostFull, full)
	diskDefuzz := fuzzy.NewWeightedAverageDefuzzificator(empty, almostEmpty, halfEmpty, almostFull, full)

	idle := fuzzy.NewTriangularFunc(0.0, 0.40)
	loaded := fuzzy.NewTriangularFunc(0.3, 0.7)
	overloaded := fuzzy.NewTriangularFunc(0.6, 1.0)

	cpuSelector := func(m fuzzy.Metric) float64 {
		return m.CPU
	}
	cpuFuzz := fuzzy.NewFuzzificator(cpuSelector, idle, loaded, overloaded)
	cpuDefuzz := fuzzy.NewWeightedAverageDefuzzificator(idle, loaded, overloaded)

	ai := fuzzy.NewInstanceRegulator(diskFuzz, cpuFuzz, diskDefuzz, cpuDefuzz)

	metrics := make(chan fuzzy.Metric)
	instances := ai.Start(metrics)
	instanceCount := 10
	for i := 0; i < 10000; i++ {
		time.Sleep(time.Millisecond * 500)
		metrics <- fuzzy.Metric{
			Instances: instanceCount,
			Disk:      rand.Float64(),
			CPU:       rand.Float64(),
		}
		instanceCount = <-instances
	}
}
