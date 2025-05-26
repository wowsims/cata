package frost

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/mage"
)

func (frostMage *FrostMage) ApplyMastery() {
	//These aren't technically spells but I'm not sure how else to create them

	frostMage.icicleCast = frostMage.Mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 148022},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: mage.MageSpellIcicle,
		MissileSpeed:   20,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, 0, spell.OutcomeMagicHit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})

	frostMage.icicleDamageMod = frostMage.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.0,
		ClassMask:  mage.MageSpellIcicle,
	})

	// leaving this as a stub as I still have to redo the water elemental code.
	waterElementalDamageMod := frostMage.Mage.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: frostMage.GetMasteryBonus(),
	})

	frostMage.Mage.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMasteryRating, newMasteryRating float64) {
		waterElementalDamageMod.UpdateFloatValue(frostMage.GetMasteryBonus())
	})

	core.MakePermanent(frostMage.Mage.RegisterAura(core.Aura{
		Label: "Mastery: Icicles",
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			waterElementalDamageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			waterElementalDamageMod.Deactivate()
		},
	}))
}

func (frostMage *FrostMage) castIcicleWithDamage(sim *core.Simulation, target *core.Unit, damage float64) {
	frostMage.icicleDamageMod.UpdateFloatValue(damage)
	frostMage.icicleCast.Cast(sim, target)
	frostMage.icicleDamageMod.UpdateFloatValue(0.0)
}

func (frostMage *FrostMage) handleIcicleGeneration(sim *core.Simulation, target *core.Unit, baseDamage float64) {
	numIcicles := len(frostMage.icicles)
	if numIcicles == 5 {
		frostMage.castIcicleWithDamage(sim, target, frostMage.icicles[0])
		frostMage.icicles = frostMage.icicles[1:]
	}
	frostMage.icicles = append(frostMage.icicles, baseDamage*frostMage.GetMasteryBonus())
}

func (frostMage *FrostMage) handleUseAllIcicles(sim *core.Simulation, target *core.Unit) {
	for i := int32(0); i < int32(len(frostMage.icicles)); i++ {
		frostMage.castIcicleWithDamage(sim, target, frostMage.icicles[i])
	}
	frostMage.icicles = make([]float64, 0) // Note this only really works if the code executes purely sequentially.
}
