package combat

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/rogue"
)

var AdrenalineRushActionID = core.ActionID{SpellID: 13750}

func (comRogue *CombatRogue) registerAdrenalineRushCD() {
	speedBonus := 1.2
	inverseBonus := 1 / speedBonus

	// Reduces the GCD of Sinister Strike, Revealing Strike, Eviscerate, Slice and Dice, and Rupture by 0.2 sec
	gcdReduction := comRogue.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_GlobalCooldown_Flat,
		ClassMask: rogue.RogueSpellRupture | rogue.RogueSpellEviscerate | rogue.RogueSpellSliceAndDice | rogue.RogueSpellRevealingStrike | rogue.RogueSpellSinisterStrike,
		TimeValue: time.Millisecond * -200,
	})

	comRogue.AdrenalineRushAura = comRogue.RegisterAura(core.Aura{
		Label:    "Adrenaline Rush",
		ActionID: AdrenalineRushActionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			comRogue.ApplyAdditiveEnergyRegenBonus(sim, 1.0)
			comRogue.MultiplyMeleeSpeed(sim, speedBonus)
			gcdReduction.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			comRogue.ApplyAdditiveEnergyRegenBonus(sim, -1.0)
			comRogue.MultiplyMeleeSpeed(sim, inverseBonus)
			gcdReduction.Deactivate()
		},
	})

	comRogue.AdrenalineRush = comRogue.RegisterSpell(core.SpellConfig{
		ActionID:       AdrenalineRushActionID,
		ClassSpellMask: rogue.RogueSpellAdrenalineRush,
		Flags:          core.SpellFlagReadinessTrinket,

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
	})
}
