package mage

import (
	"github.com/wowsims/mop/sim/core"
)

func (mage *Mage) ApplyMastery() {
	//These aren't technically spells but I'm not sure how else to create them

	mage.Icicle = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 148022},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellIcicle,
		MissileSpeed:   20,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, 0, spell.OutcomeMagicHit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})

	// leaving this as a stub as I still have to redo the water elemental code.
	waterElementalDamageMod := mage.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: mage.GetFrostMasteryBonus(),
	})

	mage.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMasteryRating, newMasteryRating float64) {
		waterElementalDamageMod.UpdateFloatValue(mage.GetFrostMasteryBonus())
	})

	core.MakePermanent(mage.RegisterAura(core.Aura{
		Label: "Mastery: Icicles",
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			waterElementalDamageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			waterElementalDamageMod.Deactivate()
		},
	}))
}

func (mage *Mage) castIcicleWithDamage(sim *core.Simulation, target *core.Unit, damage float64) {
	mage.Icicle.DamageMultiplier *= damage
	mage.Icicle.Cast(sim, target)
	mage.Icicle.DamageMultiplier /= damage
}

func (mage *Mage) HandleIcicleGeneration(sim *core.Simulation, target *core.Unit, baseDamage float64) {
	numIcicles := len(mage.icicles)
	if numIcicles == 5 {
		mage.castIcicleWithDamage(sim, target, mage.icicles[0])
		mage.icicles = mage.icicles[1:]
	}
	mage.icicles = append(mage.icicles, baseDamage*mage.GetFrostMasteryBonus())
}

func (mage *Mage) HandleUseAllIcicles(sim *core.Simulation, target *core.Unit) {
	for _, icicle := range mage.icicles {
		mage.castIcicleWithDamage(sim, target, icicle)
	}
	mage.icicles = make([]float64, 0) // Note this only really works if the code executes purely sequentially.
}
