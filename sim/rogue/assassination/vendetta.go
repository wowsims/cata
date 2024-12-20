package assassination

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/rogue"
)

func (sinRogue *AssassinationRogue) registerVendetta() {
	if !sinRogue.Talents.Vendetta {
		return
	}

	actionID := core.ActionID{SpellID: 79140}
	hasGlyph := sinRogue.HasPrimeGlyph(proto.RoguePrimeGlyph_GlyphOfVendetta)
	getDuration := func() time.Duration {
		return time.Duration((30.0+core.TernaryFloat64(sinRogue.Has4pcT13(), 3.0, 0))*core.TernaryFloat64(hasGlyph, 1.2, 1.0)) * time.Second
	}

	vendettaAura := sinRogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Vendetta",
			ActionID: actionID,
			Duration: getDuration(),
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				sinRogue.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier *= 1.2
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				sinRogue.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier /= 1.2
			},
		})
	})

	sinRogue.Vendetta = sinRogue.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		Flags:          core.SpellFlagAPL | core.SpellFlagMCD,
		ClassSpellMask: rogue.RogueSpellVendetta,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    sinRogue.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			aura := vendettaAura.Get(target)
			aura.Duration = getDuration()
			aura.Activate(sim)
		},
	})

	sinRogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    sinRogue.Vendetta,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityDefault,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return sinRogue.ComboPoints() >= 4
		},
	})
}
