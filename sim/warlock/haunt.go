package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

// TODO: SpellFlagHauntSE
func (warlock *Warlock) registerHauntSpell() {
	if !warlock.Talents.Haunt {
		return
	}

	actionID := core.ActionID{SpellID: 48181}
	debuffMult := core.TernaryFloat64(warlock.HasPrimeGlyph(proto.WarlockPrimeGlyph_GlyphOfHaunt), 1.23, 1.2)

	warlock.HauntDebuffAuras = warlock.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Haunt-" + warlock.Label,
			ActionID: actionID,
			Duration: time.Second * 12,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				warlock.AttackTables[aura.Unit.UnitIndex].HauntSEDamageTakenMultiplier *= debuffMult
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				warlock.AttackTables[aura.Unit.UnitIndex].HauntSEDamageTakenMultiplier /= debuffMult
			},
		})
	})

	warlock.Haunt = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellHaunt,
		MissileSpeed:   20,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.12,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * 8,
			},
		},

		DamageMultiplierAdditive: 1 + warlock.GrandFirestoneBonus(),
		CritMultiplier:           warlock.DefaultSpellCritMultiplier(),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 1078 + (spell.SpellPower()*0.5577)*1.25
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
				if result.Landed() {
					warlock.HauntDebuffAuras.Get(result.Target).Activate(sim)
					//TODO: EverlastingAffliction
					//warlock.everlastingAfflictionRefresh(sim, target)
				}
			})
		},
	})
}
