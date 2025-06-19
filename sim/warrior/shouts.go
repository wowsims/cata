package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

const ShoutExpirationThreshold = time.Second * 3

func (warrior *Warrior) MakeShoutSpellHelper(actionID core.ActionID, spellMask int64, allyAuras core.AuraArray) *core.Spell {
	hasGlyph := warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfHoarseVoice)
	shoutMetrics := warrior.NewRageMetrics(actionID)
	rageGen := core.TernaryFloat64(hasGlyph, 10, 20)
	duration := core.TernaryDuration(hasGlyph, time.Second*30, time.Minute*1)

	return warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
		ClassSpellMask: spellMask,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.sharedShoutsCD,
				Duration: duration,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			warrior.AddRage(sim, rageGen, shoutMetrics)
			for _, aura := range allyAuras {
				if aura != nil {
					aura.Activate(sim)
				}
			}
		},

		RelatedAuraArrays: allyAuras.ToMap(),
	})
}

func (warrior *Warrior) registerShouts() {
	warrior.BattleShout = warrior.MakeShoutSpellHelper(core.ActionID{SpellID: 6673}, SpellMaskBattleShout, warrior.NewAllyAuraArray(func(unit *core.Unit) *core.Aura {
		if unit.Type == core.PetUnit {
			return nil
		}
		return core.BattleShoutAura(unit, false)
	}))

	warrior.CommandingShout = warrior.MakeShoutSpellHelper(core.ActionID{SpellID: 469}, SpellMaskCommandingShout, warrior.NewAllyAuraArray(func(unit *core.Unit) *core.Aura {
		if unit.Type == core.PetUnit {
			return nil
		}
		return core.CommandingShoutAura(unit, false)
	}))
}
