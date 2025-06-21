package survival

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/hunter"
)

func (survival *SurvivalHunter) ApplyTalents() {
	survival.applyLNL()
	survival.ApplyMods()
	survival.Hunter.ApplyTalents()
}

func (survival *SurvivalHunter) ApplyMods() {
	survival.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  hunter.HunterSpellSerpentSting,
		FloatValue: 0.5,
	})
}

// Todo: Should we support precasting freezing/ice trap?
func (survival *SurvivalHunter) applyLNL() {
	actionID := core.ActionID{SpellID: 56343}
	procChance := core.TernaryFloat64(survival.CouldHaveSetBonus(hunter.YaunGolSlayersBattlegear, 4), 0.40, 0.20)

	icd := core.Cooldown{
		Timer:    survival.NewTimer(),
		Duration: time.Second * 10,
	}

	lnlAura := survival.RegisterAura(core.Aura{
		Icd:       &icd,
		Label:     "Lock and Load Proc",
		ActionID:  actionID,
		Duration:  time.Second * 12,
		MaxStacks: 2,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if survival.ExplosiveShot != nil {
				survival.ExplosiveShot.Cost.PercentModifier -= 100
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if survival.ExplosiveShot != nil {
				survival.ExplosiveShot.Cost.PercentModifier += 100
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == survival.ExplosiveShot {
				survival.ExplosiveShot.CD.Reset()
				// Weird check but..
				if !aura.Unit.HasActiveAura("Burning Adrenaline") {
					aura.RemoveStack(sim)
				}
			}
		},
	})

	survival.RegisterAura(core.Aura{
		Label:    "Lock and Load",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Matches(hunter.HunterSpellBlackArrow) {
				return
			}

			if !icd.IsReady(sim) {
				return
			}

			if sim.RandomFloat("Lock and Load") < procChance {
				icd.Use(sim)
				lnlAura.Activate(sim)
				lnlAura.SetStacks(sim, 2)
				if survival.ExplosiveShot != nil {
					survival.ExplosiveShot.CD.Reset()
				}
			}
		},
	})
}
