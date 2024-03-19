package arms

import (
	"math"
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

const (
	StrikesOfOpportunityHitID int32 = 76858
	StrikesOfOpportunityTag   int32 = 3 // 1 and 2 are MH and OH hits, respectively
)

func (war *ArmsWarrior) CalcMasteryPoints() float64 {
	return math.Floor(war.GetStat(stats.Mastery) / core.MasteryRatingPerMasteryPoint)
}

func (war *ArmsWarrior) GetMasteryProcChance() float64 {
	return (17.6 + 2.2*war.CalcMasteryPoints()) / 100
}

func (war *ArmsWarrior) RegisterMastery() {
	// TODO:
	//	test what things the extra attack can proc
	//	does the extra attack use the same hit table
	//
	// 4.3.3 simcraft implements SoO as a standard autoattack
	procAttackConfig := *war.AutoAttacks.MHConfig()
	procAttackConfig.ActionID = core.ActionID{SpellID: StrikesOfOpportunityHitID, Tag: procAttackConfig.ActionID.Tag}
	procAttack := war.RegisterSpell(procAttackConfig)

	icd := core.Cooldown{
		Timer:    war.NewTimer(),
		Duration: time.Millisecond * 500, // From 4.3.3 simcraft, needs to be tested
	}

	war.RegisterAura(core.Aura{
		Label:    "Strikes of Opportunity",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ActionID.IsOtherAction(proto.OtherAction_OtherActionAttack) {
				return
			}

			if !icd.IsReady(sim) {
				return
			}

			proc := war.GetMasteryProcChance()
			if sim.RandomFloat(aura.Label) < proc {
				icd.Use(sim)
				procAttack.Cast(sim, result.Target)
			}
		},
	})
}
