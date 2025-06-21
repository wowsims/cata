package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (mage *Mage) registerArmorSpells() {

	mageArmorEffectCategory := "MageArmors"

	moltenArmor := mage.RegisterAura(core.Aura{
		Label:    "Molten Armor",
		ActionID: core.ActionID{SpellID: 30482},
		Duration: core.NeverExpires,
	}).AttachStatBuff(stats.SpellCritPercent, 5)

	moltenArmor.NewExclusiveEffect(mageArmorEffectCategory, true, core.ExclusiveEffect{})

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 30482},
		SpellSchool:    core.SpellSchoolFire,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
		ClassSpellMask: MageSpellMoltenArmor,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 3,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !moltenArmor.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			moltenArmor.Activate(sim)
		},
	})

	mageArmor := mage.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 6117},
		Label:    "Mage Armor",
		Duration: core.NeverExpires,
	}).AttachStatBuff(stats.MasteryRating, 3000.0)

	mageArmor.NewExclusiveEffect(mageArmorEffectCategory, true, core.ExclusiveEffect{})

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 6117},
		SpellSchool:    core.SpellSchoolArcane,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
		ClassSpellMask: MageSpellMageArmor,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 3,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !mageArmor.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mageArmor.Activate(sim)
		},
	})

	frostArmor := mage.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 7302},
		Label:    "Frost Armor",
		Duration: core.NeverExpires,
	}).AttachMultiplyCastSpeed(1.07)

	frostArmor.NewExclusiveEffect(mageArmorEffectCategory, true, core.ExclusiveEffect{})

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 7302},
		SpellSchool:    core.SpellSchoolFrost,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
		ClassSpellMask: MageSpellFrostArmor,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 3,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !frostArmor.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			frostArmor.Activate(sim)
		},
	})
}
