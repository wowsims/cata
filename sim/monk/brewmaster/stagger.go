package brewmaster

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/monk"
)

func (bm *BrewmasterMonk) registerStagger() {
	actionId := core.ActionID{SpellID: 124255}

	bm.Stagger = bm.RegisterSpell(core.SpellConfig{
		ActionID:         actionId,
		SpellSchool:      core.SpellSchoolPhysical,
		ProcMask:         core.ProcMaskSpellProc,
		ClassSpellMask:   monk.MonkSpellStagger,
		Flags:            core.SpellFlagNoMetrics | core.SpellFlagIgnoreModifiers | core.SpellFlagNoSpellMods | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Hot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Stagger" + bm.Label,
				ActionID:  actionId.WithTag(1),
				MaxStacks: math.MaxInt32,
			},
			SelfOnly:            true,
			NumberOfTicks:       10,
			TickLength:          1 * time.Second,
			AffectedByCastSpeed: false,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				if !dot.Aura.IsActive() {
					return
				}

				damage := max(0, dot.SnapshotBaseDamage)
				target.RemoveHealth(sim, damage)

				if sim.Log != nil && dot.Aura.IsActive() {
					bm.Log(sim, "[DEBUG] Stagger ticked for %0.0f Damage", damage)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.SelfHot().Apply(sim)
		},
	})

	bm.RefreshStagger = func(sim *core.Simulation, target *core.Unit, damagePerTick float64) {
		dot := bm.Stagger.SelfHot()
		if damagePerTick <= 0 {
			dot.Deactivate(sim)
			if sim.Log != nil {
				bm.Log(sim, "[DEBUG] Stagger reset")
			}
		} else {
			oldDamagePerTick := dot.SnapshotBaseDamage
			dot.SnapshotBaseDamage = damagePerTick
			bm.Stagger.Cast(sim, target)
			newStaggerValue := int32(damagePerTick)
			if newStaggerValue < 0 {
				panic("Stagger is above 2.147 billion. Please check your Rotation/Encounter settings.")
			}

			dot.Aura.SetStacks(sim, int32(damagePerTick))

			if sim.Log != nil && dot.Aura.IsActive() {
				bm.Log(sim, "[DEBUG] Stagger tick refreshed from: %0.0f -> %0.0f", oldDamagePerTick, damagePerTick)
			}
		}

	}

	bm.AddDynamicDamageTakenModifier(func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		if !result.Landed() {
			return
		}

		avertHarmIsActive := bm.AvertHarmAura.IsActive()
		// By default Stagger only works with physical abilities
		// unless Avert Harm is active, then Magic abilities will be staggered as well (without the 20% bonus)
		if !spell.SpellSchool.Matches(core.SpellSchoolPhysical) && !avertHarmIsActive {
			return
		}

		target := result.Target
		dot := bm.Stagger.SelfHot()
		outstandingDamage := dot.OutstandingDmg()

		// Avert Harm will only gain 20% Stagger from Melee abilities (non-auto attacks)
		avertHarmMultiplier := core.TernaryFloat64(avertHarmIsActive && !spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) && spell.SpellSchool.Matches(core.SpellSchoolPhysical), 0.2, 0)
		shuffleMultiplier := core.TernaryFloat64(bm.ShuffleAura.IsActive(), 0.2, 0)
		fortifyingBrewMultiplier := core.TernaryFloat64(bm.FortifyingBrewAura.IsActive(), 0.2, 0)
		t15Brewmaster2P := core.TernaryFloat64(bm.T15Brewmaster2P != nil && bm.T15Brewmaster2P.IsActive(), 0.06, 0)

		staggerMultiplier := min(1, bm.GetMasteryBonus()) + shuffleMultiplier + fortifyingBrewMultiplier + avertHarmMultiplier + t15Brewmaster2P
		staggeredDamage := result.Damage * staggerMultiplier
		result.Damage -= staggeredDamage

		newOutstandingDamage := outstandingDamage + staggeredDamage
		newTickCount := dot.BaseTickCount
		damagePerTick := newOutstandingDamage / float64(newTickCount)

		if sim.Log != nil {
			bm.Log(sim, "[DEBUG] Stagger (%0.2f%%) mitigated %0.0f Damage - New outstanding Damage %0.0f", staggerMultiplier*100, staggeredDamage, newOutstandingDamage)
		}

		bm.RefreshStagger(sim, target, damagePerTick)

		// Dampen Harm
		// This is applied after other DR and Stagger
		if bm.DampenHarmAura.IsActive() && result.Damage > result.Target.MaxHealth()*0.2 {
			bm.DampenHarmAura.RemoveStack(sim)
			result.Damage /= 2
		}
	})

}
