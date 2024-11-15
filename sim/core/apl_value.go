package core

import (
	"reflect"
	"time"

	"github.com/wowsims/cata/sim/core/proto"
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
		value = rot.newValueConst(config.GetConst())
	case *proto.APLValue_And:
		value = rot.newValueAnd(config.GetAnd())
	case *proto.APLValue_Or:
		value = rot.newValueOr(config.GetOr())
	case *proto.APLValue_Not:
		value = rot.newValueNot(config.GetNot())
	case *proto.APLValue_Cmp:
		value = rot.newValueCompare(config.GetCmp())
	case *proto.APLValue_Math:
		value = rot.newValueMath(config.GetMath())
	case *proto.APLValue_Max:
		value = rot.newValueMax(config.GetMax())
	case *proto.APLValue_Min:
		value = rot.newValueMin(config.GetMin())

	// Encounter
	case *proto.APLValue_CurrentTime:
		value = rot.newValueCurrentTime(config.GetCurrentTime())
	case *proto.APLValue_CurrentTimePercent:
		value = rot.newValueCurrentTimePercent(config.GetCurrentTimePercent())
	case *proto.APLValue_RemainingTime:
		value = rot.newValueRemainingTime(config.GetRemainingTime())
	case *proto.APLValue_RemainingTimePercent:
		value = rot.newValueRemainingTimePercent(config.GetRemainingTimePercent())
	case *proto.APLValue_IsExecutePhase:
		value = rot.newValueIsExecutePhase(config.GetIsExecutePhase())
	case *proto.APLValue_NumberTargets:
		value = rot.newValueNumberTargets(config.GetNumberTargets())

	// Boss
	case *proto.APLValue_BossSpellIsCasting:
		value = rot.newValueBossSpellIsCasting(config.GetBossSpellIsCasting())
	case *proto.APLValue_BossSpellTimeToReady:
		value = rot.newValueBossSpellTimeToReady(config.GetBossSpellTimeToReady())

	// Resources
	case *proto.APLValue_CurrentHealth:
		value = rot.newValueCurrentHealth(config.GetCurrentHealth())
	case *proto.APLValue_CurrentHealthPercent:
		value = rot.newValueCurrentHealthPercent(config.GetCurrentHealthPercent())
	case *proto.APLValue_CurrentMana:
		value = rot.newValueCurrentMana(config.GetCurrentMana())
	case *proto.APLValue_CurrentManaPercent:
		value = rot.newValueCurrentManaPercent(config.GetCurrentManaPercent())
	case *proto.APLValue_CurrentRage:
		value = rot.newValueCurrentRage(config.GetCurrentRage())
	case *proto.APLValue_CurrentEnergy:
		value = rot.newValueCurrentEnergy(config.GetCurrentEnergy())
	case *proto.APLValue_CurrentFocus:
		value = rot.newValueCurrentFocus(config.GetCurrentFocus())
	case *proto.APLValue_CurrentComboPoints:
		value = rot.newValueCurrentComboPoints(config.GetCurrentComboPoints())
	case *proto.APLValue_CurrentRunicPower:
		value = rot.newValueCurrentRunicPower(config.GetCurrentRunicPower())

	// Resources Runes
	case *proto.APLValue_CurrentRuneCount:
		value = rot.newValueCurrentRuneCount(config.GetCurrentRuneCount())
	case *proto.APLValue_CurrentNonDeathRuneCount:
		value = rot.newValueCurrentNonDeathRuneCount(config.GetCurrentNonDeathRuneCount())
	case *proto.APLValue_CurrentRuneActive:
		value = rot.newValueCurrentRuneActive(config.GetCurrentRuneActive())
	case *proto.APLValue_CurrentRuneDeath:
		value = rot.newValueCurrentRuneDeath(config.GetCurrentRuneDeath())
	case *proto.APLValue_RuneCooldown:
		value = rot.newValueRuneCooldown(config.GetRuneCooldown())
	case *proto.APLValue_NextRuneCooldown:
		value = rot.newValueNextRuneCooldown(config.GetNextRuneCooldown())
	case *proto.APLValue_RuneSlotCooldown:
		value = rot.newValueRuneSlotCooldown(config.GetRuneSlotCooldown())

	//Unit
	case *proto.APLValue_UnitIsMoving:
		value = rot.newValueCharacterIsMoving(config.GetUnitIsMoving())

	// GCD
	case *proto.APLValue_GcdIsReady:
		value = rot.newValueGCDIsReady(config.GetGcdIsReady())
	case *proto.APLValue_GcdTimeToReady:
		value = rot.newValueGCDTimeToReady(config.GetGcdTimeToReady())

	// Auto attacks
	case *proto.APLValue_AutoTimeToNext:
		value = rot.newValueAutoTimeToNext(config.GetAutoTimeToNext())

	// Spells
	case *proto.APLValue_SpellIsKnown:
		value = rot.newValueSpellIsKnown(config.GetSpellIsKnown())
	case *proto.APLValue_SpellCanCast:
		value = rot.newValueSpellCanCast(config.GetSpellCanCast())
	case *proto.APLValue_SpellIsReady:
		value = rot.newValueSpellIsReady(config.GetSpellIsReady())
	case *proto.APLValue_SpellTimeToReady:
		value = rot.newValueSpellTimeToReady(config.GetSpellTimeToReady())
	case *proto.APLValue_SpellCastTime:
		value = rot.newValueSpellCastTime(config.GetSpellCastTime())
	case *proto.APLValue_SpellTravelTime:
		value = rot.newValueSpellTravelTime(config.GetSpellTravelTime())
	case *proto.APLValue_SpellCpm:
		value = rot.newValueSpellCPM(config.GetSpellCpm())
	case *proto.APLValue_SpellIsChanneling:
		value = rot.newValueSpellIsChanneling(config.GetSpellIsChanneling())
	case *proto.APLValue_SpellChanneledTicks:
		value = rot.newValueSpellChanneledTicks(config.GetSpellChanneledTicks())
	case *proto.APLValue_SpellCurrentCost:
		value = rot.newValueSpellCurrentCost(config.GetSpellCurrentCost())

	// Auras
	case *proto.APLValue_AuraIsKnown:
		value = rot.newValueAuraIsKnown(config.GetAuraIsKnown())
	case *proto.APLValue_AuraIsActive:
		value = rot.newValueAuraIsActive(config.GetAuraIsActive())
	case *proto.APLValue_AuraIsActiveWithReactionTime:
		value = rot.newValueAuraIsActiveWithReactionTime(config.GetAuraIsActiveWithReactionTime())
	case *proto.APLValue_AuraIsInactiveWithReactionTime:
		value = rot.newValueAuraIsInactiveWithReactionTime(config.GetAuraIsInactiveWithReactionTime())
	case *proto.APLValue_AuraRemainingTime:
		value = rot.newValueAuraRemainingTime(config.GetAuraRemainingTime())
	case *proto.APLValue_AuraNumStacks:
		value = rot.newValueAuraNumStacks(config.GetAuraNumStacks())
	case *proto.APLValue_AuraInternalCooldown:
		value = rot.newValueAuraInternalCooldown(config.GetAuraInternalCooldown())
	case *proto.APLValue_AuraIcdIsReadyWithReactionTime:
		value = rot.newValueAuraICDIsReadyWithReactionTime(config.GetAuraIcdIsReadyWithReactionTime())
	case *proto.APLValue_AuraShouldRefresh:
		value = rot.newValueAuraShouldRefresh(config.GetAuraShouldRefresh())

	// Aura sets
	case *proto.APLValue_AllTrinketStatProcsActive:
		value = rot.newValueAllTrinketStatProcsActive(config.GetAllTrinketStatProcsActive())
	case *proto.APLValue_AnyTrinketStatProcsActive:
		value = rot.newValueAnyTrinketStatProcsActive(config.GetAnyTrinketStatProcsActive())
	case *proto.APLValue_TrinketProcsMinRemainingTime:
		value = rot.newValueTrinketProcsMinRemainingTime(config.GetTrinketProcsMinRemainingTime())
	case *proto.APLValue_TrinketProcsMaxRemainingIcd:
		value = rot.newValueTrinketProcsMaxRemainingICD(config.GetTrinketProcsMaxRemainingIcd())
	case *proto.APLValue_NumEquippedStatProcTrinkets:
		value = rot.newValueNumEquippedStatProcTrinkets(config.GetNumEquippedStatProcTrinkets())

	// Dots
	case *proto.APLValue_DotIsActive:
		value = rot.newValueDotIsActive(config.GetDotIsActive())
	case *proto.APLValue_DotRemainingTime:
		value = rot.newValueDotRemainingTime(config.GetDotRemainingTime())
	case *proto.APLValue_DotTickFrequency:
		value = rot.newValueDotTickFrequency(config.GetDotTickFrequency())

	// Sequences
	case *proto.APLValue_SequenceIsComplete:
		value = rot.newValueSequenceIsComplete(config.GetSequenceIsComplete())
	case *proto.APLValue_SequenceIsReady:
		value = rot.newValueSequenceIsReady(config.GetSequenceIsReady())
	case *proto.APLValue_SequenceTimeToReady:
		value = rot.newValueSequenceTimeToReady(config.GetSequenceTimeToReady())

	// Properties
	case *proto.APLValue_ChannelClipDelay:
		value = rot.newValueChannelClipDelay(config.GetChannelClipDelay())
	case *proto.APLValue_InputDelay:
		value = rot.newValueInputDelay(config.GetInputDelay())

	default:
		value = nil
	}

	if value != nil {
		reflect.ValueOf(value).Elem().FieldByName("Uuid").Set(reflect.ValueOf(config.Uuid))
	}

	return value
}

// Default implementation of Agent.NewAPLValue so each spec doesn't need this boilerplate.
func (unit *Unit) NewAPLValue(rot *APLRotation, config *proto.APLValue) APLValue {
	return nil
}
