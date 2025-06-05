package shared

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type ProcStatBonusEffect struct {
	Name      string
	ItemID    int32
	EnchantID int32
	Callback  core.AuraCallback
	ProcMask  core.ProcMask
	Outcome   core.HitOutcome
	Harmful   bool

	// Any other custom proc conditions not covered by the above fields.
	CustomProcCondition core.CustomStatBuffProcCondition
}

type DamageEffect struct {
	SpellID          int32
	School           core.SpellSchool
	MinDmg           float64
	MaxDmg           float64
	BonusCoefficient float64
	IsMelee          bool
	ProcMask         core.ProcMask
	Outcome          OutcomeType
}

type ExtraSpellInfo struct {
	Spell   *core.Spell
	Trigger func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult)
}

type ItemVariant struct {
	ItemID   int32
	ItemName string
}

type CustomProcHandler func(sim *core.Simulation, procAura *core.StatBuffAura)

func NewProcStatBonusEffectWithDamageProc(config ProcStatBonusEffect, damage DamageEffect) {
	procMask := core.ProcMaskEmpty
	if damage.ProcMask != core.ProcMaskUnknown {
		procMask = damage.ProcMask
	}

	factory_StatBonusEffect(config, func(agent core.Agent, _ proto.ItemLevelState) ExtraSpellInfo {
		character := agent.GetCharacter()

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:                 core.ActionID{SpellID: damage.SpellID},
			SpellSchool:              damage.School,
			ProcMask:                 procMask,
			Flags:                    core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,
			DamageMultiplier:         1,
			CritMultiplier:           character.DefaultCritMultiplier(),
			DamageMultiplierAdditive: 1,
			ThreatMultiplier:         1,
			BonusCoefficient:         damage.BonusCoefficient,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(damage.MinDmg, damage.MaxDmg), GetOutcome(spell, damage.Outcome))
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

func factory_StatBonusEffect(config ProcStatBonusEffect, extraSpell func(agent core.Agent, _ proto.ItemLevelState) ExtraSpellInfo) {
	isEnchant := config.EnchantID != 0

	// Ignore empty dummy implementations
	if config.Callback == core.CallbackEmpty {
		return
	}

	var effectFn func(id int32, effect core.ApplyEffect)
	var effectID int32
	var triggerActionID core.ActionID
	if isEnchant {
		effectID = config.EnchantID
		effectFn = core.NewEnchantEffect
		triggerActionID = core.ActionID{SpellID: effectID}
	} else {
		effectID = config.ItemID
		effectFn = core.NewItemEffect
		triggerActionID = core.ActionID{ItemID: effectID}
	}

	effectFn(effectID, func(agent core.Agent, itemLevelState proto.ItemLevelState) {

		// Soft fail to allow for overrides for bad effects
		if core.HasItemEffect(effectID) {
			return
		}

		character := agent.GetCharacter()

		var eligibleSlots []proto.ItemSlot
		var procEffect *proto.ItemEffect
		if isEnchant {
			eligibleSlots = character.ItemSwap.EligibleSlotsForEffect(effectID)
			ench := core.GetEnchantByEffectID(effectID)
			if ench.EnchantEffect.GetProc() != nil {
				procEffect = ench.EnchantEffect
			}
		} else {
			eligibleSlots = character.ItemSwap.EligibleSlotsForItem(effectID)

			item := core.GetItemByID(effectID)
			if item.ItemEffect != nil {
				if item.ItemEffect.GetProc() != nil {
					procEffect = item.ItemEffect
				}
			}
		}

		if procEffect == nil {
			err, _ := fmt.Printf("Error getting proc effect for item/enchant %v", effectID)
			panic(err)
		}

		proc := procEffect.GetProc()
		procAction := core.ActionID{SpellID: procEffect.BuffId}
		procAura := character.NewTemporaryStatsAura(
			config.Name+" Proc",
			procAction,
			stats.FromProtoMap(procEffect.ScalingOptions[int32(itemLevelState)].Stats),
			time.Millisecond*time.Duration(procEffect.EffectDurationMs),
		)

		var dpm *core.DynamicProcManager
		if (proc.Ppm != 0) && (config.ProcMask == core.ProcMaskUnknown) {
			if isEnchant {
				dpm = character.AutoAttacks.NewDynamicProcManagerForEnchant(effectID, proc.Ppm, 0)
			} else {
				dpm = character.AutoAttacks.NewDynamicProcManagerForWeaponEffect(effectID, proc.Ppm, 0)
			}
		}

		procAura.CustomProcCondition = config.CustomProcCondition
		var customHandler CustomProcHandler
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
			procSpell = extraSpell(agent, itemLevelState)
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

		if config.ProcMask.Matches(core.ProcMaskEmpty) &&
			(core.GetItemByID(config.ItemID).WeaponType > 0 || core.GetItemByID(config.ItemID).RangedWeaponType > 0) {
			config.ProcMask = *character.GetDynamicProcMaskForWeaponEffect(config.ItemID)
		}

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID:   triggerActionID,
			Name:       config.Name,
			Callback:   config.Callback,
			ProcMask:   config.ProcMask,
			Outcome:    config.Outcome,
			Harmful:    config.Harmful,
			ProcChance: proc.ProcChance,
			PPM:        proc.Ppm,
			DPM:        dpm,
			ICD:        time.Millisecond * time.Duration(proc.IcdMs),
			Handler:    handler,
		})
		if proc.IcdMs != 0 {
			procAura.Icd = triggerAura.Icd
		}
		if isEnchant {
			character.ItemSwap.RegisterEnchantProcWithSlots(effectID, triggerAura, eligibleSlots)
		} else {
			character.ItemSwap.RegisterProcWithSlots(effectID, triggerAura, eligibleSlots)
		}

		character.AddStatProcBuff(effectID, procAura, isEnchant, eligibleSlots)
	})
}

func NewProcStatBonusEffectWithVariants(config ProcStatBonusEffect, variants []ItemVariant) {
	for _, variant := range variants {
		config.Name = variant.ItemName
		config.ItemID = variant.ItemID
		NewProcStatBonusEffect(config)
	}
}

func NewProcStatBonusEffect(config ProcStatBonusEffect) {
	factory_StatBonusEffect(config, nil)
}

func NewSimpleStatActive(itemID int32) {

	// Soft fail to allow for overrides for bad effects
	if core.HasItemEffect(itemID) {
		return
	}

	core.NewItemEffect(itemID, func(agent core.Agent, scalingSelector proto.ItemLevelState) {
		item := core.GetItemByID(itemID)
		if item == nil {
			panic(fmt.Sprintf("No item with ID: %d", itemID))
		}

		itemEffect := item.ItemEffect // Assuming it can be collapsed to one relevant effect per item in pre-processing
		if itemEffect == nil {
			panic(fmt.Sprintf("No effect data for item with ID: %d", itemID))
		}

		onUseData := itemEffect.GetOnUse()
		if onUseData == nil {
			panic(fmt.Sprintf("Item effect for item with ID: %d is not an active effect!", itemID))
		}

		spellConfig := core.SpellConfig{
			ActionID: core.ActionID{ItemID: itemID},
		}

		character := agent.GetCharacter()
		spellConfig.Cast.CD = core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Duration(onUseData.CooldownMs) * time.Millisecond,
		}
		// if SpellCategoryID is 0 we seemingly do not share cd with anything
		// Say Darkmoon Card: Earthquake and Ruthless Gladiator's Emblem of Cruelty even though tooltip shows as such
		if onUseData.CategoryId > 0 {
			sharedCDDuration := time.Duration(onUseData.CategoryCooldownMs) * time.Millisecond
			if sharedCDDuration == 0 {
				sharedCDDuration = time.Millisecond * time.Duration(itemEffect.EffectDurationMs)
			}

			sharedCDTimer := character.GetOrInitSpellCategoryTimer(onUseData.CategoryId)
			spellConfig.Cast.SharedCD = core.Cooldown{
				Timer:    sharedCDTimer,
				Duration: sharedCDDuration,
			}
		}

		core.RegisterTemporaryStatsOnUseCD(character, itemEffect.BuffName, stats.FromProtoMap(itemEffect.ScalingOptions[int32(scalingSelector)].Stats), time.Millisecond*time.Duration(itemEffect.EffectDurationMs), spellConfig)
	})
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
	core.NewItemEffect(config.ID, func(agent core.Agent, _ proto.ItemLevelState) {
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
	ItemID     int32
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
	core.NewItemEffect(config.ItemID, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()

		eligibleSlotsForItem := character.ItemSwap.EligibleSlotsForItem(config.ItemID)

		auraID := core.ActionID{SpellID: config.AuraID}
		if auraID.IsEmptyAction() {
			auraID = core.ActionID{ItemID: config.ItemID}
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

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID:   core.ActionID{ItemID: config.ItemID},
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

		character.AddStatProcBuff(config.ItemID, procAura, false, eligibleSlotsForItem)
		character.ItemSwap.RegisterProcWithSlots(config.ItemID, triggerAura, eligibleSlotsForItem)
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
	OutcomeRangedCanCrit
)

type ProcDamageEffect struct {
	ItemID    int32
	SpellID   int32
	EnchantID int32
	Trigger   core.ProcTrigger

	School  core.SpellSchool
	MinDmg  float64
	MaxDmg  float64
	IsMelee bool
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
	case OutcomeRangedCanCrit:
		return spell.OutcomeRangedHitAndCrit
	default:
		return spell.OutcomeMagicHitAndCrit
	}
}

func NewProcDamageEffect(config ProcDamageEffect) {
	isEnchant := config.EnchantID != 0

	var effectFn func(id int32, effect core.ApplyEffect)
	var effectID int32
	var triggerActionID core.ActionID

	if isEnchant {
		effectID = config.EnchantID
		effectFn = core.NewEnchantEffect
		triggerActionID = core.ActionID{SpellID: config.SpellID}
	} else {
		effectID = config.ItemID
		effectFn = core.NewItemEffect
		triggerActionID = core.ActionID{ItemID: config.ItemID}
	}

	effectFn(effectID, func(agent core.Agent, _ proto.ItemLevelState) {
		character := agent.GetCharacter()

		minDmg := config.MinDmg
		maxDmg := config.MaxDmg

		if core.ActionID.IsEmptyAction(config.Trigger.ActionID) {
			config.Trigger.ActionID = triggerActionID
		}

		damageSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: config.SpellID},
			SpellSchool: config.School,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       config.Flags,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(minDmg, maxDmg), GetOutcome(spell, config.Outcome))
			},
		})

		triggerConfig := config.Trigger
		triggerConfig.Handler = func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			damageSpell.Cast(sim, character.CurrentTarget)
		}
		triggerAura := core.MakeProcTriggerAura(&character.Unit, triggerConfig)

		if isEnchant {
			character.ItemSwap.RegisterEnchantProc(effectID, triggerAura)
		} else {
			character.ItemSwap.RegisterProc(effectID, triggerAura)
		}
	})
}
