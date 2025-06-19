package core

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core/proto"
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

type APLValueDotPercentIncrease struct {
	DefaultAPLValueImpl
	dot *Dot
}

func (rot *APLRotation) newValueDotPercentIncrease(config *proto.APLValueDotPercentIncrease, _ *proto.UUID) APLValue {
	dot := rot.GetAPLDot(rot.GetTargetUnit(config.TargetUnit), config.SpellId)
	if dot == nil || dot.Spell.expectedTickDamageInternal == nil {
		return nil
	}

	return &APLValueDotPercentIncrease{
		dot: dot,
	}
}

func (value *APLValueDotPercentIncrease) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}

func (value *APLValueDotPercentIncrease) GetFloat(sim *Simulation) float64 {

	// If DoT is not active expected damage increase is always 100%
	if !value.dot.IsActive() {
		return 1
	}

	// calculate haste based increase
	hasteFactor := 1.0
	if value.dot.affectedByCastSpeed {
		tickPeriod := value.dot.Spell.Unit.ApplyCastSpeedForSpell(value.dot.BaseTickLength, value.dot.Spell).Round(time.Millisecond)
		hasteFactor = float64(value.dot.tickPeriod) / float64(tickPeriod)
	}

	return hasteFactor*value.dot.Spell.ExpectedTickDamage(sim, value.dot.Unit)/
		value.dot.Spell.ExpectedTickDamageFromCurrentSnapshot(sim, value.dot.Unit) - 1
}

func (value *APLValueDotPercentIncrease) String() string {
	return fmt.Sprintf("Dot Percent Increase (%s)", value.dot.Spell.ActionID)
}
