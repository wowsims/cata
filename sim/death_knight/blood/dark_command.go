package blood

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

var DarkCommandActionID = core.ActionID{SpellID: 56222}

// Commands the target to attack you, and increases threat that you generate against the target by 200% for 3 sec.
func (bdk *BloodDeathKnight) registerDarkCommand() {
	tdbcHandler := func(sim *core.Simulation, spell *core.Spell, attackTable *core.AttackTable) float64 {
		return 2.0
	}

	darkCommandAuras := bdk.NewEnemyAuraArray(func(unit *core.Unit) *core.Aura {
		return unit.RegisterAura(core.Aura{
			Label:    "Dark Command" + unit.Label,
			ActionID: DarkCommandActionID,
			Duration: time.Second * 3,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				core.EnableThreatDoneByCaster(TDBC_DarkCommand, TDBC_Total, bdk.AttackTables[aura.Unit.UnitIndex], tdbcHandler)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				core.DisableThreatDoneByCaster(TDBC_DarkCommand, bdk.AttackTables[aura.Unit.UnitIndex])
			},
		})
	})

	bdk.RegisterSpell(core.SpellConfig{
		ActionID:       DarkCommandActionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: death_knight.DeathKnightSpellDarkCommand,

		MaxRange: 30,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    bdk.NewTimer(),
				Duration: time.Second * 8,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeAlwaysHit)

			darkCommandAuras.Get(target).Activate(sim)

			spell.DealOutcome(sim, result)
		},
	})
}
