package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (warlock *Warlock) ApplyDestructionTalents() {
	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask: WarlockSpellShadowBolt | WarlockSpellChaosBolt | WarlockSpellImmolate,
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: -1 * time.Duration([]int{0, 100, 300, 500}[warlock.Talents.Bane]) * time.Millisecond,
	})

	//TODO: Add/Mult?
	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask:  WarlockSpellShadowBolt | WarlockSpellChaosBolt,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: []float64{0.0, 1.04, 1.08, 1.12}[warlock.Talents.ShadowAndFlame],
	})

	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask:  WarlockSpellImmolate,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: []float64{1.0, 1.1, 1.2}[warlock.Talents.ImprovedImmolate],
	})

	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask: WarlockSpellSoulFire,
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: -1 * time.Duration([]int{0, 500, 1000}[warlock.Talents.Emberstorm]) * time.Millisecond,
	})

	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask: WarlockSpellIncinerate,
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: -1 * time.Duration([]int{0, 130, 250}[warlock.Talents.Emberstorm]) * time.Millisecond,
	})

	//TODO: Improved Searing Pain

	//TODO: Improved Soul Fire

	//TODO: Backdraft
	// warlock.setupBackdraft()

	//TODO: Shadowburn

	//TODO: Burning Embers

	//TODO: Soul Leech
	// warlock.setupImprovedSoulLeech()

	//BACKLASH NA
	//NETHERWARD NA

	//TODO: FireAndBrimstone inc & chaos bolt damage is done in immolate, not sure how to do this as a spell mod
	warlock.AddStaticMod(core.SpellModConfig{
		ClassMask:  WarlockSpellConflagrate,
		Kind:       core.SpellMod_BonusCrit_Rating,
		FloatValue: 5.0 * float64(warlock.Talents.FireAndBrimstone) * core.CritRatingPerCritChance,
	})

	// NETHER PROTECTION NA

	// TODO: EMPOWERED IMP
	// warlock.setupEmpoweredImp()

	// TODO: BANE OF HAVOC

	if warlock.Talents.ChaosBolt {
		warlock.registerChaosBoltSpell()
	}
}

// func (warlock *Warlock) setupEmpoweredImp() {
// 	if warlock.Talents.EmpoweredImp <= 0 || warlock.Options.Summon != proto.WarlockOptions_Imp {
// 		return
// 	}

// 	warlock.Pet.PseudoStats.DamageDealtMultiplier *= 1.0 + 0.1*float64(warlock.Talents.EmpoweredImp)

// 	var affectedSpells []*core.Spell
// 	warlock.EmpoweredImpAura = warlock.RegisterAura(core.Aura{
// 		Label:    "Empowered Imp Proc Aura",
// 		ActionID: core.ActionID{SpellID: 47283},
// 		Duration: time.Second * 8,
// 		OnInit: func(aura *core.Aura, sim *core.Simulation) {
// 			affectedSpells = core.FilterSlice([]*core.Spell{
// 				warlock.Immolate,
// 				warlock.ShadowBolt,
// 				warlock.Incinerate,
// 				warlock.Shadowburn,
// 				warlock.SoulFire,
// 				warlock.ChaosBolt,
// 				warlock.SearingPain,
// 				// missing: shadowfury, shadowflame, seed explosion (not dot)
// 				//          rain of fire (consumes proc on cast start, but doesn't increase crit, ticks
// 				//          also consume the proc but do seem to benefit from the increaesed crit)
// 			}, func(spell *core.Spell) bool { return spell != nil })
// 		},
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			for _, spell := range affectedSpells {
// 				spell.BonusCritRating += 100 * core.CritRatingPerCritChance
// 			}
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			for _, spell := range affectedSpells {
// 				spell.BonusCritRating -= 100 * core.CritRatingPerCritChance
// 			}
// 		},
// 		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
// 			if slices.Contains(affectedSpells, spell) {
// 				aura.Deactivate(sim)
// 			}
// 		},
// 	})

// 	warlock.Pet.RegisterAura(core.Aura{
// 		Label:    "Empowered Imp Hidden Aura",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if result.DidCrit() {
// 				warlock.EmpoweredImpAura.Activate(sim)
// 			}
// 		},
// 	})
// }

// func (warlock *Warlock) setupBackdraft() {
// 	if warlock.Talents.Backdraft <= 0 {
// 		return
// 	}

// 	castTimeModifier := 0.1 * float64(warlock.Talents.Backdraft)
// 	var affectedSpells []*core.Spell

// 	warlock.BackdraftAura = warlock.RegisterAura(core.Aura{
// 		Label:     "Backdraft Proc Aura",
// 		ActionID:  core.ActionID{SpellID: 54277},
// 		Duration:  time.Second * 15,
// 		MaxStacks: 3,
// 		OnInit: func(aura *core.Aura, sim *core.Simulation) {
// 			affectedSpells = core.FilterSlice([]*core.Spell{
// 				warlock.Incinerate,
// 				warlock.SoulFire,
// 				warlock.ShadowBolt,
// 				warlock.ChaosBolt,
// 				warlock.Immolate,
// 				warlock.SearingPain,
// 			}, func(spell *core.Spell) bool { return spell != nil })
// 		},
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			for _, destroSpell := range affectedSpells {
// 				destroSpell.CastTimeMultiplier -= castTimeModifier
// 				destroSpell.DefaultCast.GCD = time.Duration(float64(destroSpell.DefaultCast.GCD) * (1 - castTimeModifier))
// 			}
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			for _, destroSpell := range affectedSpells {
// 				destroSpell.CastTimeMultiplier += castTimeModifier
// 				destroSpell.DefaultCast.GCD = time.Duration(float64(destroSpell.DefaultCast.GCD) / (1 - castTimeModifier))
// 			}
// 		},
// 		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
// 			if slices.Contains(affectedSpells, spell) {
// 				aura.RemoveStack(sim)
// 			}
// 		},
// 	})

// 	warlock.RegisterAura(core.Aura{
// 		Label:    "Backdraft Hidden Aura",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if spell == warlock.Conflagrate && result.Landed() {
// 				warlock.BackdraftAura.Activate(sim)
// 				warlock.BackdraftAura.SetStacks(sim, 3)
// 			}
// 		},
// 	})
// }

// func (warlock *Warlock) setupImprovedSoulLeech() {
// 	if warlock.Talents.ImprovedSoulLeech <= 0 {
// 		return
// 	}

// 	soulLeechProcChance := 0.1 * float64(warlock.Talents.SoulLeech)
// 	impSoulLeechProcChance := float64(warlock.Talents.ImprovedSoulLeech) / 2.
// 	actionID := core.ActionID{SpellID: 54118}
// 	impSoulLeechManaMetric := warlock.NewManaMetrics(actionID)
// 	var impSoulLeechPetManaMetric *core.ResourceMetrics
// 	if warlock.Pet != nil {
// 		impSoulLeechPetManaMetric = warlock.Pet.NewManaMetrics(actionID)
// 	}
// 	replSrc := warlock.Env.Raid.NewReplenishmentSource(core.ActionID{SpellID: 54118})

// 	warlock.RegisterAura(core.Aura{
// 		Label:    "Improved Soul Leech Hidden Aura",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if (spell == warlock.Conflagrate || spell == warlock.ShadowBolt || spell == warlock.ChaosBolt ||
// 				spell == warlock.SoulFire || spell == warlock.Incinerate) && result.Landed() {
// 				if !sim.Proc(soulLeechProcChance, "SoulLeech") {
// 					return
// 				}

// 				restorePct := float64(warlock.Talents.ImprovedSoulLeech) / 100
// 				warlock.AddMana(sim, warlock.MaxMana()*restorePct, impSoulLeechManaMetric)
// 				pet := warlock.Pet
// 				if pet != nil {
// 					pet.AddMana(sim, pet.MaxMana()*restorePct, impSoulLeechPetManaMetric)
// 				}

// 				if sim.Proc(impSoulLeechProcChance, "ImprovedSoulLeech") {
// 					warlock.Env.Raid.ProcReplenishment(sim, replSrc)
// 				}
// 			}
// 		},
// 	})
// }
// }
