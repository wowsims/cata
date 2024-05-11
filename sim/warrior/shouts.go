package warrior

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

const ShoutExpirationThreshold = time.Second * 3

func (warrior *Warrior) MakeShoutSpellHelper(actionID core.ActionID, spellMask int64, allyAuras core.AuraArray) *core.Spell {

	shoutMetrics := warrior.NewRageMetrics(actionID)
	rageGen := 20.0 + 5.0*float64(warrior.Talents.BoomingVoice)
	return warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL | core.SpellFlagHelpful,
		ClassSpellMask: spellMask,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(), // TODO: double-check that BS and CS don't share CDs
				Duration: time.Minute,
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

		RelatedAuras: []core.AuraArray{allyAuras},
	})
}

func (warrior *Warrior) RegisterShouts() {
	warrior.BattleShout = warrior.MakeShoutSpellHelper(core.ActionID{SpellID: 6673}, SpellMaskBattleShout, warrior.NewAllyAuraArray(func(unit *core.Unit) *core.Aura {
		if unit.Type == core.PetUnit {
			return nil
		}
		return core.BattleShoutAura(unit, false, warrior.HasMinorGlyph(proto.WarriorMinorGlyph_GlyphOfBattle))
	}))

	warrior.CommandingShout = warrior.MakeShoutSpellHelper(core.ActionID{SpellID: 469}, SpellMaskCommandingShout, warrior.NewAllyAuraArray(func(unit *core.Unit) *core.Aura {
		if unit.Type == core.PetUnit {
			return nil
		}
		return core.CommandingShoutAura(unit, false, warrior.HasMinorGlyph(proto.WarriorMinorGlyph_GlyphOfCommand))
	}))
}
