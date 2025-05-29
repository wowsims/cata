package core

import (
	"reflect"
	"time"

	"github.com/wowsims/mop/sim/core/proto"
)

type APLValue interface {
	// Returns all inner APLValues.
	GetInnerValues() []APLValue

	// The type of value that will be returned.
	Type() proto.APLValueType

	// Gets the value, assuming it is a particular type. Usually only one of
	// these should be implemented in each class.
	GetBool(*Simulation) bool
	GetInt(*Simulation) int32
	GetFloat(*Simulation) float64
	GetDuration(*Simulation) time.Duration
	GetString(*Simulation) string

	// Performs optional post-processing.
	Finalize(*APLRotation)

	// Pretty-print string for debugging.
	String() string
}

// Provides empty implementations for the GetX() value interface functions.
type DefaultAPLValueImpl struct {
	Uuid *proto.UUID
}

func (impl DefaultAPLValueImpl) GetInnerValues() []APLValue { return nil }
func (impl DefaultAPLValueImpl) Finalize(*APLRotation)      {}

func (impl DefaultAPLValueImpl) GetBool(sim *Simulation) bool {
	panic("Unimplemented GetBool")
}
func (impl DefaultAPLValueImpl) GetInt(sim *Simulation) int32 {
	panic("Unimplemented GetInt")
}
func (impl DefaultAPLValueImpl) GetFloat(sim *Simulation) float64 {
	panic("Unimplemented GetFloat")
}
func (impl DefaultAPLValueImpl) GetDuration(sim *Simulation) time.Duration {
	panic("Unimplemented GetDuration")
}
func (impl DefaultAPLValueImpl) GetString(sim *Simulation) string {
	panic("Unimplemented GetString")
}

func (rot *APLRotation) newAPLValue(config *proto.APLValue) APLValue {
	if config == nil {
		return nil
	}

	customValue := rot.unit.Env.GetAgentFromUnit(rot.unit).NewAPLValue(rot, config)
	if customValue != nil {
		return customValue
	}

	var value APLValue
	switch config.Value.(type) {
	// Operators
	case *proto.APLValue_Const:
		value = rot.newValueConst(config.GetConst(), config.Uuid)
	case *proto.APLValue_And:
		value = rot.newValueAnd(config.GetAnd(), config.Uuid)
	case *proto.APLValue_Or:
		value = rot.newValueOr(config.GetOr(), config.Uuid)
	case *proto.APLValue_Not:
		value = rot.newValueNot(config.GetNot(), config.Uuid)
	case *proto.APLValue_Cmp:
		value = rot.newValueCompare(config.GetCmp(), config.Uuid)
	case *proto.APLValue_Math:
		value = rot.newValueMath(config.GetMath(), config.Uuid)
	case *proto.APLValue_Max:
		value = rot.newValueMax(config.GetMax(), config.Uuid)
	case *proto.APLValue_Min:
		value = rot.newValueMin(config.GetMin(), config.Uuid)

	// Encounter
	case *proto.APLValue_CurrentTime:
		value = rot.newValueCurrentTime(config.GetCurrentTime(), config.Uuid)
	case *proto.APLValue_CurrentTimePercent:
		value = rot.newValueCurrentTimePercent(config.GetCurrentTimePercent(), config.Uuid)
	case *proto.APLValue_RemainingTime:
		value = rot.newValueRemainingTime(config.GetRemainingTime(), config.Uuid)
	case *proto.APLValue_RemainingTimePercent:
		value = rot.newValueRemainingTimePercent(config.GetRemainingTimePercent(), config.Uuid)
	case *proto.APLValue_IsExecutePhase:
		value = rot.newValueIsExecutePhase(config.GetIsExecutePhase(), config.Uuid)
	case *proto.APLValue_NumberTargets:
		value = rot.newValueNumberTargets(config.GetNumberTargets(), config.Uuid)

	// Boss
	case *proto.APLValue_BossSpellIsCasting:
		value = rot.newValueBossSpellIsCasting(config.GetBossSpellIsCasting(), config.Uuid)
	case *proto.APLValue_BossSpellTimeToReady:
		value = rot.newValueBossSpellTimeToReady(config.GetBossSpellTimeToReady(), config.Uuid)

	// Resources
	case *proto.APLValue_CurrentHealth:
		value = rot.newValueCurrentHealth(config.GetCurrentHealth(), config.Uuid)
	case *proto.APLValue_CurrentHealthPercent:
		value = rot.newValueCurrentHealthPercent(config.GetCurrentHealthPercent(), config.Uuid)
	case *proto.APLValue_CurrentMana:
		value = rot.newValueCurrentMana(config.GetCurrentMana(), config.Uuid)
	case *proto.APLValue_CurrentManaPercent:
		value = rot.newValueCurrentManaPercent(config.GetCurrentManaPercent(), config.Uuid)
	case *proto.APLValue_CurrentRage:
		value = rot.newValueCurrentRage(config.GetCurrentRage(), config.Uuid)
	case *proto.APLValue_CurrentEnergy:
		value = rot.newValueCurrentEnergy(config.GetCurrentEnergy(), config.Uuid)
	case *proto.APLValue_CurrentFocus:
		value = rot.newValueCurrentFocus(config.GetCurrentFocus(), config.Uuid)
	case *proto.APLValue_CurrentComboPoints:
		value = rot.newValueCurrentComboPoints(config.GetCurrentComboPoints(), config.Uuid)
	case *proto.APLValue_CurrentRunicPower:
		value = rot.newValueCurrentRunicPower(config.GetCurrentRunicPower(), config.Uuid)
	case *proto.APLValue_MaxHealth:
		value = rot.newValueMaxHealth(config.GetMaxHealth(), config.Uuid)
	case *proto.APLValue_MaxComboPoints:
		value = rot.newValueMaxComboPoints(config.GetMaxComboPoints(), config.Uuid)
	case *proto.APLValue_MaxEnergy:
		value = rot.newValueMaxEnergy(config.GetMaxEnergy(), config.Uuid)
	case *proto.APLValue_MaxFocus:
		value = rot.newValueMaxFocus(config.GetMaxFocus(), config.Uuid)
	case *proto.APLValue_MaxRunicPower:
		value = rot.newValueMaxRunicPower(config.GetMaxRunicPower(), config.Uuid)
	case *proto.APLValue_EnergyRegenPerSecond:
		value = rot.newValueEnergyRegenPerSecond(config.GetEnergyRegenPerSecond(), config.Uuid)
	case *proto.APLValue_FocusRegenPerSecond:
		value = rot.newValueFocusRegenPerSecond(config.GetFocusRegenPerSecond(), config.Uuid)
	case *proto.APLValue_EnergyTimeToTarget:
		value = rot.newValueEnergyTimeToTarget(config.GetEnergyTimeToTarget(), config.Uuid)
	case *proto.APLValue_FocusTimeToTarget:
		value = rot.newValueFocusTimeToTarget(config.GetFocusTimeToTarget(), config.Uuid)
	case *proto.APLValue_CurrentGenericResource:
		value = rot.newValueCurrentGenericResource(config.GetCurrentGenericResource(), config.Uuid)

	// Resources Runes
	case *proto.APLValue_CurrentRuneCount:
		value = rot.newValueCurrentRuneCount(config.GetCurrentRuneCount(), config.Uuid)
	case *proto.APLValue_CurrentNonDeathRuneCount:
		value = rot.newValueCurrentNonDeathRuneCount(config.GetCurrentNonDeathRuneCount(), config.Uuid)
	case *proto.APLValue_CurrentRuneActive:
		value = rot.newValueCurrentRuneActive(config.GetCurrentRuneActive(), config.Uuid)
	case *proto.APLValue_CurrentRuneDeath:
		value = rot.newValueCurrentRuneDeath(config.GetCurrentRuneDeath(), config.Uuid)
	case *proto.APLValue_RuneCooldown:
		value = rot.newValueRuneCooldown(config.GetRuneCooldown(), config.Uuid)
	case *proto.APLValue_NextRuneCooldown:
		value = rot.newValueNextRuneCooldown(config.GetNextRuneCooldown(), config.Uuid)
	case *proto.APLValue_RuneSlotCooldown:
		value = rot.newValueRuneSlotCooldown(config.GetRuneSlotCooldown(), config.Uuid)

	//Unit
	case *proto.APLValue_UnitIsMoving:
		value = rot.newValueCharacterIsMoving(config.GetUnitIsMoving(), config.Uuid)

	// GCD
	case *proto.APLValue_GcdIsReady:
		value = rot.newValueGCDIsReady(config.GetGcdIsReady(), config.Uuid)
	case *proto.APLValue_GcdTimeToReady:
		value = rot.newValueGCDTimeToReady(config.GetGcdTimeToReady(), config.Uuid)

	// Auto attacks
	case *proto.APLValue_AutoTimeToNext:
		value = rot.newValueAutoTimeToNext(config.GetAutoTimeToNext(), config.Uuid)

	// Spells
	case *proto.APLValue_SpellIsKnown:
		value = rot.newValueSpellIsKnown(config.GetSpellIsKnown(), config.Uuid)
	case *proto.APLValue_SpellCanCast:
		value = rot.newValueSpellCanCast(config.GetSpellCanCast(), config.Uuid)
	case *proto.APLValue_SpellIsReady:
		value = rot.newValueSpellIsReady(config.GetSpellIsReady(), config.Uuid)
	case *proto.APLValue_SpellTimeToReady:
		value = rot.newValueSpellTimeToReady(config.GetSpellTimeToReady(), config.Uuid)
	case *proto.APLValue_SpellCastTime:
		value = rot.newValueSpellCastTime(config.GetSpellCastTime(), config.Uuid)
	case *proto.APLValue_SpellTravelTime:
		value = rot.newValueSpellTravelTime(config.GetSpellTravelTime(), config.Uuid)
	case *proto.APLValue_SpellCpm:
		value = rot.newValueSpellCPM(config.GetSpellCpm(), config.Uuid)
	case *proto.APLValue_SpellIsChanneling:
		value = rot.newValueSpellIsChanneling(config.GetSpellIsChanneling(), config.Uuid)
	case *proto.APLValue_SpellChanneledTicks:
		value = rot.newValueSpellChanneledTicks(config.GetSpellChanneledTicks(), config.Uuid)
	case *proto.APLValue_SpellCurrentCost:
		value = rot.newValueSpellCurrentCost(config.GetSpellCurrentCost(), config.Uuid)
	case *proto.APLValue_SpellNumCharges:
		value = rot.newValueSpellNumCharges(config.GetSpellNumCharges(), config.Uuid)
	case *proto.APLValue_SpellTimeToCharge:
		value = rot.newValueSpellTimeToCharge(config.GetSpellTimeToCharge(), config.Uuid)

	// Auras
	case *proto.APLValue_AuraIsKnown:
		value = rot.newValueAuraIsKnown(config.GetAuraIsKnown(), config.Uuid)
	case *proto.APLValue_AuraIsActive:
		value = rot.newValueAuraIsActive(config.GetAuraIsActive(), config.Uuid)
	case *proto.APLValue_AuraIsActiveWithReactionTime:
		value = rot.newValueAuraIsActiveWithReactionTime(config.GetAuraIsActiveWithReactionTime(), config.Uuid)
	case *proto.APLValue_AuraIsInactiveWithReactionTime:
		value = rot.newValueAuraIsInactiveWithReactionTime(config.GetAuraIsInactiveWithReactionTime(), config.Uuid)
	case *proto.APLValue_AuraRemainingTime:
		value = rot.newValueAuraRemainingTime(config.GetAuraRemainingTime(), config.Uuid)
	case *proto.APLValue_AuraNumStacks:
		value = rot.newValueAuraNumStacks(config.GetAuraNumStacks(), config.Uuid)
	case *proto.APLValue_AuraInternalCooldown:
		value = rot.newValueAuraInternalCooldown(config.GetAuraInternalCooldown(), config.Uuid)
	case *proto.APLValue_AuraIcdIsReadyWithReactionTime:
		value = rot.newValueAuraICDIsReadyWithReactionTime(config.GetAuraIcdIsReadyWithReactionTime(), config.Uuid)
	case *proto.APLValue_AuraShouldRefresh:
		value = rot.newValueAuraShouldRefresh(config.GetAuraShouldRefresh(), config.Uuid)

	// Aura sets
	case *proto.APLValue_AllTrinketStatProcsActive:
		value = rot.newValueAllItemStatProcsActive(config.GetAllTrinketStatProcsActive(), config.Uuid)
	case *proto.APLValue_AnyTrinketStatProcsActive:
		value = rot.newValueAnyTrinketStatProcsActive(config.GetAnyTrinketStatProcsActive(), config.Uuid)
	case *proto.APLValue_TrinketProcsMinRemainingTime:
		value = rot.newValueItemProcsMinRemainingTime(config.GetTrinketProcsMinRemainingTime(), config.Uuid)
	case *proto.APLValue_TrinketProcsMaxRemainingIcd:
		value = rot.newValueItemsProcsMaxRemainingICD(config.GetTrinketProcsMaxRemainingIcd(), config.Uuid)
	case *proto.APLValue_NumEquippedStatProcTrinkets:
		value = rot.newValueNumEquippedStatProcItems(config.GetNumEquippedStatProcTrinkets(), config.Uuid)
	case *proto.APLValue_NumStatBuffCooldowns:
		value = rot.newValueNumStatBuffCooldowns(config.GetNumStatBuffCooldowns(), config.Uuid)

	// Dots
	case *proto.APLValue_DotIsActive:
		value = rot.newValueDotIsActive(config.GetDotIsActive(), config.Uuid)
	case *proto.APLValue_DotRemainingTime:
		value = rot.newValueDotRemainingTime(config.GetDotRemainingTime(), config.Uuid)
	case *proto.APLValue_DotTickFrequency:
		value = rot.newValueDotTickFrequency(config.GetDotTickFrequency(), config.Uuid)

	// Sequences
	case *proto.APLValue_SequenceIsComplete:
		value = rot.newValueSequenceIsComplete(config.GetSequenceIsComplete(), config.Uuid)
	case *proto.APLValue_SequenceIsReady:
		value = rot.newValueSequenceIsReady(config.GetSequenceIsReady(), config.Uuid)
	case *proto.APLValue_SequenceTimeToReady:
		value = rot.newValueSequenceTimeToReady(config.GetSequenceTimeToReady(), config.Uuid)

	// Properties
	case *proto.APLValue_ChannelClipDelay:
		value = rot.newValueChannelClipDelay(config.GetChannelClipDelay(), config.Uuid)
	case *proto.APLValue_InputDelay:
		value = rot.newValueInputDelay(config.GetInputDelay(), config.Uuid)

	default:
		value = nil
	}

	if value != nil {
		// The APLValue type doesn't embed APLValueDefaultImpl,
		// but all of the concrete subtypes do (e.g. APLValueConst)
		// Using reflection, we can set the field on the concrete type without casting.
		reflect.ValueOf(value).Elem().FieldByName("Uuid").Set(reflect.ValueOf(config.Uuid))
	}

	return value
}

// Default implementation of Agent.NewAPLValue so each spec doesn't need this boilerplate.
func (unit *Unit) NewAPLValue(rot *APLRotation, config *proto.APLValue) APLValue {
	return nil
}
