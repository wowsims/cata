package affliction

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (warlock *AfflictionWarlock) NewAPLValue(rot *core.APLRotation, config *proto.APLValue) core.APLValue {
	switch config.Value.(type) {
	case *proto.APLValue_WarlockHauntInFlight:
		return warlock.newValueWarlockHauntInFlight(rot, config.GetWarlockHauntInFlight())
	default:
		return warlock.Warlock.NewAPLValue(rot, config)
	}
}

type APLValueWarlockHauntInFlight struct {
	core.DefaultAPLValueImpl
	warlock *AfflictionWarlock
}

func (warlock *AfflictionWarlock) newValueWarlockHauntInFlight(_ *core.APLRotation, _ *proto.APLValueWarlockHauntInFlight) core.APLValue {
	return &APLValueWarlockHauntInFlight{
		warlock: warlock,
	}
}
func (value *APLValueWarlockHauntInFlight) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueWarlockHauntInFlight) GetBool(sim *core.Simulation) bool {
	warlock := value.warlock
	return warlock.HauntImpactTime > 0 && sim.CurrentTime < warlock.HauntImpactTime
}
func (value *APLValueWarlockHauntInFlight) String() string {
	return "Warlock Haunt in Flight()"
}
