package brewmaster

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/monk"
)

func (bm *BrewmasterMonk) registerStagger() {
	actionId := core.ActionID{SpellID: 124255}

	staggerSpell := bm.RegisterSpell(core.SpellConfig{
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
				newHealth := max(0, dot.SnapshotBaseDamage)
				target.RemoveHealth(sim, newHealth)
				dot.Aura.Activate(sim)
				dot.Aura.SetStacks(sim, int32(newHealth))
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.SelfHot().Apply(sim)
		},
	})

	refreshStagger := func(sim *core.Simulation, target *core.Unit, damagePerTick float64) {
		dot := staggerSpell.SelfHot()
		dot.SnapshotBaseDamage = damagePerTick
		staggerSpell.Cast(sim, target)
	}

	bm.AddDynamicDamageTakenModifier(func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		if spell.SpellSchool != core.SpellSchoolPhysical {
			return
		}

		target := result.Target
		dot := staggerSpell.SelfHot()
		outstandingDamage := dot.OutstandingDmg()
		staggerMultiplier := min(1, bm.GetMasteryBonus()) + core.TernaryFloat64(bm.ShuffleAura.IsActive(), 0.2, 0)
		newStaggeredDamage := result.Damage * staggerMultiplier
		result.Damage -= newStaggeredDamage

		totalDamage := outstandingDamage + newStaggeredDamage
		newTickCount := dot.BaseTickCount + core.TernaryInt32(dot.IsActive(), 1, 0)
		damagePerTick := totalDamage / float64(newTickCount)

		refreshStagger(sim, target, damagePerTick)
	})

}
