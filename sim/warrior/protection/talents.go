package protection

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
	"github.com/wowsims/cata/sim/warrior"
)

func (war *ProtectionWarrior) ApplyTalents() {
	// Vigilance is not implemented as it requires a second friendly target
	// We can probably fake it or make it configurable or something, but I expect it wouldn't
	// make much of a difference as I think tanks getting hit by bosses will be sitting at or near
	// their max vengeance bonus pretty much all the time

	war.Warrior.ApplyCommonTalents()

	war.RegisterConcussionBlow()
	war.RegisterDevastate()
	war.RegisterLastStand()
	war.RegisterShockwave()

	war.applyBastionOfDefense()
	war.applyHeavyRepercussions()
	war.applyImpendingVictory()
	war.applyImprovedRevenge()
	war.applySwordAndBoard()
	war.applyThunderstruck()

	war.ApplyGlyphs()
}

func (war *ProtectionWarrior) applyBastionOfDefense() {
	if war.Talents.BastionOfDefense == 0 {
		return
	}

	damageDealtMultiplier := 1.0 + 0.05*float64(war.Talents.BastionOfDefense)
	enrageChance := 0.1 * float64(war.Talents.BastionOfDefense)

	enrageAura := war.GetOrRegisterAura(core.Aura{
		Label:    "Enrage",
		ActionID: core.ActionID{SpellID: 57516},
		Duration: 12 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= damageDealtMultiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= damageDealtMultiplier
		},
	})

	core.MakePermanent(war.GetOrRegisterAura(core.Aura{
		Label: "Enrage Trigger",
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Outcome.Matches(core.OutcomeBlock | core.OutcomeDodge | core.OutcomeParry) {
				if sim.Proc(enrageChance, "Enrage Trigger Chance") {
					enrageAura.Activate(sim)
				}
			}
		},
	}))
}

func (war *ProtectionWarrior) applyImprovedRevenge() {
	if war.Talents.ImprovedRevenge == 0 {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskRevenge,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.3 * float64(war.Talents.ImprovedRevenge),
	})

	// extra hit is implemented inside of revenge
}

func (war *ProtectionWarrior) applyImpendingVictory() {
	if war.Talents.ImpendingVictory == 0 {
		return
	}

	const vrReady = "Impending Victory"

	enableVRAura := war.RegisterAura(core.Aura{
		Label:    "Victorious",
		ActionID: core.ActionID{SpellID: 82368},
		Tag:      vrReady,

		Duration: 20 * time.Second,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if (spell.ClassSpellMask & warrior.SpellMaskVictoryRush) != 0 {
				aura.Deactivate(sim)
			}
		},
	})

	procChance := 0.25 * float64(war.Talents.ImpendingVictory)
	core.MakePermanent(war.RegisterAura(core.Aura{
		Label: "Impending Victory Trigger",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if (spell.ClassSpellMask&warrior.SpellMaskDevastate) == 0 || !result.Landed() {
				return
			}

			if spell.Unit.CurrentHealthPercent() <= 0.2 && sim.Proc(procChance, "Impending Victory") {
				enableVRAura.Activate(sim)
			}
		},
	}))

	// We register Victory Rush in here as this talent is the only way it can be used rotationally
	war.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 34428},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAPL | core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,
		ClassSpellMask: warrior.SpellMaskVictoryRush | warrior.SpellMaskSpecialAttack,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return war.HasActiveAuraWithTag(vrReady)
		},

		DamageMultiplier: 1.0,
		CritMultiplier:   war.DefaultMeleeCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.MeleeAttackPower() * 0.56
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}

func (war *ProtectionWarrior) applyThunderstruck() {
	if war.Talents.Thunderstruck == 0 {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskRend | warrior.SpellMaskCleave | warrior.SpellMaskThunderClap,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.03 * float64(war.Talents.Thunderstruck),
	})

	shockwaveBuff := war.AddDynamicMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskShockwave,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.0,
	})

	shockwaveBonus := 0.05 * float64(war.Talents.Thunderstruck)
	shockwaveAura := war.RegisterAura(core.Aura{
		Label:     "Thunderstruck",
		ActionID:  core.ActionID{SpellID: 87096},
		Duration:  20 * time.Second,
		MaxStacks: 3,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			if newStacks != 0 {
				bonus := shockwaveBonus * float64(newStacks)
				shockwaveBuff.UpdateFloatValue(bonus)
				shockwaveBuff.Activate()
			} else {
				shockwaveBuff.Deactivate()
			}
		},
	})

	core.MakePermanent(war.RegisterAura(core.Aura{
		Label: "Thunderstruck Trigger",

		// The shockwave buff is gained after any cast of Thunder Clap, even if it doesn't hit any targets
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if (spell.ClassSpellMask & warrior.SpellMaskThunderClap) == 0 {
				return
			}

			shockwaveAura.Activate(sim)
			shockwaveAura.AddStack(sim)
		},
	}))
}

func (war *ProtectionWarrior) applyHeavyRepercussions() {
	if war.Talents.HeavyRepercussions == 0 {
		return
	}

	damageMod := war.AddDynamicMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskShieldSlam,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.5 * float64(war.Talents.HeavyRepercussions),
	})

	buff := war.RegisterAura(core.Aura{
		Label:    "Heavy Repercussions",
		Duration: 10 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Deactivate()
		},
	})

	core.MakePermanent(war.RegisterAura(core.Aura{
		Label: "Heavy Repercussions Trigger",
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if (spell.ClassSpellMask & warrior.SpellMaskShieldBlock) != 0 {
				buff.Activate(sim)
			}
		},
	}))
}

func (war *ProtectionWarrior) applySwordAndBoard() {
	if war.Talents.SwordAndBoard == 0 {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskDevastate,
		Kind:       core.SpellMod_BonusCrit_Rating,
		FloatValue: 5 * float64(war.Talents.SwordAndBoard) * core.CritRatingPerCritChance,
	})

	costMod := war.AddDynamicMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskShieldSlam,
		Kind:       core.SpellMod_PowerCost_Pct,
		FloatValue: -1.0,
	})

	buffAura := war.RegisterAura(core.Aura{
		Label:    "Sword and Board",
		ActionID: core.ActionID{SpellID: 50227},
		Duration: 5 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			costMod.Activate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if (spell.ClassSpellMask & warrior.SpellMaskShieldSlam) != 0 {
				aura.Deactivate(sim)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			costMod.Deactivate()
		},
	})

	procChance := 0.1 * float64(war.Talents.SwordAndBoard)
	core.MakePermanent(war.RegisterAura(core.Aura{
		Label: "Sword and Board Trigger",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || (spell.ClassSpellMask&(warrior.SpellMaskDevastate|warrior.SpellMaskRevenge)) == 0 {
				return
			}

			if sim.Proc(procChance, "Sword and Board") {
				war.shieldSlam.CD.Reset()
				buffAura.Activate(sim)
			}
		},
	}))
}
