package rogue

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (rogue *Rogue) registerTricksOfTheTradeSpell() {
	hasGlyph := rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfTricksOfTheTrade)

	actionID := core.ActionID{SpellID: 57934}
	energyMetrics := rogue.NewEnergyMetrics(actionID)
	hasShadowblades := rogue.HasSetBonus(Tier10, 2)
	energyCost := core.TernaryFloat64(hasGlyph || hasShadowblades, 0, 15)

	var tottTarget *core.Unit
	if rogue.Options.TricksOfTheTradeTarget != nil {
		tottTarget = rogue.GetUnit(rogue.Options.TricksOfTheTradeTarget)
	}

	tricksOfTheTradeThreatTransferAura := rogue.GetOrRegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 59628},
		Label:    "TricksOfTheTradeThreatTransfer",
		Duration: core.TernaryDuration(hasGlyph, time.Second*10, time.Second*6),
	})

	tricksOfTheTradeDamageAura := rogue.NewAllyAuraArray(func(unit *core.Unit) *core.Aura {
		if unit.Type == core.PetUnit {
			return nil
		}
		return core.TricksOfTheTradeAura(unit, rogue.Index, hasGlyph)
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
				}
				aura.Deactivate(sim)
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.TricksOfTheTrade.CD.Set(core.NeverExpires)
			rogue.UpdateMajorCooldowns()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.TricksOfTheTrade.CD.Set(sim.CurrentTime + time.Second*time.Duration(30))
			rogue.UpdateMajorCooldowns()
		},
	})

	rogue.TricksOfTheTrade = rogue.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL | core.SpellFlagHelpful,

		EnergyCost: core.EnergyCostOptions{
			Cost: energyCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * time.Duration(30), // CD is handled by application aura
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if tottTarget != nil {
				castTarget = tottTarget
			} else if target.Type == core.PlayerUnit && target != &rogue.Unit { // Cant cast on ourself
				castTarget = target
			}
			tricksOfTheTradeApplicationAura.Activate(sim)
			if hasShadowblades {
				rogue.AddEnergy(sim, 15, energyMetrics)
			}
		},
	})

	rogue.AddMajorCooldown(core.MajorCooldown{
		Spell: rogue.TricksOfTheTrade,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
			return tottTarget != nil
		},
	})
}
