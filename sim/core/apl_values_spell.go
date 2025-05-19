package core

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core/proto"
)

type APLValueSpellIsKnown struct {
	DefaultAPLValueImpl
	spell *Spell
}

func (rot *APLRotation) newValueSpellIsKnown(config *proto.APLValueSpellIsKnown, _ *proto.UUID) APLValue {
	spell := rot.GetAPLSpell(config.SpellId)
	return &APLValueSpellIsKnown{
		spell: spell,
	}
}
func (value *APLValueSpellIsKnown) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueSpellIsKnown) GetBool(sim *Simulation) bool {
	return value.spell != nil
}
func (value *APLValueSpellIsKnown) String() string {
	return fmt.Sprintf("Is Known(%s)", value.spell.ActionID)
}

type APLValueSpellCanCast struct {
	DefaultAPLValueImpl
	spell *Spell
}

func (rot *APLRotation) newValueSpellCanCast(config *proto.APLValueSpellCanCast, _ *proto.UUID) APLValue {
	spell := rot.GetAPLSpell(config.SpellId)
	if spell == nil {
		return nil
	}
	return &APLValueSpellCanCast{
		spell: spell,
	}
}
func (value *APLValueSpellCanCast) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueSpellCanCast) GetBool(sim *Simulation) bool {
	return value.spell.CanCastOrQueue(sim, value.spell.Unit.CurrentTarget)
}
func (value *APLValueSpellCanCast) String() string {
	return fmt.Sprintf("Can Cast(%s)", value.spell.ActionID)
}

type APLValueSpellIsReady struct {
	DefaultAPLValueImpl
	spell *Spell
}

func (rot *APLRotation) newValueSpellIsReady(config *proto.APLValueSpellIsReady, _ *proto.UUID) APLValue {
	spell := rot.GetAPLSpell(config.SpellId)
	if spell == nil {
		return nil
	}
	return &APLValueSpellIsReady{
		spell: spell,
	}
}
func (value *APLValueSpellIsReady) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueSpellIsReady) GetBool(sim *Simulation) bool {
	return value.spell.IsReady(sim) || (value.spell.TimeToReady(sim) <= MaxSpellQueueWindow)
}
func (value *APLValueSpellIsReady) String() string {
	return fmt.Sprintf("Is Ready(%s)", value.spell.ActionID)
}

type APLValueSpellTimeToReady struct {
	DefaultAPLValueImpl
	spell *Spell
}

func (rot *APLRotation) newValueSpellTimeToReady(config *proto.APLValueSpellTimeToReady, _ *proto.UUID) APLValue {
	spell := rot.GetAPLSpell(config.SpellId)
	if spell == nil {
		return nil
	}
	return &APLValueSpellTimeToReady{
		spell: spell,
	}
}
func (value *APLValueSpellTimeToReady) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueSpellTimeToReady) GetDuration(sim *Simulation) time.Duration {
	return value.spell.TimeToReady(sim)
}
func (value *APLValueSpellTimeToReady) GetFloat(sim *Simulation) float64 {
	return value.spell.TimeToReady(sim).Seconds()
}
func (value *APLValueSpellTimeToReady) String() string {
	return fmt.Sprintf("Time To Ready(%s)", value.spell.ActionID)
}

type APLValueSpellCastTime struct {
	DefaultAPLValueImpl
	spell *Spell
}

func (rot *APLRotation) newValueSpellCastTime(config *proto.APLValueSpellCastTime, _ *proto.UUID) APLValue {
	spell := rot.GetAPLSpell(config.SpellId)
	if spell == nil {
		return nil
	}
	return &APLValueSpellCastTime{
		spell: spell,
	}
}
func (value *APLValueSpellCastTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueSpellCastTime) GetDuration(_ *Simulation) time.Duration {
	return value.spell.CastTime()
}
func (value *APLValueSpellCastTime) String() string {
	return fmt.Sprintf("Cast Time(%s)", value.spell.ActionID)
}

type APLValueSpellTravelTime struct {
	DefaultAPLValueImpl
	spell *Spell
}

func (rot *APLRotation) newValueSpellTravelTime(config *proto.APLValueSpellTravelTime, _ *proto.UUID) APLValue {
	spell := rot.GetAPLSpell(config.SpellId)
	if spell == nil {
		return nil
	}
	return &APLValueSpellTravelTime{
		spell: spell,
	}
}
func (value *APLValueSpellTravelTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueSpellTravelTime) GetDuration(_ *Simulation) time.Duration {
	return value.spell.TravelTime()
}
func (value *APLValueSpellTravelTime) String() string {
	return fmt.Sprintf("Travel Time(%s)", value.spell.ActionID)
}

type APLValueSpellCPM struct {
	DefaultAPLValueImpl
	spell *Spell
}

func (rot *APLRotation) newValueSpellCPM(config *proto.APLValueSpellCPM, _ *proto.UUID) APLValue {
	spell := rot.GetAPLSpell(config.SpellId)
	if spell == nil {
		return nil
	}
	return &APLValueSpellCPM{
		spell: spell,
	}
}
func (value *APLValueSpellCPM) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueSpellCPM) GetFloat(sim *Simulation) float64 {
	return value.spell.CurCPM(sim)
}
func (value *APLValueSpellCPM) String() string {
	return fmt.Sprintf("CPM(%s)", value.spell.ActionID)
}

type APLValueSpellIsChanneling struct {
	DefaultAPLValueImpl
	spell *Spell
}

func (rot *APLRotation) newValueSpellIsChanneling(config *proto.APLValueSpellIsChanneling, _ *proto.UUID) APLValue {
	spell := rot.GetAPLSpell(config.SpellId)
	if spell == nil {
		return nil
	}
	return &APLValueSpellIsChanneling{
		spell: spell,
	}
}
func (value *APLValueSpellIsChanneling) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueSpellIsChanneling) GetBool(_ *Simulation) bool {
	return value.spell.Unit.ChanneledDot != nil && value.spell.Unit.ChanneledDot.Spell == value.spell
}
func (value *APLValueSpellIsChanneling) String() string {
	return fmt.Sprintf("IsChanneling(%s)", value.spell.ActionID)
}

type APLValueSpellChanneledTicks struct {
	DefaultAPLValueImpl
	spell *Spell
}

func (rot *APLRotation) newValueSpellChanneledTicks(config *proto.APLValueSpellChanneledTicks, _ *proto.UUID) APLValue {
	spell := rot.GetAPLSpell(config.SpellId)
	if spell == nil {
		return nil
	}
	return &APLValueSpellChanneledTicks{
		spell: spell,
	}
}
func (value *APLValueSpellChanneledTicks) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}
func (value *APLValueSpellChanneledTicks) GetInt(_ *Simulation) int32 {
	channeledDot := value.spell.Unit.ChanneledDot
	if channeledDot == nil {
		return 0
	} else {
		return channeledDot.TickCount()
	}
}
func (value *APLValueSpellChanneledTicks) String() string {
	return fmt.Sprintf("ChanneledTicks(%s)", value.spell.ActionID)
}

type APLValueSpellCurrentCost struct {
	DefaultAPLValueImpl
	spell *Spell
}

func (rot *APLRotation) newValueSpellCurrentCost(config *proto.APLValueSpellCurrentCost, _ *proto.UUID) APLValue {
	spell := rot.GetAPLSpell(config.SpellId)
	if spell == nil {
		return nil
	}
	return &APLValueSpellCurrentCost{
		spell: spell,
	}
}
func (value *APLValueSpellCurrentCost) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueSpellCurrentCost) GetFloat(_ *Simulation) float64 {
	spell := value.spell
	if spell.Cost == nil {
		return 0
	}
	return spell.Cost.GetCurrentCost()
}
func (value *APLValueSpellCurrentCost) String() string {
	return fmt.Sprintf("CurrentCost(%s)", value.spell.ActionID)
}

// Spell Charges
type APLValueSpellNumCharges struct {
	DefaultAPLValueImpl
	spell *Spell
}

func (rot *APLRotation) newValueSpellNumCharges(config *proto.APLValueSpellNumCharges, _ *proto.UUID) APLValue {
	spell := rot.GetAPLSpell(config.SpellId)
	if spell == nil {
		return nil
	}
	return &APLValueSpellNumCharges{
		spell: spell,
	}
}

func (value *APLValueSpellNumCharges) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}

func (value *APLValueSpellNumCharges) GetInt(_ *Simulation) int32 {
	return int32(value.spell.GetNumCharges())
}

func (value *APLValueSpellNumCharges) String() string {
	return fmt.Sprintf("SpellNumCharges(%s)", value.spell.ActionID)
}

type APLValueSpellTimeToCharge struct {
	DefaultAPLValueImpl
	spell *Spell
}

func (rot *APLRotation) newValueSpellTimeToCharge(config *proto.APLValueSpellTimeToCharge, _ *proto.UUID) APLValue {
	spell := rot.GetAPLSpell(config.SpellId)
	if spell == nil {
		return nil
	}
	return &APLValueSpellTimeToCharge{
		spell: spell,
	}
}

func (value *APLValueSpellTimeToCharge) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}

func (value *APLValueSpellTimeToCharge) GetDuration(sim *Simulation) time.Duration {
	return value.spell.NextChargeIn(sim)
}

func (value *APLValueSpellTimeToCharge) GetFloat(sim *Simulation) float64 {
	return value.GetDuration(sim).Seconds()
}

func (value *APLValueSpellTimeToCharge) String() string {
	return fmt.Sprintf("SpellTimeToCharge(%s)", value.spell.ActionID)
}
