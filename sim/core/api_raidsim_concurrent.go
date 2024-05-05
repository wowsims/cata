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

type raidSimResultCombiner struct {
	combined *proto.RaidSimResult
}

func (rsrc *raidSimResultCombiner) newDistMetrics() *proto.DistributionMetrics {
	return &proto.DistributionMetrics{
		Min:       math.MaxFloat64,
		MinSeed:   math.MaxInt64,
		Hist:      make(map[int32]int32),
		AllValues: make([]float64, 0),
	}
}

func (rsrc *raidSimResultCombiner) newUnitMetrics(baseUnit *proto.UnitMetrics) *proto.UnitMetrics {
	newUm := &proto.UnitMetrics{
		Name:      baseUnit.Name,
		UnitIndex: baseUnit.UnitIndex,
		Dps:       rsrc.newDistMetrics(),
		Dpasp:     rsrc.newDistMetrics(),
		Threat:    rsrc.newDistMetrics(),
		Dtps:      rsrc.newDistMetrics(),
		Tmi:       rsrc.newDistMetrics(),
		Hps:       rsrc.newDistMetrics(),
		Tto:       rsrc.newDistMetrics(),
		Actions:   make([]*proto.ActionMetrics, 0, len(baseUnit.Actions)),
		Auras:     make([]*proto.AuraMetrics, 0, len(baseUnit.Auras)),
		Resources: make([]*proto.ResourceMetrics, 0, len(baseUnit.Resources)),
		Pets:      make([]*proto.UnitMetrics, 0, len(baseUnit.Pets)),
	}

	for i, pet := range baseUnit.Pets {
		newUm.Pets[i] = rsrc.newUnitMetrics(pet)
	}

	return newUm
}

func (rsrc *raidSimResultCombiner) newPartyMetrics(baseParty *proto.PartyMetrics) *proto.PartyMetrics {
	newPm := &proto.PartyMetrics{
		Dps:     rsrc.newDistMetrics(),
		Hps:     rsrc.newDistMetrics(),
		Players: make([]*proto.UnitMetrics, len(baseParty.Players)),
	}

	for i, player := range baseParty.Players {
		newPm.Players[i] = rsrc.newUnitMetrics(player)
	}

	return newPm
}

func (rsrc *raidSimResultCombiner) combineDistMetrics(base *proto.DistributionMetrics, add *proto.DistributionMetrics, isLast bool, weight float64) {
	base.Avg += add.Avg * weight
	base.Stdev += (add.Stdev * add.Stdev) * weight
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
}

func (rsrc *raidSimResultCombiner) addActionMetrics(unit *proto.UnitMetrics, add *proto.ActionMetrics) {
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

func (rsrc *raidSimResultCombiner) addAuraMetrics(unit *proto.UnitMetrics, add *proto.AuraMetrics, isLast bool, weight float64) {
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

	am.UptimeSecondsAvg += add.UptimeSecondsAvg * weight
	am.ProcsAvg += add.ProcsAvg * weight
	am.UptimeSecondsStdev += (add.UptimeSecondsStdev * add.UptimeSecondsStdev) * weight
	if isLast {
		am.UptimeSecondsStdev = math.Sqrt(am.UptimeSecondsStdev)
	}
}

func (rsrc *raidSimResultCombiner) addResourceMetrics(unit *proto.UnitMetrics, add *proto.ResourceMetrics) {
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

func (rsrc *raidSimResultCombiner) combineUnitMetrics(base *proto.UnitMetrics, add *proto.UnitMetrics, isLast bool, weight float64) {
	if base.Name != add.Name {
		panic("Names do not match?!")
	}

	if base.UnitIndex != add.UnitIndex {
		panic("UnitIndices do not match?!")
	}

	rsrc.combineDistMetrics(base.Dps, add.Dps, isLast, weight)
	rsrc.combineDistMetrics(base.Dpasp, add.Dpasp, isLast, weight)
	rsrc.combineDistMetrics(base.Threat, add.Threat, isLast, weight)
	rsrc.combineDistMetrics(base.Dtps, add.Dtps, isLast, weight)
	rsrc.combineDistMetrics(base.Tmi, add.Tmi, isLast, weight)
	rsrc.combineDistMetrics(base.Hps, add.Hps, isLast, weight)
	rsrc.combineDistMetrics(base.Tto, add.Tto, isLast, weight)

	base.SecondsOomAvg += add.SecondsOomAvg * weight
	base.ChanceOfDeath += add.ChanceOfDeath * weight

	for _, addAction := range add.Actions {
		rsrc.addActionMetrics(base, addAction)
	}

	for _, addAura := range add.Auras {
		rsrc.addAuraMetrics(base, addAura, isLast, weight)
	}

	for _, addResource := range add.Resources {
		rsrc.addResourceMetrics(base, addResource)
	}

	for i, addPet := range add.Pets {
		rsrc.combineUnitMetrics(base.Pets[i], addPet, isLast, weight)
	}
}

func (rsrc *raidSimResultCombiner) addResult(result *proto.RaidSimResult, isLast bool, weight float64) {
	rsrc.combineDistMetrics(rsrc.combined.RaidMetrics.Dps, result.RaidMetrics.Dps, isLast, weight)
	rsrc.combineDistMetrics(rsrc.combined.RaidMetrics.Hps, result.RaidMetrics.Hps, isLast, weight)

	for partyIdx, party := range result.RaidMetrics.Parties {
		baseParty := rsrc.combined.RaidMetrics.Parties[partyIdx]
		rsrc.combineDistMetrics(baseParty.Dps, party.Dps, isLast, weight)
		rsrc.combineDistMetrics(baseParty.Hps, party.Hps, isLast, weight)
		for playerIdx, player := range party.Players {
			rsrc.combineUnitMetrics(baseParty.Players[playerIdx], player, isLast, weight)
		}
	}

	for i, tar := range result.EncounterMetrics.Targets {
		rsrc.combineUnitMetrics(rsrc.combined.EncounterMetrics.Targets[i], tar, isLast, weight)
	}

	rsrc.combined.AvgIterationDuration += result.AvgIterationDuration * weight
}

func (rsrc *raidSimResultCombiner) setBaseResult(baseRsr *proto.RaidSimResult) {
	newRsr := &proto.RaidSimResult{
		RaidMetrics: &proto.RaidMetrics{
			Dps:     rsrc.newDistMetrics(),
			Hps:     rsrc.newDistMetrics(),
			Parties: make([]*proto.PartyMetrics, len(baseRsr.RaidMetrics.Parties)),
		},
		EncounterMetrics: &proto.EncounterMetrics{
			Targets: make([]*proto.UnitMetrics, len(baseRsr.EncounterMetrics.Targets)),
		},
		Logs:                   baseRsr.Logs,
		FirstIterationDuration: baseRsr.FirstIterationDuration,
	}

	for i, party := range baseRsr.RaidMetrics.Parties {
		newRsr.RaidMetrics.Parties[i] = rsrc.newPartyMetrics(party)
	}

	for i, tar := range baseRsr.EncounterMetrics.Targets {
		newRsr.EncounterMetrics.Targets[i] = rsrc.newUnitMetrics(tar)
	}

	rsrc.combined = newRsr
}

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

func (csd *concurrentSimData) GetCombinedFinalResult() *proto.RaidSimResult {
	if csd.Concurrency == 1 {
		return csd.FinalResults[0]
	}

	rsrc := raidSimResultCombiner{}
	rsrc.setBaseResult(csd.FinalResults[0])
	for i, result := range csd.FinalResults {
		resultWeight := float64(csd.IterationsDone[i]) / float64(csd.IterationsTotal)
		rsrc.addResult(result, i == len(csd.FinalResults)-1, resultWeight)
	}

	return rsrc.combined
}

// Run sim on multiple threads concurrently by splitting interations over multiple sims, transparently combining results into the progress channel.
func runConcurrentSim(request *proto.RaidSimRequest, progress chan *proto.ProgressMetrics) {
	concurrency := runtime.NumCPU()
	substituteChannels := make([]chan *proto.ProgressMetrics, concurrency)
	substituteCases := make([]reflect.SelectCase, concurrency)
	running := concurrency
	csd := concurrentSimData{
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
			FinalRaidResult:     csd.GetCombinedFinalResult(),
		}
	}()
}
