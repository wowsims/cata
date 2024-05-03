package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (warlock *Warlock) registerSummonDoomguardSpell(timer *core.Timer) {
	duration := time.Second * time.Duration(45+10*warlock.Talents.AncientGrimoire)

	summonDoomguardAura := warlock.RegisterAura(core.Aura{
		Label:    "Summon Doomguard",
		ActionID: core.ActionID{SpellID: 18540},
		Duration: duration,
	})

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 18540},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellSummonDoomguard,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.8,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 0,
			},
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Minute * 10,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			warlock.Doomguard.EnableWithTimeout(sim, warlock.Doomguard, duration)
			summonDoomguardAura.Activate(sim)
		},
	})
}

type DoomguardPet struct {
	core.Pet

	owner *Warlock

	DoomBolt *core.Spell

	CurseOfGuldanDebuffs core.AuraArray
}

func (warlock *Warlock) NewDoomguardPet() *DoomguardPet {
	stats := stats.Stats{
		stats.Strength:  297,
		stats.Agility:   79,
		stats.Stamina:   118,
		stats.Intellect: 369,
		stats.Spirit:    367,
		stats.Mana:      1174,
		stats.MP5:       270, // rough guess, unclear if it's affected by other stats
		stats.MeleeCrit: 3.454 * core.CritRatingPerCritChance,
		stats.SpellCrit: 0.9075 * core.CritRatingPerCritChance,
	}

	wp := &DoomguardPet{
		Pet:   core.NewPet("Doomguard", &warlock.Character, stats, warlock.MakeStatInheritance(), false, true),
		owner: warlock,
	}

	//TODO: Power Modifier
	wp.EnableManaBarWithModifier(0.33)

	//wp.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	//wp.AddStat(stats.AttackPower, -20)

	warlock.AddPet(wp)

	return wp
}

func (doomguard *DoomguardPet) GetPet() *core.Pet {
	return &doomguard.Pet
}

func (doomguard *DoomguardPet) Initialize() {
	doomguard.registerDoomBolt()
}

func (doomguard *DoomguardPet) Reset(_ *core.Simulation) {
}

func (doomguard *DoomguardPet) ExecuteCustomRotation(sim *core.Simulation) {
	doomguard.DoomBolt.Cast(sim, doomguard.CurrentTarget)
}

func (doomguard *DoomguardPet) registerDoomBolt() {
	doomguard.DoomBolt = doomguard.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 85692},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: WarlockSpellDoomguardDoomBolt,
		//TODO: Same as shadowbolt? How to get this
		MissileSpeed: 20,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
				//TODO: I think it's 3 seconds
				CastTime: time.Millisecond * 3000,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           2,
		ThreatMultiplier:         1,
		BonusCoefficient:         1.2855,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, 1017, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
