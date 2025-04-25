package core

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

// Parent struct for all APL values that perform checks on the set of buff auras
// associated with equipped proc items.
type APLValueItemStatProcCheck struct {
	DefaultAPLValueImpl

	baseName        string
	includeWarnings bool

	statTypesToMatch []stats.Stat
	matchingAuras    []*StatBuffAura
}

func (rot *APLRotation) newItemStatProcValue(valueName string, statType1 int32, statType2 int32, statType3 int32, minIcdSeconds float64, requireMatch bool, uuid *proto.UUID) *APLValueItemStatProcCheck {
	statTypesToMatch := stats.IntTupleToStatsList(statType1, statType2, statType3)
	minIcd := DurationFromSeconds(minIcdSeconds)
	matchingAuras := rot.GetAPLItemProcAuras(statTypesToMatch, minIcd, requireMatch, uuid)

	if (len(matchingAuras) == 0) && requireMatch {
		return nil
	}

	return &APLValueItemStatProcCheck{
		baseName:         valueName,
		includeWarnings:  requireMatch,
		statTypesToMatch: statTypesToMatch,
		matchingAuras:    matchingAuras,
	}
}
func (value *APLValueItemStatProcCheck) String() string {
	return fmt.Sprintf("%s(%s)", value.baseName, StringFromStatTypes(value.statTypesToMatch))
}
func (value *APLValueItemStatProcCheck) Finalize(rot *APLRotation) {
	if !value.includeWarnings {
		return
	}

	validAuras := FilterSlice(value.matchingAuras, func(aura *StatBuffAura) bool {
		return !aura.IsSwapped
	})
	actionIDs := MapSlice(validAuras, func(aura *StatBuffAura) ActionID {
		return aura.ActionID
	})

	rot.ValidationMessageByUUID(value.Uuid, proto.LogLevel_Information, "%s will check the following auras: %s", value, StringFromActionIDs(actionIDs))
}

type APLValueAllItemStatProcsActive struct {
	*APLValueItemStatProcCheck
}

func (rot *APLRotation) newValueAllItemStatProcsActive(config *proto.APLValueAllTrinketStatProcsActive, uuid *proto.UUID) APLValue {
	parentImpl := rot.newItemStatProcValue("AllItemStatProcsActive", config.StatType1, config.StatType2, config.StatType3, config.MinIcdSeconds, true, uuid)

	if parentImpl == nil {
		return nil
	}

	return &APLValueAllItemStatProcsActive{
		APLValueItemStatProcCheck: parentImpl,
	}
}
func (value *APLValueAllItemStatProcsActive) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueAllItemStatProcsActive) GetBool(sim *Simulation) bool {
	for _, aura := range value.matchingAuras {
		if (!aura.IsActive() || (aura.GetStacks() < aura.MaxStacks)) && aura.CanProc(sim) {
			return false
		}
	}

	return true
}

type APLValueAnyItemStatProcsActive struct {
	*APLValueItemStatProcCheck
}

func (rot *APLRotation) newValueAnyTrinketStatProcsActive(config *proto.APLValueAnyTrinketStatProcsActive, uuid *proto.UUID) APLValue {
	parentImpl := rot.newItemStatProcValue("AnyItemStatProcsActive", config.StatType1, config.StatType2, config.StatType3, config.MinIcdSeconds, true, uuid)

	if parentImpl == nil {
		return nil
	}

	return &APLValueAnyItemStatProcsActive{
		APLValueItemStatProcCheck: parentImpl,
	}
}
func (value *APLValueAnyItemStatProcsActive) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueAnyItemStatProcsActive) GetBool(sim *Simulation) bool {
	for _, aura := range value.matchingAuras {
		if aura.IsActive() && (aura.GetStacks() == aura.MaxStacks) {
			return true
		}
	}

	return false
}

type APLValueItemProcsMinRemainingTime struct {
	*APLValueItemStatProcCheck
}

func (rot *APLRotation) newValueItemProcsMinRemainingTime(config *proto.APLValueTrinketProcsMinRemainingTime, uuid *proto.UUID) APLValue {
	parentImpl := rot.newItemStatProcValue("ItemProcsMinRemainingTime", config.StatType1, config.StatType2, config.StatType3, config.MinIcdSeconds, true, uuid)

	if parentImpl == nil {
		return nil
	}

	return &APLValueItemProcsMinRemainingTime{
		APLValueItemStatProcCheck: parentImpl,
	}
}
func (value *APLValueItemProcsMinRemainingTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueItemProcsMinRemainingTime) GetDuration(sim *Simulation) time.Duration {
	minRemainingTime := NeverExpires

	for _, aura := range value.matchingAuras {
		if aura.IsActive() {
			minRemainingTime = min(minRemainingTime, aura.RemainingDuration(sim))
		}
	}

	return minRemainingTime
}

type APLValueItemProcsMaxRemainingICD struct {
	*APLValueItemStatProcCheck
}

func (rot *APLRotation) newValueItemsProcsMaxRemainingICD(config *proto.APLValueTrinketProcsMaxRemainingICD, uuid *proto.UUID) APLValue {
	parentImpl := rot.newItemStatProcValue("ItemProcsMaxRemainingICD", config.StatType1, config.StatType2, config.StatType3, config.MinIcdSeconds, true, uuid)

	if parentImpl == nil {
		return nil
	}

	return &APLValueItemProcsMaxRemainingICD{
		APLValueItemStatProcCheck: parentImpl,
	}
}
func (value *APLValueItemProcsMaxRemainingICD) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueItemProcsMaxRemainingICD) GetDuration(sim *Simulation) time.Duration {
	var maxRemainingICD time.Duration

	for _, aura := range value.matchingAuras {
		if aura.CanProc(sim) && !aura.IsActive() && (aura.Icd != nil) {
			maxRemainingICD = max(maxRemainingICD, aura.Icd.TimeToReady(sim))
		}
	}

	return maxRemainingICD
}

type APLValueNumEquippedStatProcItems struct {
	*APLValueItemStatProcCheck
}

func (rot *APLRotation) newValueNumEquippedStatProcItems(config *proto.APLValueNumEquippedStatProcTrinkets, uuid *proto.UUID) APLValue {
	parentImpl := rot.newItemStatProcValue("NumEquippedStatProcItems", config.StatType1, config.StatType2, config.StatType3, config.MinIcdSeconds, false, uuid)

	return &APLValueNumEquippedStatProcItems{
		APLValueItemStatProcCheck: parentImpl,
	}
}
func (value *APLValueNumEquippedStatProcItems) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}
func (value *APLValueNumEquippedStatProcItems) GetInt(sim *Simulation) int32 {
	return int32(len(FilterSlice(value.matchingAuras, func(aura *StatBuffAura) bool {
		return aura.CanProc(sim)
	})))
}

type APLValueNumStatBuffCooldowns struct {
	DefaultAPLValueImpl

	statTypesToMatch []stats.Stat
	matchingSpells   []*Spell
}

func (rot *APLRotation) newValueNumStatBuffCooldowns(config *proto.APLValueNumStatBuffCooldowns, _ *proto.UUID) APLValue {
	unit := rot.unit
	character := unit.Env.Raid.GetPlayerFromUnit(unit).GetCharacter()
	statTypesToMatch := stats.IntTupleToStatsList(config.StatType1, config.StatType2, config.StatType3)
	matchingSpells := character.GetMatchingStatBuffSpells(statTypesToMatch)

	return &APLValueNumStatBuffCooldowns{
		statTypesToMatch: statTypesToMatch,
		matchingSpells:   matchingSpells,
	}
}
func (value *APLValueNumStatBuffCooldowns) String() string {
	return fmt.Sprintf("NumStatBuffCooldowns(%s)", StringFromStatTypes(value.statTypesToMatch))
}
func (value *APLValueNumStatBuffCooldowns) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}
func (value *APLValueNumStatBuffCooldowns) GetInt(_ *Simulation) int32 {
	validSpellCount := int32(0)
	for _, spell := range value.matchingSpells {
		if !spell.Flags.Matches(SpellFlagSwapped) {
			validSpellCount++
		}
	}
	return validSpellCount
}
func (value *APLValueNumStatBuffCooldowns) Finalize(rot *APLRotation) {
	actionIDs := MapSlice(value.matchingSpells, func(spell *Spell) ActionID {
		return spell.ActionID
	})

	rot.ValidationMessageByUUID(value.Uuid, proto.LogLevel_Information, "%s will count the currently equipped subset of: %s", value, StringFromActionIDs(actionIDs))
}
