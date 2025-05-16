package monk

import (
	"time"

	"github.com/wowsims/mop/sim/common/cata"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

var combatConditioningActionID = core.ActionID{SpellID: 100784}.WithTag(2) // actual 128531

func combatConditioningConfig(unit *core.Unit, isSEFClone bool) cata.IgniteConfig {
	config := cata.IgniteConfig{
		ActionID:           combatConditioningActionID,
		DotAuraLabel:       "Blackout Kick (DoT)" + unit.Label,
		DisableCastMetrics: true,
		IncludeAuraDelay:   true,
		SpellSchool:        core.SpellSchoolPhysical,
		NumberOfTicks:      4,
		TickLength:         time.Second,

		ProcTrigger: core.ProcTrigger{
			Name:           "Combat Conditioning" + unit.Label,
			Callback:       core.CallbackOnSpellHitDealt,
			ClassSpellMask: MonkSpellBlackoutKick,
			Outcome:        core.OutcomeLanded,
		},

		DamageCalculator: func(result *core.SpellResult) float64 {
			return result.Damage * 0.2
		},
	}

	if isSEFClone {
		config.ActionID = config.ActionID.WithTag(config.ActionID.Tag + SEFSpellID)
	}

	return config
}

func (monk *Monk) registerCombatConditioning() {
	if monk.Spec != proto.Spec_SpecWindwalkerMonk || (!monk.HasMinorGlyph(proto.MonkMinorGlyph_GlyphOfBlackoutKick) || monk.PseudoStats.InFrontOfTarget) {
		return
	}

	cata.RegisterIgniteEffect(&monk.Unit, combatConditioningConfig(&monk.Unit, false))
}

func (pet *StormEarthAndFirePet) registerSEFCombatConditioning() {
	if pet.owner.Spec != proto.Spec_SpecWindwalkerMonk || !pet.owner.HasMinorGlyph(proto.MonkMinorGlyph_GlyphOfBlackoutKick) {
		return
	}

	cata.RegisterIgniteEffect(&pet.Unit, combatConditioningConfig(&pet.Unit, true))
}
