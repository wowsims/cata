package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

// TODO: Check damage and coefficients
func (warlock *Warlock) registerDrainSoulSpell() {
	soulSiphonMultiplier := 0.03 * float64(warlock.Talents.SoulSiphon)

	calcSoulSiphonMult := func(target *core.Unit) float64 {
		// missing: death coil
		afflictionDots := WarlockSpellUnstableAffliction | WarlockSpellCorruption |
			WarlockSpellSeedOfCorruption | WarlockSpellBaneOfAgony | WarlockSpellBaneOfDoom

		afflictionAuras := WarlockSpellHaunt | WarlockSpellCurseOfElements | WarlockSpellCurseOfWeakness |
			WarlockSpellCurseOfTongues

		numActive := 0
		for _, spell := range warlock.Spellbook {
			if (spell.Matches(afflictionDots) && spell.Dot(target).IsActive()) ||
				(spell.Matches(afflictionAuras) && spell.RelatedAuras[0].Get(target).IsActive()) {
				numActive++
			}
		}
		return 1.0 + float64(min(3, numActive))*soulSiphonMultiplier
	}

	ds := warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 1120},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagChanneled | core.SpellFlagHauntSE | core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellDrainSoul,

		ManaCost: core.ManaCostOptions{BaseCost: 0.14},
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
				baseDmg := 385/5 + 0.378*dot.Spell.SpellPower()
				dot.SnapshotBaseDamage = baseDmg * calcSoulSiphonMult(target)
				dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex], true)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--
				dot := spell.Dot(target)
				dot.Apply(sim)
				dot.UpdateExpires(dot.ExpiresAt())
			}
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, spell.OutcomeExpectedMagicCrit)
			} else {
				baseDmg := (385/5 + 0.378*spell.SpellPower()) * calcSoulSiphonMult(target)
				return spell.CalcPeriodicDamage(sim, target, baseDmg, spell.OutcomeExpectedMagicCrit)
			}
		},
	})

	drainSoulExecuteAura := warlock.RegisterAura(core.Aura{
		Label:    "Drain Soul Execute",
		ActionID: core.ActionID{SpellID: 1120},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			ds.DamageMultiplier *= 2.0
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			ds.DamageMultiplier /= 2.0
		},
	})

	warlock.RegisterResetEffect(func(sim *core.Simulation) {
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int32) {
			if isExecute == 25 {
				drainSoulExecuteAura.Activate(sim)
			}
		})
	})
}
