package protection

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/paladin"
)

/*
Reduce damage taken by 20% for 10 sec.
While Ardent Defender is active, the next attack that would otherwise kill you will instead cause you to be healed for 15% of your maximum health.
*/
func (prot *ProtectionPaladin) registerArdentDefender() {
	actionID := core.ActionID{SpellID: 31850}

	adAura := prot.RegisterAura(core.Aura{
		Label:    "Ardent Defender" + prot.Label,
		ActionID: actionID,
		Duration: time.Second * 10,
	}).AttachMultiplicativePseudoStatBuff(&prot.PseudoStats.DamageTakenMultiplier, 0.8)

	ardentDefender := prot.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
		SpellSchool:    core.SpellSchoolHoly,
		ClassSpellMask: paladin.SpellMaskArdentDefender,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    prot.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
		},

		RelatedSelfBuff: adAura,
	})

	adHealAmount := 0.0

	// Spell to heal you when AD has procced; fire this before fatal damage so that a Death is not detected
	adHeal := prot.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 66235},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		CritMultiplier:   1,
		ThreatMultiplier: 0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealHealing(sim, &prot.Unit, adHealAmount, spell.OutcomeHealing)
		},
	})

	// >= 15% hp, hit gets reduced so we end up at 15% without heal
	// < 15% hp, hit gets reduced to 0 and we heal the remaining health up to 15%
	prot.AddDynamicDamageTakenModifier(func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult, isPeriodic bool) {
		if adAura.IsActive() && result.Damage >= prot.CurrentHealth() {
			maxHealth := prot.MaxHealth()
			currentHealth := prot.CurrentHealth()
			incomingDamage := result.Damage

			if currentHealth/maxHealth >= 0.15 {
				// Incoming attack gets reduced so we end up at 15% hp
				// TODO: Overkill counted as absorb but not as healing in logs
				result.Damage = currentHealth - maxHealth*0.15
				if sim.Log != nil {
					prot.Log(sim, "Ardent Defender absorbed %.1f damage", incomingDamage-result.Damage)
				}
			} else {
				// Incoming attack gets reduced to 0
				// Heal up to 15% hp
				// TODO: Overkill counted as absorb but not as healing in logs
				result.Damage = 0
				adHealAmount = maxHealth*0.15 - currentHealth
				adHeal.Cast(sim, &prot.Unit)
				if sim.Log != nil {
					prot.Log(sim, "Ardent Defender absorbed %.1f damage and healed for %.1f", incomingDamage, adHealAmount)
				}
			}

			adAura.Deactivate(sim)
		}
	})

	prot.AddDefensiveCooldownAura(adAura)
	prot.AddMajorCooldown(core.MajorCooldown{
		Spell:    ardentDefender,
		Type:     core.CooldownTypeSurvival,
		Priority: core.CooldownPriorityLow + 10,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return !prot.AnyActiveDefensiveCooldown()
		},
	})
}
