package warlock

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (warlock *Warlock) calcSoulSiphonMult(target *core.Unit) float64 {
	soulSiphonMultiplier := 0.03 * float64(warlock.Talents.SoulSiphon)

	// missing: death coil
	afflictionDots := WarlockSpellUnstableAffliction | WarlockSpellCorruption |
		WarlockSpellSeedOfCorruption | WarlockSpellBaneOfAgony | WarlockSpellBaneOfDoom

	afflictionAuras := WarlockSpellHaunt | WarlockSpellCurseOfElements | WarlockSpellCurseOfWeakness |
		WarlockSpellCurseOfTongues

	numActive := 0
	for _, spell := range warlock.Spellbook {
		if (spell.Matches(afflictionDots) && spell.Dot(target).IsActive()) ||
			(spell.Matches(afflictionAuras) && spell.RelatedAuraArrays.AnyActive(target)) {
			numActive++
		}
	}
	return 1.0 + float64(min(3, numActive))*soulSiphonMultiplier
}

// TODO: Check damage and coefficients
func (warlock *Warlock) registerDrainSoul() {
	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 1120},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagChanneled | core.SpellFlagHauntSE | core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellDrainSoul,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 14},
		Cast:     core.CastConfig{DefaultCast: core.Cast{GCD: core.GCDDefault}},

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		CritMultiplier:           warlock.DefaultSpellCritMultiplier(),
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Drain Soul",
			},
			NumberOfTicks:       5,
			TickLength:          3 * time.Second,
			AffectedByCastSpeed: true,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				baseDmg := warlock.CalcScalingSpellDmg(0.07999999821) + 0.37799999118*dot.Spell.SpellPower()
				dot.SnapshotBaseDamage = baseDmg * warlock.calcSoulSiphonMult(target)
				dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex], true)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)
			if result.Landed() {
				dot := spell.Dot(target)
				dot.Apply(sim)
				dot.UpdateExpires(dot.ExpiresAt())
				spell.DealOutcome(sim, result)
			}
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, spell.OutcomeExpectedMagicCrit)
			} else {
				baseDmg := warlock.CalcScalingSpellDmg(0.07999999821) + 0.37799999118*spell.SpellPower()
				baseDmg *= warlock.calcSoulSiphonMult(target)
				return spell.CalcPeriodicDamage(sim, target, baseDmg, spell.OutcomeExpectedMagicCrit)
			}
		},
	})

	executeMod := warlock.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Pct,
		ClassMask: WarlockSpellDrainSoul,
		// we have to correct for death's embrace here, since they stack additively together,
		// but multiplicatively with everything else
		FloatValue: 1.0 / (1 + 0.04*float64(warlock.Talents.DeathsEmbrace)),
	})

	warlock.RegisterResetEffect(func(sim *core.Simulation) {
		executeMod.Deactivate()
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int32) {
			if isExecute == 25 {
				executeMod.Activate()
			}
		})
	})
}
