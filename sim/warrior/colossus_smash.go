package warrior

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (warrior *Warrior) RegisterColossusSmash() {
	actionID := core.ActionID{SpellID: 86346}
	aura := warrior.RegisterAura(core.Aura{
		Label:    "Colossus Smash",
		ActionID: actionID,
		Duration: time.Second * 6,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.AttackTables[aura.Unit.UnitIndex].IgnoreArmor = true
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.AttackTables[aura.Unit.UnitIndex].IgnoreArmor = false
		},
	})

	warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagAPL | core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,
		RageCost: core.RageCostOptions{
			Cost:   20,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 20,
			},
		},
		CritMultiplier: warrior.DefaultMeleeCritMultiplier(),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			aura.Activate(sim)
			if warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfColossusSmash) {
				warrior.TryApplySunderArmorEffect(sim, target)
			}
		},
	})
}
