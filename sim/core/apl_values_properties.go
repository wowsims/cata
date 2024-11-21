package core

import (
	"time"

	"github.com/wowsims/cata/sim/core/proto"
)

type APLValueChannelClipDelay struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueChannelClipDelay(config *proto.APLValueChannelClipDelay, uuid *proto.UUID) APLValue {
	return &APLValueChannelClipDelay{
		unit: rot.unit,
	}
}
func (value *APLValueChannelClipDelay) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueChannelClipDelay) GetDuration(sim *Simulation) time.Duration {
	return value.unit.ChannelClipDelay
}
func (value *APLValueChannelClipDelay) String() string {
	return "Channel Clip Delay()"
}

type APLValueInputDelay struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueInputDelay(config *proto.APLValueInputDelay, uuid *proto.UUID) APLValue {
	return &APLValueInputDelay{
		unit: rot.unit,
	}
}
func (value *APLValueInputDelay) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueInputDelay) GetDuration(sim *Simulation) time.Duration {
	return value.unit.ReactionTime
}
func (value *APLValueInputDelay) String() string {
	return "Channel Clip Delay()"
}

type APLValueFrontOfTarget struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueFrontOfTarget(config *proto.APLValueFrontOfTarget, uuid *proto.UUID) APLValue {
	return &APLValueFrontOfTarget{
		unit: rot.unit,
	}
}
func (value *APLValueFrontOfTarget) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueFrontOfTarget) GetBool(sim *Simulation) bool {
	return value.unit.PseudoStats.InFrontOfTarget
}
func (value *APLValueFrontOfTarget) String() string {
	return "Front of Target()"
}
