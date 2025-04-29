package demonology

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

func (demonology *DemonologyWarlock) CurseOfGuldanDebuffAura(target *core.Unit) *core.Aura {
	// TODO: the talent tooltip says this applies to "any Warlock demons". It's unclear if this means
	// any warlock pet or just the ones belonging to the caster
	return target.GetOrRegisterAura(core.Aura{
		Label:    "Curse of Guldan-" + demonology.Label,
		ActionID: core.ActionID{SpellID: 86000},
		Duration: 15 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, pet := range demonology.Pets {
				pet.AttackTables[aura.Unit.UnitIndex].BonusSpellCritPercent += 10
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, pet := range demonology.Pets {
				pet.AttackTables[aura.Unit.UnitIndex].BonusSpellCritPercent -= 10
			}
		},
	})
}

func (demonology *DemonologyWarlock) registerHandOfGuldan() {
	if !demonology.Talents.HandOfGuldan {
		return
	}

	curseOfGuldanAuras := demonology.NewEnemyAuraArray(demonology.CurseOfGuldanDebuffAura)
	demonology.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 71521},
		SpellSchool:    core.SpellSchoolFire | core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: warlock.WarlockSpellHandOfGuldan,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 7},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 2 * time.Second,
			},
			CD: core.Cooldown{
				Timer:    demonology.NewTimer(),
				Duration: 12 * time.Second,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   demonology.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 0.96799999475,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := demonology.CalcAndRollDamageRange(sim, 1.59300005436, 0.16599999368)
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				for _, target := range sim.Encounter.TargetUnits {
					curseOfGuldanAuras.Get(target).Activate(sim)
				}
			}
		},
	})
}
