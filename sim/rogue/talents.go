package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func (rogue *Rogue) ApplyTalents() {
	rogue.ApplyArmorSpecializationEffect(stats.Agility, proto.ArmorType_ArmorTypeLeather, 87504)

	// Hotfix Passive: https://www.wowhead.com/mop-classic/spell=137034/hotfix-passive
	rogue.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  RogueSpellAmbush,
		FloatValue: 0.12,
	})

	// Nightstalker
	if rogue.Talents.Nightstalker {
		rogue.NightstalkerMod = rogue.AddDynamicMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Pct,
			ClassMask:  RogueSpellsAll,
			FloatValue: 0.5,
		})
	}

	// Subterfuge
	if rogue.Talents.Subterfuge {
		rogue.SubterfugeAura = rogue.RegisterAura(core.Aura{
			Label:    "Subterfuge",
			Duration: time.Second * 3,
			ActionID: core.ActionID{SpellID: 108208},
		})
	}

	// Shadow Focus
	if rogue.Talents.ShadowFocus {
		rogue.ShadowFocusMod = rogue.AddDynamicMod(core.SpellModConfig{
			Kind:       core.SpellMod_PowerCost_Pct,
			ClassMask:  RogueSpellsAll,
			FloatValue: -0.75,
		})
	}

	// Marked for Death
	if rogue.Talents.MarkedForDeath {
		mfdMetrics := rogue.NewComboPointMetrics(core.ActionID{SpellID: 137619})

		mfdSpell := rogue.RegisterSpell(core.SpellConfig{
			ActionID:       core.ActionID{SpellID: 137619},
			Flags:          core.SpellFlagAPL,
			ClassSpellMask: RogueSpellMarkedForDeath,

			Cast: core.CastConfig{
				IgnoreHaste: true,
				CD: core.Cooldown{
					Timer:    rogue.NewTimer(),
					Duration: time.Minute * 1,
				},
			},
			ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
				rogue.AddComboPoints(sim, 5, mfdMetrics)
			},
		})

		rogue.AddMajorCooldown(core.MajorCooldown{
			Spell:    mfdSpell,
			Type:     core.CooldownTypeDPS,
			Priority: core.CooldownPriorityDefault,
		})
	}

	// Anticipation
	if rogue.Talents.Anticipation {
		action := core.ActionID{SpellID: 114015}
		antiMetrics := rogue.NewComboPointMetrics(action)

		rogue.AnticipationAura = rogue.RegisterAura(core.Aura{
			Label:     "Anticipation",
			ActionID:  action,
			Duration:  time.Second * 15,
			MaxStacks: 5,

			// Adding stacks is driven by rogue.AddComboPointsOrAnticipation()

			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Landed() && spell.Flags.Matches(SpellFlagFinisher) {
					rogue.AddComboPoints(sim, aura.GetStacks(), antiMetrics)
					aura.SetStacks(sim, 0)
					aura.Deactivate(sim)
				}
			},
		})
	}
}

func (rogue *Rogue) ApplyCutToTheChase(sim *core.Simulation) {
	if rogue.Spec == proto.Spec_SpecAssassinationRogue && rogue.SliceAndDiceAura.IsActive() {
		rogue.SliceAndDiceAura.Duration = rogue.sliceAndDiceDurations[5]
		rogue.SliceAndDiceAura.Activate(sim)
	}
}
