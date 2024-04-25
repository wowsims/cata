package demonology

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/warlock"
)

func (demonology *DemonologyWarlock) CurseOfGuldanDebuffAura(target *core.Unit) *core.Aura {
	return target.GetOrRegisterAura(core.Aura{
		Label:    "CurseOfGuldan-" + demonology.Label,
		ActionID: core.ActionID{SpellID: 86000},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			//TODO: Implement Crit rating for pet vs this target only
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			//TODO: Implement Crit rating for pet vs this target only
		},
	})
}

// TODO: Since we are assuming hand of guldan will hit all targets can we just give the pets 10% crit for the duration?
func (demonology *DemonologyWarlock) registerHandOfGuldanSpell() {
	if !demonology.Talents.Haunt {
		return
	}

	demonology.RegisterAura(core.Aura{
		Label:    "Hand of Guldan",
		ActionID: core.ActionID{SpellID: 47197},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
		},
	})

	//TODO: Damage is shadow flame... How to make this shadow and fire schools?
	demonology.HandOfGuldan = demonology.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 71521},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: warlock.WarlockSpellHandOfGuldan,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.07,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 2,
			},
			CD: core.Cooldown{
				Timer:    demonology.NewTimer(),
				Duration: time.Second * 12,
			},
		},

		CritMultiplier:   demonology.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 0.968,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealDamage(sim, target, 1793, spell.OutcomeMagicHitAndCrit)
			if !result.Landed() {
				return
			}
		},
	})
}
