package blood

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/death_knight"
)

func RegisterBloodDeathKnight() {
	core.RegisterAgentFactory(
		proto.Player_BloodDeathKnight{},
		proto.Spec_SpecBloodDeathKnight,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewBloodDeathKnight(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_BloodDeathKnight)
			if !ok {
				panic("Invalid spec value for Blood Death Knight!")
			}
			player.Spec = playerSpec
		},
	)
}

type BloodDeathKnight struct {
	*death_knight.DeathKnight

	vengeance *core.VengeanceTracker
}

func NewBloodDeathKnight(character *core.Character, options *proto.Player) *BloodDeathKnight {
	dkOptions := options.GetBloodDeathKnight()

	bdk := &BloodDeathKnight{
		DeathKnight: death_knight.NewDeathKnight(character, death_knight.DeathKnightInputs{
			IsDps:              false,
			StartingRunicPower: dkOptions.Options.ClassOptions.StartingRunicPower,
			Spec:               proto.Spec_SpecBloodDeathKnight,
		}, options.TalentsString, 50034),
		vengeance: &core.VengeanceTracker{},
	}

	return bdk
}

func (bdk *BloodDeathKnight) GetDeathKnight() *death_knight.DeathKnight {
	return bdk.DeathKnight
}

func (bdk *BloodDeathKnight) Initialize() {
	bdk.DeathKnight.Initialize()

	// bdk.registerHeartStrikeSpell()

	// TODO: Fix this to work with the new talent system.
	// if bdk.Talents.DancingRuneWeapon {
	// 	bdk.RuneWeapon.AddCopySpell(HeartStrikeActionID, bdk.registerDrwHeartStrikeSpell())
	// }
}

func (bdk BloodDeathKnight) getBloodShieldMasteryBonus() float64 {
	return 0.5 + 0.0625*bdk.GetMasteryPoints()
}

func (bdk *BloodDeathKnight) ApplyTalents() {
	bdk.DeathKnight.ApplyTalents()
	bdk.ApplyArmorSpecializationEffect(stats.Stamina, proto.ArmorType_ArmorTypePlate, 86524)

	// Veteran of the Third War
	bdk.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		ClassMask: death_knight.DeathKnightSpellOutbreak,
		TimeValue: time.Second * -30,
	})
	bdk.MultiplyStat(stats.Stamina, 1.09)
	bdk.AddStat(stats.ExpertiseRating, 6*core.ExpertisePerQuarterPercentReduction)
	core.MakePermanent(bdk.GetOrRegisterAura(core.Aura{
		Label:    "Veteran of the Third War",
		ActionID: core.ActionID{SpellID: 50029},
	}))

	// Vengeance
	core.ApplyVengeanceEffect(&bdk.Character, bdk.vengeance, 93099)

	// Mastery: Blood Shield
	shieldAmount := 0.0
	currentShield := 0.0
	shieldSpell := bdk.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 77535},
		ProcMask:    core.ProcMaskSpellHealing,
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Shield: core.ShieldConfig{
			SelfOnly: true,
			Aura: core.Aura{
				Label:    "Blood Shield",
				Duration: core.NeverExpires,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if currentShield < bdk.MaxHealth() {
				shieldAmount = min(shieldAmount, bdk.MaxHealth()-currentShield)
				currentShield += shieldAmount
				spell.SelfShield().Apply(sim, shieldAmount)
			}
		},
	})
	core.MakePermanent(bdk.GetOrRegisterAura(core.Aura{
		Label:    "Mastery: Blood Shield",
		ActionID: core.ActionID{SpellID: 77513},
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			shieldAmount = 0.0
			currentShield = 0.0
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.SpellSchool.Matches(core.SpellSchoolPhysical) {
				return
			}

			if currentShield <= 0 || result.Damage <= 0 {
				return
			}

			damageReduced := min(result.Damage, currentShield)
			currentShield -= damageReduced

			bdk.GainHealth(sim, damageReduced, shieldSpell.HealthMetrics(result.Target))

			if currentShield <= 0 {
				shieldSpell.SelfShield().Deactivate(sim)
			}
		},
		OnHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ClassSpellMask&death_knight.DeathKnightSpellDeathStrikeHeal == 0 {
				return
			}

			shieldAmount = result.Damage * bdk.getBloodShieldMasteryBonus()
			shieldSpell.Cast(sim, result.Target)
		},
	}))

}

func (bdk *BloodDeathKnight) Reset(sim *core.Simulation) {
	bdk.DeathKnight.Reset(sim)
}
