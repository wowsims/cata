package paladin

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

/*
Fills you with Holy Light, increasing your casting speed by 10%, improving healing spells by 5% and giving melee attacks a chance to heal

-- Glyph of the Battle Healer --

the most wounded member of your party or raid

-- else --

you

----------

for

-- Holy --

(0.15 * <AP> + 0.15 * <SP>) * 1.25
and restore 4% of base mana when striking a target outside of an Arena or Battleground

-- else --

(0.15 * <AP> + 0.15 * <SP>)

----------
*/
func (paladin *Paladin) registerSealOfInsight() {
	hasGlyphOfTheBattleHealer := paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfTheBattleHealer)
	isHoly := paladin.Spec == proto.Spec_SpecHolyPaladin
	actionID := core.ActionID{SpellID: 20167}
	manaMetrics := paladin.NewManaMetrics(actionID)

	// Seal of Insight on-hit proc
	onSpecialOrSwingProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagHelpful | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,
		ClassSpellMask: SpellMaskSealOfInsight,

		DamageMultiplier: 1,
		CritMultiplier:   0,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			heal := 0.15*spell.SpellPower() + 0.15*spell.MeleeAttackPower()
			spell.CalcAndDealHealing(sim, target, heal, spell.OutcomeHealing)

			if isHoly {
				// Beta changes 2025-06-13: https://www.wowhead.com/mop-classic/news/additional-holy-priest-and-paladin-changes-coming-to-mists-of-pandaria-classic-377264
				// - Seal of Insight now has a chance to restore 4% of the Holy Paladin’s base mana when striking a target in a non-PvP environment. [5.4 Revert].
				// Beta changes 2025-06-16: https://www.wowhead.com/mop-classic/news/blood-death-knights-buffed-and-even-more-class-balance-adjustments-mists-of-377292
				// - Seal of Insight’s Base Mana restoration increased to 6% of base mana (was 4%). [New]
				paladin.AddMana(sim, paladin.BaseMana*0.06, manaMetrics)
			}
		},
	})

	dpm := paladin.NewLegacyPPMManager(15, core.ProcMaskMeleeMH)
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

			// SoI only procs on white hits, CS, HoW, ShotR and TV
			if (spell.ProcMask&core.ProcMaskMeleeWhiteHit == 0 &&
				!spell.Matches(SpellMaskCanTriggerSealOfInsight)) ||
				!dpm.Proc(sim, spell.ProcMask, "Seal of Insight"+paladin.Label) {
				return
			}

			if hasGlyphOfTheBattleHealer {
				onSpecialOrSwingProc.Cast(sim, sim.Raid.GetLowestHealthAllyUnit())
			} else {
				onSpecialOrSwingProc.Cast(sim, &paladin.Unit)
			}
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  SpellMaskModifiedBySealOfInsight,
		FloatValue: 0.05,
	}).AttachMultiplyCastSpeed(1.1)

	// Seal of Insight self-buff.
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
			paladin.CurrentSeal = paladin.SealOfInsightAura
			paladin.CurrentSeal.Activate(sim)
		},
	})
}
