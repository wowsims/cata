package shaman

import (
	"github.com/wowsims/mop/sim/core"
)

func (shaman *Shaman) StormlashActionID() core.ActionID {
	return core.ActionID{
		SpellID: 120668,
		Tag:     shaman.Index,
	}
}

func (shaman *Shaman) registerStormlashCD() {
	actionID := shaman.StormlashActionID()

	slAuras := []*core.Aura{}
	for _, party := range shaman.Env.Raid.Parties {
		for _, partyMember := range party.Players {
			slAuras = append(slAuras, core.StormLashAura(partyMember.GetCharacter(), actionID.Tag))
		}
	}

	spell := shaman.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskStormlashTotem,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 5.9,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: core.StormLashCD,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !shaman.HasActiveAuraWithTag("StormLash")
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			for _, slAura := range slAuras {
				slAura.Activate(sim)
			}
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell:    spell,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityDefault,
	})
}
