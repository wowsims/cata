package arcane

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/mage"
)

func (arcane *ArcaneMage) registerArcaneBlastSpell() {

	//https://wago.tools/db2/SpellEffect?build=5.5.0.60802&filter%5BSpellID%5D=30451
	arcaneBlastVariance := .15
	arcaneBlastCoefficient := .78
	arcaneBlastScaling := .78

	arcane.RegisterSpell(core.SpellConfig{

		ActionID:       core.ActionID{SpellID: 30451},
		SpellSchool:    core.SpellSchoolArcane,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: mage.MageSpellArcaneBlast,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 1.5, // Arcane Blast mana cost lowered by 10% to 1.5% of base mana (was 1.666%) -  https://eu.forums.blizzard.com/en/wow/t/mists-of-pandaria-classic-development-notes-updated-20-june/571162/13
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2000,
			},
		},

		DamageMultiplier: 1 * 1.37, // Arcane Blast damage increased by 37% -  https://www.wowhead.com/mop-classic/news/guardian-druid-and-arcane-mage-buffed-additional-mists-of-pandaria-class-changes-377468
		CritMultiplier:   arcane.DefaultCritMultiplier(),
		BonusCoefficient: arcaneBlastCoefficient,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := arcane.CalcAndRollDamageRange(sim, arcaneBlastScaling, arcaneBlastVariance)
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				arcane.ArcaneChargesAura.Activate(sim)
				arcane.ArcaneChargesAura.AddStack(sim)
			}
		},
	})
}
