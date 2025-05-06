package core

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core/proto"
)

type APLValueCurrentHealth struct {
	DefaultAPLValueImpl
	unit UnitReference
}

func (rot *APLRotation) newValueCurrentHealth(config *proto.APLValueCurrentHealth, uuid *proto.UUID) APLValue {
	unit := rot.GetSourceUnit(config.SourceUnit)
	if unit.Get() == nil {
		return nil
	}
	if !unit.Get().HasHealthBar() {
		rot.ValidationMessageByUUID(uuid, proto.LogLevel_Warning, "%s does not use Health", unit.Get().Label)
		return nil
	}
	return &APLValueCurrentHealth{
		unit: unit,
	}
}
func (value *APLValueCurrentHealth) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueCurrentHealth) GetFloat(sim *Simulation) float64 {
	return value.unit.Get().CurrentHealth()
}
func (value *APLValueCurrentHealth) String() string {
	return "Current Health"
}

type APLValueCurrentHealthPercent struct {
	DefaultAPLValueImpl
	unit UnitReference
}

func (rot *APLRotation) newValueCurrentHealthPercent(config *proto.APLValueCurrentHealthPercent, uuid *proto.UUID) APLValue {
	unit := rot.GetSourceUnit(config.SourceUnit)
	if unit.Get() == nil {
		return nil
	}
	if !unit.Get().HasHealthBar() {
		rot.ValidationMessageByUUID(uuid, proto.LogLevel_Warning, "%s does not use Health", unit.Get().Label)
		return nil
	}
	return &APLValueCurrentHealthPercent{
		unit: unit,
	}
}
func (value *APLValueCurrentHealthPercent) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueCurrentHealthPercent) GetFloat(sim *Simulation) float64 {
	return value.unit.Get().CurrentHealthPercent()
}
func (value *APLValueCurrentHealthPercent) String() string {
	return fmt.Sprintf("Current Health %%")
}

type APLValueCurrentMana struct {
	DefaultAPLValueImpl
	unit UnitReference
}

func (rot *APLRotation) newValueCurrentMana(config *proto.APLValueCurrentMana, uuid *proto.UUID) APLValue {
	unit := rot.GetSourceUnit(config.SourceUnit)
	if unit.Get() == nil {
		return nil
	}
	if !unit.Get().HasManaBar() {
		rot.ValidationMessageByUUID(uuid, proto.LogLevel_Warning, "%s does not use Mana", unit.Get().Label)
		return nil
	}
	return &APLValueCurrentMana{
		unit: unit,
	}
}
func (value *APLValueCurrentMana) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueCurrentMana) GetFloat(sim *Simulation) float64 {
	return value.unit.Get().CurrentMana()
}
func (value *APLValueCurrentMana) String() string {
	return "Current Mana"
}

type APLValueCurrentManaPercent struct {
	DefaultAPLValueImpl
	unit UnitReference
}

func (rot *APLRotation) newValueCurrentManaPercent(config *proto.APLValueCurrentManaPercent, uuid *proto.UUID) APLValue {
	unit := rot.GetSourceUnit(config.SourceUnit)
	if unit.Get() == nil {
		return nil
	}
	if !unit.Get().HasManaBar() {
		rot.ValidationMessageByUUID(uuid, proto.LogLevel_Warning, "%s does not use Mana", unit.Get().Label)
		return nil
	}
	return &APLValueCurrentManaPercent{
		unit: unit,
	}
}
func (value *APLValueCurrentManaPercent) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueCurrentManaPercent) GetFloat(sim *Simulation) float64 {
	return value.unit.Get().CurrentManaPercent()
}
func (value *APLValueCurrentManaPercent) String() string {
	return fmt.Sprintf("Current Mana %%")
}

type APLValueCurrentRage struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueCurrentRage(config *proto.APLValueCurrentRage, uuid *proto.UUID) APLValue {
	unit := rot.unit
	if !unit.HasRageBar() {
		rot.ValidationMessageByUUID(uuid, proto.LogLevel_Warning, "%s does not use Rage", unit.Label)
		return nil
	}
	return &APLValueCurrentRage{
		unit: unit,
	}
}
func (value *APLValueCurrentRage) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueCurrentRage) GetFloat(sim *Simulation) float64 {
	return value.unit.CurrentRage()
}
func (value *APLValueCurrentRage) String() string {
	return "Current Rage"
}

type APLValueCurrentFocus struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueCurrentFocus(config *proto.APLValueCurrentFocus, uuid *proto.UUID) APLValue {
	unit := rot.unit
	if !unit.HasFocusBar() {
		rot.ValidationMessageByUUID(uuid, proto.LogLevel_Warning, "%s does not use Focus", unit.Label)
		return nil
	}
	return &APLValueCurrentFocus{
		unit: unit,
	}
}

func (value *APLValueCurrentFocus) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}

func (value *APLValueCurrentFocus) GetFloat(sim *Simulation) float64 {
	return value.unit.CurrentFocus()
}

func (value *APLValueCurrentFocus) String() string {
	return "Current Focus"
}

type APLValueMaxFocus struct {
	DefaultAPLValueImpl
	maxFocus float64
}

func (rot *APLRotation) newValueMaxFocus(_ *proto.APLValueMaxFocus, uuid *proto.UUID) APLValue {
	unit := rot.unit
	if !unit.HasFocusBar() {
		rot.ValidationMessageByUUID(uuid, proto.LogLevel_Error, "%s does not use Focus", unit.Label)
		return nil
	}
	return &APLValueMaxFocus{
		maxFocus: unit.MaximumFocus(),
	}
}
func (value *APLValueMaxFocus) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueMaxFocus) GetFloat(sim *Simulation) float64 {
	return value.maxFocus
}
func (value *APLValueMaxFocus) String() string {
	return fmt.Sprintf("Max Focus(%f)", value.maxFocus)
}

type APLValueFocusRegenPerSecond struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueFocusRegenPerSecond(_ *proto.APLValueFocusRegenPerSecond, uuid *proto.UUID) APLValue {
	unit := rot.unit
	if !unit.HasFocusBar() {
		rot.ValidationMessageByUUID(uuid, proto.LogLevel_Warning, "%s does not use Focus", unit.Label)
		return nil
	}
	return &APLValueFocusRegenPerSecond{
		unit: unit,
	}
}
func (value *APLValueFocusRegenPerSecond) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueFocusRegenPerSecond) GetFloat(sim *Simulation) float64 {
	return value.unit.FocusRegenPerSecond()
}
func (value *APLValueFocusRegenPerSecond) String() string {
	return "Focus Regen Per Second"
}

type APLValueFocusTimeToTarget struct {
	DefaultAPLValueImpl
	unit        *Unit
	targetFocus APLValue
}

func (rot *APLRotation) newValueFocusTimeToTarget(config *proto.APLValueFocusTimeToTarget, uuid *proto.UUID) APLValue {
	unit := rot.unit
	if !unit.HasFocusBar() {
		rot.ValidationMessageByUUID(uuid, proto.LogLevel_Warning, "%s does not use Focus", unit.Label)
		return nil
	}

	targetFocus := rot.coerceTo(rot.newAPLValue(config.TargetFocus), proto.APLValueType_ValueTypeFloat)
	if targetFocus == nil {
		return nil
	}

	return &APLValueFocusTimeToTarget{
		unit:        unit,
		targetFocus: targetFocus,
	}
}
func (value *APLValueFocusTimeToTarget) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueFocusTimeToTarget) GetDuration(sim *Simulation) time.Duration {
	return value.unit.TimeToTargetFocus(value.targetFocus.GetFloat(sim))
}
func (value *APLValueFocusTimeToTarget) String() string {
	return "Estimated Time To Target Focus"
}

type APLValueCurrentEnergy struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueCurrentEnergy(config *proto.APLValueCurrentEnergy, uuid *proto.UUID) APLValue {
	unit := rot.unit
	if !unit.HasEnergyBar() {
		rot.ValidationMessageByUUID(uuid, proto.LogLevel_Warning, "%s does not use Energy", unit.Label)
		return nil
	}
	return &APLValueCurrentEnergy{
		unit: unit,
	}
}
func (value *APLValueCurrentEnergy) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueCurrentEnergy) GetFloat(sim *Simulation) float64 {
	return value.unit.CurrentEnergy()
}
func (value *APLValueCurrentEnergy) String() string {
	return "Current Energy"
}

type APLValueMaxEnergy struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueMaxEnergy(_ *proto.APLValueMaxEnergy, uuid *proto.UUID) APLValue {
	unit := rot.unit
	if !unit.HasEnergyBar() {
		rot.ValidationMessageByUUID(uuid, proto.LogLevel_Error, "%s does not use Energy", unit.Label)
		return nil
	}
	return &APLValueMaxEnergy{
		unit: unit,
	}
}
func (value *APLValueMaxEnergy) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueMaxEnergy) GetFloat(sim *Simulation) float64 {
	return value.unit.MaximumEnergy()
}
func (value *APLValueMaxEnergy) String() string {
	return "Max Energy"
}

type APLValueEnergyRegenPerSecond struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueEnergyRegenPerSecond(_ *proto.APLValueEnergyRegenPerSecond, uuid *proto.UUID) APLValue {
	unit := rot.unit
	if !unit.HasEnergyBar() {
		rot.ValidationMessageByUUID(uuid, proto.LogLevel_Warning, "%s does not use Energy", unit.Label)
		return nil
	}
	return &APLValueEnergyRegenPerSecond{
		unit: unit,
	}
}
func (value *APLValueEnergyRegenPerSecond) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueEnergyRegenPerSecond) GetFloat(sim *Simulation) float64 {
	return value.unit.EnergyRegenPerSecond()
}
func (value *APLValueEnergyRegenPerSecond) String() string {
	return "Energy Regen Per Second"
}

type APLValueEnergyTimeToTarget struct {
	DefaultAPLValueImpl
	unit         *Unit
	targetEnergy APLValue
}

func (rot *APLRotation) newValueEnergyTimeToTarget(config *proto.APLValueEnergyTimeToTarget, uuid *proto.UUID) APLValue {
	unit := rot.unit
	if !unit.HasEnergyBar() {
		rot.ValidationMessageByUUID(uuid, proto.LogLevel_Warning, "%s does not use Energy", unit.Label)
		return nil
	}

	targetEnergy := rot.coerceTo(rot.newAPLValue(config.TargetEnergy), proto.APLValueType_ValueTypeFloat)
	if targetEnergy == nil {
		return nil
	}

	return &APLValueEnergyTimeToTarget{
		unit:         unit,
		targetEnergy: targetEnergy,
	}
}
func (value *APLValueEnergyTimeToTarget) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueEnergyTimeToTarget) GetDuration(sim *Simulation) time.Duration {
	return value.unit.TimeToTargetEnergy(value.targetEnergy.GetFloat(sim))
}
func (value *APLValueEnergyTimeToTarget) String() string {
	return "Estimated Time To Target Energy"
}

type APLValueCurrentComboPoints struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueCurrentComboPoints(config *proto.APLValueCurrentComboPoints, uuid *proto.UUID) APLValue {
	unit := rot.unit
	if !unit.HasEnergyBar() {
		rot.ValidationMessageByUUID(uuid, proto.LogLevel_Warning, "%s does not use Combo Points", unit.Label)
		return nil
	}
	return &APLValueCurrentComboPoints{
		unit: unit,
	}
}
func (value *APLValueCurrentComboPoints) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}
func (value *APLValueCurrentComboPoints) GetInt(sim *Simulation) int32 {
	return value.unit.ComboPoints()
}
func (value *APLValueCurrentComboPoints) String() string {
	return "Current Combo Points"
}

type APLValueMaxComboPoints struct {
	DefaultAPLValueImpl
	maxComboPoints int32
}

func (rot *APLRotation) newValueMaxComboPoints(_ *proto.APLValueMaxComboPoints, uuid *proto.UUID) APLValue {
	unit := rot.unit
	if !unit.HasEnergyBar() {
		rot.ValidationMessageByUUID(uuid, proto.LogLevel_Error, "%s does not use Combo Points", unit.Label)
		return nil
	}
	return &APLValueMaxComboPoints{
		maxComboPoints: unit.MaxComboPoints(),
	}
}
func (value *APLValueMaxComboPoints) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}
func (value *APLValueMaxComboPoints) GetInt(sim *Simulation) int32 {
	return value.maxComboPoints
}
func (value *APLValueMaxComboPoints) String() string {
	return fmt.Sprintf("Max Combo Points(%d)", value.maxComboPoints)
}

type APLValueCurrentRunicPower struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueCurrentRunicPower(config *proto.APLValueCurrentRunicPower, uuid *proto.UUID) APLValue {
	unit := rot.unit
	if !unit.HasRunicPowerBar() {
		rot.ValidationMessageByUUID(uuid, proto.LogLevel_Warning, "%s does not use Runic Power", unit.Label)
		return nil
	}
	return &APLValueCurrentRunicPower{
		unit: unit,
	}
}
func (value *APLValueCurrentRunicPower) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}
func (value *APLValueCurrentRunicPower) GetInt(sim *Simulation) int32 {
	return int32(value.unit.CurrentRunicPower())
}
func (value *APLValueCurrentRunicPower) String() string {
	return "Current Runic Power"
}

type APLValueMaxRunicPower struct {
	DefaultAPLValueImpl
	maxRunicPower int32
}

func (rot *APLRotation) newValueMaxRunicPower(_ *proto.APLValueMaxRunicPower, uuid *proto.UUID) APLValue {
	unit := rot.unit
	if !unit.HasRunicPowerBar() {
		rot.ValidationMessageByUUID(uuid, proto.LogLevel_Error, "%s does not use Runic Power", unit.Label)
		return nil
	}
	return &APLValueMaxRunicPower{
		maxRunicPower: int32(unit.MaximumRunicPower()),
	}
}
func (value *APLValueMaxRunicPower) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}
func (value *APLValueMaxRunicPower) GetInt(sim *Simulation) int32 {
	return value.maxRunicPower
}
func (value *APLValueMaxRunicPower) String() string {
	return fmt.Sprintf("Max Runic Power(%d)", value.maxRunicPower)
}

type APLValueCurrentGenericResource struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueCurrentGenericResource(_ *proto.APLValueCurrentGenericResource, uuid *proto.UUID) APLValue {
	unit := rot.unit
	if unit.secondaryResourceBar == nil {
		rot.ValidationMessageByUUID(uuid, proto.LogLevel_Warning, "%s does not have secondary resource", unit.Label)
		return nil
	}
	return &APLValueCurrentGenericResource{
		unit: unit,
	}
}
func (value *APLValueCurrentGenericResource) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}
func (value *APLValueCurrentGenericResource) GetInt(sim *Simulation) int32 {
	return value.unit.secondaryResourceBar.Value()
}
func (value *APLValueCurrentGenericResource) String() string {
	return "Current {GENERIC_RESOURCE}"
}
