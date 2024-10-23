package elemental

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/shaman"
)

func (elemental *ElementalShaman) registerThunderstormSpell() {
	actionID := core.ActionID{SpellID: 51490}
	manaMetrics := elemental.NewManaMetrics(actionID)

	manaRestore := 0.08
	if elemental.HasMinorGlyph(proto.ShamanMinorGlyph_GlyphOfThunderstorm) {
		manaRestore = 0.02
	}

	elemental.Thunderstorm = elemental.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL | shaman.SpellFlagFocusable,
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: shaman.SpellMaskThunderstorm,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    elemental.NewTimer(),
				Duration: time.Second * 45,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   elemental.DefaultSpellCritMultiplier(),
		BonusCoefficient: 0.571,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			elemental.AddMana(sim, elemental.MaxMana()*manaRestore, manaMetrics)

			if elemental.Shaman.ThunderstormInRange {
				results := make([]*core.SpellResult, elemental.Env.GetNumTargets())
				baseDamage := elemental.GetShaman().CalcAndRollDamageRange(sim, 1.62999999523, 0.13300000131)
				aoeMult := sim.Encounter.AOECapMultiplier()
				spell.DamageMultiplier *= aoeMult
				for i, aoeTarget := range sim.Encounter.TargetUnits {
					results[i] = spell.CalcDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
				}
				for i, _ := range sim.Encounter.TargetUnits {
					spell.DealDamage(sim, results[i])
				}
				spell.DamageMultiplier /= aoeMult
			}
		},
	})
}
