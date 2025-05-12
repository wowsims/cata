package subtlety

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (subRogue *SubtletyRogue) registerHonorAmongThieves() {
	// When anyone in your group critically hits with a damage or healing spell or ability,
	// you have a 100% chance to gain a combo point on your current target.
	// This effect cannot occur more than once per 2 seconds.
	// Cannot trigger before combat starts

	procChance := 1.0
	comboMetrics := subRogue.NewComboPointMetrics(core.ActionID{SpellID: 51701})
	honorAmongThievesID := core.ActionID{SpellID: 51701}

	icd := core.Cooldown{
		Timer:    subRogue.NewTimer(),
		Duration: time.Second * 2,
	}

	maybeProc := func(sim *core.Simulation) {
		if icd.IsReady(sim) && sim.Proc(procChance, "Honor Among Thieves") {
			subRogue.AddComboPointsOrAnticipation(sim, 1, comboMetrics)

			if subRogue.T16EnergyAura != nil {
				subRogue.T16EnergyAura.Activate(sim)
				subRogue.T16EnergyAura.AddStack(sim)
			}

			icd.Use(sim)
		}
	}

	subRogue.HonorAmongThieves = core.MakePermanent(subRogue.RegisterAura(core.Aura{
		Label:    "Honor Among Thieves Combo Point Aura",
		ActionID: honorAmongThievesID,
		Icd:      &icd,
		OnGain: func(_ *core.Aura, sim *core.Simulation) {
			// In an ideal party, you'd probably get up to 6 ability crits/s (Rate = 600).
			//  Survival Hunters, Enhancement Shamans, and Assassination Rogues are particularly good.
			if subRogue.SubtletyOptions.HonorAmongThievesCritRate <= 0 {
				return
			}

			if subRogue.SubtletyOptions.HonorAmongThievesCritRate > 2000 {
				subRogue.SubtletyOptions.HonorAmongThievesCritRate = 2000 // limited, so performance doesn't suffer
			}

			rateToDuration := float64(time.Second) * 100 / float64(subRogue.SubtletyOptions.HonorAmongThievesCritRate)

			pa := &core.PendingAction{}
			pa.OnAction = func(sim *core.Simulation) {
				maybeProc(sim)
				pa.NextActionAt = sim.CurrentTime + time.Duration(sim.RandomExpFloat("next party crit")*rateToDuration)
				sim.AddPendingAction(pa)
			}
			pa.NextActionAt = sim.CurrentTime + time.Duration(sim.RandomExpFloat("next party crit")*rateToDuration)
			sim.AddPendingAction(pa)
		},
		OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidCrit() && !spell.ProcMask.Matches(core.ProcMaskMeleeMHAuto|core.ProcMaskMeleeOHAuto|core.ProcMaskRangedAuto) {
				maybeProc(sim)
			}
		},
		OnPeriodicDamageDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidCrit() {
				maybeProc(sim)
			}
		},
	}))
}
