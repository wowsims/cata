package monk

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

/*
Tooltip:
The Monk attunes $G himself:herself; differently depending on the weapon type.

One-handed weapons / Dual-wield one-handed weapons:
Autoattack damage increased by 40%.

Two-handed weapons:
Melee attack speed increased by 40%.
*/
func (monk *Monk) registerWayOfTheMonk() {
	aura := core.MakePermanent(monk.RegisterAura(core.Aura{
		Label:      "Way of the Monk" + monk.Label,
		ActionID:   core.ActionID{SpellID: 120277},
		BuildPhase: core.CharacterBuildPhaseBase,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if monk.HandType == proto.HandType_HandTypeTwoHand {
				monk.MultiplyMeleeSpeed(sim, 1.4)
			} else {
				if monk.MHAutoSpell != nil {
					monk.MHAutoSpell.DamageMultiplier *= 1.4
				}
				if monk.OHAutoSpell != nil {
					monk.OHAutoSpell.DamageMultiplier *= 1.4
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if monk.HandType == proto.HandType_HandTypeTwoHand {
				monk.MultiplyMeleeSpeed(sim, 1/1.4)
			} else {
				if monk.MHAutoSpell != nil {
					monk.MHAutoSpell.DamageMultiplier /= 1.4
				}
				if monk.OHAutoSpell != nil {
					monk.OHAutoSpell.DamageMultiplier /= 1.4
				}
			}
		},
	}))

	monk.RegisterItemSwapCallback(core.MeleeWeaponSlots(), func(sim *core.Simulation, slot proto.ItemSlot) {
		aura.Deactivate(sim)
		aura.Activate(sim)
	})

}

/*
Tooltip:
Increases your chance to parry by 5%.

Whenever you parry an attack, you reflexively strike back at the enemy for ${0.3*$<low>} to ${0.3*$<high>} damage. This strike has a 1 sec cooldown.

$stnc=$?a103985[${1.2*7.5}][${1.0*7.5}]
*/
func (monk *Monk) registerSwiftReflexes() {
	swiftReflexesAttack := monk.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 124335},
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell,
		MaxRange:    core.MaxMeleeRange,

		DamageMultiplier: 0.3 * 7.5,
		ThreatMultiplier: 1,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := monk.CalculateMonkStrikeDamage(sim, spell)

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})

	aura := core.MakeProcTriggerAura(&monk.Unit, core.ProcTrigger{
		Name:     "Swift Reflexes" + monk.Label,
		ActionID: core.ActionID{SpellID: 124334},
		Outcome:  core.OutcomeParry,
		Callback: core.CallbackOnSpellHitTaken,
		ICD:      time.Second,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			swiftReflexesAttack.Cast(sim, spell.Unit)
		},
	})
	aura.BuildPhase = core.CharacterBuildPhaseBase
	aura.AttachAdditivePseudoStatBuff(&monk.PseudoStats.BaseParryChance, 0.05)
}
