package concurrency

import (
	"fmt"
	"math"
	"runtime/debug"

	"github.com/wowsims/cata/sim/core/proto"
)

type raidSimResultCombiner struct {
	Debug    bool
	Combined *proto.RaidSimResult
}

func (rsrc *raidSimResultCombiner) newDistMetrics() *proto.DistributionMetrics {
	return &proto.DistributionMetrics{
		Min:            math.MaxFloat64,
		MinSeed:        math.MaxInt64,
		Hist:           make(map[int32]int32),
		AllValues:      make([]float64, 0),
		AggregatorData: &proto.AggregatorData{},
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
		Auras:     make([]*proto.AuraMetrics, len(baseUnit.Auras)),
		Resources: make([]*proto.ResourceMetrics, 0, len(baseUnit.Resources)),
		Pets:      make([]*proto.UnitMetrics, len(baseUnit.Pets)),
	}

	for i, aura := range baseUnit.Auras {
		newUm.Auras[i] = &proto.AuraMetrics{
			Id:             aura.Id,
			AggregatorData: &proto.AggregatorData{},
		}
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

	if add.Max > base.Max {
		base.Max = add.Max
		base.MaxSeed = add.MaxSeed
	}

	if add.Min == 0 || add.Min < base.Min {
		base.Min = add.Min
		base.MinSeed = add.MinSeed
	} else if add.Min == base.Min {
		base.MinSeed = add.MinSeed
	}

	for idx, val := range add.Hist {
		base.Hist[idx] += val
	}

	base.AllValues = append(base.AllValues, add.AllValues...)

	base.AggregatorData.N += add.AggregatorData.N
	base.AggregatorData.SumSq += add.AggregatorData.SumSq
	if isLast {
		base.Stdev = math.Sqrt(base.AggregatorData.SumSq/float64(base.AggregatorData.N) - base.Avg*base.Avg)
	}
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

func (rsrc *raidSimResultCombiner) combineAuraMetrics(base *proto.AuraMetrics, add *proto.AuraMetrics, weight float64, isLast bool) {
	base.UptimeSecondsAvg += add.UptimeSecondsAvg * weight
	base.ProcsAvg += add.ProcsAvg * weight

	base.AggregatorData.N += add.AggregatorData.N
	base.AggregatorData.SumSq += add.AggregatorData.SumSq
	if isLast {
		base.UptimeSecondsStdev = math.Sqrt(base.AggregatorData.SumSq/float64(base.AggregatorData.N) - base.UptimeSecondsAvg*base.UptimeSecondsAvg)
	}
}

func (rsrc *raidSimResultCombiner) addResourceMetrics(unit *proto.UnitMetrics, add *proto.ResourceMetrics) {
	var rm *proto.ResourceMetrics

	rkey := func(r *proto.ResourceMetrics) string {
		return fmt.Sprintf("%s-%d", r.Id.String(), r.Type)
	}

	for _, baseResource := range unit.Resources {
		if rkey(baseResource) == rkey(add) {
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

	for i, addAura := range add.Auras {
		rsrc.combineAuraMetrics(base.Auras[i], addAura, weight, isLast)
	}

	for _, addResource := range add.Resources {
		rsrc.addResourceMetrics(base, addResource)
	}

	for i, addPet := range add.Pets {
		rsrc.combineUnitMetrics(base.Pets[i], addPet, isLast, weight)
	}
}

func (rsrc *raidSimResultCombiner) AddResult(result *proto.RaidSimResult, isLast bool, weight float64) {
	rsrc.combineDistMetrics(rsrc.Combined.RaidMetrics.Dps, result.RaidMetrics.Dps, isLast, weight)
	rsrc.combineDistMetrics(rsrc.Combined.RaidMetrics.Hps, result.RaidMetrics.Hps, isLast, weight)

	for partyIdx, party := range result.RaidMetrics.Parties {
		baseParty := rsrc.Combined.RaidMetrics.Parties[partyIdx]
		rsrc.combineDistMetrics(baseParty.Dps, party.Dps, isLast, weight)
		rsrc.combineDistMetrics(baseParty.Hps, party.Hps, isLast, weight)
		for playerIdx, player := range party.Players {
			rsrc.combineUnitMetrics(baseParty.Players[playerIdx], player, isLast, weight)
		}
	}

	for i, tar := range result.EncounterMetrics.Targets {
		rsrc.combineUnitMetrics(rsrc.Combined.EncounterMetrics.Targets[i], tar, isLast, weight)
	}

	rsrc.Combined.AvgIterationDuration += result.AvgIterationDuration * weight
	rsrc.Combined.IterationsDone += result.IterationsDone

	if rsrc.Debug {
		rsrc.Combined.Logs += "-SIMSTART-\n" + result.Logs
	}
}

func (rsrc *raidSimResultCombiner) SetBaseResult(baseRsr *proto.RaidSimResult) {
	newRsr := &proto.RaidSimResult{
		RaidMetrics: &proto.RaidMetrics{
			Dps:     rsrc.newDistMetrics(),
			Hps:     rsrc.newDistMetrics(),
			Parties: make([]*proto.PartyMetrics, len(baseRsr.RaidMetrics.Parties)),
		},
		EncounterMetrics: &proto.EncounterMetrics{
			Targets: make([]*proto.UnitMetrics, len(baseRsr.EncounterMetrics.Targets)),
		},
		FirstIterationDuration: baseRsr.FirstIterationDuration,
	}

	if !rsrc.Debug {
		newRsr.Logs = baseRsr.Logs
	}

	for i, party := range baseRsr.RaidMetrics.Parties {
		newRsr.RaidMetrics.Parties[i] = rsrc.newPartyMetrics(party)
	}

	for i, tar := range baseRsr.EncounterMetrics.Targets {
		newRsr.EncounterMetrics.Targets[i] = rsrc.newUnitMetrics(tar)
	}

	rsrc.Combined = newRsr
}

func CombineConcurrentResults(results []*proto.RaidSimResult, isDebug bool) *proto.RaidSimResult {
	numResults := len(results)

	if numResults == 0 {
		panic("Result set is empty!")
	}

	if numResults == 1 {
		return results[0]
	}

	var totalIterations int32 = 0
	for _, req := range results {
		totalIterations += req.IterationsDone
	}

	rsrc := raidSimResultCombiner{Debug: isDebug}
	rsrc.SetBaseResult(results[0])
	for i, result := range results {
		resultWeight := float64(results[i].IterationsDone) / float64(totalIterations)
		rsrc.AddResult(result, i == numResults-1, resultWeight)
	}

	return rsrc.Combined
}

func CombineConcurrentResultsAsProto(results []*proto.RaidSimResult, isDebug bool) (rsrCombResult *proto.RaidSimResultCombinationResult) {
	defer func() {
		if err := recover(); err != nil {
			errStr := ""
			switch errt := err.(type) {
			case string:
				errStr = errt
			case error:
				errStr = errt.Error()
			}
			errStr += "\nStack Trace:\n" + string(debug.Stack())
			rsrCombResult = &proto.RaidSimResultCombinationResult{ErrorResult: errStr}
		}
	}()
	rsrCombResult = &proto.RaidSimResultCombinationResult{CombinedResult: CombineConcurrentResults(results, isDebug)}
	return rsrCombResult
}
