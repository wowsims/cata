package monk

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (monk *Monk) NewAPLValue(_ *core.APLRotation, config *proto.APLValue) core.APLValue {
	switch config.Value.(type) {
	case *proto.APLValue_MonkCurrentChi:
		return monk.newValueCurrentChi(config.GetMonkCurrentChi(), config.Uuid)
	case *proto.APLValue_MonkMaxChi:
		return monk.newValueMaxChi(config.GetMonkMaxChi(), config.Uuid)
	case *proto.APLValue_MonkNextChiBrewRecharge:
		return monk.newValueNextChiBrewRecharge(config.GetMonkNextChiBrewRecharge(), config.Uuid)
	default:
		return nil
	}
}

type APLValueCurrentChi struct {
	core.DefaultAPLValueImpl
	monk *Monk
}

func (monk *Monk) newValueCurrentChi(_ *proto.APLValueMonkCurrentChi, _ *proto.UUID) core.APLValue {
	return &APLValueCurrentChi{
		monk: monk,
	}
}
func (value *APLValueCurrentChi) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}
func (value *APLValueCurrentChi) GetInt(_ *core.Simulation) int32 {
	return value.monk.ComboPoints()
}
func (value *APLValueCurrentChi) String() string {
	return "Current Chi"
}

type APLValueMaxChi struct {
	core.DefaultAPLValueImpl
	maxChi int32
}

func (monk *Monk) newValueMaxChi(_ *proto.APLValueMonkMaxChi, _ *proto.UUID) core.APLValue {
	return &APLValueMaxChi{
		maxChi: monk.MaxComboPoints(),
	}
}
func (value *APLValueMaxChi) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}
func (value *APLValueMaxChi) GetInt(_ *core.Simulation) int32 {
	return value.maxChi
}
func (value *APLValueMaxChi) String() string {
	return fmt.Sprintf("Max Chi(%d)", value.maxChi)
}

type APLValueNextChiBrewRecharge struct {
	core.DefaultAPLValueImpl
	monk *Monk
}

func (monk *Monk) newValueNextChiBrewRecharge(_ *proto.APLValueMonkNextChiBrewRecharge, _ *proto.UUID) core.APLValue {
	return &APLValueNextChiBrewRecharge{
		monk: monk,
	}
}
func (value *APLValueNextChiBrewRecharge) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueNextChiBrewRecharge) GetDuration(sim *core.Simulation) time.Duration {
	if value.monk.chiBrewRecharge != nil && value.monk.chiBrewRecharge.NextActionAt > sim.CurrentTime {
		return value.monk.chiBrewRecharge.NextActionAt - sim.CurrentTime
	}

	return 0
}
func (value *APLValueNextChiBrewRecharge) String() string {
	return "Time To Next Chi Brew Recharge"
}
