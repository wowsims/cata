package core

import (
	"math"
	"time"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/simsignals"
	"github.com/wowsims/cata/sim/core/stats"
	googleProto "google.golang.org/protobuf/proto"
)

const DTPSReferenceStat = stats.Armor

type UnitStats struct {
	Stats       stats.Stats
	PseudoStats []float64
}

func NewUnitStats() UnitStats {
	return UnitStats{
		PseudoStats: make([]float64, stats.PseudoStatsLen),
	}
}
func (s *UnitStats) AddStat(stat stats.UnitStat, value float64) {
	if stat.IsStat() {
		s.Stats[stat.StatIdx()] += value
	} else {
		s.PseudoStats[stat.PseudoStatIdx()] += value
	}
}
func (s *UnitStats) Get(stat stats.UnitStat) float64 {
	if stat.IsStat() {
		return s.Stats[stat.StatIdx()]
	} else {
		return s.PseudoStats[stat.PseudoStatIdx()]
	}
}

func (s *UnitStats) ToProto() *proto.UnitStats {
	return &proto.UnitStats{
		Stats:       s.Stats[:],
		PseudoStats: s.PseudoStats,
	}
}

type StatWeightValues struct {
	Weights       UnitStats
	WeightsStdev  UnitStats
	EpValues      UnitStats
	EpValuesStdev UnitStats
}

func NewStatWeightValues() StatWeightValues {
	return StatWeightValues{
		Weights:       NewUnitStats(),
		WeightsStdev:  NewUnitStats(),
		EpValues:      NewUnitStats(),
		EpValuesStdev: NewUnitStats(),
	}
}

func (swv *StatWeightValues) ToProto() *proto.StatWeightValues {
	return &proto.StatWeightValues{
		Weights:       swv.Weights.ToProto(),
		WeightsStdev:  swv.WeightsStdev.ToProto(),
		EpValues:      swv.EpValues.ToProto(),
		EpValuesStdev: swv.EpValuesStdev.ToProto(),
	}
}

type StatWeightsResult struct {
	Dps    StatWeightValues
	Hps    StatWeightValues
	Tps    StatWeightValues
	Dtps   StatWeightValues
	Tmi    StatWeightValues
	PDeath StatWeightValues
}

func NewStatWeightsResult() *StatWeightsResult {
	return &StatWeightsResult{
		Dps:    NewStatWeightValues(),
		Hps:    NewStatWeightValues(),
		Tps:    NewStatWeightValues(),
		Dtps:   NewStatWeightValues(),
		Tmi:    NewStatWeightValues(),
		PDeath: NewStatWeightValues(),
	}
}

func (swr *StatWeightsResult) ToProto() *proto.StatWeightsResult {
	return &proto.StatWeightsResult{
		Dps:    swr.Dps.ToProto(),
		Hps:    swr.Hps.ToProto(),
		Tps:    swr.Tps.ToProto(),
		Dtps:   swr.Dtps.ToProto(),
		Tmi:    swr.Tmi.ToProto(),
		PDeath: swr.PDeath.ToProto(),
	}
}

func buildStatWeightRequests(swr *proto.StatWeightsRequest) *proto.StatWeightRequestsData {
	if swr.Player.BonusStats == nil {
		swr.Player.BonusStats = &proto.UnitStats{}
	}
	if swr.Player.BonusStats.Stats == nil {
		swr.Player.BonusStats.Stats = make([]float64, stats.Len)
	}
	if swr.Player.BonusStats.PseudoStats == nil {
		swr.Player.BonusStats.PseudoStats = make([]float64, stats.PseudoStatsLen)
	}

	raidProto := SinglePlayerRaidProto(swr.Player, swr.PartyBuffs, swr.RaidBuffs, swr.Debuffs)
	raidProto.Tanks = swr.Tanks

	swr.SimOptions.SaveAllValues = true

	// Cut in half since we're doing above and below separately.
	// This number needs to be the same for the baseline sim too, so that RNG lines up perfectly.
	swr.SimOptions.Iterations /= 2

	// Make sure an RNG seed is always set because it gives more consistent results.
	// When there is no user-supplied seed it needs to be a randomly-selected seed
	// though, so that run-run differences still exist.
	if swr.SimOptions.RandomSeed == 0 {
		swr.SimOptions.RandomSeed = time.Now().UnixNano()
	}

	// Reduce variance even more by using test-level RNG controls.
	swr.SimOptions.UseLabeledRands = true

	swBaseResponse := &proto.StatWeightRequestsData{
		BaseRequest: &proto.RaidSimRequest{
			RequestId:  swr.RequestId,
			Raid:       raidProto,
			Encounter:  swr.Encounter,
			SimOptions: swr.SimOptions,
		},
		EpReferenceStat: swr.EpReferenceStat,
		StatSimRequests: []*proto.StatWeightStatRequestData{},
	}

	// Do half the iterations with a positive, and half with a negative value for better accuracy.
	const defaultStatMod = 40.0 // match to the impact of a single gem
	statModsLow := make([]float64, stats.UnitStatsLen)
	statModsHigh := make([]float64, stats.UnitStatsLen)

	// Make sure reference stat is included.
	statModsLow[swr.EpReferenceStat] = -defaultStatMod
	statModsHigh[swr.EpReferenceStat] = defaultStatMod

	statsToWeigh := stats.ProtoArrayToStatsList(swr.StatsToWeigh)
	for _, s := range statsToWeigh {
		stat := stats.UnitStatFromStat(s)
		statMod := defaultStatMod
		if stat.EqualsStat(stats.Armor) || stat.EqualsStat(stats.BonusArmor) {
			statMod = defaultStatMod * 10
		}
		statModsHigh[stat] = statMod
		statModsLow[stat] = -statMod
	}
	for _, s := range swr.PseudoStatsToWeigh {
		stat := stats.UnitStatFromPseudoStat(s)
		statMod := defaultStatMod * 0.5
		statModsHigh[stat] = statMod
		statModsLow[stat] = -statMod
	}

	for i := range statModsLow {
		stat := stats.UnitStatFromIdx(i)
		if statModsLow[stat] == 0 {
			continue
		}

		lowSimRequest := googleProto.Clone(swBaseResponse.BaseRequest).(*proto.RaidSimRequest)
		stat.AddToStatsProto(lowSimRequest.Raid.Parties[0].Players[0].BonusStats, statModsLow[stat])

		highSimRequest := googleProto.Clone(swBaseResponse.BaseRequest).(*proto.RaidSimRequest)
		stat.AddToStatsProto(highSimRequest.Raid.Parties[0].Players[0].BonusStats, statModsHigh[stat])

		swBaseResponse.StatSimRequests = append(swBaseResponse.StatSimRequests, &proto.StatWeightStatRequestData{
			UnitStat:    int32(stat),
			RequestLow:  lowSimRequest,
			RequestHigh: highSimRequest,
			ModLow:      statModsLow[stat],
			ModHigh:     statModsHigh[stat],
		})
	}

	return swBaseResponse
}

func computeStatWeights(swcr *proto.StatWeightsCalcRequest) *proto.StatWeightsResult {
	haveRefStat := false
	for _, statResult := range swcr.StatSimResults {
		if statResult.UnitStat == int32(swcr.EpReferenceStat) {
			haveRefStat = true
			break
		}
	}
	if !haveRefStat {
		return &proto.StatWeightsResult{ErrorResult: "No result for reference stat exists!"}
	}

	result := NewStatWeightsResult()
	for _, statResult := range swcr.StatSimResults {
		stat := stats.UnitStatFromIdx(int(statResult.UnitStat))

		baselinePlayer := swcr.BaseResult.RaidMetrics.Parties[0].Players[0]
		modPlayerLow := statResult.ResultLow.RaidMetrics.Parties[0].Players[0]
		modPlayerHigh := statResult.ResultHigh.RaidMetrics.Parties[0].Players[0]

		// Check for hard caps. Hard caps will have results identical to the baseline because RNG is fixed.
		// When we find a hard-capped stat, just skip it (will return 0).
		if modPlayerHigh.Dps.Avg == baselinePlayer.Dps.Avg && modPlayerHigh.Hps.Avg == baselinePlayer.Hps.Avg && modPlayerHigh.Tmi.Avg == baselinePlayer.Tmi.Avg {
			continue
		}

		calcWeightResults := func(baselineMetrics *proto.DistributionMetrics, modLowMetrics *proto.DistributionMetrics, modHighMetrics *proto.DistributionMetrics, weightResults *StatWeightValues) {
			var lo, hi aggregator
			for i := 0; i < len(baselineMetrics.AllValues); i++ {
				lo.add(modLowMetrics.AllValues[i] - baselineMetrics.AllValues[i])
			}
			lo.scale(1 / statResult.ModLow)
			for i := 0; i < len(baselineMetrics.AllValues); i++ {
				hi.add(modHighMetrics.AllValues[i] - baselineMetrics.AllValues[i])
			}
			hi.scale(1 / statResult.ModHigh)

			mean, stdev := lo.merge(&hi).meanAndStdDev()
			weightResults.Weights.AddStat(stat, mean)
			weightResults.WeightsStdev.AddStat(stat, stdev)
		}

		calcWeightResults(baselinePlayer.Dps, modPlayerLow.Dps, modPlayerHigh.Dps, &result.Dps)
		calcWeightResults(baselinePlayer.Hps, modPlayerLow.Hps, modPlayerHigh.Hps, &result.Hps)
		calcWeightResults(baselinePlayer.Threat, modPlayerLow.Threat, modPlayerHigh.Threat, &result.Tps)
		calcWeightResults(baselinePlayer.Dtps, modPlayerLow.Dtps, modPlayerHigh.Dtps, &result.Dtps)
		calcWeightResults(baselinePlayer.Tmi, modPlayerLow.Tmi, modPlayerHigh.Tmi, &result.Tmi)
		meanLow := (modPlayerLow.ChanceOfDeath - baselinePlayer.ChanceOfDeath) / statResult.ModLow
		meanHigh := (modPlayerHigh.ChanceOfDeath - baselinePlayer.ChanceOfDeath) / statResult.ModHigh
		result.PDeath.Weights.AddStat(stat, (meanLow+meanHigh)/2)
		result.PDeath.WeightsStdev.AddStat(stat, 0)
	}

	referenceStat := stats.Stat(swcr.EpReferenceStat)

	// Compute EP results.
	for _, statData := range swcr.StatSimResults {
		stat := stats.UnitStatFromIdx(int(statData.UnitStat))

		calcEpResults := func(weightResults *StatWeightValues, refStat stats.Stat) {
			if weightResults.Weights.Stats[refStat] == 0 {
				return
			}
			mean := weightResults.Weights.Get(stat) / weightResults.Weights.Stats[refStat]
			stdev := weightResults.WeightsStdev.Get(stat) / math.Abs(weightResults.Weights.Stats[refStat])
			weightResults.EpValues.AddStat(stat, mean)
			weightResults.EpValuesStdev.AddStat(stat, stdev)
		}

		calcEpResults(&result.Dps, referenceStat)
		calcEpResults(&result.Hps, referenceStat)
		calcEpResults(&result.Tps, referenceStat)
		calcEpResults(&result.Dtps, DTPSReferenceStat)
		calcEpResults(&result.Tmi, DTPSReferenceStat)
		calcEpResults(&result.PDeath, DTPSReferenceStat)
	}

	return result.ToProto()
}

// Run stat weight sims and compute weights.
func runStatWeights(request *proto.StatWeightsRequest, progress chan *proto.ProgressMetrics, signals simsignals.Signals) *proto.StatWeightsResult {
	requestData := buildStatWeightRequests(request)

	var iterationsTotal int32 = requestData.BaseRequest.SimOptions.Iterations
	var iterationsDone int32 = 0
	var simsTotal int32 = 1
	var simsCompleted int32 = 0

	for _, reqData := range requestData.StatSimRequests {
		iterationsTotal += reqData.RequestLow.SimOptions.Iterations
		iterationsTotal += reqData.RequestHigh.SimOptions.Iterations
		simsTotal += 2
	}

	waitForResult := func(srcProgressChannel chan *proto.ProgressMetrics) *proto.RaidSimResult {
		var lastCompleted int32 = 0
		for metrics := range srcProgressChannel {
			iterationsDone += metrics.CompletedIterations - lastCompleted
			lastCompleted = metrics.CompletedIterations

			if progress != nil {
				progress <- &proto.ProgressMetrics{
					TotalIterations:     iterationsTotal,
					CompletedIterations: iterationsDone,
					CompletedSims:       simsCompleted,
					TotalSims:           simsTotal,
				}
			}

			if metrics.FinalRaidResult != nil {
				simsCompleted++
				return metrics.FinalRaidResult
			}
		}
		return nil
	}

	simFunc := runSimConcurrent
	if IsRunningInWasm() {
		simFunc = RunSim
	}

	baseProgress := make(chan *proto.ProgressMetrics, 100)
	go simFunc(requestData.BaseRequest, baseProgress, signals)
	baselineResult := waitForResult(baseProgress)
	if baselineResult.ErrorResult != "" {
		return &proto.StatWeightsResult{ErrorResult: baselineResult.ErrorResult}
	}

	statResults := []*proto.StatWeightStatResultData{}

	for _, reqData := range requestData.StatSimRequests {
		lowProgress := make(chan *proto.ProgressMetrics, 100)
		go simFunc(reqData.RequestLow, lowProgress, signals)
		lowRes := waitForResult(lowProgress)
		if lowRes.ErrorResult != "" {
			return &proto.StatWeightsResult{ErrorResult: lowRes.ErrorResult}
		}

		highProgress := make(chan *proto.ProgressMetrics, 100)
		go simFunc(reqData.RequestHigh, highProgress, signals)
		highRes := waitForResult(highProgress)
		if highRes.ErrorResult != "" {
			return &proto.StatWeightsResult{ErrorResult: highRes.ErrorResult}
		}

		statResults = append(statResults, &proto.StatWeightStatResultData{
			UnitStat:   reqData.UnitStat,
			ResultLow:  lowRes,
			ResultHigh: highRes,
			ModLow:     reqData.ModLow,
			ModHigh:    reqData.ModHigh,
		})
	}

	return computeStatWeights(&proto.StatWeightsCalcRequest{
		BaseResult:      baselineResult,
		EpReferenceStat: requestData.EpReferenceStat,
		StatSimResults:  statResults,
	})
}
