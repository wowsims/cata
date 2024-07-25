// Proto-based function interface for the simulator
package core

import (
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/simsignals"
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
	return runStatWeights(request, nil, simsignals.CreateSignals())
}

func StatWeightsAsync(request *proto.StatWeightsRequest, progress chan *proto.ProgressMetrics, requestId string) {
	signals, err := simsignals.RegisterWithId(requestId)
	if err != nil {
		progress <- &proto.ProgressMetrics{
			FinalWeightResult: &proto.StatWeightsResult{
				Error: &proto.ErrorOutcome{
					Message: "Couldn't register for signal API: " + err.Error(),
				},
			},
		}
		return
	}
	go func() {
		defer simsignals.UnregisterId(requestId)
		result := runStatWeights(request, progress, signals)
		progress <- &proto.ProgressMetrics{
			FinalWeightResult: result,
		}
	}()
}

// Get data for all requests needed for stat weights.
func StatWeightRequests(request *proto.StatWeightsRequest) *proto.StatWeightRequestsData {
	return buildStatWeightRequests(request)
}

func StatWeightCompute(request *proto.StatWeightsCalcRequest) *proto.StatWeightsResult {
	return computeStatWeights(request)
}

/**
 * Runs multiple iterations of the sim with a full raid.
 */
func RunRaidSim(request *proto.RaidSimRequest) *proto.RaidSimResult {
	return RunSim(request, nil, simsignals.CreateSignals())
}

func RunRaidSimAsync(request *proto.RaidSimRequest, progress chan *proto.ProgressMetrics, requestId string) {
	signals, err := simsignals.RegisterWithId(requestId)
	if err != nil {
		progress <- &proto.ProgressMetrics{
			FinalRaidResult: &proto.RaidSimResult{
				Error: &proto.ErrorOutcome{
					Message: "Couldn't register for signal API: " + err.Error(),
				},
			},
		}
		return
	}
	go func() {
		defer simsignals.UnregisterId(requestId)
		RunSim(request, progress, signals)
	}()
}

// Threading does not work in WASM!
func RunRaidSimConcurrent(request *proto.RaidSimRequest) *proto.RaidSimResult {
	return runSimConcurrent(request, nil, simsignals.CreateSignals())
}

// Threading does not work in WASM!
func RunRaidSimConcurrentAsync(request *proto.RaidSimRequest, progress chan *proto.ProgressMetrics, requestId string) {
	signals, err := simsignals.RegisterWithId(requestId)
	if err != nil {
		progress <- &proto.ProgressMetrics{
			FinalRaidResult: &proto.RaidSimResult{
				Error: &proto.ErrorOutcome{
					Message: "Couldn't register for signal API: " + err.Error(),
				},
			},
		}
		return
	}
	go func() {
		defer simsignals.UnregisterId(requestId)
		runSimConcurrent(request, progress, signals)
	}()
}

func RunBulkSim(request *proto.BulkSimRequest) *proto.BulkSimResult {
	return BulkSim(simsignals.CreateSignals(), request, nil)
}

func RunBulkSimAsync(request *proto.BulkSimRequest, progress chan *proto.ProgressMetrics, requestId string) {
	signals, err := simsignals.RegisterWithId(requestId)
	if err != nil {
		progress <- &proto.ProgressMetrics{
			FinalBulkResult: &proto.BulkSimResult{
				Error: &proto.ErrorOutcome{
					Message: "Couldn't register for signal API: " + err.Error(),
				},
			},
		}
		return
	}
	go func() {
		defer simsignals.UnregisterId(requestId)
		BulkSim(signals, request, progress)
	}()
}

var runningInWasm = false

func SetRunningInWasm() {
	runningInWasm = true
}

func IsRunningInWasm() bool {
	return runningInWasm
}

func RunBulkCombos(request *proto.BulkSimCombosRequest) *proto.BulkSimCombosResult {
	bulkSimReq := &proto.BulkSimCombosRequest{
		BaseSettings: request.BaseSettings,
		BulkSettings: request.BulkSettings,
	}
	return BulkSimCombos(simsignals.CreateSignals(), bulkSimReq)
}
