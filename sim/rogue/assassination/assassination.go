package assassination

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
	"github.com/wowsims/cata/sim/rogue"
)

func RegisterAssassinationRogue() {
	core.RegisterAgentFactory(
		proto.Player_AssassinationRogue{},
		proto.Spec_SpecAssassinationRogue,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewAssassinationRogue(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_AssassinationRogue)
			if !ok {
				panic("Invalid spec value for Assassination Rogue!")
			}
			player.Spec = playerSpec
		},
	)
}

func (sinRogue *AssassinationRogue) Initialize() {
	sinRogue.Rogue.Initialize()

	sinRogue.registerMutilateSpell()
	sinRogue.registerOverkill()
	sinRogue.registerColdBloodCD()
	sinRogue.applySealFate()
	sinRogue.registerVenomousWounds()

	sinRogue.applyMastery()
}

func NewAssassinationRogue(character *core.Character, options *proto.Player) *AssassinationRogue {
	sinOptions := options.GetAssassinationRogue().Options

	sinRogue := &AssassinationRogue{
		Rogue: rogue.NewRogue(character, sinOptions.ClassOptions, options.TalentsString),
	}
	sinRogue.AssassinationOptions = sinOptions

	return sinRogue
}

type AssassinationRogue struct {
	*rogue.Rogue

	masteryAura *core.Aura
}

func (sinRogue *AssassinationRogue) GetRogue() *rogue.Rogue {
	return sinRogue.Rogue
}

func (sinRogue *AssassinationRogue) Reset(sim *core.Simulation) {
	sinRogue.Rogue.Reset(sim)
}

var OverkillActionID = core.ActionID{SpellID: 58427}

func (sinRogue *AssassinationRogue) registerOverkill() {
	if !sinRogue.Talents.Overkill {
		return
	}

	effectDuration := time.Second * 20
	if sinRogue.StealthAura.IsActive() {
		effectDuration = core.NeverExpires
	}

	sinRogue.OverkillAura = sinRogue.RegisterAura(core.Aura{
		Label:    "Overkill",
		ActionID: OverkillActionID,
		Duration: effectDuration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			sinRogue.ApplyEnergyTickMultiplier(0.3)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			sinRogue.ApplyEnergyTickMultiplier(-0.3)
		},
	})
}

func (sinRogue *AssassinationRogue) registerColdBloodCD() {
	if !sinRogue.Talents.ColdBlood {
		return
	}

	actionID := core.ActionID{SpellID: 14177}

	coldBloodAura := sinRogue.RegisterAura(core.Aura{
		Label:    "Cold Blood",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range sinRogue.Spellbook {
				if spell.Flags.Matches(rogue.SpellFlagColdBlooded) {
					spell.BonusCritRating += 100 * core.CritRatingPerCritChance
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range sinRogue.Spellbook {
				if spell.Flags.Matches(rogue.SpellFlagColdBlooded) {
					spell.BonusCritRating -= 100 * core.CritRatingPerCritChance
				}
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// for Fan of Knives and Mutilate, the offhand hit comes first and is ignored, so the aura doesn't fade too early
			if spell.Flags.Matches(rogue.SpellFlagColdBlooded) && spell.ProcMask.Matches(core.ProcMaskMeleeMH) {
				aura.Deactivate(sim)
			}
		},
	})

	sinRogue.ColdBlood = sinRogue.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    sinRogue.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			coldBloodAura.Activate(sim)
		},
	})

	sinRogue.AddMajorCooldown(core.MajorCooldown{
		Spell: sinRogue.ColdBlood,
		Type:  core.CooldownTypeDPS,
	})
}

func (sinRogue *AssassinationRogue) applySealFate() {
	if sinRogue.Talents.SealFate == 0 {
		return
	}

	procChance := 0.2 * float64(sinRogue.Talents.SealFate)
	cpMetrics := sinRogue.NewComboPointMetrics(core.ActionID{SpellID: 14190})

	icd := core.Cooldown{
		Timer:    sinRogue.NewTimer(),
		Duration: 500 * time.Millisecond,
	}

	sinRogue.RegisterAura(core.Aura{
		Label:    "Seal Fate",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Flags.Matches(rogue.SpellFlagBuilder) {
				return
			}

			if !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}

			if icd.IsReady(sim) && sim.Proc(procChance, "Seal Fate") {
				sinRogue.AddComboPoints(sim, 1, cpMetrics)
				icd.Use(sim)
			}
		},
	})
}

func (sinRogue *AssassinationRogue) registerVenomousWounds() {
	if sinRogue.Talents.VenomousWounds == 0 {
		return
	}

	vwSpellID := 79132 + sinRogue.Talents.VenomousWounds
	vwActionID := core.ActionID{SpellID: vwSpellID}

	// https://web.archive.org/web/20111128070437/http://elitistjerks.com/f78/t105429-cataclysm_mechanics_testing/  Ctrl-F "Venomous Wounds"
	vwBaseTickDamage := 675.0
	vwMetrics := sinRogue.NewEnergyMetrics(vwActionID)

	sinRogue.VenomousWounds = sinRogue.RegisterSpell(core.SpellConfig{
		ActionID:         vwActionID,
		SpellSchool:      core.SpellSchoolNature,
		ProcMask:         core.ProcMaskSpellDamage,
		CritMultiplier:   sinRogue.MeleeCritMultiplier(false),
		ThreatMultiplier: 1,
		DamageMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			vwDamage := vwBaseTickDamage + 0.176*spell.MeleeAttackPower()
			spell.CalcAndDealDamage(sim, target, vwDamage, spell.OutcomeAlwaysHit)
			sinRogue.AddEnergy(sim, 10, vwMetrics)
		},
	})
}

func (sinRogue *AssassinationRogue) applyMastery() {
	const damagePerPercent = .035
	const baseEffect = .28
	sinRogue.masteryAura = sinRogue.RegisterAura(core.Aura{
		Label:    "Mastery: Potent Poisons",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			masteryPercent := sinRogue.GetStat(stats.Mastery) / core.MasteryRatingPerMasteryPercent
			masteryEffect := baseEffect + masteryPercent*damagePerPercent
			for _, spell := range sinRogue.InstantPoison {
				spell.DamageMultiplier += masteryEffect
			}
			for _, spell := range sinRogue.WoundPoison {
				spell.DamageMultiplier += masteryEffect
			}
			sinRogue.DeadlyPoison.DamageMultiplier += masteryEffect
			sinRogue.Envenom.DamageMultiplier += masteryEffect
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			masteryPercent := sinRogue.GetStat(stats.Mastery) / core.MasteryRatingPerMasteryPercent
			masteryEffect := baseEffect + masteryPercent*damagePerPercent
			for _, spell := range sinRogue.InstantPoison {
				spell.DamageMultiplier -= masteryEffect
			}
			for _, spell := range sinRogue.WoundPoison {
				spell.DamageMultiplier -= masteryEffect
			}
			sinRogue.DeadlyPoison.DamageMultiplier -= masteryEffect
			sinRogue.Envenom.DamageMultiplier -= masteryEffect
		},
	})
}
