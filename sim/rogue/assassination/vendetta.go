package assassination

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/rogue"
)

func (sinRogue *AssassinationRogue) registerVendetta() {
	actionID := core.ActionID{SpellID: 79140}
	hasGlyph := sinRogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfVendetta)
	duration := time.Second * time.Duration(core.TernaryFloat64(hasGlyph, 30, 20))

	vendettaAuras := sinRogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Vendetta",
			ActionID: actionID,
			Duration: duration,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				core.EnableDamageDoneByCaster(DDBC_Vendetta, DDBC_Total, sinRogue.AttackTables[aura.Unit.UnitIndex], func(sim *core.Simulation, spell *core.Spell, attackTable *core.AttackTable) float64 {
					if spell.Matches(rogue.RogueSpellsAll) || spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
						return core.TernaryFloat64(hasGlyph, 1.3, 1.25)
					}
					return 1.0
				})
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				core.DisableDamageDoneByCaster(DDBC_Vendetta, sinRogue.AttackTables[aura.Unit.UnitIndex])
			},
		})
	})

	sinRogue.Vendetta = sinRogue.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		Flags:          core.SpellFlagAPL | core.SpellFlagMCD | core.SpellFlagReadinessTrinket,
		ClassSpellMask: rogue.RogueSpellVendetta,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    sinRogue.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			aura := vendettaAuras.Get(target)
			aura.Activate(sim)
		},
		RelatedAuraArrays: vendettaAuras.ToMap(),
	})

	sinRogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    sinRogue.Vendetta,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityDefault,
	})
}
