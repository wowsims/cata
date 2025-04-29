package combat

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/rogue"
)

var AdrenalineRushActionID = core.ActionID{SpellID: 13750}

func (comRogue *CombatRogue) registerAdrenalineRushCD() {
	if !comRogue.Talents.AdrenalineRush {
		return
	}

	speedBonus := 1.2
	inverseBonus := 1 / speedBonus

	comRogue.AdrenalineRushAura = comRogue.RegisterAura(core.Aura{
		Label:    "Adrenaline Rush",
		ActionID: AdrenalineRushActionID,
		Duration: core.TernaryDuration(comRogue.HasPrimeGlyph(proto.RoguePrimeGlyph_GlyphOfAdrenalineRush), time.Second*20, time.Second*15),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			comRogue.ApplyAdditiveEnergyRegenBonus(sim, 1.0)
			comRogue.MultiplyMeleeSpeed(sim, speedBonus)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			comRogue.ApplyAdditiveEnergyRegenBonus(sim, -1.0)
			comRogue.MultiplyMeleeSpeed(sim, inverseBonus)
		},
	})

	comRogue.AdrenalineRush = comRogue.RegisterSpell(core.SpellConfig{
		ActionID:       AdrenalineRushActionID,
		ClassSpellMask: rogue.RogueSpellAdrenalineRush,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    comRogue.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			comRogue.BreakStealth(sim)
			spell.RelatedSelfBuff.Activate(sim)
		},
		RelatedSelfBuff: comRogue.AdrenalineRushAura,
	})

	comRogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    comRogue.AdrenalineRush,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityDefault,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			thresh := 40.0
			return comRogue.CurrentEnergy() <= thresh && !comRogue.KillingSpree.IsReady(sim)
		},
	})
}
