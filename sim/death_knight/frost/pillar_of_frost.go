package frost

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/death_knight"
)

/*
Calls upon the power of Frost to increase the Death Knight's Strength by 20%.
Icy crystals hang heavy upon the Death Knight's body, providing immunity against external movement

-- Glyph of Pillar of Frost --

and all effects that cause loss of control, but also reducing movement speed by 70% while active.

-- else --

such as knockbacks.

----------

Lasts 20 sec.
*/
func (fdk *FrostDeathKnight) registerPillarOfFrost() {
	actionID := core.ActionID{SpellID: 51271}

	pillarOfFrostAura := fdk.RegisterAura(core.Aura{
		Label:    "Pillar of Frost" + fdk.Label,
		ActionID: actionID,
		Duration: time.Second * 20,
	}).AttachStatDependency(
		fdk.NewDynamicMultiplyStat(stats.Strength, 1.2),
	)

	spell := fdk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,
		ClassSpellMask: death_knight.DeathKnightSpellPillarOfFrost,

		RuneCost: core.RuneCostOptions{
			FrostRuneCost:  1,
			RunicPowerGain: 10,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    fdk.NewTimer(),
				Duration: time.Minute * 1,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
		},

		RelatedSelfBuff: pillarOfFrostAura,
	})

	fdk.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}
