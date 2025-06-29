package shared

import (
	"fmt"
	"math"
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

	// Soft fail to allow for overrides for bad effects
	if isEnchant {
		if core.HasEnchantEffect(config.EnchantID) {
			return
		}
	} else {
		if core.HasItemEffect(config.ItemID) {
			return
		}
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
		if proc.GetRppm() != nil {
			dpm = character.NewRPPMProcManager(effectID, isEnchant, config.ProcMask, core.RppmConfigFromProcEffectProto(proc))
		} else if proc.GetPpm() > 0 {
			if config.ProcMask == core.ProcMaskUnknown {
				if isEnchant {
					dpm = character.NewDynamicLegacyProcForEnchant(effectID, proc.GetPpm(), 0)
				} else {
					dpm = character.NewDynamicLegacyProcForWeapon(effectID, proc.GetPpm(), 0)
				}
			} else {
				dpm = character.NewLegacyPPMManager(proc.GetPpm(), config.ProcMask)
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

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID:   triggerActionID,
			Name:       config.Name,
			Callback:   config.Callback,
			ProcMask:   config.ProcMask,
			Outcome:    config.Outcome,
			Harmful:    config.Harmful,
			ProcChance: proc.GetProcChance(),
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
	var maxItemID int32

	for _, variant := range variants {
		maxItemID = max(maxItemID, variant.ItemID)
	}

	for _, variant := range variants {
		config.Name = variant.ItemName
		config.ItemID = variant.ItemID
		core.AddEffectsToTest = (config.ItemID == maxItemID)
		NewProcStatBonusEffect(config)
	}

	core.AddEffectsToTest = true
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
	Rppm        core.RPPMConfig

	// The stacks will only be granted as long as the trinket is active
	TrinketLimitsDuration bool
}

// Creates a new stacking stats bonus aura based on the configuration. If Bonus is not given, the ItemEffect of the item will be used
// to determine the correct values.
func NewStackingStatBonusCD(config StackingStatBonusCD) {
	core.NewItemEffect(config.ID, func(agent core.Agent, state proto.ItemLevelState) {
		character := agent.GetCharacter()

		auraID := core.ActionID{SpellID: config.AuraID}
		if auraID.IsEmptyAction() {
			auraID = core.ActionID{ItemID: config.ID}
		}

		// If we do not get a manual stat, overwrite it with scaling stats
		if config.Bonus.Equals(stats.Stats{}) {
			item := core.GetItemByID(config.ID)
			if item == nil || item.ItemEffect == nil {
				panic("Unsupported Item-/Effect")
			}

			config.Bonus = stats.FromProtoMap(item.ItemEffect.ScalingOptions[int32(state)].Stats)
		}

		var dpm *core.DynamicProcManager
		if config.Rppm.PPM > 0 {
			dpm = character.NewRPPMProcManager(config.ID, false, config.ProcMask, config.Rppm)
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
			DPM:        dpm,
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
	Rppm       core.RPPMConfig
	SpellFlags core.SpellFlag
	Outcome    core.HitOutcome
	Harmful    bool
	ProcChance float64
	Icd        time.Duration
}

func NewStackingStatBonusEffect(config StackingStatBonusEffect) {
	core.NewItemEffect(config.ItemID, func(agent core.Agent, state proto.ItemLevelState) {
		character := agent.GetCharacter()

		eligibleSlotsForItem := character.ItemSwap.EligibleSlotsForItem(config.ItemID)

		auraID := core.ActionID{SpellID: config.AuraID}
		if auraID.IsEmptyAction() {
			auraID = core.ActionID{ItemID: config.ItemID}
		}

		// If we do not get a manual stat, overwrite it with scaling stats
		if config.Bonus.Equals(stats.Stats{}) {
			item := core.GetItemByID(config.ItemID)
			if item == nil || item.ItemEffect == nil {
				panic("Unsupported Item-/Effect")
			}

			config.Bonus = stats.FromProtoMap(item.ItemEffect.ScalingOptions[int32(state)].Stats)
		}

		var dpm *core.DynamicProcManager
		if config.Rppm.PPM > 0 {
			dpm = character.NewRPPMProcManager(config.ItemID, false, config.ProcMask, config.Rppm)
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
			DPM:        dpm,
			ICD:        config.Icd,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				procAura.Activate(sim)
				procAura.AddStack(sim)
			},
		})

		procAura.Icd = triggerAura.Icd
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
	ItemID     int32
	SpellID    int32
	EnchantID  int32
	Trigger    core.ProcTrigger
	TriggerDPM func(*core.Character) *core.DynamicProcManager
	School     core.SpellSchool
	MinDmg     float64
	MaxDmg     float64
	IsMelee    bool
	Flags      core.SpellFlag
	Outcome    OutcomeType
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

		if config.TriggerDPM != nil {
			config.Trigger.DPM = config.TriggerDPM(character)
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

// Takes in the SpellResult for the triggering spell, and returns the total damage
// of a *fresh* Ignite triggered by that spell. Roll-over damage
// calculations for existing Ignites are handled internally.
type IgniteDamageCalculator func(result *core.SpellResult) float64

type IgniteConfig struct {
	ActionID           core.ActionID
	ClassSpellMask     int64
	SpellSchool        core.SpellSchool
	DisableCastMetrics bool
	DotAuraLabel       string
	DotAuraTag         string
	ProcTrigger        core.ProcTrigger // Ignores the Handler field and creates a custom one, but uses all others.
	DamageCalculator   IgniteDamageCalculator
	IncludeAuraDelay   bool // "munching" and "free roll-over" interactions
	NumberOfTicks      int32
	TickLength         time.Duration
	ParentAura         *core.Aura
}

func RegisterIgniteEffect(unit *core.Unit, config IgniteConfig) *core.Spell {
	spellFlags := core.SpellFlagIgnoreModifiers | core.SpellFlagNoSpellMods | core.SpellFlagNoOnCastComplete

	if config.DisableCastMetrics {
		spellFlags |= core.SpellFlagPassiveSpell
	}

	if config.SpellSchool == 0 {
		config.SpellSchool = core.SpellSchoolFire
	}

	if config.NumberOfTicks == 0 {
		config.NumberOfTicks = 2
	}

	if config.TickLength == 0 {
		config.TickLength = time.Second * 2
	}

	igniteSpell := unit.RegisterSpell(core.SpellConfig{
		ActionID:         config.ActionID,
		SpellSchool:      config.SpellSchool,
		ProcMask:         core.ProcMaskSpellProc,
		ClassSpellMask:   config.ClassSpellMask,
		Flags:            spellFlags,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     config.DotAuraLabel,
				Tag:       config.DotAuraTag,
				MaxStacks: math.MaxInt32,
			},

			NumberOfTicks:       config.NumberOfTicks,
			TickLength:          config.TickLength,
			AffectedByCastSpeed: false,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.Spell.CalcAndDealPeriodicDamage(sim, target, dot.SnapshotBaseDamage, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})

	refreshIgnite := func(sim *core.Simulation, target *core.Unit, damagePerTick float64) {
		// Cata Ignite
		// 1st ignite application = 4s, split into 2 ticks (2s, 0s)
		// Ignite refreshes: Duration = 4s + MODULO(remaining duration, 2), max 6s. Split damage over 3 ticks at 4s, 2s, 0s.
		dot := igniteSpell.Dot(target)
		dot.SnapshotBaseDamage = damagePerTick
		igniteSpell.Cast(sim, target)
		dot.Aura.SetStacks(sim, int32(dot.SnapshotBaseDamage))
	}

	var scheduledRefresh *core.PendingAction
	procTrigger := config.ProcTrigger
	procTrigger.Handler = func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
		target := result.Target
		dot := igniteSpell.Dot(target)
		outstandingDamage := dot.OutstandingDmg()
		newDamage := config.DamageCalculator(result)
		totalDamage := outstandingDamage + newDamage
		newTickCount := dot.BaseTickCount + core.TernaryInt32(dot.IsActive(), 1, 0)
		damagePerTick := totalDamage / float64(newTickCount)

		if config.IncludeAuraDelay {
			// Rough 2-bucket model for the aura update delay distribution based
			// on PTR measurements. Most updates occur on either the same or very
			// next spell batch after the proc, and can therefore be modeled by a
			// 0-10 ms random draw. But a reasonable minority fraction take ~10x
			// longer than this to fire. The origin of these longer delays is
			// likely not actually random in reality, but can be treated that way
			// in practice since the player cannot play around them.
			var delaySeconds float64

			if sim.Proc(0.75, "Aura Delay") {
				delaySeconds = 0.010 * sim.RandomFloat("Aura Delay")
			} else {
				delaySeconds = 0.090 + 0.020*sim.RandomFloat("Aura Delay")
			}

			applyDotAt := sim.CurrentTime + core.DurationFromSeconds(delaySeconds)

			// Cancel any prior aura updates already in the queue
			if (scheduledRefresh != nil) && (scheduledRefresh.NextActionAt > sim.CurrentTime) {
				scheduledRefresh.Cancel(sim)

				if sim.Log != nil {
					unit.Log(sim, "Previous %s proc was munched due to server aura delay", config.DotAuraLabel)
				}
			}

			// Schedule a delayed refresh of the DoT with cached damagePerTick value (allowing for "free roll-overs")
			if sim.Log != nil {
				unit.Log(sim, "Schedule travel (%0.1f ms) for %s", delaySeconds*1000, config.DotAuraLabel)

				if dot.IsActive() && (dot.NextTickAt() < applyDotAt) {
					unit.Log(sim, "%s rolled with %0.3f damage both ticking and rolled into next", config.DotAuraLabel, outstandingDamage)
				}
			}

			scheduledRefresh = core.StartDelayedAction(sim, core.DelayedActionOptions{
				DoAt:     applyDotAt,
				Priority: core.ActionPriorityDOT,

				OnAction: func(_ *core.Simulation) {
					refreshIgnite(sim, target, damagePerTick)
				},
			})
		} else {
			refreshIgnite(sim, target, damagePerTick)
		}
	}

	if config.ParentAura != nil {
		config.ParentAura.AttachProcTrigger(procTrigger)
	} else {
		core.MakeProcTriggerAura(unit, procTrigger)
	}

	return igniteSpell
}

type ItemVersion byte

const (
	ItemVersionLFR = iota
	ItemVersionNormal
	ItemVersionHeroic
	ItemVersionThunderforged
	ItemVersionHeroicThunderforged
	ItemVersionWarforged
	ItemVersionHeroicWarforged
	ItemVersionFlexible
)

type ItemVersionMap map[ItemVersion]int32
type ItemVersionFactory func(version ItemVersion, id int32, versionLabel string)

func (version ItemVersion) GetLabel() string {
	switch version {
	case ItemVersionLFR:
		return "(Celestial)"
	case ItemVersionHeroic:
		return "(Heroic)"
	case ItemVersionThunderforged:
		return "(Thunderforged)"
	case ItemVersionHeroicThunderforged:
		return "(Heroic Thunderforged)"
	case ItemVersionWarforged:
		return "(Warforged)"
	case ItemVersionHeroicWarforged:
		return "(Heroic Warforged)"
	case ItemVersionFlexible:
		return "(Flex)"
	}
	return ""
}

func (versions ItemVersionMap) RegisterAll(fac ItemVersionFactory) {
	var maxItemID int32

	for _, id := range versions {
		maxItemID = max(maxItemID, id)
	}

	for version, id := range versions {
		core.AddEffectsToTest = (id == maxItemID)
		fac(version, id, version.GetLabel())
	}

	core.AddEffectsToTest = true
}

func RegisterRiposteEffect(character *core.Character, auraSpellID int32, triggerSpellID int32) {
	riposteAura := character.RegisterAura(core.Aura{
		Label:     "Riposte" + character.Label,
		ActionID:  core.ActionID{SpellID: auraSpellID},
		Duration:  time.Second * 20,
		MaxStacks: math.MaxInt32,

		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			character.AddStatDynamic(sim, stats.CritRating, float64(newStacks-oldStacks))
		},
	})

	core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
		Name:     "Riposte Trigger" + character.Label,
		ActionID: core.ActionID{SpellID: triggerSpellID},
		Callback: core.CallbackOnSpellHitTaken,
		Outcome:  core.OutcomeDodge | core.OutcomeParry,
		ICD:      time.Second * 1,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			bonusCrit := math.Round((character.GetStat(stats.DodgeRating) + character.GetParryRatingWithoutStrength()) * 0.75)
			riposteAura.Activate(sim)
			riposteAura.SetStacks(sim, int32(bonusCrit))
		},
	})
}
