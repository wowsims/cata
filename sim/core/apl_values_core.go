package core

import (
	"fmt"
	"time"

	"github.com/wowsims/cata/sim/core/proto"
)

type APLValueDotIsActive struct {
	DefaultAPLValueImpl
	dot *Dot
}

func (rot *APLRotation) newValueDotIsActive(config *proto.APLValueDotIsActive, _ *proto.UUID) APLValue {
	dot := rot.GetAPLDot(rot.GetTargetUnit(config.TargetUnit), config.SpellId)
	if dot == nil {
		return nil
	}
	return &APLValueDotIsActive{
		dot: dot,
	}
}
func (value *APLValueDotIsActive) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueDotIsActive) GetBool(sim *Simulation) bool {
	return value.dot.IsActive()
}
func (value *APLValueDotIsActive) String() string {
	return fmt.Sprintf("Dot Is Active(%s)", value.dot.Spell.ActionID)
}

type APLValueDotRemainingTime struct {
	DefaultAPLValueImpl
	dot *Dot
}

func (rot *APLRotation) newValueDotRemainingTime(config *proto.APLValueDotRemainingTime, _ *proto.UUID) APLValue {
	dot := rot.GetAPLDot(rot.GetTargetUnit(config.TargetUnit), config.SpellId)
	if dot == nil {
		return nil
	}
	return &APLValueDotRemainingTime{
		dot: dot,
	}
}
func (value *APLValueDotRemainingTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueDotRemainingTime) GetDuration(sim *Simulation) time.Duration {
	return TernaryDuration(value.dot.IsActive(), value.dot.RemainingDuration(sim), 0)
}
func (value *APLValueDotRemainingTime) String() string {
	return fmt.Sprintf("Dot Remaining Time(%s)", value.dot.Spell.ActionID)
}

type APLValueDotTickFrequency struct {
	DefaultAPLValueImpl
	dot *Dot
}

func (rot *APLRotation) newValueDotTickFrequency(config *proto.APLValueDotTickFrequency, _ *proto.UUID) APLValue {
	dot := rot.GetAPLDot(rot.GetTargetUnit(config.TargetUnit), config.SpellId)
	if dot == nil {
		return nil
	}
	return &APLValueDotTickFrequency{
		dot: dot,
	}
}

func (value *APLValueDotTickFrequency) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueDotTickFrequency) GetDuration(_ *Simulation) time.Duration {
	return value.dot.tickPeriod
}
func (value *APLValueDotTickFrequency) String() string {
	return fmt.Sprintf("Dot Tick Frequency(%s)", value.dot.tickPeriod)
}
