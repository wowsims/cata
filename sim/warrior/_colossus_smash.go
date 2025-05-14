package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (warrior *Warrior) RegisterColossusSmash() {
	actionID := core.ActionID{SpellID: 86346}

	warrior.ColossusSmashAuras = warrior.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Colossus Smash",
			ActionID: actionID,
			Duration: time.Second * 6,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				warrior.AttackTables[aura.Unit.UnitIndex].IgnoreArmor = true
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				warrior.AttackTables[aura.Unit.UnitIndex].IgnoreArmor = false
			},
		})
	})

	warrior.ColossusSmash = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAPL | core.SpellFlagMeleeMetrics,
		ClassSpellMask: SpellMaskColossusSmash | SpellMaskSpecialAttack,
		MaxRange:       core.MaxMeleeRange,

		RageCost: core.RageCostOptions{
			Cost:   20,
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 20,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1.5,
		CritMultiplier:   warrior.DefaultCritMultiplier(),

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 120.0 + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			if !result.Landed() {
				spell.IssueRefund(sim)
			} else {
				aura := warrior.ColossusSmashAuras.Get(target)
				aura.Activate(sim)
				if warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfColossusSmash) {
					warrior.TryApplySunderArmorEffect(sim, target)
				}
			}
		},
	})
}
