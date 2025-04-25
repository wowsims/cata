package arcane

import (
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	mop.CreateDTRClassConfig(proto.Spec_SpecArcaneMage, 1.0/12)
}
