package protection

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *ProtectionWarrior) ApplyTalents() {
	war.ApplyArmorSpecializationEffect(stats.Stamina, proto.ArmorType_ArmorTypePlate, 86526)

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
	actionID := core.ActionID{SpellID: 57516}
	enrageAura := war.GetOrRegisterAura(core.Aura{
		Label:    "Enrage",
		ActionID: actionID,
		Duration: 12 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= damageDealtMultiplier
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= damageDealtMultiplier
		},
	})

	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:       "Enrage Trigger",
		ActionID:   actionID,
		Callback:   core.CallbackOnSpellHitTaken,
		Outcome:    core.OutcomeBlock | core.OutcomeDodge | core.OutcomeParry,
		ProcChance: enrageChance,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			enrageAura.Activate(sim)
		},
	})
}

func (war *ProtectionWarrior) applyImprovedRevenge() {
	if war.Talents.ImprovedRevenge == 0 {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskRevenge,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.3 * float64(war.Talents.ImprovedRevenge),
	})

	// extra hit is implemented inside of revenge
}

func (war *ProtectionWarrior) applyImpendingVictory() {
	if war.Talents.ImpendingVictory == 0 {
		return
	}

	const vrReady = "Impending Victory"
	actionID := core.ActionID{SpellID: 82368}
	enableVRAura := war.RegisterAura(core.Aura{
		Label:    "Victorious",
		ActionID: actionID,
		Tag:      vrReady,

		Duration: 20 * time.Second,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if (spell.ClassSpellMask & warrior.SpellMaskVictoryRush) != 0 {
				aura.Deactivate(sim)
			}
		},
	})

	procChance := 0.25 * float64(war.Talents.ImpendingVictory)
	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:           "Impending Victory Trigger",
		ActionID:       actionID,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ProcChance:     procChance,
		ClassSpellMask: warrior.SpellMaskDevastate,
		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			return spell.Unit.CurrentHealthPercent() <= 0.2
		},
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			enableVRAura.Activate(sim)
		},
	})

	// We register Victory Rush in here as this talent is the only way it can be used rotationally
	war.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 34428},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagAPL | core.SpellFlagMeleeMetrics,
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
		CritMultiplier:   war.DefaultCritMultiplier(),

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
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.03 * float64(war.Talents.Thunderstruck),
	})

	shockwaveBuff := war.AddDynamicMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskShockwave,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.0,
	})

	actionID := core.ActionID{SpellID: 87096}
	shockwaveBonus := 0.05 * float64(war.Talents.Thunderstruck)
	shockwaveAura := war.RegisterAura(core.Aura{
		Label:     "Thunderstruck",
		ActionID:  actionID,
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

	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:           "Thunderstruck Trigger",
		ActionID:       actionID,
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: warrior.SpellMaskThunderClap,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			shockwaveAura.Activate(sim)
			shockwaveAura.AddStack(sim)
		},
	})
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

	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:           "Heavy Repercussions Trigger",
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: warrior.SpellMaskShieldBlock,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			buff.Activate(sim)
		},
	})
}

func (war *ProtectionWarrior) applySwordAndBoard() {
	if war.Talents.SwordAndBoard == 0 {
		return
	}

	war.AddStaticMod(core.SpellModConfig{
		ClassMask:  warrior.SpellMaskDevastate,
		Kind:       core.SpellMod_BonusCrit_Percent,
		FloatValue: 5 * float64(war.Talents.SwordAndBoard),
	})

	costMod := war.AddDynamicMod(core.SpellModConfig{
		ClassMask: warrior.SpellMaskShieldSlam,
		Kind:      core.SpellMod_PowerCost_Pct,
		IntValue:  -100,
	})

	actionID := core.ActionID{SpellID: 50227}
	buffAura := war.RegisterAura(core.Aura{
		Label:    "Sword and Board",
		ActionID: actionID,
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
	core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:           "Heavy Repercussions Trigger",
		ActionID:       actionID,
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: warrior.SpellMaskDevastate | warrior.SpellMaskRevenge,
		Outcome:        core.OutcomeLanded,
		ProcChance:     procChance,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			war.shieldSlam.CD.Reset()
			buffAura.Activate(sim)
		},
	})
}
