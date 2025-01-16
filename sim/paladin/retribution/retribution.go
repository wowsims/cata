package retribution

import (
	"math"
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
	"github.com/wowsims/cata/sim/paladin"
)

func RegisterRetributionPaladin() {
	core.RegisterAgentFactory(
		proto.Player_RetributionPaladin{},
		proto.Spec_SpecRetributionPaladin,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewRetributionPaladin(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_RetributionPaladin) // I don't really understand this line
			if !ok {
				panic("Invalid spec value for Retribution Paladin!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewRetributionPaladin(character *core.Character, options *proto.Player) *RetributionPaladin {
	retOptions := options.GetRetributionPaladin()

	ret := &RetributionPaladin{
		Paladin: paladin.NewPaladin(character, options.TalentsString, retOptions.Options.ClassOptions),
	}
	ret.StartingHolyPower = retOptions.Options.StartingHolyPower

	return ret
}

type RetributionPaladin struct {
	*paladin.Paladin

	HoLDamage float64
}

func (ret *RetributionPaladin) GetPaladin() *paladin.Paladin {
	return ret.Paladin
}

func (ret *RetributionPaladin) Initialize() {
	ret.Paladin.Initialize()
	ret.RegisterSpecializationEffects()
	ret.RegisterTemplarsVerdict()
}

func (ret *RetributionPaladin) ApplyTalents() {
	ret.Paladin.ApplyTalents()
	ret.ApplyArmorSpecializationEffect(stats.Strength, proto.ArmorType_ArmorTypePlate, 86525)
}

func (ret *RetributionPaladin) Reset(sim *core.Simulation) {
	ret.Paladin.Reset(sim)
}

func (ret *RetributionPaladin) RegisterSpecializationEffects() {
	ret.RegisterMastery()

	// Sheath of Light
	ret.AddStatDependency(stats.AttackPower, stats.SpellPower, 0.3)
	ret.AddStat(stats.SpellHitPercent, 8)

	// Two-Handed Weapon Specialization
	mhWeapon := ret.GetMHWeapon()
	if mhWeapon != nil && mhWeapon.HandType == proto.HandType_HandTypeTwoHand {
		ret.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.25
		ret.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Pct,
			ClassMask:  paladin.SpellMaskModifiedByTwoHandedSpec,
			FloatValue: 0.25,
		})
	}

	// Judgements of the Bold
	ret.ApplyJudgmentsOfTheBold()
}

func (ret *RetributionPaladin) RegisterMastery() {
	actionId := core.ActionID{SpellID: 76672}

	// Hand of Light
	ret.HandOfLight = ret.RegisterSpell(core.SpellConfig{
		ActionID:       actionId,
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagIgnoreModifiers | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,
		ClassSpellMask: paladin.SpellMaskHandOfLight,

		DamageMultiplier: 1.0,
		ThreatMultiplier: 1.0,
		CritMultiplier:   0.0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := ret.HoLDamage
			if target.HasActiveAuraWithTag(core.SpellDamageEffectAuraTag) {
				baseDamage *= 1.08
			}
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)
		},
	})

	core.MakeProcTriggerAura(&ret.Unit, core.ProcTrigger{
		Name:           "Hand of Light" + ret.Label,
		ActionID:       actionId,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ProcMask:       core.ProcMaskMeleeSpecial,
		ClassSpellMask: paladin.SpellMaskCrusaderStrike | paladin.SpellMaskDivineStorm | paladin.SpellMaskTemplarsVerdict,
		ProcChance:     1.0,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			ret.HoLDamage = ((16.8 + 2.1*ret.GetMasteryPoints()) / 100.0) * result.Damage
			ret.HandOfLight.Cast(sim, result.Target)
		},
	})
}

func (ret *RetributionPaladin) ApplyJudgmentsOfTheBold() {
	actionID := core.ActionID{SpellID: 89901}
	manaMetrics := ret.NewManaMetrics(actionID)

	// It's 25% of base mana over 10 seconds, with haste adding ticks.
	manaPerTick := math.Round(0.025 * ret.BaseMana)

	jotb := ret.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagHelpful | core.SpellFlagNoMetrics | core.SpellFlagNoLogs,

		Hot: core.DotConfig{
			SelfOnly: true,
			Aura: core.Aura{
				Label: "Judgements of the Bold" + ret.Label,
			},
			NumberOfTicks:        10,
			TickLength:           time.Second * 1,
			AffectedByCastSpeed:  true,
			HasteReducesDuration: false,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				ret.AddMana(sim, manaPerTick, manaMetrics)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.SelfHot().Apply(sim)
		},
	})

	core.MakeProcTriggerAura(&ret.Unit, core.ProcTrigger{
		Name:           "Judgements of the Bold Trigger" + ret.Label,
		ActionID:       actionID,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: paladin.SpellMaskJudgement,
		ProcChance:     1.0,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			jotb.Cast(sim, &ret.Unit)
		},
	})
}
