package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

// TODO: Check damage and coefficients
func (warlock *Warlock) registerDrainSoulSpell() {
	soulSiphonMultiplier := 0.03 * float64(warlock.Talents.SoulSiphon)

	calcSoulSiphonMult := func(target *core.Unit) float64 {
		auras := []*core.Aura{
			warlock.UnstableAffliction.Dot(target).Aura,
			warlock.Corruption.Dot(target).Aura,
			warlock.Seed.Dot(target).Aura,
			warlock.BaneOfAgony.Dot(target).Aura,
			warlock.BaneOfDoom.Dot(target).Aura,
			warlock.CurseOfElementsAuras.Get(target),
			warlock.CurseOfWeaknessAuras.Get(target),
			warlock.CurseOfTonguesAuras.Get(target),
			warlock.ShadowEmbraceDebuffAura(target),
			// missing: death coil
		}
		if warlock.HauntDebuffAuras != nil {
			auras = append(auras, warlock.HauntDebuffAuras.Get(target))
		}
		numActive := 0
		for _, aura := range auras {
			if aura.IsActive() {
				numActive++
			}
		}
		return 1.0 + float64(min(3, numActive))*soulSiphonMultiplier
	}

	warlock.DrainSoul = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 1120},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagChanneled | core.SpellFlagHauntSE | core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellDrainSoul,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.14,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
				// ChannelTime: channelTime,
			},
		},

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
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex], true)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickCounted)
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
				return dot.CalcSnapshotDamage(sim, target, spell.OutcomeExpectedMagicAlwaysHit)
			} else {
				baseDmg := (385/5 + 0.378*spell.SpellPower()) * calcSoulSiphonMult(target)
				return spell.CalcPeriodicDamage(sim, target, baseDmg, spell.OutcomeExpectedMagicAlwaysHit)
			}
		},
	})

	drainSoulExecuteAura := warlock.RegisterAura(core.Aura{
		Label:    "Drain Soul Execute",
		ActionID: core.ActionID{SpellID: 1120},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.DrainSoul.DamageMultiplier *= 2.0
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.DrainSoul.DamageMultiplier /= 2.0
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
