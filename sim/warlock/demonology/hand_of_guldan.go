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

// TODO: Curse
func (demonology *DemonologyWarlock) registerHandOfGuldanSpell() {
	if !demonology.Talents.HandOfGuldan {
		return
	}

	// TODO: If you switch pets or summon a new one they won't have the attack tables do not switch when I think curse of guldan would apply to any active and future pets
	// When the curse expires will it be taken away from the current active pet or the pet we originally assigned it to?
	// demonology.ActivePet.CurseOfGuldanDebuffs = demonology.ActivePet.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
	// 	return target.GetOrRegisterAura(core.Aura{
	// 		Label:    "Curse of Guldan",
	// 		ActionID: core.ActionID{SpellID: 86000},
	// 		Duration: time.Second * 15,
	// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
	// 			demonology.ActivePet.AttackTables[aura.Unit.UnitIndex].BonusCritRating += 10.0 * core.CritRatingPerCritChance
	// 		},
	// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
	// 			demonology.ActivePet.AttackTables[aura.Unit.UnitIndex].BonusCritRating -= 10.0 * core.CritRatingPerCritChance
	// 		},
	// 	})
	// })

	demonology.HandOfGuldan = demonology.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 71521},
		SpellSchool:    core.SpellSchoolFire | core.SpellSchoolShadow,
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

			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := demonology.CalcBaseDamageWithVariance(sim, 1.593, 0.166)
				result := spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)

				if result.Landed() {
					//aura := demonology.ActivePet.CurseOfGuldanDebuffs.Get(target)
					//aura.Activate(sim)
				}
			}
		},
	})
}
