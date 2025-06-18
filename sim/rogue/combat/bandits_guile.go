package combat

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/rogue"
)

func (comRogue *CombatRogue) registerBanditsGuile() {
	attackCounter := int32(0)
	var bgDamageAuras [3]*core.Aura
	currentInsightIndex := -1

	for index := 0; index < 3; index++ {
		var label string
		var actionID core.ActionID
		switch index {
		case 0:
			label = "Shallow Insight"
			actionID = core.ActionID{SpellID: 84745}
		case 1:
			label = "Moderate Insight"
			actionID = core.ActionID{SpellID: 84746}
		case 2:
			label = "Deep Insight"
			actionID = core.ActionID{SpellID: 84747}
		}

		damageBonus := []float64{0.1, 0.2, 0.3}[index]

		bgDamageMod := comRogue.AddDynamicMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Pct,
			ClassMask:  rogue.RogueSpellsAll,
			FloatValue: damageBonus,
		})

		bgDamageAuras[index] = comRogue.RegisterAura(core.Aura{
			Label:    label,
			ActionID: actionID,
			Duration: time.Second * 15,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				comRogue.AutoAttacks.MHAuto().DamageMultiplier *= (1 + damageBonus)
				comRogue.AutoAttacks.OHAuto().DamageMultiplier *= (1 + damageBonus)
				bgDamageMod.Activate()
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				comRogue.AutoAttacks.MHAuto().DamageMultiplier /= (1 + damageBonus)
				comRogue.AutoAttacks.OHAuto().DamageMultiplier /= (1 + damageBonus)
				bgDamageMod.Deactivate()
				if currentInsightIndex == 2 {
					currentInsightIndex = -1
					attackCounter = 0
				}
			},
		})
	}

	comRogue.BanditsGuileAura = comRogue.RegisterAura(core.Aura{
		Label:     "Bandit's Guile Tracker",
		ActionID:  core.ActionID{SpellID: 84654},
		Duration:  core.NeverExpires,
		MaxStacks: 4,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			currentInsightIndex = -1
			attackCounter = 0
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if currentInsightIndex < 2 && result.Landed() && (spell == comRogue.SinisterStrike || spell == comRogue.RevealingStrike) {
				attackCounter += 1

				if attackCounter == 4 {
					attackCounter = 0
					// Deactivate previous aura
					if currentInsightIndex >= 0 {
						bgDamageAuras[currentInsightIndex].Deactivate(sim)
					}
					currentInsightIndex += 1
					// Activate next aura
					bgDamageAuras[currentInsightIndex].Activate(sim)
				} else {
					// Refresh duration of existing aura
					if currentInsightIndex >= 0 {
						bgDamageAuras[currentInsightIndex].Duration = time.Second * 15
						bgDamageAuras[currentInsightIndex].Activate(sim)
					}
				}

				comRogue.BanditsGuileAura.SetStacks(sim, attackCounter+1)
			}
		},
	})
}
