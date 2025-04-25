package combat

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/rogue"
)

func (comRogue *CombatRogue) registerBanditsGuile() {
	chanceToProc := []float64{0.0, 0.33, 0.67, 1.0}[comRogue.Talents.BanditsGuile]
	attackCounter := int32(0)
	var lastAttacked *core.Unit
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

		damageBonus := []float64{1.1, 1.2, 1.3}[index]

		bgDamageAuras[index] = comRogue.RegisterAura(core.Aura{
			Label:    label,
			ActionID: actionID,
			Duration: time.Second * 15,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				core.EnableDamageDoneByCaster(DDBC_BanditsGuile, DDBC_Total, comRogue.AttackTables[comRogue.CurrentTarget.UnitIndex], func(sim *core.Simulation, spell *core.Spell, attackTable *core.AttackTable) float64 {
					if spell.Matches(rogue.RogueSpellsAll) || spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
						return damageBonus
					}
					return 1.0
				})
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				core.DisableDamageDoneByCaster(0, comRogue.AttackTables[comRogue.CurrentTarget.UnitIndex])
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
				if lastAttacked != result.Target {
					// Reset back to no insight, no casts
					attackCounter = 0
					if currentInsightIndex >= 0 {
						bgDamageAuras[currentInsightIndex].Deactivate(sim)
					}
					currentInsightIndex = -1
				}
				lastAttacked = result.Target

				if chanceToProc == 1 || sim.Proc(chanceToProc, "Bandit's Guile") {
					attackCounter += 1
					comRogue.BanditsGuileAura.SetStacks(sim, attackCounter+1)
					if attackCounter == 4 {
						attackCounter = 0
						comRogue.BanditsGuileAura.SetStacks(sim, attackCounter+1)
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
				}

			}
		},
	})
}
