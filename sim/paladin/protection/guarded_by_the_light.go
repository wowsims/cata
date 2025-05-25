package protection

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/paladin"
)

func (prot *ProtectionPaladin) registerGuardedByTheLight() {
	actionID := core.ActionID{SpellID: 53592}
	manaMetrics := prot.NewManaMetrics(actionID)

	oldGetSpellPowerValue := prot.GetSpellPowerValue
	newGetSpellPowerValue := func(spell *core.Spell) float64 {
		return spell.MeleeAttackPower() * 0.5
	}

	core.MakePermanent(prot.RegisterAura(core.Aura{
		Label:      "Guarded by the Light" + prot.Label,
		ActionID:   actionID,
		BuildPhase: core.CharacterBuildPhaseTalents,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   time.Second * 2,
				Priority: core.ActionPriorityRegen,
				OnAction: func(*core.Simulation) {
					prot.AddMana(sim, 0.15*prot.MaxMana(), manaMetrics)
				},
			})

			prot.GetSpellPowerValue = newGetSpellPowerValue
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			prot.GetSpellPowerValue = oldGetSpellPowerValue
		},
	})).AttachStatDependency(
		prot.NewDynamicMultiplyStat(stats.Stamina, 1.25),
	).AttachMultiplicativePseudoStatBuff(
		&prot.PseudoStats.BaseBlockChance,
		0.1,
	).AttachAdditivePseudoStatBuff(
		&prot.PseudoStats.ReducedCritTakenChance,
		0.06,
	).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		ClassMask: paladin.SpellMaskCrusaderStrike,
		IntValue:  -80,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		ClassMask: paladin.SpellMaskJudgment,
		IntValue:  -40,
	})
}
