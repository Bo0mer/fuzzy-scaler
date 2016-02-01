package fuzzy

import "log"

const (
	DiskStatusEmpty = iota
	DiskStatusAlmostEmpty
	DiskStatusHalfEmpty
	DiskStatusAlmostFull
	DiskStatusFull
)

var DiskSets = []string{
	DiskStatusEmpty:       "empty",
	DiskStatusAlmostEmpty: "almost empty",
	DiskStatusHalfEmpty:   "half empty",
	DiskStatusAlmostFull:  "almost full",
	DiskStatusFull:        "full"}

const (
	CPUStatusIdle = iota
	CPUStatusLoaded
	CPUStatusOverloaded
)

var CPUSets = []string{
	CPUStatusIdle:       "idle",
	CPUStatusLoaded:     "loaded",
	CPUStatusOverloaded: "overloaded",
}

type Metric struct {
	Instances int
	Disk      float64
	CPU       float64
}

type InstanceRegulator struct {
	diskFuzz   Fuzzificator
	diskDefuzz Defuzzificator
	cpuFuzz    Fuzzificator
	cpuDefuzz  Defuzzificator
	input      <-chan Metric
	output     chan int
	done       chan struct{}
}

func NewInstanceRegulator(diskFuzz, cpuFuzz Fuzzificator, diskDefuzz, cpuDefuzz Defuzzificator) *InstanceRegulator {
	return &InstanceRegulator{
		diskFuzz:   diskFuzz,
		diskDefuzz: diskDefuzz,
		cpuFuzz:    cpuFuzz,
		cpuDefuzz:  cpuDefuzz,
		output:     make(chan int),
		done:       make(chan struct{}),
	}
}

func (r *InstanceRegulator) Start(input <-chan Metric) <-chan int {
	r.input = input
	go r.run()
	return r.output
}

func (r *InstanceRegulator) Close() {
	r.done <- struct{}{}
	<-r.done
	close(r.done)
	close(r.output)
}

func (r *InstanceRegulator) run() {
	for {
		select {
		case m := <-r.input:
			diskWeightVector := r.diskFuzz.Fuzzify(m)
			cpuWeightVector := r.cpuFuzz.Fuzzify(m)
			log.Printf("regulator: calculated disk weight vector: %v\n", diskWeightVector)
			log.Printf("regulator: calculated cpu weight vector: %v\n", cpuWeightVector)
			factor := r.processRules(diskWeightVector, cpuWeightVector)
			// TODO: figure out good scale factor
			instances := int(float64(m.Instances) * factor * 2.0)
			if instances == 0 {
				instances = m.Instances
			}

			log.Printf("regulator: desired number of instances %d\n", instances)
			r.output <- instances
		case <-r.done:
			// notify that we're really leaving
			r.done <- struct{}{}
			return
		}
	}
}

func (r *InstanceRegulator) processRules(diskwv, cpuwv []float64) float64 {
	diskVal := r.diskDefuzz.Defuzzify(diskwv)
	cpuVal := r.cpuDefuzz.Defuzzify(cpuwv)
	log.Printf("regulator: disk weighted average %f\n", diskVal)
	log.Printf("regulator: cpu weighted average %f\n", cpuVal)
	m := (diskVal*0.8 + cpuVal*0.2)
	log.Printf("regulator: output real value %f\n", m)
	return m
}
