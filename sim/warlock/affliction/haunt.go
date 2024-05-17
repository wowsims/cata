package affliction

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/warlock"
)

func (affliction *AfflictionWarlock) registerHaunt() {
	if !affliction.Talents.Haunt {
		return
	}

	actionID := core.ActionID{SpellID: 48181}
	debuffMult := core.TernaryFloat64(affliction.HasPrimeGlyph(proto.WarlockPrimeGlyph_GlyphOfHaunt), 1.23, 1.2)

	affliction.HauntDebuffAuras = affliction.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Haunt-" + affliction.Label,
			ActionID: actionID,
			Duration: 12 * time.Second,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				affliction.AttackTables[aura.Unit.UnitIndex].HauntSEDamageTakenMultiplier *= debuffMult
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				affliction.AttackTables[aura.Unit.UnitIndex].HauntSEDamageTakenMultiplier /= debuffMult
			},
		})
	})

	affliction.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: warlock.WarlockSpellHaunt,
		MissileSpeed:   20,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.12,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 1500 * time.Millisecond,
			},
			CD: core.Cooldown{
				Timer:    affliction.NewTimer(),
				Duration: 8 * time.Second,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           affliction.DefaultSpellCritMultiplier(),
		ThreatMultiplier:         1,
		BonusCoefficient:         0.5577 * 1.25,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := affliction.CalcScalingSpellDmg(warlock.Coefficient_Haunt)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
				if result.Landed() {
					affliction.HauntDebuffAuras.Get(result.Target).Activate(sim)
				}
			})

		},
		RelatedAuras: []core.AuraArray{affliction.HauntDebuffAuras},
	})
}
