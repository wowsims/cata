package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (warlock *Warlock) registerShadowBoltSpell() {
	spellCoeff := 0.857 * (1 + 0.04*float64(warlock.Talents.ShadowAndFlame))
	ISBProcChance := 0.2 * float64(warlock.Talents.ImprovedShadowBolt)

	var shadowMasteryAuras core.AuraArray
	if ISBProcChance > 0 {
		shadowMasteryAuras = warlock.NewEnemyAuraArray(core.ShadowMasteryAura)
	}

	warlock.ShadowBolt = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 47809},
		SpellSchool:  core.SpellSchoolShadow,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        core.SpellFlagAPL,
		MissileSpeed: 20,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.17,
			Multiplier: 1 -
				core.TernaryFloat64(warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfShadowBolt), 0.1, 0),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * (3000 - 100*time.Duration(warlock.Talents.Bane)),
			},
		},

		DamageMultiplierAdditive: 1 +
			warlock.GrandFirestoneBonus() +
			0.03*float64(warlock.Talents.ShadowMastery) +
			0.02*float64(warlock.Talents.ImprovedShadowBolt),
		CritMultiplier:   warlock.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(694, 775) + spellCoeff*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
				if result.Landed() {
					// ISB debuff
					if sim.Proc(ISBProcChance, "ISB") {
						shadowMasteryAuras.Get(target).Activate(sim)
					}
					warlock.everlastingAfflictionRefresh(sim, target)
				}
			})
		},
	})

	if ISBProcChance > 0 {
		warlock.ShadowBolt.RelatedAuras = append(warlock.ShadowBolt.RelatedAuras, shadowMasteryAuras)
	}
}
