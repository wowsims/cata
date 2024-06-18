package simsignals_test

import (
	"testing"
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/simsignals"
	"github.com/wowsims/cata/sim/warrior/arms"
)

func TestAbort(t *testing.T) {
	arms.RegisterArmsWarrior()

	player := &proto.Player{
		Name:      "John",
		Race:      proto.Race_RaceOrc,
		Class:     proto.Class_ClassWarrior,
		Equipment: core.GetGearSet("../../../ui/warrior/arms/gear_sets", "p1_arms_bis").GearSet,
		Rotation:  &proto.APLRotation{},
		Consumes:  &proto.Consumes{},
		Spec: &proto.Player_ArmsWarrior{
			ArmsWarrior: &proto.ArmsWarrior{
				Options: &proto.ArmsWarrior_Options{
					ClassOptions: &proto.WarriorOptions{
						StartingRage:       50,
						UseShatteringThrow: true,
						Shout:              proto.WarriorShout_WarriorShoutBattle,
					},
				},
			},
		},
		Glyphs:        &proto.Glyphs{},
		TalentsString: "",
		Buffs:         &proto.IndividualBuffs{},
	}

	rsr := &proto.RaidSimRequest{
		RequestId: "uniqueidlol",
		Raid:      core.SinglePlayerRaidProto(player, core.FullPartyBuffs, core.FullRaidBuffs, core.FullDebuffs),
		Encounter: &proto.Encounter{
			Duration: 300,
			Targets: []*proto.Target{
				core.NewDefaultTarget(),
			},
		},
		SimOptions: &proto.SimOptions{
			Iterations: 33333,
			IsTest:     true,
			RandomSeed: 123,
		},
	}

	t.Run("RunRaidSimAsync", func(t *testing.T) {
		progress := make(chan *proto.ProgressMetrics, 10)
		core.RunRaidSimAsync(rsr, progress)
		simsignals.AbortById(rsr.RequestId)
		simsignals.AbortById(rsr.RequestId)
		simsignals.AbortById(rsr.RequestId)
		for {
			msg := <-progress
			if msg.FinalRaidResult != nil {
				if msg.FinalRaidResult.Error == nil || msg.FinalRaidResult.Error.Type != proto.ErrorOutcomeType_ErrorOutcomeAborted {
					t.Fatal("Sim did not abort!")
				}
				return
			}
		}
	})

	t.Run("RunRaidSimAsyncMultiManual", func(t *testing.T) {
		rsr.RequestId += "x"
		var conc int32 = 2
		progress := make([]chan *proto.ProgressMetrics, conc)
		rsrSplits := core.SplitSimRequestForConcurrency(rsr, conc)
		for i, rsrSplit := range rsrSplits.Requests {
			progress[i] = make(chan *proto.ProgressMetrics, 10)
			core.RunRaidSimAsync(rsrSplit, progress[i])
			simsignals.AbortById(rsrSplit.RequestId)
		}

		running := conc

		for {
			for i, p := range progress {
				msg, ok := <-p
				if ok && msg.FinalRaidResult != nil {
					if msg.FinalRaidResult.Error == nil || msg.FinalRaidResult.Error.Type != proto.ErrorOutcomeType_ErrorOutcomeAborted {
						t.Fatalf("Sim instance %d did not abort!", i)
					}
					running--
					if running == 0 {
						return
					}
				}
			}
		}
	})

	t.Run("RunRaidSimConcurrentAsync", func(t *testing.T) {
		rsr.RequestId += "x"
		progress := make(chan *proto.ProgressMetrics, 10)
		core.RunRaidSimConcurrentAsync(rsr, progress)
		simsignals.AbortById(rsr.RequestId)
		for {
			msg := <-progress
			if msg.FinalRaidResult != nil {
				if msg.FinalRaidResult.Error == nil || msg.FinalRaidResult.Error.Type != proto.ErrorOutcomeType_ErrorOutcomeAborted {
					t.Fatal("Sim did not abort!")
				}
				return
			}
		}
	})

	t.Run("RunRaidSimConcurrentAsync-Delayed", func(t *testing.T) {
		rsr.RequestId += "x"
		progress := make(chan *proto.ProgressMetrics, 10)
		core.RunRaidSimConcurrentAsync(rsr, progress)
		go func() {
			time.Sleep(time.Second)
			simsignals.AbortById(rsr.RequestId)
		}()
		for {
			msg := <-progress
			if msg.FinalRaidResult != nil {
				if msg.FinalRaidResult.Error == nil || msg.FinalRaidResult.Error.Type != proto.ErrorOutcomeType_ErrorOutcomeAborted {
					t.Fatal("Sim did not abort!")
				}
				return
			}
		}
	})

	t.Run("StatWeightsAsync", func(t *testing.T) {
		swr := &proto.StatWeightsRequest{
			RequestId:  "lel",
			Player:     player,
			RaidBuffs:  core.FullRaidBuffs,
			PartyBuffs: core.FullPartyBuffs,
			Debuffs:    core.FullDebuffs,
			Encounter:  core.MakeSingleTargetEncounter(0),
			SimOptions: core.StatWeightsDefaultSimTestOptions,
			Tanks:      make([]*proto.UnitReference, 0),

			StatsToWeigh: []proto.Stat{
				proto.Stat_StatAgility,
				proto.Stat_StatAttackPower,
				proto.Stat_StatMastery,
				proto.Stat_StatMeleeHit,
				proto.Stat_StatExpertise,
			},
			EpReferenceStat: proto.Stat_StatAttackPower,
		}
		swr.SimOptions.Iterations = 9999

		progress := make(chan *proto.ProgressMetrics, 10)
		core.StatWeightsAsync(swr, progress)

		go func() {
			time.Sleep(time.Second)
			simsignals.AbortById(swr.RequestId)
		}()

		for msg := range progress {
			if msg.FinalWeightResult != nil {
				if msg.FinalWeightResult.Error == nil || msg.FinalWeightResult.Error.Type != proto.ErrorOutcomeType_ErrorOutcomeAborted {
					t.Fatalf("Sim did not abort!")
				}
				return
			}
		}
	})

	t.Run("RunBulkSimAsync", func(t *testing.T) {
		bsr := &proto.BulkSimRequest{
			RequestId:    "lel",
			BaseSettings: rsr,
			BulkSettings: &proto.BulkSettings{
				Combinations:       true,
				Items:              []*proto.ItemSpec{{Id: 77949}, {Id: 55068}},
				IterationsPerCombo: 9999,
				FastMode:           false,
			},
		}

		progress := make(chan *proto.ProgressMetrics, 10)
		core.RunBulkSimAsync(bsr, progress)

		go func() {
			time.Sleep(time.Second)
			simsignals.AbortById(bsr.RequestId)
		}()

		for msg := range progress {
			if msg.FinalBulkResult != nil {
				if msg.FinalBulkResult.Error == nil || msg.FinalBulkResult.Error.Type != proto.ErrorOutcomeType_ErrorOutcomeAborted {
					t.Fatalf("Sim did not abort!")
				}
				return
			}
		}
	})
}
