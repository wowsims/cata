package shared

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

type ProcStatBonusEffect struct {
	Name       string
	ID         int32
	AuraID     int32
	Bonus      stats.Stats
	Duration   time.Duration
	Callback   core.AuraCallback
	ProcMask   core.ProcMask
	Outcome    core.HitOutcome
	Harmful    bool
	ProcChance float64
	PPM        float64
	ICD        time.Duration

	// For ignoring a hardcoded spell.
	IgnoreSpellID int32

	// Any other custom proc conditions not covered by the above fields.
	CustomProcCondition core.CustomStatBuffProcCondition
}

type DamageEffect struct {
	SpellID          int32
	School           core.SpellSchool
	MinDmg           float64
	MaxDmg           float64
	BonusCoefficient float64
	ProcMask         core.ProcMask
}

type ExtraSpellInfo struct {
	Spell   *core.Spell
	Trigger func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult)
}

type CustomProcHandler func(sim *core.Simulation, procAura *core.StatBuffAura)

func NewProcStatBonusEffectWithDamageProc(config ProcStatBonusEffect, damage DamageEffect) {
	procMask := core.ProcMaskEmpty
	if damage.ProcMask != core.ProcMaskUnknown {
		procMask = damage.ProcMask
	}

	factory_StatBonusEffect(config, func(agent core.Agent) ExtraSpellInfo {
		character := agent.GetCharacter()
		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:                 core.ActionID{SpellID: damage.SpellID},
			SpellSchool:              damage.School,
			ProcMask:                 procMask,
			Flags:                    core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,
			DamageMultiplier:         1,
			CritMultiplier:           character.DefaultSpellCritMultiplier(),
			DamageMultiplierAdditive: 1,
			ThreatMultiplier:         1,
			BonusCoefficient:         damage.BonusCoefficient,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(damage.MinDmg, damage.MaxDmg), spell.OutcomeMagicCrit)
			},
		})

		return ExtraSpellInfo{
			Spell: procSpell,
			Trigger: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, result.Target)
			},
		}
	})
}

func factory_StatBonusEffect(config ProcStatBonusEffect, extraSpell func(agent core.Agent) ExtraSpellInfo) {
	core.NewItemEffect(config.ID, func(agent core.Agent) {
		character := agent.GetCharacter()

		procID := core.ActionID{SpellID: config.AuraID}
		if procID.IsEmptyAction() {
			procID = core.ActionID{ItemID: config.ID}
		}
		procAura := character.NewTemporaryStatsAura(config.Name+" Proc", procID, config.Bonus, config.Duration)

		var itemSwapProcCondition core.CustomStatBuffProcCondition
		if character.ItemSwap.IsEnabled() && character.ItemSwap.ItemExistsInSwapSet(config.ID) {
			itemSwapProcCondition = func(sim *core.Simulation, aura *core.Aura) bool {
				return character.ItemSwap.HasItemEquipped(config.ID)
			}
		}

		var customHandler CustomProcHandler
		if config.CustomProcCondition != nil {
			if itemSwapProcCondition != nil {
				procAura.CustomProcCondition = func(sim *core.Simulation, aura *core.Aura) bool {
					return itemSwapProcCondition(sim, aura) && config.CustomProcCondition(sim, aura)
				}
			} else {
				procAura.CustomProcCondition = config.CustomProcCondition
			}

		} else if itemSwapProcCondition != nil {
			procAura.CustomProcCondition = itemSwapProcCondition
		}

		if config.CustomProcCondition != nil {
			customHandler = func(sim *core.Simulation, procAura *core.StatBuffAura) {
				if procAura.CanProc(sim) {
					procAura.Activate(sim)
				} else {
					// reset ICD condition was not fulfilled
					if procAura.Icd != nil && procAura.Icd.Duration != 0 {
						procAura.Icd.Reset()
					}
				}
			}
		}

		var procSpell ExtraSpellInfo
		if extraSpell != nil {
			procSpell = extraSpell(agent)
		}

		handler := func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if customHandler != nil {
				customHandler(sim, procAura)
			} else {
				procAura.Activate(sim)
				if procSpell.Spell != nil {
					procSpell.Trigger(sim, spell, result)
				}
			}
		}

		if config.IgnoreSpellID != 0 {
			ignoreSpellID := config.IgnoreSpellID
			handler = func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !spell.IsSpellAction(ignoreSpellID) {
					if customHandler != nil {
						customHandler(sim, procAura)
					} else {
						procAura.Activate(sim)
						if procSpell.Spell != nil {
							procSpell.Trigger(sim, spell, result)
						}
					}
				}
			}
		}

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID:   core.ActionID{ItemID: config.ID},
			Name:       config.Name,
			Callback:   config.Callback,
			ProcMask:   config.ProcMask,
			Outcome:    config.Outcome,
			Harmful:    config.Harmful,
			ProcChance: config.ProcChance,
			PPM:        config.PPM,
			ICD:        config.ICD,
			Handler:    handler,
		})

		procAura.Icd = triggerAura.Icd
		character.TrinketProcBuffs = append(character.TrinketProcBuffs, procAura)
		character.ItemSwap.RegisterOnSwapItemForItemProcEffect(config.ID, triggerAura, core.EligibleSlotsForItem(core.GetItemByID(config.ID), false))
	})
}

func NewProcStatBonusEffect(config ProcStatBonusEffect) {
	factory_StatBonusEffect(config, nil)
}

type StatCDFactory func(itemID int32, duration time.Duration, cooldown time.Duration)

// Wraps factory functions so that only the first item is included in tests.
func testFirstOnly(factory StatCDFactory) StatCDFactory {
	first := true
	return func(itemID int32, duration time.Duration, cooldown time.Duration) {
		if first {
			first = false
			factory(itemID, duration, cooldown)
		} else {
			core.AddEffectsToTest = false
			factory(itemID, duration, cooldown)
			core.AddEffectsToTest = true
		}
	}
}

func CreateOffensiveStatActive(itemID int32, duration time.Duration, cooldown time.Duration, stats stats.Stats) {
	testFirstOnly(func(itemID int32, duration time.Duration, cooldown time.Duration) {
		core.NewSimpleStatOffensiveTrinketEffect(itemID, stats, duration, cooldown)
	})(itemID, duration, cooldown)
}

func CreateDefensiveStatActive(itemID int32, duration time.Duration, cooldown time.Duration, stats stats.Stats) {
	testFirstOnly(func(itemID int32, duration time.Duration, cooldown time.Duration) {
		core.NewSimpleStatDefensiveTrinketEffect(itemID, stats, duration, cooldown)
	})(itemID, duration, cooldown)
}

func NewStrengthActive(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
	CreateOffensiveStatActive(itemID, duration, cooldown, stats.Stats{stats.Strength: bonus})
}

func NewAgilityActive(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
	CreateOffensiveStatActive(itemID, duration, cooldown, stats.Stats{stats.Agility: bonus})
}

func NewIntActive(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
	CreateOffensiveStatActive(itemID, duration, cooldown, stats.Stats{stats.Intellect: bonus})
}

func NewSpiritActive(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
	CreateOffensiveStatActive(itemID, duration, cooldown, stats.Stats{stats.Spirit: bonus})
}

func NewCritActive(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
	CreateOffensiveStatActive(itemID, duration, cooldown, stats.Stats{stats.CritRating: bonus})
}

func NewHasteActive(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
	CreateOffensiveStatActive(itemID, duration, cooldown, stats.Stats{stats.HasteRating: bonus})
}

func NewDodgeActive(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
	CreateDefensiveStatActive(itemID, duration, cooldown, stats.Stats{stats.DodgeRating: bonus})
}

func NewSpellPowerActive(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
	CreateOffensiveStatActive(itemID, duration, cooldown, stats.Stats{stats.SpellPower: bonus})
}

func NewHealthActive(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
	CreateDefensiveStatActive(itemID, duration, cooldown, stats.Stats{stats.Health: bonus})
}

func NewParryActive(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
	CreateDefensiveStatActive(itemID, duration, cooldown, stats.Stats{stats.ParryRating: bonus})
}

func NewMasteryActive(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
	CreateOffensiveStatActive(itemID, duration, cooldown, stats.Stats{stats.MasteryRating: bonus})
}

type StackingStatBonusCD struct {
	Name        string
	ID          int32
	AuraID      int32
	Bonus       stats.Stats
	Duration    time.Duration
	MaxStacks   int32
	CD          time.Duration
	Callback    core.AuraCallback
	ProcMask    core.ProcMask
	SpellFlags  core.SpellFlag
	Outcome     core.HitOutcome
	Harmful     bool
	ProcChance  float64
	IsDefensive bool

	// The stacks will only be granted as long as the trinket is active
	TrinketLimitsDuration bool
}

func NewStackingStatBonusCD(config StackingStatBonusCD) {
	core.NewItemEffect(config.ID, func(agent core.Agent) {
		character := agent.GetCharacter()

		auraID := core.ActionID{SpellID: config.AuraID}
		if auraID.IsEmptyAction() {
			auraID = core.ActionID{ItemID: config.ID}
		}

		duration := core.TernaryDuration(config.TrinketLimitsDuration, core.NeverExpires, config.Duration)
		statAura := core.MakeStackingAura(character, core.StackingStatAura{
			Aura: core.Aura{
				Label:     config.Name + " Proc",
				ActionID:  auraID,
				Duration:  duration,
				MaxStacks: config.MaxStacks,
			},
			BonusPerStack: config.Bonus,
		})

		// If trinket limits duration create a separate proc aura
		var procAura *core.Aura = statAura.Aura
		if config.TrinketLimitsDuration {
			procAura = character.RegisterAura(core.Aura{
				Label:    config.Name + " Aura",
				ActionID: auraID,
				Duration: config.Duration,
				OnExpire: func(_ *core.Aura, sim *core.Simulation) {
					statAura.Deactivate(sim)
				},
			})
		}

		core.ApplyProcTriggerCallback(&character.Unit, procAura, core.ProcTrigger{
			Name:       config.Name,
			Callback:   config.Callback,
			ProcMask:   config.ProcMask,
			SpellFlags: config.SpellFlags,
			Outcome:    config.Outcome,
			Harmful:    config.Harmful,
			ProcChance: config.ProcChance,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				statAura.AddStack(sim)
			},
		})

		var sharedTimer *core.Timer
		if config.IsDefensive {
			sharedTimer = character.GetDefensiveTrinketCD()
		} else {
			sharedTimer = character.GetOffensiveTrinketCD()
		}

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{ItemID: config.ID},
			Flags:    core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: config.CD,
				},
				SharedCD: core.Cooldown{
					Timer:    sharedTimer,
					Duration: config.Duration,
				},
			},

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
				statAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeDPS,
		})
	})
}

type StackingStatBonusEffect struct {
	Name       string
	ID         int32
	AuraID     int32
	Bonus      stats.Stats
	Duration   time.Duration
	MaxStacks  int32
	Callback   core.AuraCallback
	ProcMask   core.ProcMask
	SpellFlags core.SpellFlag
	Outcome    core.HitOutcome
	Harmful    bool
	ProcChance float64
	Icd        time.Duration
}

func NewStackingStatBonusEffect(config StackingStatBonusEffect) {
	core.NewItemEffect(config.ID, func(agent core.Agent) {
		character := agent.GetCharacter()

		auraID := core.ActionID{SpellID: config.AuraID}
		if auraID.IsEmptyAction() {
			auraID = core.ActionID{ItemID: config.ID}
		}
		procAura := core.MakeStackingAura(character, core.StackingStatAura{
			Aura: core.Aura{
				Label:     config.Name + " Proc",
				ActionID:  auraID,
				Duration:  config.Duration,
				MaxStacks: config.MaxStacks,
			},
			BonusPerStack: config.Bonus,
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID:   core.ActionID{ItemID: config.ID},
			Name:       config.Name,
			Callback:   config.Callback,
			ProcMask:   config.ProcMask,
			SpellFlags: config.SpellFlags,
			Outcome:    config.Outcome,
			Harmful:    config.Harmful,
			ProcChance: config.ProcChance,
			ICD:        config.Icd,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				procAura.Activate(sim)
				procAura.AddStack(sim)
			},
		})

		character.TrinketProcBuffs = append(character.TrinketProcBuffs, procAura)
	})
}

type OutcomeType uint64

const (
	OutcomeDefault                  = 0
	OutcomeMeleeCanCrit OutcomeType = iota
	OutcomeMeleeNoCrit
	OutcomeMeleeNoBlockDodgeParryCrit
	OutcomeSpellCanCrit
	OutcomeSpellNoCrit
	OutcomeSpellNoMissCanCrit
)

type ProcDamageEffect struct {
	ItemID  int32
	SpellID int32
	Trigger core.ProcTrigger

	School  core.SpellSchool
	MinDmg  float64
	MaxDmg  float64
	Flags   core.SpellFlag
	Outcome OutcomeType
}

func GetOutcome(spell *core.Spell, outcome OutcomeType) core.OutcomeApplier {
	switch outcome {
	case OutcomeMeleeCanCrit:
		return spell.OutcomeMeleeSpecialHitAndCrit
	case OutcomeMeleeNoCrit:
		return spell.OutcomeMeleeSpecialHit
	case OutcomeMeleeNoBlockDodgeParryCrit:
		return spell.OutcomeMeleeSpecialNoBlockDodgeParryNoCrit
	case OutcomeSpellCanCrit:
		return spell.OutcomeMagicHitAndCrit
	case OutcomeSpellNoMissCanCrit:
		return spell.OutcomeMagicCrit
	case OutcomeSpellNoCrit:
		return spell.OutcomeMagicHit
	default:
		return spell.OutcomeMagicHitAndCrit
	}
}

func NewProcDamageEffect(config ProcDamageEffect) {
	core.NewItemEffect(config.ItemID, func(agent core.Agent) {
		character := agent.GetCharacter()

		minDmg := config.MinDmg
		maxDmg := config.MaxDmg

		if core.ActionID.IsEmptyAction(config.Trigger.ActionID) {
			config.Trigger.ActionID = core.ActionID{ItemID: config.ItemID}
		}

		damageSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: config.SpellID},
			SpellSchool: config.School,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       config.Flags,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(minDmg, maxDmg), GetOutcome(spell, config.Outcome))
			},
		})

		triggerConfig := config.Trigger
		triggerConfig.Handler = func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			damageSpell.Cast(sim, character.CurrentTarget)
		}
		core.MakeProcTriggerAura(&character.Unit, triggerConfig)
	})
}
