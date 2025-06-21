package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (war *Warrior) registerVictoryRush() {

	war.VictoryRushAura = core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:     "Victorious",
		ActionID: core.ActionID{SpellID: 32216},
		Duration: 20 * time.Second,
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask:  SpellMaskImpendingVictory,
		Kind:       core.SpellMod_PowerCost_Pct,
		FloatValue: -2,
	})

	if war.Talents.ImpendingVictory {
		return
	}

	actionID := core.ActionID{SpellID: 34428}
	healthMetrics := war.NewHealthMetrics(actionID)
	hasGlyph := war.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfVictoryRush)

	war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAPL | core.SpellFlagMeleeMetrics,
		ClassSpellMask: SpellMaskImpendingVictory,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   war.DefaultCritMultiplier(),

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return war.VictoryRushAura.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			war.VictoryRushAura.Deactivate(sim)
			baseDamage := 56 + spell.MeleeAttackPower()*0.56
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			healthMultiplier := core.TernaryFloat64(hasGlyph, 0.3, 0.2)

			if result.Landed() {
				war.GainHealth(sim, war.MaxHealth()*healthMultiplier, healthMetrics)
			}
		},
	})
}
