package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (rogue *Rogue) registerVanishSpell() {
	rogue.Vanish = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 1856},
		SpellSchool:    core.SpellSchoolPhysical,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: RogueSpellVanish,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * time.Duration(180-30*rogue.Talents.Elusiveness),
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Pause auto attacks
			rogue.AutoAttacks.CancelAutoSwing(sim)
			// Apply stealth
			rogue.StealthAura.Activate(sim)
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    rogue.Vanish,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityDrums,

		ShouldActivate: func(sim *core.Simulation, unit *core.Character) bool {
			if rogue.Talents.Overkill {
				return !(rogue.StealthAura.IsActive() || rogue.OverkillAura.IsActive()) && rogue.CurrentEnergy() > 50
			}
			if rogue.Spec == proto.Spec_SpecSubtletyRogue { // Master of Subtlety is now a Subtlety rogue passive
				if rogue.MasterOfSubtletyAura.IsActive() {
					return false // possible after preparation
				}

				wantPremed, premedCPs := checkPremediation(sim, rogue)
				if wantPremed && premedCPs == 0 {
					return false // essentially sync with premed if possible
				}

				if rogue.CurrentEnergy() < rogue.Ambush.DefaultCast.Cost {
					return false
				}

				return rogue.ComboPoints()+premedCPs <= 5 // heuristically, "<= 5" is too strict (since omitting premed is fine)
			}

			return false
		},
	})
}

func checkPremediation(sim *core.Simulation, rogue *Rogue) (bool, int32) {
	if rogue.Premeditation == nil {
		return false, 0
	}

	if !rogue.Premeditation.IsReady(sim) {
		return false, 0
	}
	return true, 2
}
