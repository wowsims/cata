package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (shaman *Shaman) registerUnleashFlame() {

	spellMask := SpellMaskLavaBurst | SpellMaskFlameShock | SpellMaskFireNova | SpellMaskElementalBlast

	unleashFlameAura := shaman.RegisterAura(core.Aura{
		Label:    "Unleash Flame",
		ActionID: core.ActionID{SpellID: 73683},
		Duration: time.Second * 8,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(spellMask) && aura.StartedAt() < (sim.CurrentTime-spell.TravelTime()) { // In case unleash element is used during LvB/EB travel time

				//Unleash flame applies to both direct damage and dot,
				//As the 2 parts are separated we wait to deactivate the aura
				pa := &core.PendingAction{
					NextActionAt: sim.CurrentTime + time.Duration(1),
					Priority:     core.ActionPriorityGCD,

					OnAction: func(sim *core.Simulation) {
						aura.Deactivate(sim)
					},
				}
				sim.AddPendingAction(pa)
			}
		},
	}).AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  spellMask,
		FloatValue: 0.3,
	})

	shaman.UnleashFlame = shaman.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 73683},
		SpellSchool:      core.SpellSchoolFire,
		ProcMask:         core.ProcMaskSpellDamage,
		CritMultiplier:   shaman.DefaultCritMultiplier(),
		ClassSpellMask:   SpellMaskUnleashFlame,
		Flags:            SpellFlagFocusable | core.SpellFlagPassiveSpell,
		DamageMultiplier: 1,
		BonusCoefficient: 0.42899999022,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := shaman.CalcAndRollDamageRange(sim, 1.11300003529, 0.17000000179)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			unleashFlameAura.Activate(sim)
		},
	})
}

func (shaman *Shaman) registerUnleashFrost() {

	shaman.UnleashFrost = shaman.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 73682},
		SpellSchool:      core.SpellSchoolFrost,
		ProcMask:         core.ProcMaskSpellDamage,
		ClassSpellMask:   SpellMaskUnleashFrost,
		Flags:            core.SpellFlagPassiveSpell,
		CritMultiplier:   shaman.DefaultCritMultiplier(),
		DamageMultiplier: 1,
		BonusCoefficient: 0.38600000739,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := shaman.CalcAndRollDamageRange(sim, 0.86900001764, 0.15000000596)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

func (shaman *Shaman) registerUnleashWind() {

	speedMultiplier := 1 + 0.5

	unleashWindAura := shaman.RegisterAura(core.Aura{
		Label:     "Unleash Wind",
		ActionID:  core.ActionID{SpellID: 73681},
		Duration:  time.Second * 12,
		MaxStacks: 6,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyMeleeSpeed(sim, speedMultiplier)
			aura.SetStacks(sim, aura.MaxStacks)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyMeleeSpeed(sim, 1/speedMultiplier)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if aura.GetStacks() > 0 && spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				aura.RemoveStack(sim)
			}
		},
	})

	shaman.UnleashWind = shaman.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 73681},
		SpellSchool:      core.SpellSchoolPhysical,
		ProcMask:         core.ProcMaskRangedSpecial,
		Flags:            core.SpellFlagPassiveSpell,
		DamageMultiplier: 0.9,
		CritMultiplier:   shaman.DefaultCritMultiplier(),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeRangedHitAndCrit)
			unleashWindAura.Activate(sim)
		},
	})
}

func (shaman *Shaman) registerUnleashLife() {
	var affectedSpells []*core.Spell

	unleashLifeAura := shaman.RegisterAura(core.Aura{
		Label:    "Unleash Life",
		ActionID: core.ActionID{SpellID: 73685},
		Duration: time.Second * 8,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells = core.FilterSlice([]*core.Spell{
				shaman.ChainHeal,
				shaman.HealingWave,
				shaman.GreaterHealingWave,
				shaman.HealingSurge,
				shaman.Riptide,
			}, func(spell *core.Spell) bool { return spell != nil })
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.DamageMultiplierAdditive += 0.2
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.DamageMultiplierAdditive -= 0.2
			}
		},
	})

	shaman.UnleashLife = shaman.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 73685},
		SpellSchool:      core.SpellSchoolNature,
		ProcMask:         core.ProcMaskSpellHealing,
		Flags:            core.SpellFlagHelpful | core.SpellFlagPassiveSpell,
		CritMultiplier:   shaman.DefaultCritMultiplier(),
		BonusCoefficient: 0.28600001335,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseHeal := shaman.CalcScalingSpellDmg(2.82999992371)
			result := spell.CalcAndDealHealing(sim, target, baseHeal, spell.OutcomeHealingCrit)

			if result.Outcome.Matches(core.OutcomeCrit) {
				shaman.ancestralHealingAmount = result.Damage * 0.3
				shaman.AncestralAwakening.Cast(sim, target)
			}
			unleashLifeAura.Activate(sim)
		},
	})
}

func (shaman *Shaman) registerUnleashElements() {
	unleashElementsTimer := shaman.NewTimer()
	shaman.registerUnleashFlame()
	shaman.registerUnleashFrost()
	shaman.registerUnleashWind()
	shaman.registerUnleashLife()

	shaman.UnleashElements = shaman.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 73680},
		Flags:          SpellFlagShamanSpell | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskUnleashElements,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 8.2,
			PercentModifier: 100,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    unleashElementsTimer,
				Duration: time.Second * 15,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			switch shaman.SelfBuffs.ImbueMH {
			case proto.ShamanImbue_FlametongueWeapon:
				shaman.UnleashFlame.Cast(sim, target)
			case proto.ShamanImbue_WindfuryWeapon:
				shaman.UnleashWind.Cast(sim, target)
			case proto.ShamanImbue_EarthlivingWeapon:
				shaman.UnleashLife.Cast(sim, target)
			case proto.ShamanImbue_FrostbrandWeapon:
				shaman.UnleashFrost.Cast(sim, target)
			}
			if shaman.SelfBuffs.ImbueOH != proto.ShamanImbue_NoImbue && shaman.SelfBuffs.ImbueOH != shaman.SelfBuffs.ImbueMH {
				switch shaman.SelfBuffs.ImbueOH {
				case proto.ShamanImbue_FlametongueWeapon:
					shaman.UnleashFlame.Cast(sim, target)
				case proto.ShamanImbue_WindfuryWeapon:
					shaman.UnleashWind.Cast(sim, target)
				case proto.ShamanImbue_EarthlivingWeapon:
					shaman.UnleashLife.Cast(sim, target)
				case proto.ShamanImbue_FrostbrandWeapon:
					shaman.UnleashFrost.Cast(sim, target)
				}
			}
		},
	})
}
