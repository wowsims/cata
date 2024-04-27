package protection

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
	"github.com/wowsims/cata/sim/warrior"
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

	core.VengeanceTracker

	shieldSlam *core.Spell
}

func NewProtectionWarrior(character *core.Character, options *proto.Player) *ProtectionWarrior {
	protOptions := options.GetProtectionWarrior().Options

	war := &ProtectionWarrior{
		Warrior: warrior.NewWarrior(character, options.TalentsString, warrior.WarriorInputs{}),
		Options: protOptions,
	}

	war.Character.PrimaryStat = stats.Stamina

	rbo := core.RageBarOptions{
		StartingRage:   protOptions.ClassOptions.StartingRage,
		RageMultiplier: 1.0,

		OnHitDealtRageGain: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult, rage float64) float64 {
			if result.Target != nil && result.Target.CurrentTarget != &war.Unit {
				return rage * 1.5 // Sentinel: Generate 50% addl rage from attacking targets not attacking the warrior
			}

			return rage
		},
	}
	if mh := war.GetMHWeapon(); mh != nil {
		rbo.MHSwingSpeed = mh.SwingSpeed
	}

	war.EnableRageBar(rbo)
	war.EnableAutoAttacks(war, core.AutoAttackOptions{
		MainHand:       war.WeaponFromMainHand(war.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
	})

	//healingModel := options.HealingModel
	//if healingModel != nil {
	//	if healingModel.InspirationUptime > 0.0 {
	//		core.ApplyInspiration(war.GetCharacter(), healingModel.InspirationUptime)
	//	}
	//}

	return war
}

func (war *ProtectionWarrior) RegisterSpecializationEffects() {
	// Critical block
	war.RegisterMastery()

	// Sentinel stat buffs
	war.MultiplyStat(stats.Stamina, 1.15)
	war.MultiplyStat(stats.Block, 1.15)

	// Vengeance
	core.ApplyVengeanceEffect(war.GetCharacter(), &war.VengeanceTracker, 93098)
}

func (war *ProtectionWarrior) RegisterMastery() {
	dummyCriticalBlockSpell := war.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 76857}, // Doesn't seem like there's an actual spell ID for the block itself, so use the mastery ID
		Flags:    core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,
	})

	// Seems to work pretty much the same as WotLK critical block
	war.AddDynamicDamageTakenModifier(func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		if result.Outcome.Matches(core.OutcomeBlock) && !result.Outcome.Matches(core.OutcomeMiss) && !result.Outcome.Matches(core.OutcomeParry) && !result.Outcome.Matches(core.OutcomeDodge) {
			procChance := war.CriticalBlockChance // Use the member and not GetCriticalBlockChance as Hold the Line may have been applied from the baseline warrior impl
			if sim.Proc(procChance, "Critical Block Roll") {
				blockValue := war.BlockValue()
				result.Damage = max(0, result.Damage-blockValue)
				dummyCriticalBlockSpell.Cast(sim, spell.Unit)
			}
		}
	})

	// Crit block mastery also applies an equal amount to regular block
	// set initial block rating from stats
	war.CriticalBlockChance = war.GetCriticalBlockChance()
	war.AddStat(stats.Block, (war.CriticalBlockChance*100.0)*core.BlockRatingPerBlockChance)

	// and keep it updated when mastery changes
	war.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery, newMastery float64) {
		oldBlockRating := (1.5 * core.MasteryRatingToMasteryPoints(oldMastery)) * core.BlockRatingPerBlockChance
		newBlockRating := (1.5 * core.MasteryRatingToMasteryPoints(newMastery)) * core.BlockRatingPerBlockChance

		war.AddStatDynamic(sim, stats.Block, -oldBlockRating+newBlockRating)
		war.CriticalBlockChance = war.GetCriticalBlockChance()
	})
}

func CalcMasteryPercent(points float64) float64 {
	return 12.0 + 1.5*points
}

func (war *ProtectionWarrior) GetCriticalBlockChance() float64 {
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

func (war *ProtectionWarrior) Reset(sim *core.Simulation) {
	war.Warrior.Reset(sim)
}
