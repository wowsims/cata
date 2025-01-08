package rogue

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (rogue *Rogue) getVendettaDuration(baseDuration float64) time.Duration {
	hasGlyph := rogue.HasPrimeGlyph(proto.RoguePrimeGlyph_GlyphOfVendetta)
	return time.Duration((30.0+baseDuration)*core.TernaryFloat64(hasGlyph, 1.2, 1.0)) * time.Second
}

func (rogue *Rogue) registerVendetta() {
	if !rogue.Talents.Vendetta {
		return
	}

	actionID := core.ActionID{SpellID: 79140}
	duration := rogue.getVendettaDuration(0)

	vendettaAuras := rogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Vendetta",
			ActionID: actionID,
			Duration: duration,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				rogue.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier *= 1.2
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				rogue.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier /= 1.2
			},
		})
	})

	rogue.Vendetta = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		Flags:          core.SpellFlagAPL | core.SpellFlagMCD,
		ClassSpellMask: RogueSpellVendetta,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			aura := vendettaAuras.Get(target)
			aura.Activate(sim)
		},
		RelatedAuras: []core.AuraArray{vendettaAuras},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    rogue.Vendetta,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityDefault,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return rogue.ComboPoints() >= 4
		},
	})
}
