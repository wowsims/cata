package core

import (
	"log"
	"reflect"
	"runtime"
	"runtime/debug"

	"github.com/wowsims/cata/sim/core/concurrency"
	"github.com/wowsims/cata/sim/core/proto"
)

type concurrentSimData struct {
	Concurrency     int32
	IterationsTotal int32
	IterationsDone  []int32

	DpsValues []float64
	HpsValues []float64

	FinalResults []*proto.RaidSimResult
}

func (csd *concurrentSimData) GetIterationsDone() int32 {
	var total int32 = 0
	for _, done := range csd.IterationsDone {
		total += done
	}
	return total
}

func (csd *concurrentSimData) GetDpsAvg() float64 {
	total := 0.0
	for _, done := range csd.DpsValues {
		total += done
	}
	return total / float64(csd.Concurrency)
}

func (csd *concurrentSimData) GetHpsAvg() float64 {
	total := 0.0
	for _, done := range csd.HpsValues {
		total += done
	}
	return total / float64(csd.Concurrency)
}

func (csd *concurrentSimData) UpdateProgress(idx int, msg *proto.ProgressMetrics) bool {
	csd.IterationsDone[idx] = msg.CompletedIterations
	csd.DpsValues[idx] = msg.Dps
	csd.HpsValues[idx] = msg.Hps

	if msg.FinalRaidResult != nil {
		csd.FinalResults[idx] = msg.FinalRaidResult
		return true
	}

	return false
}

func (csd *concurrentSimData) MakeProgressMetrics() *proto.ProgressMetrics {
	return &proto.ProgressMetrics{
		TotalIterations:     csd.IterationsTotal,
		CompletedIterations: csd.GetIterationsDone(),
		Dps:                 csd.GetDpsAvg(),
		Hps:                 csd.GetHpsAvg(),
	}
}

// Run sim on multiple threads concurrently by splitting interations over multiple sims, transparently combining results into the progress channel.
func runSimConcurrent(request *proto.RaidSimRequest, progress chan *proto.ProgressMetrics) (result *proto.RaidSimResult) {
	defer func() {
		if !request.SimOptions.IsTest {
			if err := recover(); err != nil {
				errStr := ""
				switch errt := err.(type) {
				case string:
					errStr = errt
				case error:
					errStr = errt.Error()
				}

				errStr += "\nStack Trace:\n" + string(debug.Stack())
				result = &proto.RaidSimResult{ErrorResult: errStr}

				if progress != nil {
					progress <- &proto.ProgressMetrics{FinalRaidResult: result}
				}
			}
		}

		if progress != nil {
			close(progress)
		}
	}()

	threads := TernaryInt32(request.SimOptions.IsTest, 3, int32(runtime.NumCPU()))

	splitRes := concurrency.SplitRequestForConcurrency(request, threads)
	if splitRes.ErrorResult != "" {
		panic(splitRes.ErrorResult)
	}
	threads = splitRes.SplitsDone

	substituteChannels := make([]chan *proto.ProgressMetrics, threads)
	substituteCases := make([]reflect.SelectCase, threads)
	quitChannels := make([]chan bool, threads)
	running := threads

	csd := concurrentSimData{
		Concurrency:     threads,
		IterationsTotal: request.SimOptions.Iterations,
		IterationsDone:  make([]int32, threads),
		DpsValues:       make([]float64, threads),
		HpsValues:       make([]float64, threads),
		FinalResults:    make([]*proto.RaidSimResult, threads),
	}

	for i := 0; i < int(threads); i++ {
		substituteChannels[i] = make(chan *proto.ProgressMetrics, 20)
		substituteCases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(substituteChannels[i])}
		quitChannels[i] = make(chan bool, 1)
	}

	if !request.SimOptions.IsTest {
		log.Printf("Running %d iterations on %d concurrent sims.", csd.IterationsTotal, csd.Concurrency)
	}

	defer func() {
		// Send quit signals to threads in case we returned due to an error.
		for _, quitChan := range quitChannels {
			quitChan <- true
			close(quitChan)
		}
	}()

	for i, req := range splitRes.Requests {
		go RunSim(req, substituteChannels[i], quitChannels[i])
		// Wait for first message to make sure env was constructed. Otherwise concurrent map writes to simdb will happen.
		msg := <-substituteChannels[i]
		// First message may be due to an immediate error, otherwise it can be ignored.
		if msg.FinalRaidResult != nil && msg.FinalRaidResult.ErrorResult != "" {
			if progress != nil {
				progress <- msg
			}
			log.Printf("Thread %d had an error. Cancelling all sims!", i)
			return msg.FinalRaidResult
		}
	}

	progressCounter := 0

	for running > 0 {
		i, val, ok := reflect.Select(substituteCases)

		if !ok {
			substituteCases[i].Chan = reflect.ValueOf(nil)
			running -= 1
			continue
		}

		msg := val.Interface().(*proto.ProgressMetrics)
		if csd.UpdateProgress(i, msg) {
			if msg.FinalRaidResult != nil && msg.FinalRaidResult.ErrorResult != "" {
				if progress != nil {
					progress <- msg
				}
				log.Printf("Thread %d had an error. Cancelling all sims!", i)
				return msg.FinalRaidResult
			}
			substituteCases[i].Chan = reflect.ValueOf(nil)
			running -= 1
			continue
		}

		if progress != nil {
			progressCounter++ // Don't spam progress
			if progressCounter%int(threads) == 0 {
				progress <- csd.MakeProgressMetrics()
			}
		}
	}

	for _, res := range csd.FinalResults {
		if res == nil {
			panic("Missing one or more final sim result(s)!")
		}
	}

	if !request.SimOptions.IsTest {
		log.Printf("All %d sims finished successfully.", csd.Concurrency)
	}

	result = concurrency.CombineConcurrentResults(csd.FinalResults, request.SimOptions.Debug)

	if progress != nil {
		pm := csd.MakeProgressMetrics()
		pm.FinalRaidResult = result
		progress <- pm
	}

	return result
}
