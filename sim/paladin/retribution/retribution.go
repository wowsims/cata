package retribution

import (
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
		Paladin: paladin.NewPaladin(character, options.TalentsString),
		Seal:    retOptions.Options.ClassOptions.Seal,
	}

	ret.PaladinAura = retOptions.Options.ClassOptions.Aura

	return ret
}

type RetributionPaladin struct {
	*paladin.Paladin

	Seal      proto.PaladinSeal
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
	ret.ApplyArmorSpecializationEffect(stats.Strength, proto.ArmorType_ArmorTypePlate)
}

func (ret *RetributionPaladin) Reset(sim *core.Simulation) {
	ret.Paladin.Reset(sim)
}

func (ret *RetributionPaladin) RegisterSpecializationEffects() {
	ret.RegisterMastery()

	// Sheath of Light
	ret.AddStatDependency(stats.AttackPower, stats.SpellPower, 0.3)
	ret.AddStat(stats.SpellHit, core.SpellHitRatingPerHitChance*8)

	// Two-Handed Weapon Specialization
	mhWeapon := ret.GetMHWeapon()
	if mhWeapon != nil && mhWeapon.HandType == proto.HandType_HandTypeTwoHand {
		ret.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.25
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
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagNoOnCastComplete,
		ClassSpellMask: paladin.SpellMaskHandOfLight,

		DamageMultiplier: 1.0,
		ThreatMultiplier: 1.0,
		CritMultiplier:   ret.DefaultMeleeCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			newResult := spell.CalcOutcome(sim, target, spell.OutcomeAlwaysHit)
			// TODO: this damage needs to be manually boosted by any 8% magic damage taken debuff present on the target.
			newResult.Damage = ret.HoLDamage
			newResult.Threat = spell.ThreatFromDamage(newResult.Outcome, newResult.Damage)
			spell.DealDamage(sim, newResult)
		},
	})

	core.MakeProcTriggerAura(&ret.Unit, core.ProcTrigger{
		Name:           "Hand of Light",
		ActionID:       actionId,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ProcMask:       core.ProcMaskMeleeSpecial,
		ClassSpellMask: paladin.SpellMaskCrusaderStrike | paladin.SpellMaskDivineStorm | paladin.SpellMaskTemplarsVerdict,
		ProcChance:     1.0,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			ret.HoLDamage = (16.8 + 2.1*ret.GetMasteryPoints()) / 100.0 * result.Damage
			ret.HandOfLight.Cast(sim, result.Target)
		},
	})
}

func (ret *RetributionPaladin) ApplyJudgmentsOfTheBold() {
	actionID := core.ActionID{SpellID: 89901}
	manaMetrics := ret.NewManaMetrics(actionID)
	var pa *core.PendingAction

	jotbAura := ret.RegisterAura(core.Aura{
		Label:    "Judgements of the Bold",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			pa = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period: time.Second * 2,
				OnAction: func(sim *core.Simulation) {
					ret.AddMana(sim, 0.25*ret.BaseMana, manaMetrics)
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			pa.Cancel(sim)
		},
	})

	core.MakeProcTriggerAura(&ret.Unit, core.ProcTrigger{
		Name:       "Judgements of the Bold Trigger",
		ActionID:   actionID,
		Callback:   core.CallbackOnSpellHitDealt,
		Outcome:    core.OutcomeLanded,
		ProcMask:   core.ProcMaskMeleeSpecial,
		ProcChance: 1.0,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			jotbAura.Activate(sim)
		},
	})
}
