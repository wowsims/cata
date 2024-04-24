package shaman

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (shaman *Shaman) registerUnleashFlame() {

	spellMask := SpellMaskLavaBurst | SpellMaskFlameShock | SpellMaskFireNova

	unleashFlameMod := shaman.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		ClassMask:  spellMask,
		FloatValue: 0.2 * (1 + 0.25*float64(shaman.Talents.ElementalWeapons)),
	})

	unleashFlameAura := shaman.RegisterAura(core.Aura{
		Label:    "Unleash Flame",
		ActionID: core.ActionID{SpellID: 73683},
		Duration: time.Second * 8,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			unleashFlameMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			unleashFlameMod.Deactivate()
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spellMask&spell.ClassSpellMask > 0 {
				aura.Deactivate(sim)
			}
		},
	})

	shaman.UnleashFlame = shaman.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 73683},
		SpellSchool:      core.SpellSchoolFire,
		ProcMask:         core.ProcMaskSpellDamage,
		CritMultiplier:   shaman.DefaultSpellCritMultiplier(),
		ClassSpellMask:   SpellMaskUnleashFlame,
		DamageMultiplier: 1,
		BonusCoefficient: 0.429,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, 1118, spell.OutcomeMagicHitAndCrit)
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
		CritMultiplier:   shaman.DefaultSpellCritMultiplier(),
		DamageMultiplier: 1,
		BonusCoefficient: 0.386,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, 873, spell.OutcomeMagicHitAndCrit)
		},
	})
}

func (shaman *Shaman) registerUnleashWind() {

	speedMultiplier := 1 + 0.4*(1+0.25*float64(shaman.Talents.ElementalWeapons))

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
		DamageMultiplier: 1.75,
		CritMultiplier:   shaman.DefaultMeleeCritMultiplier(),
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
				spell.CostMultiplier -= 0.2
			}
		},
	})

	shaman.UnleashLife = shaman.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 73685},
		SpellSchool:      core.SpellSchoolNature,
		ProcMask:         core.ProcMaskSpellHealing,
		Flags:            core.SpellFlagHelpful,
		CritMultiplier:   shaman.DefaultHealingCritMultiplier(),
		BonusCoefficient: 0.201,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealHealing(sim, target, 1996, spell.OutcomeHealingCrit)

			if result.Outcome.Matches(core.OutcomeCrit) {
				if shaman.Talents.AncestralAwakening > 0 {
					shaman.ancestralHealingAmount = result.Damage * 0.3
					shaman.AncestralAwakening.Cast(sim, target)
				}
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
		ActionID: core.ActionID{SpellID: 73680},
		Flags:    core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.07,
			Multiplier: 1 - 0.05*float64(shaman.Talents.Convection) - shaman.GetMentalQuicknessBonus(),
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
