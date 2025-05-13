package paladin

import (
	"math"

	"github.com/wowsims/mop/sim/core"
)

func (paladin *Paladin) registerSealOfInsight() {
	actionID := core.ActionID{SpellID: 20167}
	healthMetrics := paladin.NewHealthMetrics(actionID)
	manaMetrics := paladin.NewManaMetrics(actionID)
	// It's 4% of base mana per tick.
	manaPerTick := math.Round(0.04 * paladin.BaseMana)

	// Seal of Insight on-hit proc
	onSpecialOrSwingProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags: core.SpellFlagHelpful |
			core.SpellFlagNoLogs |
			core.SpellFlagNoMetrics |
			core.SpellFlagNoOnCastComplete,
		ClassSpellMask: SpellMaskSealOfInsight,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			heal := 0.15*spell.SpellPower() +
				0.15*spell.MeleeAttackPower()
			paladin.GainHealth(sim, heal, healthMetrics)
			paladin.AddMana(sim, manaPerTick, manaMetrics)
		},
	})

	dpm := paladin.AutoAttacks.NewPPMManager(15, core.ProcMaskMeleeMH)
	paladin.SealOfInsightAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Insight" + paladin.Label,
		Tag:      "Seal",
		ActionID: core.ActionID{SpellID: 20165},
		Duration: core.NeverExpires,
		Dpm:      dpm,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Don't proc on misses
			if !result.Landed() {
				return
			}

			// SoJ only procs on white hits, CS, TV and HoW
			if (spell.ProcMask&core.ProcMaskMeleeWhiteHit == 0 &&
				!spell.Matches(SpellMaskCanTriggerSealOfInsight)) ||
				!dpm.Proc(sim, spell.ProcMask, "Seal of Insight"+paladin.Label) {
				return
			}

			onSpecialOrSwingProc.Cast(sim, result.Target)
		},
	})

	// Seal of Insight self-buff.
	aura := paladin.SealOfInsightAura
	paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 20165},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 16.4,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
			IgnoreHaste: true,
		},

		ThreatMultiplier: 0,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			if paladin.CurrentSeal != nil {
				paladin.CurrentSeal.Deactivate(sim)
			}
			paladin.CurrentSeal = aura
			paladin.CurrentSeal.Activate(sim)
		},
	})
}
