package paladin

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (paladin *Paladin) NewAPLValue(rot *core.APLRotation, config *proto.APLValue) core.APLValue {
	switch config.Value.(type) {
	case *proto.APLValue_CurrentHolyPower:
		return paladin.newValueCurrentHolyPower(config.GetCurrentHolyPower(), config.Uuid)
	default:
		return nil
	}
}

type APLValueCurrentHolyPower struct {
	core.DefaultAPLValueImpl
	paladin *Paladin
}

func (paladin *Paladin) newValueCurrentHolyPower(_ *proto.APLValueCurrentHolyPower, uuid *proto.UUID) core.APLValue {
	if paladin.HolyPower == nil {
		return nil
	}

	return &APLValueCurrentHolyPower{
		paladin: paladin,
	}
}
func (value *APLValueCurrentHolyPower) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}
func (value *APLValueCurrentHolyPower) GetInt(sim *core.Simulation) int32 {
	return int32(value.paladin.HolyPower.Value())
}
func (value *APLValueCurrentHolyPower) String() string {
	return "Current Holy Power"
}
