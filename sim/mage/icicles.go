package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (mage *Mage) ApplyMastery() {
	//These aren't technically spells but I'm not sure how else to create them

	mage.Icicle = mage.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 148022},
		SpellSchool:      core.SpellSchoolFrost,
		ProcMask:         core.ProcMaskSpellDamage,
		Flags:            core.SpellFlagAPL,
		ClassSpellMask:   MageSpellIcicle,
		MissileSpeed:     20,
		DamageMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, 1, spell.OutcomeMagicHit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})

	mage.IciclesAura = mage.RegisterAura(core.Aura{
		Label:     "Mastery: Icicles",
		ActionID:  core.ActionID{SpellID: 148022},
		Duration:  time.Hour * 1,
		MaxStacks: 5,
	})

	mage.IciclesAura.OnStacksChange = func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
		if newStacks == 0 {
			mage.IciclesAura.Deactivate(sim)
		}
	}

	core.MakeProcTriggerAura(&mage.Unit, core.ProcTrigger{
		Name:           "Icicles - Trigger",
		ClassSpellMask: MageSpellFrostbolt | MageSpellFrostfireBolt,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			mage.IciclesAura.Activate(sim)
			mage.IciclesAura.AddStack(sim)
		},
	})
}

func (mage *Mage) castIcicleWithDamage(sim *core.Simulation, target *core.Unit, damage float64) {
	mage.Icicle.DamageMultiplier *= damage
	mage.Icicle.Cast(sim, target)
	mage.Icicle.DamageMultiplier /= damage
	if mage.IciclesAura.IsActive() {
		mage.IciclesAura.RemoveStack(sim)
	}
}

func (mage *Mage) HandleIcicleGeneration(sim *core.Simulation, target *core.Unit, baseDamage float64) {
	numIcicles := len(mage.Icicles)
	if numIcicles == 5 {
		mage.castIcicleWithDamage(sim, target, mage.Icicles[0])
		mage.Icicles = mage.Icicles[1:]
	}
	mage.Icicles = append(mage.Icicles, baseDamage*mage.GetFrostMasteryBonus())
}

func (mage *Mage) HandleUseAllIcicles(sim *core.Simulation, target *core.Unit) {
	for _, icicle := range mage.Icicles {
		mage.castIcicleWithDamage(sim, target, icicle)
	}
	mage.Icicles = make([]float64, 0) // Note this only really works if the code executes purely sequentially.
}
