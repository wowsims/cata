package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (rogue *Rogue) registerShadowBladesCD() {
	mhHit := rogue.makeShadowBladeHit(true)
	ohHit := rogue.makeShadowBladeHit(false)
	cpMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: 121471})

	sbAura := rogue.RegisterAura(core.Aura{
		Label:    "Shadow Blades",
		ActionID: core.ActionID{SpellID: 121471},
		Duration: time.Second * 12,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				if spell == rogue.AutoAttacks.MHAuto() {
					mhHit.Cast(sim, result.Target)
				} else if spell == rogue.AutoAttacks.OHAuto() {
					ohHit.Cast(sim, result.Target)
				}

				if spell.Flags.Matches(SpellFlagBuilder) {
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
			sbAura.Activate(sim)
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    rogue.ShadowBlades,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityDefault,
	})
}

func (rogue *Rogue) makeShadowBladeHit(isMH bool) *core.Spell {
	procMask := core.Ternary(isMH, core.ProcMaskMeleeMH, core.ProcMaskMeleeOH)
	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 121473},
		ClassSpellMask: RogueSpellShadowBlades,
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
