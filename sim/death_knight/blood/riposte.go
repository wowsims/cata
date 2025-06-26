package blood

import (
	"github.com/wowsims/mop/sim/common/shared"
)

func (bdk *BloodDeathKnight) registerRiposte() {
	shared.RegisterRiposteEffect(&bdk.Character, 145677, 145676)
}
