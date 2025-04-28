package paladin

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (paladin *Paladin) registerSealOfInsight() {
	// Judgement of Insight cast on Judgement
	paladin.JudgementOfInsight = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 54158},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskMeleeSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagSecondaryJudgement,
		ClassSpellMask: SpellMaskJudgementOfInsight,

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 1 +
				0.25*spell.SpellPower() +
				0.15999999642*spell.MeleeAttackPower()

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
		},
	})

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
		Duration: time.Minute * 30,
		Dpm:      dpm,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Don't proc on misses
			if !result.Landed() {
				return
			}

			// SoJ only procs on white hits, CS, TV and HoW
			if (spell.ProcMask&core.ProcMaskMeleeWhiteHit == 0 &&
				spell.ClassSpellMask&SpellMaskCanTriggerSealOfInsight == 0) ||
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
			BaseCostPercent: 14,
			PercentModifier: 100,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ThreatMultiplier: 0,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			if paladin.CurrentSeal != nil {
				paladin.CurrentSeal.Deactivate(sim)
			}
			paladin.CurrentSeal = aura
			paladin.CurrentJudgement = paladin.JudgementOfInsight
			paladin.CurrentSeal.Activate(sim)
		},
	})
}
