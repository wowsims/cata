package shaman

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

// Shared logic for all shocks.
func (shaman *Shaman) newShockSpellConfig(spellID int32, spellSchool core.SpellSchool, baseCost float64, shockTimer *core.Timer) core.SpellConfig {
	actionID := core.ActionID{SpellID: spellID}

	return core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: spellSchool,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagShock | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: baseCost,
			Multiplier: 1 -
				0.05*float64(shaman.Talents.Convection),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shockTimer,
				Duration: time.Second*6 - time.Millisecond*500*time.Duration(shaman.Talents.Reverberation),
			},
		},

		BonusHitRating:   float64(shaman.Talents.ElementalPrecision) * core.SpellHitRatingPerHitChance,
		DamageMultiplier: 1 + 0.02*float64(shaman.Talents.Concussion),
		CritMultiplier:   shaman.ElementalFuryCritMultiplier(0),
	}
}

func (shaman *Shaman) registerEarthShockSpell(shockTimer *core.Timer) {
	config := shaman.newShockSpellConfig(8042, core.SpellSchoolNature, 0.18, shockTimer)
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := 931 + 0.386*spell.SpellPower()
		spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
	}

	shaman.EarthShock = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerFlameShockSpell(shockTimer *core.Timer) {
	config := shaman.newShockSpellConfig(8050, core.SpellSchoolFire, 0.17, shockTimer)

	flameShockBaseNumberOfTicks := int32(6)

	if shaman.HasPrimeGlyph(proto.ShamanPrimeGlyph_GlyphOfFlameShock) {
		flameShockBaseNumberOfTicks += 3
	}

	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := 531 + 0.214*spell.SpellPower()
		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		if result.Landed() {
			spell.Dot(target).NumberOfTicks = flameShockBaseNumberOfTicks
			spell.Dot(target).Apply(sim)
		}
		spell.DealDamage(sim, result)
	}

	bonusPeriodicDamageMultiplier := 0 + 0.2*float64(shaman.Talents.LavaFlows)

	config.Dot = core.DotConfig{
		Aura: core.Aura{
			Label: "FlameShock",
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				shaman.LavaBurst.BonusCritRating += 100 * core.CritRatingPerCritChance
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				shaman.LavaBurst.BonusCritRating -= 100 * core.CritRatingPerCritChance
			},
		},
		NumberOfTicks:       flameShockBaseNumberOfTicks,
		TickLength:          time.Second * 3,
		AffectedByCastSpeed: true,

		//TODO: Snapshot?
		//TODO: I don't know what 834/6 is
		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
			dot.SnapshotBaseDamage = 834/6 + 0.1*dot.Spell.SpellPower()
			dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)

			dot.Spell.DamageMultiplierAdditive += bonusPeriodicDamageMultiplier
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			dot.Spell.DamageMultiplierAdditive -= bonusPeriodicDamageMultiplier
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)

			if shaman.Talents.LavaSurge > 0 {
				if sim.RandomFloat("LavaSurge") < (0.1 * float64(shaman.Talents.LavaSurge)) {
					shaman.LavaBurst.CD.Reset()
				}
			}
		},
	}

	shaman.FlameShock = shaman.RegisterSpell(config)
}

// TODO: need base damage
func (shaman *Shaman) registerFrostShockSpell(shockTimer *core.Timer) {
	config := shaman.newShockSpellConfig(8056, core.SpellSchoolFrost, 0.18, shockTimer)
	config.ThreatMultiplier *= 2
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := sim.Roll(812, 858) + 0.386*spell.SpellPower()
		spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
	}

	shaman.FrostShock = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerShocks() {
	shockTimer := shaman.NewTimer()
	shaman.registerEarthShockSpell(shockTimer)
	shaman.registerFlameShockSpell(shockTimer)
	shaman.registerFrostShockSpell(shockTimer)
}
