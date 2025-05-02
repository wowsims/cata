package affliction

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/warlock"
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

		ManaCost: core.ManaCostOptions{BaseCostPercent: 12},
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

		DamageMultiplier: 1,
		CritMultiplier:   affliction.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 0.5577,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := affliction.CalcScalingSpellDmg(0.95810002089)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
				if result.Landed() {
					affliction.HauntDebuffAuras.Get(result.Target).Activate(sim)
				}
			})
		},

		RelatedAuraArrays: affliction.HauntDebuffAuras.ToMap(),
	})
}
