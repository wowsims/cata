package combat

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/rogue"
)

const masteryChancePerPoint = .02
const masteryBaseEffect = 0.16

func (comRogue *CombatRogue) applyMastery() {

	comRogue.mainGauche = comRogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 86392},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagPassiveSpell,
		ClassSpellMask: rogue.RogueSpellMainGauche,

		DamageMultiplier: 1,
		CritMultiplier:   comRogue.MeleeCritMultiplier(false),
		ThreatMultiplier: 1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		},
	})

	comRogue.mainGaucheAura = comRogue.RegisterAura(core.Aura{
		Label:    "Mastery: Main Gauche",
		Duration: core.NeverExpires,
		// ActionID Excluded to not clog up Buffs metrics

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMeleeMH) && spell != comRogue.mainGauche && spell != comRogue.Rupture {
				masteryPoints := comRogue.GetMasteryPoints()
				mgProcChance := masteryChancePerPoint*masteryPoints + masteryBaseEffect

				if sim.Proc(mgProcChance, "Main Gauche") {
					comRogue.mainGauche.Cast(sim, result.Target)
				}
			}
		},
	})
}
