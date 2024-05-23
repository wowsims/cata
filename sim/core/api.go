// Proto-based function interface for the simulator
package core

import (
	"context"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/simsignals"
	"github.com/wowsims/cata/sim/core/stats"
)

/**
 * Returns character stats taking into account gear / buffs / consumes / etc
 */
func ComputeStats(csr *proto.ComputeStatsRequest) *proto.ComputeStatsResult {
	encounter := csr.Encounter
	if encounter == nil {
		encounter = &proto.Encounter{}
	}

	_, raidStats, encounterStats := NewEnvironment(csr.Raid, encounter, true)

	return &proto.ComputeStatsResult{
		RaidStats:      raidStats,
		EncounterStats: encounterStats,
	}
}

/**
 * Returns stat weights and EP values, with standard deviations, for all stats.
 */
func StatWeights(request *proto.StatWeightsRequest) *proto.StatWeightsResult {
	result := CalcStatWeight(request, stats.Stat(request.EpReferenceStat), nil)
	return result.ToProto()
}

func StatWeightsAsync(request *proto.StatWeightsRequest, progress chan *proto.ProgressMetrics) {
	go func() {
		result := CalcStatWeight(request, stats.Stat(request.EpReferenceStat), progress)
		progress <- &proto.ProgressMetrics{
			FinalWeightResult: result.ToProto(),
		}
	}()
}

/**
 * Runs multiple iterations of the sim with a full raid.
 */
func RunRaidSim(request *proto.RaidSimRequest) *proto.RaidSimResult {
	return RunSim(request, nil, simsignals.CreateSignals())
}

func RunRaidSimAsync(request *proto.RaidSimRequest, progress chan *proto.ProgressMetrics) {
	simId := request.Id
	signals, err := simsignals.RegisterWithId(request.Id)
	if err != nil {
		progress <- &proto.ProgressMetrics{FinalRaidResult: &proto.RaidSimResult{ErrorResult: "Couldn't register for signal API: " + err.Error()}}
		return
	}
	go func() {
		defer simsignals.UnregisterId(simId)
		RunSim(request, progress, signals)
	}()
}

// Threading does not work in WASM!
func RunRaidSimConcurrent(request *proto.RaidSimRequest) *proto.RaidSimResult {
	return runSimConcurrent(request, nil, simsignals.CreateSignals())
}

// Threading does not work in WASM!
func RunRaidSimConcurrentAsync(request *proto.RaidSimRequest, progress chan *proto.ProgressMetrics) {
	simId := request.Id
	signals, err := simsignals.RegisterWithId(request.Id)
	if err != nil {
		progress <- &proto.ProgressMetrics{FinalRaidResult: &proto.RaidSimResult{ErrorResult: "Couldn't register for signal API: " + err.Error()}}
		return
	}
	go func() {
		defer simsignals.UnregisterId(simId)
		runSimConcurrent(request, progress, signals)
	}()
}

func RunBulkSim(request *proto.BulkSimRequest) *proto.BulkSimResult {
	return BulkSim(context.Background(), request, nil)
}

func RunBulkSimAsync(ctx context.Context, request *proto.BulkSimRequest, progress chan *proto.ProgressMetrics) {
	go BulkSim(ctx, request, progress)
}
