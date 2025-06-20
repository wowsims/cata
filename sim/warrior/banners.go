package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (war *Warrior) registerBanners() {
	war.registerSkullBanner()
	war.registerDemoralizingBanner()
}

func (war *Warrior) registerSkullBanner() {
	war.SkullBannerAura = core.SkullBannerAura(&war.Character, war.Index)

	spell := war.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       core.SkullBannerActionID,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskSkullBanner,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: core.SkullBannerCD,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			war.SkullBannerAura.Activate(sim)
		},
	})

	war.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

func (war *Warrior) registerDemoralizingBanner() {
	actionID := core.ActionID{SpellID: 114030}

	war.DemoralizingBannerAuras = war.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return war.GetOrRegisterAura(core.Aura{
			Label:    "Demoralizing Banner",
			ActionID: actionID,
			Duration: 15 * time.Second,
		}).AttachMultiplicativePseudoStatBuff(&target.PseudoStats.DamageDealtMultiplier, 0.9)
	})

	spell := war.RegisterSpell(core.SpellConfig{
		ActionID:       core.SkullBannerActionID,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskDemoralizingBanner,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			for _, target := range sim.Encounter.TargetUnits {
				war.DemoralizingBannerAuras.Get(target).Activate(sim)
			}
		},
	})

	war.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeSurvival,
	})
}
