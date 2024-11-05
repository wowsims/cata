package core

import (
	"fmt"
	"time"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

// Parent struct for all APL values that perform checks on the set of buff auras
// associated with equipped proc trinkets.
type APLValueTrinketStatProcCheck struct {
	DefaultAPLValueImpl

	baseName        string
	includeWarnings bool

	statTypesToMatch []stats.Stat
	matchingAuras    []*StatBuffAura
}

func (rot *APLRotation) newTrinketStatProcValue(valueName string, statType1 int32, statType2 int32, statType3 int32, excludeStackingProcs bool, requireMatch bool) *APLValueTrinketStatProcCheck {
	statTypesToMatch := stats.IntTupleToStatsList(statType1, statType2, statType3)
	matchingAuras := rot.GetAPLTrinketProcAuras(statTypesToMatch, excludeStackingProcs, requireMatch)

	if (len(matchingAuras) == 0) && requireMatch {
		return nil
	}

	return &APLValueTrinketStatProcCheck{
		baseName:         valueName,
		includeWarnings:  requireMatch,
		statTypesToMatch: statTypesToMatch,
		matchingAuras:    matchingAuras,
	}
}
func (value *APLValueTrinketStatProcCheck) String() string {
	return fmt.Sprintf("%s(%s)", value.baseName, StringFromStatTypes(value.statTypesToMatch))
}
func (value *APLValueTrinketStatProcCheck) Finalize(rot *APLRotation) {
	if !value.includeWarnings {
		return
	}

	actionIDs := MapSlice(value.matchingAuras, func(aura *StatBuffAura) ActionID {
		return aura.ActionID
	})

	rot.ValidationWarning("%s will check the following auras: %s", value, StringFromActionIDs(actionIDs))
}

type APLValueAllTrinketStatProcsActive struct {
	*APLValueTrinketStatProcCheck
}

func (rot *APLRotation) newValueAllTrinketStatProcsActive(config *proto.APLValueAllTrinketStatProcsActive) APLValue {
	parentImpl := rot.newTrinketStatProcValue("AllTrinketStatProcsActive", config.StatType1, config.StatType2, config.StatType3, config.ExcludeStackingProcs, true)

	if parentImpl == nil {
		return nil
	}

	return &APLValueAllTrinketStatProcsActive{
		APLValueTrinketStatProcCheck: parentImpl,
	}
}
func (value *APLValueAllTrinketStatProcsActive) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueAllTrinketStatProcsActive) GetBool(sim *Simulation) bool {
	for _, aura := range value.matchingAuras {
		if (!aura.IsActive() || (aura.GetStacks() < aura.MaxStacks)) && aura.CanProc(sim) {
			return false
		}
	}

	return true
}

type APLValueAnyTrinketStatProcsActive struct {
	*APLValueTrinketStatProcCheck
}

func (rot *APLRotation) newValueAnyTrinketStatProcsActive(config *proto.APLValueAnyTrinketStatProcsActive) APLValue {
	parentImpl := rot.newTrinketStatProcValue("AnyTrinketStatProcsActive", config.StatType1, config.StatType2, config.StatType3, config.ExcludeStackingProcs, true)

	if parentImpl == nil {
		return nil
	}

	return &APLValueAnyTrinketStatProcsActive{
		APLValueTrinketStatProcCheck: parentImpl,
	}
}
func (value *APLValueAnyTrinketStatProcsActive) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueAnyTrinketStatProcsActive) GetBool(sim *Simulation) bool {
	for _, aura := range value.matchingAuras {
		if aura.IsActive() && (aura.GetStacks() == aura.MaxStacks) {
			return true
		}
	}

	return false
}

type APLValueTrinketProcsMinRemainingTime struct {
	*APLValueTrinketStatProcCheck
}

func (rot *APLRotation) newValueTrinketProcsMinRemainingTime(config *proto.APLValueTrinketProcsMinRemainingTime) APLValue {
	parentImpl := rot.newTrinketStatProcValue("TrinketProcsMinRemainingTime", config.StatType1, config.StatType2, config.StatType3, config.ExcludeStackingProcs, true)

	if parentImpl == nil {
		return nil
	}

	return &APLValueTrinketProcsMinRemainingTime{
		APLValueTrinketStatProcCheck: parentImpl,
	}
}
func (value *APLValueTrinketProcsMinRemainingTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueTrinketProcsMinRemainingTime) GetDuration(sim *Simulation) time.Duration {
	minRemainingTime := NeverExpires

	for _, aura := range value.matchingAuras {
		if aura.IsActive() {
			minRemainingTime = min(minRemainingTime, aura.RemainingDuration(sim))
		}
	}

	return minRemainingTime
}

type APLValueNumEquippedStatProcTrinkets struct {
	*APLValueTrinketStatProcCheck
}

func (rot *APLRotation) newValueNumEquippedStatProcTrinkets(config *proto.APLValueNumEquippedStatProcTrinkets) APLValue {
	parentImpl := rot.newTrinketStatProcValue("NumEquippedStatProcTrinkets", config.StatType1, config.StatType2, config.StatType3, config.ExcludeStackingProcs, false)

	return &APLValueNumEquippedStatProcTrinkets{
		APLValueTrinketStatProcCheck: parentImpl,
	}
}
func (value *APLValueNumEquippedStatProcTrinkets) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}
func (value *APLValueNumEquippedStatProcTrinkets) GetInt(sim *Simulation) int32 {
	return int32(len(FilterSlice(value.matchingAuras, func(aura *StatBuffAura) bool {
		return aura.CanProc(sim)
	})))
}
