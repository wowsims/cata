package blood

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

// Each time you heal yourself with Death Strike while in Blood Presence, you gain (50 + (<Mastery Rating>/600)*6.25)% of the amount healed as a Physical damage absorption shield.
func (bdk *BloodDeathKnight) registerMastery() {
	shieldAmount := 0.0
	currentShield := 0.0

	var shieldSpell *core.Spell
	shieldSpell = bdk.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 77535},
		ProcMask:    core.ProcMaskSpellHealing,
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagPassiveSpell | core.SpellFlagHelpful,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Shield: core.ShieldConfig{
			SelfOnly: true,
			Aura: core.Aura{
				Label:    "Blood Shield" + bdk.Label,
				Duration: time.Second * 10,

				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					shieldAmount = 0.0
					currentShield = 0.0
				},
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !spell.SpellSchool.Matches(core.SpellSchoolPhysical) || result.Damage <= 0 {
						return
					}

					if currentShield <= 0 {
						shieldSpell.SelfShield().Deactivate(sim)
						return
					}

					damageReduced := min(result.Damage, currentShield)
					currentShield -= damageReduced

					bdk.GainHealth(sim, damageReduced, shieldSpell.HealthMetrics(result.Target))

					if currentShield <= 0 {
						shieldSpell.SelfShield().Deactivate(sim)
					}
				},
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if currentShield < bdk.MaxHealth() {
				shieldAmount = min(shieldAmount, bdk.MaxHealth()-currentShield)
				currentShield += shieldAmount
				spell.SelfShield().Apply(sim, shieldAmount)
			}
		},
	})

	core.MakeProcTriggerAura(&bdk.Unit, core.ProcTrigger{
		Name:           "Mastery: Blood Shield" + bdk.Label,
		ActionID:       core.ActionID{SpellID: 77513},
		Callback:       core.CallbackOnHealDealt,
		ClassSpellMask: death_knight.DeathKnightSpellDeathStrikeHeal,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			shieldAmount = result.Damage * bdk.getMasteryPercent()
			shieldSpell.Cast(sim, result.Target)
		},
	})

}

func (bdk BloodDeathKnight) getMasteryPercent() float64 {
	return 0.5 + 0.0625*bdk.GetMasteryPoints()
}
