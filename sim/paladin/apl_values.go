package paladin

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (paladin *Paladin) NewAPLValue(rot *core.APLRotation, config *proto.APLValue) core.APLValue {
	switch config.Value.(type) {
	case *proto.APLValue_CurrentHolyPower:
		return paladin.newValueCurrentHolyPower(config.GetCurrentHolyPower())
	default:
		return nil
	}
}

type APLValueCurrentHolyPower struct {
	core.DefaultAPLValueImpl
	paladin *Paladin
}

func (paladin *Paladin) newValueCurrentHolyPower(_ *proto.APLValueCurrentHolyPower) core.APLValue {
	if !paladin.HasHolyPowerBar() {
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
	return value.paladin.CurrentHolyPower()
}
func (value *APLValueCurrentHolyPower) String() string {
	return "Current Holy Power"
}
