package protection

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/paladin"
)

/*
Increases your total Stamina by 25% and your block chance by 10%.

Reduces the chance you will be critically hit by melee attacks by 6%.

Word of Glory is no longer on the global cooldown.

Your spell power is now equal to 50% of your attack power, and you no longer benefit from other sources of spell power.

Grants 15% of your maximum mana every 2 sec.
*/
func (prot *ProtectionPaladin) registerGuardedByTheLight() {
	actionID := core.ActionID{SpellID: 53592}

	oldGetSpellPowerValue := prot.GetSpellPowerValue
	newGetSpellPowerValue := func(spell *core.Spell) float64 {
		return math.Floor(spell.MeleeAttackPower() * 0.5)
	}

	core.MakePermanent(prot.RegisterAura(core.Aura{
		Label:      "Guarded by the Light" + prot.Label,
		ActionID:   actionID,
		BuildPhase: core.CharacterBuildPhaseTalents,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			prot.GetSpellPowerValue = newGetSpellPowerValue
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			prot.GetSpellPowerValue = oldGetSpellPowerValue
		},
	})).AttachStatDependency(
		prot.NewDynamicMultiplyStat(stats.Stamina, 1.25),
	).AttachAdditivePseudoStatBuff(
		&prot.PseudoStats.BaseBlockChance, 0.1,
	).AttachAdditivePseudoStatBuff(
		&prot.PseudoStats.ReducedCritTakenChance, 0.06,
	).AttachSpellMod(core.SpellModConfig{
		// Not in tooltip: Crusader Strike costs 80% less mana
		Kind:       core.SpellMod_PowerCost_Pct,
		ClassMask:  paladin.SpellMaskCrusaderStrike,
		FloatValue: -0.80,
	}).AttachSpellMod(core.SpellModConfig{
		// Not in tooltip: Judgmentcosts 40% less mana
		Kind:       core.SpellMod_PowerCost_Pct,
		ClassMask:  paladin.SpellMaskJudgment,
		FloatValue: -0.4,
	})

	manaMetrics := prot.NewManaMetrics(actionID)
	core.MakePermanent(prot.RegisterAura(core.Aura{
		Label: "Guarded by the Light Mana Regen" + prot.Label,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   time.Second * 2,
				Priority: core.ActionPriorityRegen,
				OnAction: func(*core.Simulation) {
					prot.AddMana(sim, 0.15*prot.MaxMana(), manaMetrics)
				},
			})
		},
	}))
}
