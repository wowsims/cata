package death_knight

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (dk *DeathKnight) DiseasesAreActive(target *core.Unit) bool {
	return dk.FrostFeverSpell.Dot(target).IsActive() || dk.BloodPlagueSpell.Dot(target).IsActive()
}

func (dk *DeathKnight) CountActiveDiseases(target *core.Unit) float64 {
	count := 0
	if dk.FrostFeverSpell.Dot(target).IsActive() {
		count++
	}
	if dk.BloodPlagueSpell.Dot(target).IsActive() {
		count++
	}
	if dk.Talents.EbonPlaguebringer > 0 && dk.EbonPlagueAura.Get(target).IsActive() {
		count++
	}
	return float64(count)
}

func (dk *DeathKnight) registerOutbreak() {
	dk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 77575},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellOutbreak,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Minute,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				dk.FrostFeverSpell.Cast(sim, target)
				dk.BloodPlagueSpell.Cast(sim, target)
			}
		},
	})
}

func (dk *DeathKnight) registerFrostFever() {
	extraEffectAura := dk.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.FrostFeverAura(&dk.Unit, target, dk.Talents.BrittleBones)
	})

	dk.FrostFeverSpell = dk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 55095},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagDisease,
		ClassSpellMask: DeathKnightSpellFrostFever,

		DamageMultiplier: 1.15,
		CritMultiplier:   dk.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "FrostFever" + dk.Label,
				Tag:   "FrostFever",
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					extraEffectAura.Get(aura.Unit).Activate(sim)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					extraEffectAura.Get(aura.Unit).Deactivate(sim)
				},
			},
			NumberOfTicks: 7,
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, dot.Spell.MeleeAttackPower()*0.055+0.31999999285*core.CharacterLevel)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dot := spell.Dot(target)
			dot.Apply(sim)
		},

		RelatedAuras: []core.AuraArray{extraEffectAura},
	})
}

func (dk *DeathKnight) registerBloodPlague() {
	dk.BloodPlagueSpell = dk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 55078},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagDisease,
		ClassSpellMask: DeathKnightSpellBloodPlague,

		DamageMultiplier: 1.15,
		CritMultiplier:   dk.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "BloodPlague" + dk.Label,
				Tag:   "BloodPlague",
			},
			NumberOfTicks: 7,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, dot.Spell.MeleeAttackPower()*0.055+0.3939999938*core.CharacterLevel)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})
}
