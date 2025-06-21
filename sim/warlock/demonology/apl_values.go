package demonology

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (warlock *DemonologyWarlock) NewAPLValue(rot *core.APLRotation, config *proto.APLValue) core.APLValue {
	switch config.Value.(type) {
	case *proto.APLValue_WarlockHandOfGuldanInFlight:
		return warlock.newValueWarlockHandOfGuldanInFlight(rot, config.GetWarlockHandOfGuldanInFlight())
	default:
		return warlock.Warlock.NewAPLValue(rot, config)
	}
}

type APLValueWarlockHandOfGuldanInFlight struct {
	core.DefaultAPLValueImpl
	warlock *DemonologyWarlock
}

func (warlock *DemonologyWarlock) newValueWarlockHandOfGuldanInFlight(rot *core.APLRotation, config *proto.APLValueWarlockHandOfGuldanInFlight) core.APLValue {
	return &APLValueWarlockHandOfGuldanInFlight{
		warlock: warlock,
	}
}
func (value *APLValueWarlockHandOfGuldanInFlight) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueWarlockHandOfGuldanInFlight) GetBool(sim *core.Simulation) bool {
	warlock := value.warlock
	return warlock.HandOfGuldanImpactTime > 0 && sim.CurrentTime < warlock.HandOfGuldanImpactTime
}
func (value *APLValueWarlockHandOfGuldanInFlight) String() string {
	return "Warlock Hand of Guldan in Flight()"
}
