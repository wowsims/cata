package unholy

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

/*
Consume 5 charges of Shadow Infusion on your Ghoul to transform it into a powerful undead monstrosity for 30 sec.
The Ghoul's abilities are empowered and take on new functions while the transformation is active.
*/
func (uhdk *UnholyDeathKnight) registerDarkTransformation() {
	actionID := core.ActionID{SpellID: 63560}

	uhdk.Ghoul.DarkTransformationAura = uhdk.Ghoul.GetOrRegisterAura(core.Aura{
		Label:    "Dark Transformation" + uhdk.Ghoul.Label,
		ActionID: actionID,
		Duration: time.Second * 30,
	}).AttachMultiplicativePseudoStatBuff(
		&uhdk.Ghoul.PseudoStats.DamageDealtMultiplier, 2.0,
	).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  death_knight.GhoulSpellClaw,
		FloatValue: 0.2,
	}).AttachDependentAura(uhdk.GetOrRegisterAura(core.Aura{
		Label:    "Dark Transformation" + uhdk.Label,
		ActionID: actionID,
		Duration: time.Second * 30,
	}))

	uhdk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolShadow,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: death_knight.DeathKnightSpellDarkTransformation,

		RuneCost: core.RuneCostOptions{
			UnholyRuneCost: 1,
			RunicPowerGain: 10,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
			IgnoreHaste: true,
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return uhdk.Ghoul.ShadowInfusionAura.GetStacks() == uhdk.Ghoul.ShadowInfusionAura.MaxStacks
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			uhdk.Ghoul.ShadowInfusionAura.Deactivate(sim)
			uhdk.Ghoul.DarkTransformationAura.Activate(sim)
		},
	})
}
