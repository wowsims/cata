package cata

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

/*
Base classes for Dragonwrath, Tarecgosa's Rest support
All relevant mechanics should be tracked and documented at:
https://github.com/wowsims/cata/issues/946
*/
const (
	supressNone   int8 = 0
	supressImpact int8 = 1 << iota
	supressDoT
	supressAll = supressImpact | supressDoT
)

type SpellHandler func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult)

type DragonwrathClassConfig struct {
	procChance   float64
	spellHandler map[int32]SpellHandler
	spellConfig  map[int32]int8
}

var classConfig = map[proto.Spec]*DragonwrathClassConfig{}
var globalSpellHandler = map[int32]SpellHandler{}

func (config *DragonwrathClassConfig) getDoTHandler(spellId int32) SpellHandler {
	if val, ok := config.spellHandler[spellId]; ok {
		return val
	}

	if val, ok := globalSpellHandler[spellId]; ok {
		return val
	}

	return defaultDoTHandler
}

func (config *DragonwrathClassConfig) getImpactHandler(spellId int32) SpellHandler {
	if val, ok := config.spellHandler[spellId]; ok {
		return val
	}

	if val, ok := globalSpellHandler[spellId]; ok {
		return val
	}

	return defaultSpellHandler
}

func (config *DragonwrathClassConfig) supressInternal(spellId int32, supressionType int8) {
	if _, ok := config.spellConfig[spellId]; !ok {
		config.spellConfig[spellId] = 0
	}

	config.spellConfig[spellId] |= supressionType
}

func (config *DragonwrathClassConfig) SupressImpact(spellId int32) {
	config.supressInternal(spellId, supressImpact)
}

func (config *DragonwrathClassConfig) SupressDoT(spellId int32) {
	config.supressInternal(spellId, supressDoT)
}

func (config *DragonwrathClassConfig) SupressSpell(spellId int32) {
	config.supressInternal(spellId, supressAll)
}

func GetDragonwrathDoTSpell(unit *core.Unit) *core.Spell {
	return unit.GetSpell(core.ActionID{SpellID: 101085})
}

func AddClassConfig(spec proto.Spec, procChance float64) *DragonwrathClassConfig {
	if classConfig[spec] != nil {
		panic("Class config already registered")
	}

	classConfig[spec] = &DragonwrathClassConfig{
		procChance:   procChance,
		spellHandler: map[int32]SpellHandler{},
		spellConfig:  map[int32]int8{},
	}

	return classConfig[spec]
}

func defaultDoTHandler(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	dotSpell := GetDragonwrathDoTSpell(spell.Unit)
	dotSpell.BonusSpellPower = result.Damage
	dotSpell.Cast(sim, result.Target)
}

func defaultSpellHandler(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	copySpell := spell.Unit.GetSpell(spell.WithTag(71086))
	if copySpell == nil {
		copySpell = spell.Unit.RegisterSpell(GetDRTSpellConfig(spell))
	}

	sim.AddPendingAction(
		&core.PendingAction{
			NextActionAt: sim.CurrentTime + time.Millisecond*200, // add slight delay
			Priority:     core.ActionPriorityAuto,
			OnAction: func(sim *core.Simulation) {
				copySpell.Cast(sim, result.Target)
			},
		},
	)
}

func GetDRTSpellConfig(spell *core.Spell) core.SpellConfig {
	baseConfig := core.SpellConfig{
		ActionID:     spell.WithTag(71086),
		SpellSchool:  spell.SpellSchool,
		ProcMask:     core.ProcMaskEmpty,
		ApplyEffects: spell.ApplyEffects,
		ManaCost:     core.ManaCostOptions{},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},
		CritMultiplier:           spell.Unit.Env.Raid.GetPlayerFromUnit(spell.Unit).GetCharacter().DefaultSpellCritMultiplier(),
		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		MissileSpeed:             spell.MissileSpeed,
		ClassSpellMask:           spell.ClassSpellMask,
		BonusCoefficient:         spell.BonusCoefficient,
		Flags:                    spell.Flags &^ core.SpellFlagAPL,
	}

	// copy BonusCoefficient as spells might have spell mods applied to them
	// we dont know the original in the base spell config so they might have been applied twice
	baseConfig.BonusCoefficient = spell.BonusCoefficient
	return baseConfig
}

func init() {
	core.NewItemEffect(71086, func(a core.Agent) {
		unit := &a.GetCharacter().Unit
		registerSpells(unit)

		lastTimestamp := time.Duration(0)
		spellList := map[int32]bool{}
		core.MakePermanent(unit.RegisterAura(core.Aura{
			ActionID: core.ActionID{ItemID: 71086},
			Label:    "Dragonwrath, Tarecgosa's Rest - Handler",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {

				// Handle direct damage only and make sure we're not proccing of our own spell
				if !result.Landed() || spell.ActionID.Tag == 71086 {
					return
				}

				// there are many apply effects with 0 damage right now for spells that actually have no impact damage
				// we do not want replicate those
				if result.PreOutcomeDamage == 0 {
					return
				}

				if !spell.ProcMask.Matches(core.ProcMaskSpellOrProc) {
					return
				}

				config := classConfig[a.GetCharacter().Spec]
				if config == nil {
					// TODO: HANDLE BETTER
					panic("DTR not supported for this spec yet")
				}

				if val, ok := config.spellConfig[spell.SpellID]; ok && val&supressImpact > 0 {
					return
				}

				// make sure the same spell impact can only trigger once per timestamp (AoE Impact spells like Arcane Explosion or Mind Sear)
				if lastTimestamp != sim.CurrentTime {
					lastTimestamp = sim.CurrentTime
					spellList = map[int32]bool{}
				}

				if _, ok := spellList[spell.SpellID]; ok {
					return
				}

				// spell has not been checked yet, add it
				spellList[spell.SpellID] = true

				if !sim.Proc(config.procChance, "Dragonwrath, Tarecgosa's Rest - DoT Proc") {
					return
				}

				config.getImpactHandler(spell.SpellID)(sim, spell, result)
			},

			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				config := classConfig[a.GetCharacter().Spec]
				if config == nil {
					// TODO: HANDLE BETTER
					panic("DTR not supported for this spec yet")
				}

				if val, ok := config.spellConfig[spell.SpellID]; ok && val&supressDoT > 0 {
					return
				}

				if !sim.Proc(config.procChance, "Dragonwrath, Tarecgosa's Rest - DoT Proc") {
					return
				}

				config.getDoTHandler(spell.SpellID)(sim, spell, result)
			},
		}))
	})

	// register custom global spell handlers
	registerPulseLightningCapacitor()
}

func registerSpells(unit *core.Unit) {
	registerDotSpell(unit)
}

func registerDotSpell(unit *core.Unit) {
	unit.RegisterSpell(core.SpellConfig{
		ActionID:                 core.ActionID{SpellID: 101085},
		SpellSchool:              core.SpellSchoolArcane,
		CritMultiplier:           0,
		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		Flags:                    core.SpellFlagNoSpellMods | core.SpellFlagIgnoreModifiers,
		ProcMask:                 core.ProcMaskEmpty,
		ManaCost:                 core.ManaCostOptions{},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, spell.BonusSpellPower, spell.OutcomeAlwaysHit)
		},
	})
}

func registerPulseLightningCapacitor() {
	globalSpellHandler[96891] = func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		copySpell := spell.Unit.GetSpell(spell.WithTag(71086))
		if copySpell == nil {
			copyConfig := GetDRTSpellConfig(spell)
			copyConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, spell.BonusSpellPower, spell.OutcomeMagicHitAndCrit)
			}
			copySpell = spell.Unit.RegisterSpell(copyConfig)
		}

		sim.AddPendingAction(
			&core.PendingAction{
				NextActionAt: sim.CurrentTime + time.Millisecond*200, // add slight delay
				Priority:     core.ActionPriorityAuto,
				OnAction: func(sim *core.Simulation) {
					// for now use the same roll as the old one as we don't carry any
					// meta data of the auras
					// only use BonusDamage fields for spells with 0 spell scaling
					copySpell.BonusSpellPower = result.PreOutcomeDamage
					copySpell.Cast(sim, result.Target)
				},
			},
		)
	}
}
