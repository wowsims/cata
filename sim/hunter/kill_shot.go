package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (hunter *Hunter) registerKillShotSpell() {
	if hunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfKillShot) {
		icd := core.Cooldown{
			Timer:    hunter.NewTimer(),
			Duration: time.Second * 6,
		}
		hunter.RegisterAura(core.Aura{
			Label:    "Kill Shot Glyph",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell == hunter.KillShot {
					if icd.IsReady(sim) {
						icd.Use(sim)
						hunter.KillShot.CD.Reset()
					}
				}
			},
		})
	}
	spellData, _ := core.CurrentSpellGen().GetDBC().FetchSpell(53351)
	eff2, _ := spellData.EffectN(2)
	actionId := core.ActionID{SpellID: 53351}
	spellConfig := core.CurrentSpellGen().ParseSpellData(53351, &hunter.Unit, &actionId)
	spellConfig.ProcMask = core.ProcMaskRangedSpecial // Probably can get this for config as well
	spellConfig.Flags = core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL
	spellConfig.ExtraCastCondition = func(sim *core.Simulation, target *core.Unit) bool {
		return sim.IsExecutePhase20()
	}
	spellConfig.BonusCritRating = 0 + 5*core.CritRatingPerCritChance*float64(hunter.Talents.SniperTraining)
	spellConfig.DamageMultiplier = eff2.Min(core.CurrentSpellGen().GetDBC(), 85, 85)
	spellConfig.CritMultiplier = hunter.CritMultiplier(true, true, false)
	spellConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		// (100% weapon dmg + 45% RAP + 543) * 150%
		weaponDamage := 0.0
		// Todo: Figure out a nice way to deal with this
		// if spell.SpellData != nil {
		// 	min, max, normalized, _ := core.CurrentSpellGen().ParseEffects(spell.SpellData, &hunter.Unit)
		// 	if normalized {
		// 		weaponDamage += hunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower(target))
		// 		weaponDamage += sim.Roll(min, max)
		// 	}
		// }
		eff1, _ := spellData.EffectN(1)
		weaponDamage += hunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower(target))
		weaponDamage += eff1.Average(core.CurrentSpellGen().GetDBC(), 85, 85)

		// Todo: Figure out where exactly this part comes from, cant find an effect for it
		rapBonusDamage := spell.RangedAttackPower(target) * (0.45 * eff2.Min(core.CurrentSpellGen().GetDBC(), 85, 85))

		baseDamage := weaponDamage + rapBonusDamage
		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

		spell.WaitTravelTime(sim, func(sim *core.Simulation) {
			spell.DealDamage(sim, result)
		})
	}

	hunter.KillShot = hunter.RegisterSpell(*spellConfig)
}
