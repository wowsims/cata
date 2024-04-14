package death_knight

import (
	//"github.com/wowsims/cata/sim/core/proto"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (dk *DeathKnight) ApplyFrostTalents() {

	// Nerves Of Cold Steel
	if dk.nervesOfColdSteelActive() {
		dk.AddStat(stats.MeleeHit, core.MeleeHitRatingPerHitChance*float64(dk.Talents.NervesOfColdSteel))
		dk.AutoAttacks.OHConfig().DamageMultiplier *= dk.nervesOfColdSteelBonus()
	}

	// Icy Talons
	dk.applyIcyTalons()

	// Killing Machine
	dk.applyKillingMachine()

	// Rime
	dk.applyRime()
}

func (dk *DeathKnight) nervesOfColdSteelActive() bool {
	return dk.HasMHWeapon() && dk.HasOHWeapon()
}

func (dk *DeathKnight) nervesOfColdSteelBonus() float64 {
	return []float64{1.0, 1.08, 1.16, 1.25}[dk.Talents.NervesOfColdSteel]
}

func (dk *DeathKnight) mercilessCombatBonus(sim *core.Simulation) float64 {
	return core.TernaryFloat64(dk.Talents.MercilessCombat > 0 && sim.IsExecutePhase35(), 1.0+0.06*float64(dk.Talents.MercilessCombat), 1.0)
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

func (dk *DeathKnight) applyIcyTalons() {
	if !dk.Talents.ImprovedIcyTalons {
		return
	}

	// TODO:
}

func (dk *DeathKnight) threatOfThassarianProc(sim *core.Simulation, result *core.SpellResult, ohSpell *core.Spell) {
	if dk.Talents.ThreatOfThassarian == 0 || dk.GetOHWeapon() == nil {
		return
	}
	if sim.Proc([]float64{0.0, 0.3, 0.6, 1.0}[dk.Talents.ThreatOfThassarian], "Threat of Thassarian") {
		ohSpell.Cast(sim, result.Target)
	}
}
