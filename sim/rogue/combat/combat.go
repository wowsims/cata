package combat

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/rogue"
)

func RegisterCombatRogue() {
	core.RegisterAgentFactory(
		proto.Player_CombatRogue{},
		proto.Spec_SpecCombatRogue,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewCombatRogue(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_CombatRogue)
			if !ok {
				panic("Invalid spec value for Combat Rogue!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewCombatRogue(character *core.Character, options *proto.Player) *CombatRogue {
	combatOptions := options.GetCombatRogue().Options

	combatRogue := &CombatRogue{
		Rogue: rogue.NewRogue(character, combatOptions.ClassOptions, options.TalentsString),
	}
	combatRogue.CombatOptions = combatOptions

	return combatRogue
}

func (combatRogue *CombatRogue) Initialize() {
	combatRogue.Rogue.Initialize()

	combatRogue.AutoAttacks.OHConfig().DamageMultiplier *= 1.75
	// combatRogue.registerKillingSpreeCD()
}

type CombatRogue struct {
	*rogue.Rogue
}

func (combatRogue *CombatRogue) GetRogue() *rogue.Rogue {
	return combatRogue.Rogue
}

func (combatRogue *CombatRogue) Reset(sim *core.Simulation) {
	combatRogue.Rogue.Reset(sim)
}

// func (rogue *Rogue) applyCombatPotency() {
// 	if rogue.Talents.CombatPotency == 0 {
// 		return
// 	}

// 	const procChance = 0.2
// 	energyBonus := 3.0 * float64(rogue.Talents.CombatPotency)
// 	energyMetrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: 35553})

// 	rogue.RegisterAura(core.Aura{
// 		Label:    "Combat Potency",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			// from 3.0.3 patch notes: "Combat Potency: Now only works with auto attacks"
// 			if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeOHAuto) {
// 				return
// 			}

// 			if sim.RandomFloat("Combat Potency") < procChance {
// 				rogue.AddEnergy(sim, energyBonus, energyMetrics)
// 			}
// 		},
// 	})
// }

// var BladeFlurryActionID = core.ActionID{SpellID: 13877}
// var BladeFlurryHitID = core.ActionID{SpellID: 22482}

// func (rogue *Rogue) registerBladeFlurryCD() {
// 	if !rogue.Talents.BladeFlurry {
// 		return
// 	}

// 	var curDmg float64
// 	bfHit := rogue.RegisterSpell(core.SpellConfig{
// 		ActionID:    BladeFlurryHitID,
// 		SpellSchool: core.SpellSchoolPhysical,
// 		ProcMask:    core.ProcMaskEmpty, // No proc mask, so it won't proc itself.
// 		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreAttackerModifiers,

// 		DamageMultiplier: 1,
// 		ThreatMultiplier: 1,

// 		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 			spell.CalcAndDealDamage(sim, target, curDmg, spell.OutcomeAlwaysHit)
// 		},
// 	})

// 	const hasteBonus = 1.2
// 	const inverseHasteBonus = 1 / 1.2

// 	dur := time.Second * 15

// 	rogue.BladeFlurryAura = rogue.RegisterAura(core.Aura{
// 		Label:    "Blade Flurry",
// 		ActionID: BladeFlurryActionID,
// 		Duration: dur,
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			rogue.MultiplyMeleeSpeed(sim, hasteBonus)
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			rogue.MultiplyMeleeSpeed(sim, inverseHasteBonus)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if sim.GetNumTargets() < 2 {
// 				return
// 			}
// 			if result.Damage == 0 || !spell.ProcMask.Matches(core.ProcMaskMelee) {
// 				return
// 			}
// 			// Fan of Knives off-hand hits are not cloned
// 			if spell.IsSpellAction(FanOfKnivesSpellID) && spell.ProcMask.Matches(core.ProcMaskMeleeOH) {
// 				return
// 			}

// 			// Undo armor reduction to get the raw damage value.
// 			curDmg = result.Damage / result.ResistanceMultiplier

// 			bfHit.Cast(sim, rogue.Env.NextTargetUnit(result.Target))
// 			bfHit.SpellMetrics[result.Target.UnitIndex].Casts--
// 		},
// 	})

// 	cooldownDur := time.Minute * 2
// 	rogue.BladeFlurry = rogue.RegisterSpell(core.SpellConfig{
// 		ActionID: BladeFlurryActionID,
// 		Flags:    core.SpellFlagAPL,

// 		EnergyCost: core.EnergyCostOptions{
// 			Cost: core.TernaryFloat64(rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfBladeFlurry), 0, 25),
// 		},
// 		Cast: core.CastConfig{
// 			DefaultCast: core.Cast{
// 				GCD: time.Second,
// 			},
// 			IgnoreHaste: true,
// 			CD: core.Cooldown{
// 				Timer:    rogue.NewTimer(),
// 				Duration: cooldownDur,
// 			},
// 		},

// 		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
// 			rogue.BreakStealth(sim)
// 			rogue.BladeFlurryAura.Activate(sim)
// 		},
// 	})

// 	rogue.AddMajorCooldown(core.MajorCooldown{
// 		Spell:    rogue.BladeFlurry,
// 		Type:     core.CooldownTypeDPS,
// 		Priority: core.CooldownPriorityDefault,
// 		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
// 			if sim.GetRemainingDuration() > cooldownDur+dur {
// 				// We'll have enough time to cast another BF, so use it immediately to make sure we get the 2nd one.
// 				return true
// 			}

// 			// Since this is our last BF, wait until we have SND / procs up.
// 			sndTimeRemaining := rogue.SliceAndDiceAura.RemainingDuration(sim)
// 			// TODO: Wait for dst/mongoose procs
// 			return sndTimeRemaining >= time.Second
// 		},
// 	})
// }

// var AdrenalineRushActionID = core.ActionID{SpellID: 13750}

// func (rogue *Rogue) registerAdrenalineRushCD() {
// 	if !rogue.Talents.AdrenalineRush {
// 		return
// 	}

// 	rogue.AdrenalineRushAura = rogue.RegisterAura(core.Aura{
// 		Label:    "Adrenaline Rush",
// 		ActionID: AdrenalineRushActionID,
// 		Duration: core.TernaryDuration(rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfAdrenalineRush), time.Second*20, time.Second*15),
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			rogue.ResetEnergyTick(sim)
// 			rogue.ApplyEnergyTickMultiplier(1.0)
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			rogue.ResetEnergyTick(sim)
// 			rogue.ApplyEnergyTickMultiplier(-1.0)
// 		},
// 	})

// 	adrenalineRushSpell := rogue.RegisterSpell(core.SpellConfig{
// 		ActionID: AdrenalineRushActionID,

// 		Cast: core.CastConfig{
// 			DefaultCast: core.Cast{
// 				GCD: time.Second,
// 			},
// 			IgnoreHaste: true,
// 			CD: core.Cooldown{
// 				Timer:    rogue.NewTimer(),
// 				Duration: time.Minute * 3,
// 			},
// 		},

// 		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
// 			rogue.BreakStealth(sim)
// 			rogue.AdrenalineRushAura.Activate(sim)
// 		},
// 	})

// 	rogue.AddMajorCooldown(core.MajorCooldown{
// 		Spell:    adrenalineRushSpell,
// 		Type:     core.CooldownTypeDPS,
// 		Priority: core.CooldownPriorityBloodlust,
// 		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
// 			thresh := 45.0
// 			return rogue.CurrentEnergy() <= thresh
// 		},
// 	})
// }

// func (rogue *Rogue) registerKillingSpreeCD() {
// 	if !rogue.Talents.KillingSpree {
// 		return
// 	}
// 	rogue.registerKillingSpreeSpell()
// }
