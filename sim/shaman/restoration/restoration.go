package restoration

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/shaman"
)

func RegisterRestorationShaman() {
	core.RegisterAgentFactory(
		proto.Player_RestorationShaman{},
		proto.Spec_SpecRestorationShaman,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewRestorationShaman(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_RestorationShaman)
			if !ok {
				panic("Invalid spec value for Restoration Shaman!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewRestorationShaman(character *core.Character, options *proto.Player) *RestorationShaman {
	restoOptions := options.GetRestorationShaman().Options

	selfBuffs := shaman.SelfBuffs{
		Shield: restoOptions.ClassOptions.Shield,
	}

	resto := &RestorationShaman{
		Shaman: shaman.NewShaman(character, options.TalentsString, selfBuffs, false, restoOptions.ClassOptions.FeleAutocast),
	}

	// if resto.HasMHWeapon() {
	// 	resto.ApplyEarthlivingImbueToItem(resto.GetMHWeapon())
	// }

	return resto
}

type RestorationShaman struct {
	*shaman.Shaman
}

func (resto *RestorationShaman) GetShaman() *shaman.Shaman {
	return resto.Shaman
}

func (resto *RestorationShaman) Reset(sim *core.Simulation) {
	resto.Shaman.Reset(sim)
}
func (resto *RestorationShaman) GetMainTarget() *core.Unit {
	// TODO: make this just grab first player that isn't self.
	target := resto.Env.Raid.GetFirstTargetDummy()
	if target == nil {
		return &resto.Unit
	} else {
		return &target.Unit
	}
}

func (resto *RestorationShaman) Initialize() {
	// resto.CurrentTarget = resto.GetMainTarget()

	// // Has to be here because earthliving can cast hots and needs Env to be set to create the hots.
	// procMask := core.ProcMaskUnknown
	// if resto.HasMHWeapon() {
	// 	procMask |= core.ProcMaskMeleeMH
	// }
	// if resto.HasOHWeapon() {
	// 	procMask |= core.ProcMaskMeleeOH
	// }
	// resto.RegisterEarthlivingImbue(procMask)

	resto.Shaman.Initialize()
	resto.Shaman.RegisterHealingSpells()
}

func (resto *RestorationShaman) ApplyTalents() {
	resto.Shaman.ApplyTalents()
	resto.ApplyArmorSpecializationEffect(stats.Intellect, proto.ArmorType_ArmorTypeMail, 86529)
}
