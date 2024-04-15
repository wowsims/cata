package death_knight

import (
	//"github.com/wowsims/cata/sim/core/proto"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (dk *DeathKnight) ApplyFrostTalents() {

	// Nerves Of Cold Steel
	if dk.HasMHWeapon() && dk.HasOHWeapon() && dk.Talents.NervesOfColdSteel > 0 {
		dk.applyNervesOfColdSteel()
	}

	// Annihilation
	if dk.Talents.Annihilation > 0 {
		dk.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  DeathKnightSpellObliterate,
			FloatValue: 0.15 * float64(dk.Talents.Annihilation),
		})
	}

	// Chill of the Grave
	dk.applyChillOfTheGrave()

	// Killing Machine
	dk.applyKillingMachine()

	// Merciless Combat
	dk.applyMercilessCombat()

	// Rime
	dk.applyRime()

	// Improved Icy Talons
	if dk.Talents.ImprovedIcyTalons {
		dk.PseudoStats.MeleeSpeedMultiplier *= 1.05
	}
}

const DeathKnightChillOfTheGrave = DeathKnightSpellIcyTouch | DeathKnightSpellHowlingBlast | DeathKnightSpellObliterate

func (dk *DeathKnight) applyChillOfTheGrave() {
	if dk.Talents.ChillOfTheGrave == 0 {
		return
	}

	rpAmount := 5.0 * float64(dk.Talents.ChillOfTheGrave)
	rpMetric := dk.NewRunicPowerMetrics(core.ActionID{SpellID: 50115})
	core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
		Name:           "Chill of the Grave",
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: DeathKnightChillOfTheGrave,
		Outcome:        core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// Dont trigger on Obliterate Off hand
			if spell.ClassSpellMask == DeathKnightSpellObliterate && spell.ProcMask.Matches(core.ProcMaskMeleeOH) {
				return
			}

			dk.AddRunicPower(sim, rpAmount, rpMetric)
		},
	})
}

func (dk *DeathKnight) applyNervesOfColdSteel() {
	dk.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*float64(dk.Talents.NervesOfColdSteel))

	dk.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: []float64{0.0, 0.08, 0.16, 0.25}[dk.Talents.NervesOfColdSteel],
		ProcMask:   core.ProcMaskMeleeOH,
	})
}

const DeathKnightSpellMercilessCombat = DeathKnightSpellIcyTouch | DeathKnightSpellObliterate | DeathKnightSpellFrostStrike | DeathKnightSpellHowlingBlast

func (dk *DeathKnight) mercilessCombatMultiplier(sim *core.Simulation, spell *core.Spell, _ *core.AttackTable) float64 {
	if spell.ClassSpellMask&(DeathKnightSpellMercilessCombat) == 0 {
		return 1.0
	}
	return 1.0 + 0.06*float64(dk.Talents.MercilessCombat)
}

func (dk *DeathKnight) applyMercilessCombat() {
	if dk.Talents.MercilessCombat == 0 {
		return
	}

	debuffs := dk.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		aura := target.GetOrRegisterAura(core.Aura{
			Label:    "Merciless Combat" + dk.Label,
			ActionID: core.ActionID{SpellID: 49538},
			Duration: core.NeverExpires,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				dk.AttackTables[aura.Unit.UnitIndex].DamageDoneByCasterMultiplier = dk.mercilessCombatMultiplier
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				dk.AttackTables[aura.Unit.UnitIndex].DamageDoneByCasterMultiplier = nil
			},
		})
		return aura
	})

	core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
		Name:           "Merciless Combat Proc",
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: DeathKnightSpellMercilessCombat,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if sim.IsExecutePhase35() {
				debuffs.Get(result.Target).Activate(sim)
			}
		},
	})
}

func (dk *DeathKnight) applyRime() {
	if dk.Talents.Rime == 0 {
		return
	}

	// TODO:
}

func (dk *DeathKnight) applyKillingMachine() {
	if dk.Talents.KillingMachine == 0 {
		return
	}

	// TODO:
}

func (dk *DeathKnight) ThreatOfThassarianProc(sim *core.Simulation, result *core.SpellResult, ohSpell *core.Spell) {
	if dk.Talents.ThreatOfThassarian == 0 || dk.GetOHWeapon() == nil {
		return
	}
	if sim.Proc([]float64{0.0, 0.3, 0.6, 1.0}[dk.Talents.ThreatOfThassarian], "Threat of Thassarian") {
		ohSpell.Cast(sim, result.Target)
	}
}
