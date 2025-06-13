package guardian

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/druid"
)

func (bear *GuardianDruid) registerIncarnation() {
	if !bear.Talents.Incarnation {
		return
	}

	actionID := core.ActionID{SpellID: 102558}

	var affectedSpells []*druid.DruidSpell
	var cdReductions []time.Duration

	bear.SonOfUrsocAura = bear.RegisterAura(core.Aura{
		Label:    "Incarnation: Son of Ursoc",
		ActionID: actionID,
		Duration: time.Second * 30,

		OnInit: func(_ *core.Aura, _ *core.Simulation) {
			affectedSpells = []*druid.DruidSpell{bear.SwipeBear, bear.Lacerate, bear.MangleBear, bear.ThrashBear, bear.Maul}
			cdReductions = make([]time.Duration, len(affectedSpells))
		},

		OnGain: func(_ *core.Aura, _ *core.Simulation) {
			for idx, spell := range affectedSpells {
				cdReductions[idx] = spell.CD.Duration - core.GCDDefault
				spell.CD.Duration -= cdReductions[idx]
				spell.CD.Reset()
			}
		},

		OnExpire: func(_ *core.Aura, _ *core.Simulation) {
			for idx, spell := range affectedSpells {
				spell.CD.Duration += cdReductions[idx]
			}
		},
	})

	bear.SonOfUrsoc = bear.RegisterSpell(druid.Any, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},

			CD: core.Cooldown{
				Timer:    bear.NewTimer(),
				Duration: time.Minute * 3,
			},

			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if !bear.InForm(druid.Bear) {
				bear.BearFormAura.Activate(sim)
			}

			bear.SonOfUrsocAura.Activate(sim)
		},
	})

	bear.AddMajorCooldown(core.MajorCooldown{
		Spell: bear.SonOfUrsoc.Spell,
		Type:  core.CooldownTypeDPS,

		ShouldActivate: func(sim *core.Simulation, _ *core.Character) bool {
			return !bear.BerserkBearAura.IsActive() && !bear.Berserk.IsReady(sim)
		},
	})
}

func (bear *GuardianDruid) registerForceOfNature() {
	if !bear.Talents.ForceOfNature {
		return
	}

	bear.ForceOfNature = bear.RegisterSpell(druid.Any, core.SpellConfig{
		ActionID: core.ActionID{SpellID: 102706},
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    bear.NewTimer(),
				Duration: time.Second * 20,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			bear.Treants.Enable(sim)
		},
	})

	bear.AddMajorCooldown(core.MajorCooldown{
		Spell: bear.ForceOfNature.Spell,
		Type:  core.CooldownTypeDPS,
	})
}
