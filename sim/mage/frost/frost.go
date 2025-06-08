package frost

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/mage"
)

func RegisterFrostMage() {
	core.RegisterAgentFactory(
		proto.Player_FrostMage{},
		proto.Spec_SpecFrostMage,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewFrostMage(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_FrostMage)
			if !ok {
				panic("Invalid spec value for Frost Mage!")
			}
			player.Spec = playerSpec
		},
	)
}

type FrostMage struct {
	*mage.Mage

	waterElemental             *WaterElemental
	frozenOrb                  *FrozenOrb
	frostfireFrozenCritBuffMod *core.SpellMod
	iceLanceFrozenCritBuffMod  *core.SpellMod
}

func NewFrostMage(character *core.Character, options *proto.Player) *FrostMage {
	frostOptions := options.GetFrostMage().Options

	frostMage := &FrostMage{
		Mage: mage.NewMage(character, options, frostOptions.ClassOptions),
	}
	frostMage.waterElemental = frostMage.NewWaterElemental()
	frostMage.frozenOrb = frostMage.NewFrozenOrb()

	return frostMage
}

func (frostMage *FrostMage) GetMage() *mage.Mage {
	return frostMage.Mage
}

func (frostMage *FrostMage) Reset(sim *core.Simulation) {
	frostMage.Mage.Reset(sim)
}

func (frostMage *FrostMage) Initialize() {
	frostMage.Mage.Initialize()

	frostMage.registerSummonWaterElementalSpell()
	frostMage.registerFingersOfFrost()
	frostMage.registerBrainFreeze()
	frostMage.registerFrostboltSpell()
	frostMage.registerFrozenOrbSpell()
}

func (frostMage *FrostMage) ApplyTalents() {
	frostMage.Mage.ApplyTalents()

	frostMage.frostfireFrozenCritBuffMod = frostMage.Mage.AddDynamicMod(core.SpellModConfig{
		FloatValue: frostMage.GetStat(stats.SpellCritPercent)*2 + 50,
		ClassMask:  mage.MageSpellFrostfireBolt,
		Kind:       core.SpellMod_CritMultiplier_Flat,
	})

	frostMage.iceLanceFrozenCritBuffMod = frostMage.Mage.AddDynamicMod(core.SpellModConfig{
		FloatValue: frostMage.GetStat(stats.SpellCritPercent)*2 + 50,
		ClassMask:  mage.MageSpellIceLance,
		Kind:       core.SpellMod_CritMultiplier_Flat,
	})

	frostMage.AddOnTemporaryStatsChange(func(sim *core.Simulation, buffAura *core.Aura, statsChangeWithoutDeps stats.Stats) {
		frostMage.frostfireFrozenCritBuffMod.UpdateFloatValue(frostMage.GetStat(stats.SpellCritPercent)*2 + 50)
		frostMage.iceLanceFrozenCritBuffMod.UpdateFloatValue(frostMage.GetStat(stats.SpellCritPercent)*2 + 50)
	})

	frostMasteryMod := frostMage.waterElemental.AddDynamicMod(core.SpellModConfig{
		ClassMask:  mage.MageWaterElementalSpellWaterBolt,
		FloatValue: frostMage.GetFrostMasteryBonus(),
		Kind:       core.SpellMod_DamageDone_Pct,
	})

	frostMage.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery, newMastery float64) {
		masteryBonus := frostMage.GetFrostMasteryBonus()
		frostMasteryMod.UpdateFloatValue(masteryBonus)
	})

	core.MakePermanent(frostMage.RegisterAura(core.Aura{
		Label: "Mastery: Icicles - Water Elemental",
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			frostMasteryMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			frostMasteryMod.Deactivate()
		},
	}))

}
