package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (shaman *Shaman) ApplyTalents() {
	shaman.ApplyElementalMastery()
	shaman.ApplyAncestralSwiftness()
	shaman.ApplyEchoOfTheElements()
	shaman.ApplyUnleashedFury()
	shaman.ApplyPrimalElementalist()
	shaman.ApplyElementalBlast()

	shaman.ApplyGlyphs()
}

func (shaman *Shaman) ApplyElementalMastery() {
	if !shaman.Talents.ElementalMastery {
		return
	}

	eleMasterActionID := core.ActionID{SpellID: 16166}

	buffAura := shaman.RegisterAura(core.Aura{
		Label:    "Elemental Mastery Buff",
		ActionID: core.ActionID{SpellID: 64701},
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyCastSpeed(1.30)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyCastSpeed(1 / 1.30)
		},
	})

	eleMastSpell := shaman.RegisterSpell(core.SpellConfig{
		ActionID:       eleMasterActionID,
		ClassSpellMask: SpellMaskElementalMastery,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 90,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			buffAura.Activate(sim)
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: eleMastSpell,
		Type:  core.CooldownTypeDPS,
	})
}

func (shaman *Shaman) ApplyAncestralSwiftness() {
	if !shaman.Talents.AncestralSwiftness {
		return
	}

	asCdTimer := shaman.NewTimer()
	asCd := time.Second * 90

	affectedSpells := SpellMaskLightningBolt | SpellMaskChainLightning | SpellMaskEarthShock | SpellMaskElementalBlast
	AncestralSwiftnessBuffAura := core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label:    "Ancestral swiftness Buff",
		ActionID: core.ActionID{SpellID: 64701},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyCastSpeed(1.05)
			shaman.MultiplyMeleeSpeed(sim, 1.1)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyCastSpeed(1 / 1.30)
			shaman.MultiplyMeleeSpeed(sim, 1/1.05)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ClassSpellMask&affectedSpells == 0 {
				return
			}
			asCdTimer.Set(sim.CurrentTime + asCd)
			shaman.UpdateMajorCooldowns()
			aura.Deactivate(sim)
		},
	}))

	asSpell := shaman.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 16188},
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			//TODO goes on cd only when buff is consumed
			CD: core.Cooldown{
				Timer:    asCdTimer,
				Duration: asCd,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			AncestralSwiftnessBuffAura.Activate(sim)
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: asSpell,
		Type:  core.CooldownTypeDPS,
	})
}

func (shaman *Shaman) ApplyEchoOfTheElements() {
	if !shaman.Talents.EchoOfTheElements {
		return
	}

	//TODO like dtr 5%
}

func (shaman *Shaman) ApplyUnleashedFury() {
	if !shaman.Talents.UnleashedFury {
		return
	}

	//TODO
}

func (shaman *Shaman) ApplyPrimalElementalist() {
	if !shaman.Talents.PrimalElementalist {
		return
	}

	//TODO
}

func (shaman *Shaman) ApplyElementalBlast() {
	if !shaman.Talents.ElementalBlast {
		return
	}
	shaman.RegisterElementalBlastSpell()
}
