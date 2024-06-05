package simsignals_test

import (
	"testing"
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/simsignals"
	"github.com/wowsims/cata/sim/warrior/arms"
)

func getTestRsr() *proto.RaidSimRequest {
	arms.RegisterArmsWarrior()
	return &proto.RaidSimRequest{
		RequestId: "uniqueidlol",
		Raid: core.SinglePlayerRaidProto(
			&proto.Player{
				Race:      proto.Race_RaceOrc,
				Class:     proto.Class_ClassWarrior,
				Equipment: &proto.EquipmentSpec{},
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
			},
			core.FullPartyBuffs,
			core.FullRaidBuffs,
			core.FullDebuffs),
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
}

func TestAbort(t *testing.T) {
	rsr := getTestRsr()

	t.Run("RunRaidSimAsync", func(t *testing.T) {
		progress := make(chan *proto.ProgressMetrics, 10)
		core.RunRaidSimAsync(rsr, progress)
		simsignals.AbortById(rsr.RequestId)
		simsignals.AbortById(rsr.RequestId)
		simsignals.AbortById(rsr.RequestId)
		for {
			msg := <-progress
			if msg.FinalRaidResult != nil {
				if msg.FinalRaidResult.ErrorResult != "aborted" {
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
					if msg.FinalRaidResult.ErrorResult != "aborted" {
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
				if msg.FinalRaidResult.ErrorResult != "aborted" {
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
				if msg.FinalRaidResult.ErrorResult != "aborted" {
					t.Fatal("Sim did not abort!")
				}
				return
			}
		}
	})

	/* t.Run("StatWeightsAsync", func(t *testing.T) {
		swr := &proto.StatWeightsRequest{
			Id:         "lel",
			Player:     getTestPlayerFeralCat(),
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
			core.AbortSimById(swr.Id)
		}()

		for msg := range progress {
			if msg.FinalWeightResult != nil {
				if msg.FinalWeightResult.ErrorResult != "aborted" {
					t.Fatalf("Sim did not abort! %s", msg.FinalWeightResult.ErrorResult)
				}
				return
			}
		}
	}) */
}
