package combat

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/rogue"
)

var BladeFlurryActionID = core.ActionID{SpellID: 13877}
var BladeFlurryHitID = core.ActionID{SpellID: 22482}

func (comRogue *CombatRogue) registerBladeFlurry() {
	var curDmg float64
	bfHit := comRogue.RegisterSpell(core.SpellConfig{
		ActionID:    BladeFlurryHitID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty, // No proc mask, so it won't proc itself.
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreModifiers | core.SpellFlagIgnoreArmor,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, curDmg, spell.OutcomeAlwaysHit)
		},
	})

	energyReduction := -0.2

	comRogue.BladeFlurryAura = comRogue.RegisterAura(core.Aura{
		Label:    "Blade Flurry",
		ActionID: BladeFlurryActionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			comRogue.ApplyAdditiveEnergyRegenBonus(sim, energyReduction)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			comRogue.ApplyAdditiveEnergyRegenBonus(sim, -energyReduction)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if sim.GetNumTargets() < 2 {
				return
			}
			if result.Damage == 0 || !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}
			// Fan of Knives is not cloned
			if spell.IsSpellAction(comRogue.FanOfKnives.SpellID) {
				return
			}

			curDmg = result.Damage * 0.4
			numHits := 0

			for enemyIndex := 0; enemyIndex < int(comRogue.Env.GetNumTargets()) && numHits < 4; enemyIndex++ {
				bfTarget := comRogue.Env.GetTargetUnit(int32(enemyIndex))
				if bfTarget != comRogue.CurrentTarget {
					numHits++
					bfHit.Cast(sim, bfTarget)
				}
			}
		},
	})

	comRogue.BladeFlurry = comRogue.RegisterSpell(core.SpellConfig{
		ActionID:       BladeFlurryActionID,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: rogue.RogueSpellBladeFlurry,

		Cast: core.CastConfig{
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    comRogue.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			comRogue.BladeFlurryAura.Activate(sim)
		},
	})
}
