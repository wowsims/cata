package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (rogue *Rogue) registerShadowBladesCD() {
	mhHit := rogue.makeShadowBladeHit(true)
	ohHit := rogue.makeShadowBladeHit(false)
	cpMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: 121471})

	rogue.ShadowBladesAura = rogue.RegisterAura(core.Aura{
		Label:    "Shadow Blades",
		ActionID: core.ActionID{SpellID: 121471},
		Duration: time.Second * 12,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			// Make auto attacks deal 0 damage for the duration of SB
			// This allows for anything tied to autos (poisons, main gauche, etc) to still fire
			rogue.AutoAttacks.MHAuto().DamageMultiplier = 0
			rogue.AutoAttacks.OHAuto().DamageMultiplier = 0
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.AutoAttacks.MHAuto().DamageMultiplier = 1
			rogue.AutoAttacks.OHAuto().DamageMultiplier = 1
		},

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				if spell == rogue.AutoAttacks.MHAuto() {
					mhHit.Cast(sim, result.Target)
				} else if spell == rogue.AutoAttacks.OHAuto() {
					ohHit.Cast(sim, result.Target)
				} else if spell.Flags.Matches(SpellFlagBuilder) {
					rogue.AddComboPoints(sim, 1, cpMetrics)
				}
			}
		},
	})

	rogue.ShadowBlades = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 121471},
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: RogueSpellShadowBlades,

		Cast: core.CastConfig{
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Minute * 3,
			},
		},
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
		},
		RelatedSelfBuff: rogue.ShadowBladesAura,
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    rogue.ShadowBlades,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityDefault,
	})
}

func (rogue *Rogue) makeShadowBladeHit(isMH bool) *core.Spell {
	procMask := core.Ternary(isMH, core.ProcMaskMeleeMH, core.ProcMaskMeleeOH)
	tag := core.TernaryInt32(isMH, 1, 2)
	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 121471, Tag: tag},
		ClassSpellMask: RogueSpellShadowBladesHit,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       procMask,
		Flags:          core.SpellFlagMeleeMetrics,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           rogue.CritMultiplier(true),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var baseDamage float64
			if isMH {
				baseDamage = spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			} else {
				baseDamage = spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			}

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialBlockAndCrit)
		},
	})
}
