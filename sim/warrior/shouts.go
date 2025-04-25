package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

const ShoutExpirationThreshold = time.Second * 3

func (warrior *Warrior) MakeShoutSpellHelper(actionID core.ActionID, spellMask int64, allyAuras core.AuraArray) *core.Spell {

	shoutMetrics := warrior.NewRageMetrics(actionID)
	rageGen := 20.0 + 5.0*float64(warrior.Talents.BoomingVoice)
	talentReduction := time.Duration(warrior.Talents.BoomingVoice*15) * time.Second

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
				Timer:    warrior.shoutsCD,
				Duration: time.Minute - talentReduction,
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
