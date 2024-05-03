//go:build !wasm

package core

import (
	"log"
	"reflect"
	"runtime"

	"github.com/wowsims/cata/sim/core/proto"
	googleProto "google.golang.org/protobuf/proto"
)

type simResultCombiner struct {
	results []*proto.RaidSimResult
}

func (comb simResultCombiner) combineDistMetrics(base *proto.DistributionMetrics, add *proto.DistributionMetrics) {
	base.Avg = (base.Avg + add.Avg) / 2
	base.Stdev = (base.Stdev + add.Stdev) / 2

	base.Max = max(base.Max, add.Max)
	base.MaxSeed = max(base.MaxSeed, add.MaxSeed)

	base.Min = min(base.Min, add.Min)
	base.MinSeed = min(base.MinSeed, add.MinSeed)

	for idx, val := range add.Hist {
		base.Hist[idx] += val
	}

	base.AllValues = append(base.AllValues, add.AllValues...)
}

func (comb simResultCombiner) combineTargetActionMetrics(base *proto.TargetedActionMetrics, add *proto.TargetedActionMetrics) {
	if base.UnitIndex != add.UnitIndex {
		panic("Unitidx doesn't match?!")
	}

	base.Casts += add.Casts
	base.Hits += add.Hits
	base.Crits += add.Crits
	base.Misses += add.Misses
	base.Dodges += add.Dodges
	base.Parries += add.Parries
	base.Blocks += add.Blocks
	base.Glances += add.Glances
	base.Damage += add.Damage
	base.Threat += add.Threat
	base.Healing += add.Healing
	base.Shielding += add.Shielding
	base.CastTimeMs += add.CastTimeMs
}

func (comb simResultCombiner) combineActionMetrics(unit *proto.UnitMetrics, add *proto.ActionMetrics) {
	addKey := add.Id.String()
	for _, baseAction := range unit.Actions {
		if baseAction.Id.String() == addKey {
			for i, tgt := range baseAction.Targets {
				comb.combineTargetActionMetrics(tgt, add.Targets[i])
			}
			return
		}
	}
	unit.Actions = append(unit.Actions, add)
}

func (comb simResultCombiner) combineAuraMetrics(unit *proto.UnitMetrics, add *proto.AuraMetrics) {
	addKey := add.Id.String()
	for _, baseAura := range unit.Auras {
		if baseAura.Id.String() == addKey {
			baseAura.UptimeSecondsAvg = (baseAura.UptimeSecondsAvg + add.UptimeSecondsAvg) / 2
			baseAura.UptimeSecondsStdev = (baseAura.UptimeSecondsStdev + add.UptimeSecondsStdev) / 2
			baseAura.ProcsAvg = (baseAura.ProcsAvg + add.ProcsAvg) / 2
			return
		}
	}
	unit.Auras = append(unit.Auras, add)
}

func (comb simResultCombiner) combineResourceMetrics(unit *proto.UnitMetrics, add *proto.ResourceMetrics) {
	addKey := add.Id.String()
	for _, baseResource := range unit.Resources {
		if baseResource.Id.String() == addKey {
			baseResource.Events += add.Events
			baseResource.Gain += add.Gain
			baseResource.ActualGain += add.ActualGain
			return
		}
	}
	unit.Resources = append(unit.Resources, add)
}

func (comb simResultCombiner) combineUnitMetrics(base *proto.UnitMetrics, add *proto.UnitMetrics) {
	if base.Name != add.Name {
		panic("Names do not match?!")
	}

	comb.combineDistMetrics(base.Dps, add.Dps)
	comb.combineDistMetrics(base.Dpasp, add.Dpasp)
	comb.combineDistMetrics(base.Threat, add.Threat)
	comb.combineDistMetrics(base.Dtps, add.Dtps)
	comb.combineDistMetrics(base.Tmi, add.Tmi)
	comb.combineDistMetrics(base.Hps, add.Hps)
	comb.combineDistMetrics(base.Tto, add.Tto)

	base.SecondsOomAvg = (base.SecondsOomAvg + add.SecondsOomAvg) / 2
	base.ChanceOfDeath = (base.ChanceOfDeath + add.ChanceOfDeath) / 2

	for _, addAction := range add.Actions {
		comb.combineActionMetrics(base, addAction)
	}

	for _, addAura := range add.Auras {
		comb.combineAuraMetrics(base, addAura)
	}

	for _, addResource := range add.Resources {
		comb.combineResourceMetrics(base, addResource)
	}

	for i, pet := range base.Pets {
		comb.combineUnitMetrics(pet, add.Pets[i])
	}
}

func (comb simResultCombiner) combinePartyMetrics(base *proto.PartyMetrics, add *proto.PartyMetrics) {
	comb.combineDistMetrics(base.Dps, add.Dps)
	comb.combineDistMetrics(base.Hps, add.Hps)
	for i, player := range base.Players {
		comb.combineUnitMetrics(player, add.Players[i])
	}
}

func (comb simResultCombiner) combineRaidMetrics(base *proto.RaidMetrics, add *proto.RaidMetrics) {
	comb.combineDistMetrics(base.Dps, add.Dps)
	comb.combineDistMetrics(base.Hps, add.Hps)
	for i, party := range base.Parties {
		comb.combinePartyMetrics(party, add.Parties[i])
	}
}

func (comb simResultCombiner) combineEncounterMetrics(base *proto.EncounterMetrics, add *proto.EncounterMetrics) {
	for i, target := range base.Targets {
		comb.combineUnitMetrics(target, add.Targets[i])
	}
}

func (comb simResultCombiner) combineRaidResults(base *proto.RaidSimResult, add *proto.RaidSimResult) {
	comb.combineRaidMetrics(base.RaidMetrics, add.RaidMetrics)
	comb.combineEncounterMetrics(base.EncounterMetrics, add.EncounterMetrics)
	base.AvgIterationDuration = (base.AvgIterationDuration + add.AvgIterationDuration) / 2
}

func (comb simResultCombiner) Combine() *proto.RaidSimResult {
	finalResult := comb.results[0]
	for i := 1; i < len(comb.results); i++ {
		comb.combineRaidResults(finalResult, comb.results[i])
	}
	return finalResult
}

type concurrentSimData struct {
	Concurrency     int32
	IterationsTotal int32
	IterationsDone  []int32

	DpsValues []float64
	HpsValues []float64

	FinalResults []*proto.RaidSimResult
}

func (csd concurrentSimData) GetIterationsDone() int32 {
	var total int32 = 0
	for _, done := range csd.IterationsDone {
		total += done
	}
	return total
}

func (csd concurrentSimData) GetDpsAvg() float64 {
	total := 0.0
	for _, done := range csd.DpsValues {
		total += done
	}
	return total / float64(len(csd.DpsValues))
}

func (csd concurrentSimData) GetHpsAvg() float64 {
	total := 0.0
	for _, done := range csd.HpsValues {
		total += done
	}
	return total / float64(len(csd.HpsValues))
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

func (csd *concurrentSimData) GetFinalResult() *proto.RaidSimResult {
	comb := simResultCombiner{
		results: csd.FinalResults,
	}
	return comb.Combine()
}

func runConcurrentSim(request *proto.RaidSimRequest, progress chan *proto.ProgressMetrics) {
	concurrency := runtime.NumCPU()
	substituteChannels := make([]chan *proto.ProgressMetrics, concurrency)
	substituteCases := make([]reflect.SelectCase, concurrency)
	running := concurrency
	csd := &concurrentSimData{
		Concurrency:     int32(concurrency),
		IterationsTotal: request.SimOptions.Iterations,
		IterationsDone:  make([]int32, concurrency),
		DpsValues:       make([]float64, concurrency),
		HpsValues:       make([]float64, concurrency),
		FinalResults:    make([]*proto.RaidSimResult, concurrency),
	}

	for i := 0; i < concurrency; i++ {
		substituteChannels[i] = make(chan *proto.ProgressMetrics, cap(progress))
		substituteCases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(substituteChannels[i])}
	}

	log.Printf("Running %d iterations on %d concurrent sims.", csd.IterationsTotal, csd.Concurrency)

	go func() {
		remainder := request.SimOptions.Iterations % int32(concurrency)
		request.SimOptions.Iterations /= int32(concurrency)

		for i := 0; i < concurrency; i++ {
			if i == 0 && remainder > 0 {
				requestRemainderIterations := googleProto.Clone(request).(*proto.RaidSimRequest)
				requestRemainderIterations.SimOptions.Iterations += remainder
				go RunSim(requestRemainderIterations, substituteChannels[i])
			} else {
				go RunSim(request, substituteChannels[i])
			}

			// Wait for first message to make sure env was constructed. Otherwise concurrent map writes to simdb will happen.
			msg := <-substituteChannels[i]
			// First message may be due to an immediate error, otherwise it can be ignored.
			if msg.FinalRaidResult != nil && msg.FinalRaidResult.ErrorResult != "" {
				progress <- msg
				return
			}
		}

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
					progress <- msg
					// TODO: cancel still running routines on error?
					return
				}
				continue
			}

			progress <- &proto.ProgressMetrics{
				TotalIterations:     csd.IterationsTotal,
				CompletedIterations: csd.GetIterationsDone(),
				Dps:                 csd.GetDpsAvg(),
				Hps:                 csd.GetHpsAvg(),
			}
		}

		log.Printf("All %d sims finished successfully.", csd.Concurrency)

		progress <- &proto.ProgressMetrics{
			TotalIterations:     csd.IterationsTotal,
			CompletedIterations: csd.GetIterationsDone(),
			Dps:                 csd.GetDpsAvg(),
			Hps:                 csd.GetHpsAvg(),
			FinalRaidResult:     csd.GetFinalResult(),
		}
	}()
}
