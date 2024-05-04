//go:build !wasm

package core

import (
	"log"
	"math"
	"reflect"
	"runtime"

	"github.com/wowsims/cata/sim/core/proto"
	googleProto "google.golang.org/protobuf/proto"
)

type simResultCombiner struct {
	count   int32
	results []*proto.RaidSimResult
}

func (comb simResultCombiner) addDistMetrics(base *proto.DistributionMetrics, add *proto.DistributionMetrics, isLast bool) *proto.DistributionMetrics {
	if base == nil {
		base = &proto.DistributionMetrics{
			Min:     math.MaxFloat64,
			MinSeed: math.MaxInt64,
			Hist:    make(map[int32]int32),
		}
	}

	base.Avg += add.Avg / float64(comb.count)
	base.Stdev += (add.Stdev * add.Stdev) / float64(comb.count)
	if isLast {
		base.Stdev = math.Sqrt(base.Stdev)
	}

	base.Max = max(base.Max, add.Max)
	base.MaxSeed = max(base.MaxSeed, add.MaxSeed)

	base.Min = min(base.Min, add.Min)
	base.MinSeed = min(base.MinSeed, add.MinSeed)

	for idx, val := range add.Hist {
		base.Hist[idx] += val
	}

	base.AllValues = append(base.AllValues, add.AllValues...)

	return base
}

func (comb simResultCombiner) addctionMetrics(unit *proto.UnitMetrics, add *proto.ActionMetrics) {
	var am *proto.ActionMetrics

	addKey := add.Id.String()
	for _, baseAction := range unit.Actions {
		if baseAction.Id.String() == addKey {
			am = baseAction
			break
		}
	}

	if am == nil {
		am = &proto.ActionMetrics{
			Id:      add.Id,
			IsMelee: add.IsMelee,
			Targets: make([]*proto.TargetedActionMetrics, len(add.Targets)),
		}
		for i, addTgt := range add.Targets {
			am.Targets[i] = &proto.TargetedActionMetrics{
				UnitIndex: addTgt.UnitIndex,
			}
		}
		unit.Actions = append(unit.Actions, am)
	}

	for i, baseTgt := range am.Targets {
		addTgt := add.Targets[i]
		if baseTgt.UnitIndex != addTgt.UnitIndex {
			panic("Unitidx doesn't match?!")
		}
		baseTgt.Casts += addTgt.Casts
		baseTgt.Hits += addTgt.Hits
		baseTgt.Crits += addTgt.Crits
		baseTgt.Misses += addTgt.Misses
		baseTgt.Dodges += addTgt.Dodges
		baseTgt.Parries += addTgt.Parries
		baseTgt.Blocks += addTgt.Blocks
		baseTgt.Glances += addTgt.Glances
		baseTgt.Damage += addTgt.Damage
		baseTgt.Threat += addTgt.Threat
		baseTgt.Healing += addTgt.Healing
		baseTgt.Shielding += addTgt.Shielding
		baseTgt.CastTimeMs += addTgt.CastTimeMs
	}
}

func (comb simResultCombiner) addAuraMetrics(unit *proto.UnitMetrics, add *proto.AuraMetrics, isLast bool) {
	var am *proto.AuraMetrics

	addKey := add.Id.String()
	for _, baseAura := range unit.Auras {
		if baseAura.Id.String() == addKey {
			am = baseAura
			break
		}
	}

	if am == nil {
		am = &proto.AuraMetrics{
			Id: add.Id,
		}
		unit.Auras = append(unit.Auras, am)
	}

	am.UptimeSecondsAvg += add.UptimeSecondsAvg / float64(comb.count)
	am.ProcsAvg += add.ProcsAvg / float64(comb.count)
	am.UptimeSecondsStdev += (add.UptimeSecondsStdev * add.UptimeSecondsStdev) / float64(comb.count)
	if isLast {
		am.UptimeSecondsStdev = math.Sqrt(am.UptimeSecondsStdev)
	}
}

func (comb simResultCombiner) addResourceMetrics(unit *proto.UnitMetrics, add *proto.ResourceMetrics) {
	var rm *proto.ResourceMetrics

	addKey := add.Id.String()
	for _, baseResource := range unit.Resources {
		if baseResource.Id.String() == addKey {
			rm = baseResource
			break
		}
	}

	if rm == nil {
		rm = &proto.ResourceMetrics{
			Id:   add.Id,
			Type: add.Type,
		}
		unit.Resources = append(unit.Resources, rm)
	}

	rm.Events += add.Events
	rm.Gain += add.Gain
	rm.ActualGain += add.ActualGain
}

func (comb simResultCombiner) addUnitMetrics(base *proto.UnitMetrics, add *proto.UnitMetrics, isLast bool) *proto.UnitMetrics {
	if base == nil {
		base = &proto.UnitMetrics{
			Name:      add.Name,
			UnitIndex: add.UnitIndex,
			Pets:      make([]*proto.UnitMetrics, len(add.Pets)),
		}
	}

	if base.Name != add.Name {
		panic("Names do not match?!")
	}

	base.Dps = comb.addDistMetrics(base.Dps, add.Dps, isLast)
	base.Dpasp = comb.addDistMetrics(base.Dpasp, add.Dpasp, isLast)
	base.Threat = comb.addDistMetrics(base.Threat, add.Threat, isLast)
	base.Dtps = comb.addDistMetrics(base.Dtps, add.Dtps, isLast)
	base.Tmi = comb.addDistMetrics(base.Tmi, add.Tmi, isLast)
	base.Hps = comb.addDistMetrics(base.Hps, add.Hps, isLast)
	base.Tto = comb.addDistMetrics(base.Tto, add.Tto, isLast)

	base.SecondsOomAvg += add.SecondsOomAvg / float64(comb.count)
	base.ChanceOfDeath += add.ChanceOfDeath / float64(comb.count)

	if base.Actions == nil {
		base.Actions = make([]*proto.ActionMetrics, 0, len(add.Actions))
	}
	for _, addAction := range add.Actions {
		comb.addctionMetrics(base, addAction)
	}

	if base.Auras == nil {
		base.Auras = make([]*proto.AuraMetrics, 0, len(add.Auras))
	}
	for _, addAura := range add.Auras {
		comb.addAuraMetrics(base, addAura, isLast)
	}

	if base.Resources == nil {
		base.Resources = make([]*proto.ResourceMetrics, 0, len(add.Resources))
	}
	for _, addResource := range add.Resources {
		comb.addResourceMetrics(base, addResource)
	}

	for i, addPet := range add.Pets {
		base.Pets[i] = comb.addUnitMetrics(base.Pets[i], addPet, isLast)
	}

	return base
}

func (comb simResultCombiner) addRaidMetrics(finalRsr *proto.RaidSimResult, add *proto.RaidMetrics, isLast bool) {
	if finalRsr.RaidMetrics == nil {
		finalRsr.RaidMetrics = &proto.RaidMetrics{
			Parties: make([]*proto.PartyMetrics, len(add.Parties)),
		}
	}

	finalRsr.RaidMetrics.Dps = comb.addDistMetrics(finalRsr.RaidMetrics.Dps, add.Dps, isLast)
	finalRsr.RaidMetrics.Hps = comb.addDistMetrics(finalRsr.RaidMetrics.Hps, add.Hps, isLast)

	for i, addParty := range add.Parties {
		if finalRsr.RaidMetrics.Parties[i] == nil {
			finalRsr.RaidMetrics.Parties[i] = &proto.PartyMetrics{
				Players: make([]*proto.UnitMetrics, len(addParty.Players)),
			}
		}

		base := finalRsr.RaidMetrics.Parties[i]

		base.Dps = comb.addDistMetrics(base.Dps, add.Dps, isLast)
		base.Hps = comb.addDistMetrics(base.Hps, add.Hps, isLast)

		for i, addPlayer := range addParty.Players {
			base.Players[i] = comb.addUnitMetrics(base.Players[i], addPlayer, isLast)
		}
	}
}

func (comb simResultCombiner) addEncounterMetrics(finalRsr *proto.RaidSimResult, add *proto.EncounterMetrics, isLast bool) {
	if finalRsr.EncounterMetrics == nil {
		finalRsr.EncounterMetrics = &proto.EncounterMetrics{
			Targets: make([]*proto.UnitMetrics, len(add.Targets)),
		}
	}

	for i, addTarget := range add.Targets {
		finalRsr.EncounterMetrics.Targets[i] = comb.addUnitMetrics(finalRsr.EncounterMetrics.Targets[i], addTarget, isLast)
	}
}

func (comb simResultCombiner) Combine() *proto.RaidSimResult {
	finalResult := &proto.RaidSimResult{
		Logs:                   comb.results[0].Logs,
		FirstIterationDuration: comb.results[0].FirstIterationDuration,
	}

	for i, addResult := range comb.results {
		isLast := i+1 == int(comb.count)
		comb.addRaidMetrics(finalResult, addResult.RaidMetrics, isLast)
		comb.addEncounterMetrics(finalResult, addResult.EncounterMetrics, isLast)
		finalResult.AvgIterationDuration += addResult.AvgIterationDuration / float64(comb.count)
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
	if csd.Concurrency == 1 {
		return csd.FinalResults[0]
	}

	comb := simResultCombiner{
		count:   csd.Concurrency,
		results: csd.FinalResults,
	}
	finalRes := comb.Combine()
	return finalRes
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
		nextStartSeed := request.SimOptions.RandomSeed // Sims increment their seed each iteration.

		for i := 0; i < concurrency; i++ {
			requestCopy := googleProto.Clone(request).(*proto.RaidSimRequest)
			if i == 0 {
				requestCopy.SimOptions.Iterations += remainder
			}
			requestCopy.SimOptions.RandomSeed = nextStartSeed
			nextStartSeed += int64(requestCopy.SimOptions.Iterations)

			go RunSim(requestCopy, substituteChannels[i])

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
