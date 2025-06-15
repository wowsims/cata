package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (druid *Druid) registerTranquilityCD() {
	// Only implemented for individual sims at present, so we don't care
	// about the target selection logic.
	targets := druid.Env.Raid.GetFirstNPlayersOrPets(5)

	// First register the stacking HoT spell that gets triggered by the main channel.
	tranquilityHot := druid.RegisterSpell(Humanoid|Tree, core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 44203},
		SpellSchool:      core.SpellSchoolNature,
		ProcMask:         core.ProcMaskSpellHealing,
		Flags:            core.SpellFlagHelpful | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   druid.DefaultCritMultiplier(),

		// TODO: Healing value calculations are very likely incorrect and will need a closer look if we care about
		// modeling the actual healing output from the spell. Right now this is just a placeholder for pre-pull use in
		// DPS sims to proc healing trinkets.
		Hot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Tranquility (HoT)",
				MaxStacks: 3,
			},

			NumberOfTicks:        4,
			TickLength:           time.Second * 2,
			AffectedByCastSpeed:  true,
			HasteReducesDuration: true,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = 0.068 * dot.Spell.HealingPower(target) * float64(dot.Aura.GetStacks())
				dot.SnapshotAttackerMultiplier = dot.CasterPeriodicHealingMultiplier()
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotHealing(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			hot := spell.Hot(target)

			if !hot.IsActive() {
				hot.Apply(sim)
			}

			hot.AddStack(sim)
			hot.TakeSnapshot(sim, false)
		},
	})

	// Then register the primary channel spell.
	druid.RegisterSpell(Humanoid|Tree, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 740},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful | core.SpellFlagAPL | core.SpellFlagChanneled,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 32,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},

			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Minute * 8,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   druid.DefaultCritMultiplier(),

		// TODO: Healing value calculations are very likely incorrect and will need a closer look if we
		// care about modeling the actual healing output from the spell. Right now this is just a
		// placeholder for pre-pull use in DPS sims to proc healing trinkets.
		Hot: core.DotConfig{
			Aura: core.Aura{
				Label: "Tranquility",
			},

			NumberOfTicks:        4,
			TickLength:           time.Second * 2,
			AffectedByCastSpeed:  true,
			HasteReducesDuration: true,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = 3882. + 0.398*dot.Spell.HealingPower(target)
				dot.SnapshotAttackerMultiplier = dot.Spell.CasterHealingMultiplier()
			},

			OnTick: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range targets {
					dot.CalcAndDealPeriodicSnapshotHealing(sim, aoeTarget, dot.OutcomeTick)
					tranquilityHot.Cast(sim, aoeTarget)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.Hot(&druid.Unit).Apply(sim)
		},
	})
}
