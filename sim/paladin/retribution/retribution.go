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

	pal := paladin.NewPaladin(character, options.TalentsString)

	ret := &RetributionPaladin{
		Paladin: pal,
		Seal:    retOptions.Options.ClassOptions.Seal,
	}

	ret.PaladinAura = retOptions.Options.ClassOptions.Aura

	ret.EnableAutoAttacks(ret, core.AutoAttackOptions{
		MainHand:       ret.WeaponFromMainHand(ret.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
	})

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
	//ret.RegisterAvengingWrathCD()
}

func (ret *RetributionPaladin) ApplyTalents() {
	ret.Paladin.ApplyTalents()
	ret.ApplyArmorSpecializationEffect(stats.Strength, proto.ArmorType_ArmorTypePlate)
}

func (ret *RetributionPaladin) Reset(sim *core.Simulation) {
	ret.Paladin.Reset(sim)

	// switch ret.Seal {
	// case proto.PaladinSeal_Vengeance:
	// 	ret.CurrentSeal = ret.SealOfVengeanceAura
	// 	ret.SealOfVengeanceAura.Activate(sim)
	// case proto.PaladinSeal_Command:
	// 	ret.CurrentSeal = ret.SealOfCommandAura
	// 	ret.SealOfCommandAura.Activate(sim)
	// case proto.PaladinSeal_Righteousness:
	// 	ret.CurrentSeal = ret.SealOfRighteousnessAura
	// 	ret.SealOfRighteousnessAura.Activate(sim)
	// }
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

}

func (ret *RetributionPaladin) RegisterMastery() {
	actionId := core.ActionID{SpellID: 76672}

	// Hand of Light
	handOfLight := ret.RegisterSpell(core.SpellConfig{
		ActionID:         actionId,
		SpellSchool:      core.SpellSchoolHoly,
		ProcMask:         core.ProcMaskMeleeMHSpecial,
		Flags:            core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagNoOnCastComplete,
		DamageMultiplier: 1.0,
		ThreatMultiplier: 1.0,
		CritMultiplier:   ret.DefaultMeleeCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			new_result := spell.CalcOutcome(sim, target, spell.OutcomeAlwaysHit)
			new_result.Damage = ret.HoLDamage
			new_result.Threat = spell.ThreatFromDamage(new_result.Outcome, new_result.Damage)
			spell.DealDamage(sim, new_result)
		},
	})

	core.MakeProcTriggerAura(&ret.Unit, core.ProcTrigger{
		Name:           "Hand of Light",
		ActionID:       actionId,
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ProcMask:       core.ProcMaskMeleeSpecial,
		ClassSpellMask: paladin.SpellMaskCrusaderStrike | paladin.SpellMaskDivineStorm | paladin.SpellMaskTemplarsVerdict,

		ProcChance: 1.0,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			ret.HoLDamage = (16.8 + 2.1*ret.GetMasteryPoints()) / 100.0 * result.Damage
			handOfLight.Cast(sim, result.Target)
		},
	})
}

func (ret *RetributionPaladin) ApplyJudymentOfTheBold() {
	// Judgement of the Bold
	actionID := core.ActionID{SpellID: 89901}
	manaMetrics := ret.NewManaMetrics(actionID)
	var manaPA *core.PendingAction

	jotbAura := ret.RegisterAura(core.Aura{
		Label:    "Judgement of the Bold",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			manaPA = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period: time.Second * 2,
				OnAction: func(sim *core.Simulation) {
					ret.AddMana(sim, 0.25*ret.BaseMana, manaMetrics)
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			manaPA.Cancel(sim)
		},
	})

	core.MakeProcTriggerAura(&ret.Unit, core.ProcTrigger{
		Name:           "Judgement of the Bold",
		ActionID:       actionID,
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: paladin.SpellMaskJudgement,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			jotbAura.Activate(sim)
		},
	})
}
