package protection

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/warrior"
)

func RegisterProtectionWarrior() {
	core.RegisterAgentFactory(
		proto.Player_ProtectionWarrior{},
		proto.Spec_SpecProtectionWarrior,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewProtectionWarrior(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_ProtectionWarrior)
			if !ok {
				panic("Invalid spec value for Protection Warrior!")
			}
			player.Spec = playerSpec
		},
	)
}

type ProtectionWarrior struct {
	*warrior.Warrior

	Options *proto.ProtectionWarrior_Options

	shieldSlam *core.Spell
}

func NewProtectionWarrior(character *core.Character, options *proto.Player) *ProtectionWarrior {
	protOptions := options.GetProtectionWarrior().Options

	war := &ProtectionWarrior{
		Warrior: warrior.NewWarrior(character, options.TalentsString, warrior.WarriorInputs{}),
		Options: protOptions,
	}

	rbo := core.RageBarOptions{
		StartingRage:       protOptions.ClassOptions.StartingRage,
		BaseRageMultiplier: 1,
	}

	war.EnableRageBar(rbo)
	war.EnableAutoAttacks(war, core.AutoAttackOptions{
		MainHand:       war.WeaponFromMainHand(war.DefaultCritMultiplier()),
		AutoSwingMelee: true,
	})

	return war
}

func (war *ProtectionWarrior) RegisterSpecializationEffects() {
	// Critical block
	war.RegisterMastery()

	// Sentinel stat buffs
	war.MultiplyStat(stats.Stamina, 1.15)
	war.AddStat(stats.BlockPercent, 15)

	// Vengeance
	war.RegisterVengeance(93098, war.DefensiveStanceAura)
}

func (war *ProtectionWarrior) RegisterMastery() {

	dummyCriticalBlockSpell := war.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 76857}, // Doesn't seem like there's an actual spell ID for the block itself, so use the mastery ID
		Flags:    core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,
	})

	war.Blockhandler = func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		if !spell.SpellSchool.Matches(core.SpellSchoolPhysical) {
			return
		}

		if result.Outcome.Matches(core.OutcomeBlock) && !result.Outcome.Matches(core.OutcomeMiss) && !result.Outcome.Matches(core.OutcomeParry) && !result.Outcome.Matches(core.OutcomeDodge) {
			procChance := war.GetCriticalBlockChance()
			if sim.Proc(procChance, "Critical Block Roll") {
				result.Damage = result.Damage * (1 - war.BlockDamageReduction()*2)
				dummyCriticalBlockSpell.Cast(sim, spell.Unit)
				return
			}
			result.Damage = result.Damage * (1 - war.BlockDamageReduction())
		}
	}

	// Crit block mastery also applies an equal amount to regular block
	// set initial block % from stats
	war.CriticalBlockChance[0] = war.CalculateCriticalBlockChance()
	war.AddStat(stats.BlockPercent, war.CriticalBlockChance[0]*100.0)

	// and keep it updated when mastery changes
	war.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMasteryRating float64, newMasteryRating float64) {
		war.AddStatDynamic(sim, stats.BlockPercent, 1.5*core.MasteryRatingToMasteryPoints(newMasteryRating-oldMasteryRating))
		war.CriticalBlockChance[0] = war.CalculateCriticalBlockChance()
	})

}

func CalcMasteryPercent(points float64) float64 {
	return 12.0 + 1.5*points
}

func (war *ProtectionWarrior) CalculateCriticalBlockChance() float64 {
	return CalcMasteryPercent(war.GetMasteryPoints()) / 100.0
}

func (war *ProtectionWarrior) GetWarrior() *warrior.Warrior {
	return war.Warrior
}

func (war *ProtectionWarrior) Initialize() {
	war.Warrior.Initialize()
	war.RegisterSpecializationEffects()
	war.RegisterShieldSlam()
}

func (war *ProtectionWarrior) ApplyTalents() {}

func (war *ProtectionWarrior) Reset(sim *core.Simulation) {
	war.Warrior.Reset(sim)
}
