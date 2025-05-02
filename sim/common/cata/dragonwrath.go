package cata

import (
	"log"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

/*
Base classes for Dragonwrath, Tarecgosa's Rest support
All relevant mechanics should be tracked and documented at:
https://github.com/wowsims/mop/issues/946
*/
const (
	supressNone   int8 = 0
	supressImpact int8 = 1 << iota
	supressDoT
	supressAll = supressImpact | supressDoT
)

type SpellHandler func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult)
type SpellCopyHandler func(unit *core.Unit, spell *core.Spell)

type DragonwrathSpellConfig struct {
	spellHandler     SpellHandler
	spellCopyHandler SpellCopyHandler
	procPerCast      bool
	tickIsCast       bool // some spells deal periodic damage but should be treated as casts
	castIsTick       bool // some spells are implemented as casts but should be treated as periodic (Treat as Periodic flag)
	isAoESpell       bool // AoE Spells use reduced proc chance
	supress          int8
}

type DragonwrathClassConfig struct {
	procChance  float64
	spellConfig map[int32]DragonwrathSpellConfig
}

// Should be used in all places where copy spells are delayed to have it consistent
var DTRDelay = time.Millisecond * 0

var classConfig = map[proto.Spec]*DragonwrathClassConfig{}
var globalSpellConfig = map[int32]DragonwrathSpellConfig{}

func (config *DragonwrathClassConfig) getDoTHandler(spellId int32) SpellHandler {
	if val, ok := config.spellConfig[spellId]; ok && val.spellHandler != nil {
		return val.spellHandler
	}

	if val, ok := globalSpellConfig[spellId]; ok && val.spellHandler != nil {
		return val.spellHandler
	}

	return defaultDoTHandler
}

func (config *DragonwrathClassConfig) getImpactHandler(spellId int32) SpellHandler {
	if val, ok := config.spellConfig[spellId]; ok && val.spellHandler != nil {
		return val.spellHandler
	}

	if val, ok := globalSpellConfig[spellId]; ok && val.spellHandler != nil {
		return val.spellHandler
	}

	return defaultSpellHandler
}

func (config DragonwrathSpellConfig) TreatCastAsTick() DragonwrathSpellConfig {
	config.castIsTick = true
	config.tickIsCast = false
	return config
}

func (config DragonwrathSpellConfig) SupressImpact() DragonwrathSpellConfig {
	config.supress |= supressImpact
	return config
}

func (config DragonwrathSpellConfig) SupressDoT() DragonwrathSpellConfig {
	config.supress |= supressDoT
	return config
}

func (config DragonwrathSpellConfig) SupressSpell() DragonwrathSpellConfig {
	config.supress = supressAll
	return config
}

func (config DragonwrathSpellConfig) WithSpellHandler(handler SpellHandler) DragonwrathSpellConfig {
	config.spellHandler = handler
	return config
}

func (config DragonwrathSpellConfig) WithCustomSpell(handler SpellCopyHandler) DragonwrathSpellConfig {
	config.spellCopyHandler = handler
	return config
}

func (config DragonwrathSpellConfig) IsAoESpell() DragonwrathSpellConfig {
	config.isAoESpell = true
	return config
}

func (config DragonwrathSpellConfig) TreatTickAsCast() DragonwrathSpellConfig {
	config.tickIsCast = true
	config.castIsTick = false
	return config
}

func (config DragonwrathSpellConfig) ProcPerCast() DragonwrathSpellConfig {
	config.procPerCast = true
	return config
}

func NewDragonwrathSpellConfig() DragonwrathSpellConfig {
	return DragonwrathSpellConfig{}
}

func GetDragonwrathDoTSpell(unit *core.Unit) *core.Spell {
	return unit.GetSpell(core.ActionID{SpellID: 101085})
}

func CreateDTRClassConfig(spec proto.Spec, procChance float64) *DragonwrathClassConfig {
	if classConfig[spec] != nil {
		panic("Class config already registered")
	}

	classConfig[spec] = &DragonwrathClassConfig{
		procChance:  procChance,
		spellConfig: map[int32]DragonwrathSpellConfig{},
	}

	return classConfig[spec]
}

func (config *DragonwrathClassConfig) AddSpell(id int32, spellConfig DragonwrathSpellConfig) *DragonwrathClassConfig {
	config.spellConfig[id] = spellConfig
	return config
}

func defaultDoTHandler(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	dotSpell := GetDragonwrathDoTSpell(spell.Unit)
	dotSpell.BonusSpellPower = result.Damage
	dotSpell.Cast(sim, result.Target)
}

func defaultSpellHandler(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	// Add DTR item ID as tag for all duplicated spells
	copySpell := spell.Unit.GetSpell(spell.WithTag(spell.Tag + 71086))
	if copySpell == nil {
		copySpell = spell.Unit.RegisterSpell(GetDRTSpellConfig(spell))

		// copy BonusCoefficient as spells might have spell mods applied to them
		// we dont know the original in the base spell config so they might have been applied twice
		copySpell.BonusCoefficient = spell.BonusCoefficient
	}

	CastDTRSpell(sim, copySpell, result.Target)
}

func GetDRTSpellConfig(spell *core.Spell) core.SpellConfig {
	// create a copy with 0 tag as we don't use spell split metrics for dtr copies
	oldTag := spell.ActionID.Tag
	if spell.GetMetricSplitCount() > 1 {
		spell.ActionID.Tag = 0
	}
	baseConfig := core.SpellConfig{
		ActionID:                 spell.WithTag(spell.Tag + 71086),
		SpellSchool:              spell.SpellSchool,
		ProcMask:                 core.ProcMaskSpellProc,
		ApplyEffects:             spell.ApplyEffects,
		ManaCost:                 core.ManaCostOptions{},
		CritMultiplier:           spell.Unit.Env.Raid.GetPlayerFromUnit(spell.Unit).GetCharacter().DefaultCritMultiplier(),
		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		MissileSpeed:             spell.MissileSpeed,
		ClassSpellMask:           spell.ClassSpellMask,
		BonusCoefficient:         spell.BonusCoefficient,
		Flags:                    spell.Flags &^ core.SpellFlagAPL,
		RelatedDotSpell:          spell.RelatedDotSpell,
	}

	// instant spells will not refresh dots if they're duplicated
	if baseConfig.MissileSpeed == 0 && baseConfig.Cast.DefaultCast.CastTime == 0 {
		baseConfig.Flags |= core.SpellFlagSupressDoTApply
	}

	spell.ActionID.Tag = oldTag
	return baseConfig
}

func init() {
	core.NewItemEffect(71086, func(a core.Agent) {
		character := a.GetCharacter()
		unit := &character.Unit
		registerSpells(unit)
		var config *DragonwrathClassConfig

		unit.OnSpellRegistered(func(spell *core.Spell) {
			if val, ok := classConfig[character.Spec]; ok {
				for id, c := range val.spellConfig {
					if c.spellCopyHandler != nil && id == spell.SpellID && spell.ActionID.Tag < 71086 {
						c.spellCopyHandler(unit, spell)
					}
				}
			}
		})

		lastTimestamp := time.Duration(0)
		spellList := map[*core.Spell]bool{}

		aura := core.MakePermanent(unit.RegisterAura(core.Aura{
			ActionID: core.ActionID{ItemID: 71086},
			Label:    "Dragonwrath, Tarecgosa's Rest - Handler",
			OnInit: func(aura *core.Aura, sim *core.Simulation) {
				config = classConfig[character.Spec]
				if config == nil {

					// Create an empty config for this spell
					config = CreateDTRClassConfig(character.Spec, 0.0)
					classConfig[character.Spec] = config

					log.Printf("Using DTR for spec %s which is not implemented. Using default config", character.Spec)
				}
			},
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				lastTimestamp = time.Duration(0)
				spellList = map[*core.Spell]bool{}
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				lastTimestamp = sim.CurrentTime
				spellList = map[*core.Spell]bool{}
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {

				// Handle direct damage only and make sure we're not proccing of our own spell
				if !result.Landed() || spell.ActionID.Tag >= 71086 {
					return
				}

				// there are many apply effects with 0 damage right now for spells that actually have no impact damage
				// we do not want replicate those
				if result.PreOutcomeDamage == 0 {
					return
				}

				// only proc damage spells or class spells
				if !spell.ProcMask.Matches(core.ProcMaskSpellDamage) && spell.ClassSpellMask == 0 {
					return
				}

				// for now make it generic, might change this later
				// this rule should disable impact replication of related dot spells
				if spell.ActionID.Tag > 0 && spell.CurDot() != nil {
					return
				}

				if lastTimestamp != sim.CurrentTime {
					lastTimestamp = sim.CurrentTime
					spellList = map[*core.Spell]bool{}
				}

				procChance := config.procChance
				isAoESpell := false
				if val, ok := config.spellConfig[spell.SpellID]; ok {
					if val.supress&supressImpact > 0 {
						return
					}

					if val.castIsTick {
						aura.OnPeriodicDamageDealt(aura, sim, spell, result)
						return
					}

					if _, ok := spellList[spell]; ok {
						return
					}

					// make sure the same spell impact can only trigger once per timestamp (AoE Impact spells like Arcane Explosion or Mind Sear)
					if val.procPerCast {
						if lastTimestamp != sim.CurrentTime {
							lastTimestamp = sim.CurrentTime
							spellList = map[*core.Spell]bool{}
						}

						if _, ok := spellList[spell]; ok {
							return
						}

						// spell has not been checked yet, add it
						spellList[spell] = true
					}

					// reduce proc chance for AoE Spells
					if val.isAoESpell {
						isAoESpell = true
						procChance *= 2.0 / 9.0
					}
				}

				if !sim.Proc(procChance, "Dragonwrath, Tarecgosa's Rest - DoT Proc") {
					return
				}

				// AoE spells can only proc once per round
				if isAoESpell {
					// spell has not been checked yet, add it
					spellList[spell] = true
				}

				config.getImpactHandler(spell.SpellID)(sim, spell, result)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if val, ok := config.spellConfig[spell.SpellID]; ok {
					if val.isAoESpell {
						spellList = map[*core.Spell]bool{}
						return
					}
				}
				if _, ok := spellList[spell]; ok {
					spellList[spell] = false
				}
			},
			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if val, ok := config.spellConfig[spell.SpellID]; ok {
					if val.supress&supressDoT > 0 {
						return
					}

					if val.tickIsCast {
						aura.OnSpellHitDealt(aura, sim, spell, result)
						return
					}
				}

				if !sim.Proc(config.procChance, "Dragonwrath, Tarecgosa's Rest - DoT Proc") {
					return
				}

				config.getDoTHandler(spell.SpellID)(sim, spell, result)
			},
		}))

		character.ItemSwap.RegisterProc(71086, aura)
	})

	// register custom global spell handlers
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
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, spell.BonusSpellPower, spell.OutcomeAlwaysHit)
		},
	})
}

func CastDTRSpell(sim *core.Simulation, spell *core.Spell, target *core.Unit) {
	// instant spells will not refresh dots if they're duplicated
	if spell.Flags&core.SpellFlagSupressDoTApply > 0 && spell.RelatedDotSpell != nil {
		oldFlags := spell.RelatedDotSpell.Flags
		spell.RelatedDotSpell.Flags |= core.SpellFlagSupressDoTApply
		spell.Cast(sim, target)
		spell.RelatedDotSpell.Flags = oldFlags
		return
	}

	sim.AddPendingAction(&core.PendingAction{
		NextActionAt: sim.CurrentTime,
		Priority:     core.ActionPriorityAuto,
		OnAction: func(sim *core.Simulation) {
			spell.Cast(sim, target)
		},
	})
}
