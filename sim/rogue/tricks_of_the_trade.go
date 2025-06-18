package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (rogue *Rogue) registerTricksOfTheTradeSpell() {
	hasGlyph := rogue.HasMinorGlyph(proto.RogueMinorGlyph_GlyphOfTricksOfTheTrade)
	damageMult := core.TernaryFloat64(hasGlyph, 1.0, 1.15)
	actionID := core.ActionID{SpellID: 57934}

	var tottTarget *core.Unit
	if rogue.Options.TricksOfTheTradeTarget != nil {
		tottTarget = rogue.GetUnit(rogue.Options.TricksOfTheTradeTarget)
	}

	tricksOfTheTradeThreatTransferAura := rogue.GetOrRegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 59628},
		Label:    "TricksOfTheTradeThreatTransfer",
		Duration: time.Second * 6,
	})

	// Bogus Tricks threat "cast" for hooking T12/T13 set bonuses
	totThreatTransferSpell := rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 59628},
		ClassSpellMask: RogueSpellTricksOfTheTradeThreat,
	})

	tricksOfTheTradeDamageAura := rogue.NewAllyAuraArray(func(unit *core.Unit) *core.Aura {
		if unit.Type == core.PetUnit {
			return nil
		}
		return core.TricksOfTheTradeAura(unit, rogue.Index, damageMult)
	})

	var castTarget *core.Unit
	tricksOfTheTradeApplicationAura := rogue.GetOrRegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 57934},
		Label:    "TricksOfTheTradeApplication",
		Duration: 30 * time.Second,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				tricksOfTheTradeThreatTransferAura.Activate(sim)
				if castTarget != nil {
					tricksOfTheTradeDamageAura.Get(castTarget).Activate(sim)
					totThreatTransferSpell.Cast(sim, castTarget)
				} else {
					totThreatTransferSpell.Cast(sim, &rogue.Unit)
				}
				aura.Deactivate(sim)
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.TricksOfTheTrade.CD.Set(core.NeverExpires)
			rogue.UpdateMajorCooldowns()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.TricksOfTheTrade.CD.Set(sim.CurrentTime + time.Second*30)
			rogue.UpdateMajorCooldowns()
		},
	})

	rogue.TricksOfTheTrade = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
		ClassSpellMask: RogueSpellTricksOfTheTrade,

		EnergyCost: core.EnergyCostOptions{
			Cost: 15,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * 30, // CD is handled by application aura
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if tottTarget != nil {
				castTarget = tottTarget
			} else if target.Type == core.PlayerUnit && target != &rogue.Unit { // Cant cast on ourself
				castTarget = target
			}
			tricksOfTheTradeApplicationAura.Activate(sim)
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell: rogue.TricksOfTheTrade,
		Type:  core.CooldownTypeDPS,
	})
}
