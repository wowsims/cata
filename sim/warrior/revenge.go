package warrior

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (warrior *Warrior) RegisterRevengeSpell() {
	actionID := core.ActionID{SpellID: 6572}

	warrior.revengeProcAura = warrior.RegisterAura(core.Aura{
		Label:    "Revenge Ready",
		Duration: 5 * time.Second,
		ActionID: actionID,
	})

	core.MakePermanent(warrior.RegisterAura(core.Aura{
		Label: "Revenge Trigger",
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Outcome.Matches(core.OutcomeBlock | core.OutcomeDodge | core.OutcomeParry) {
				warrior.revengeProcAura.Activate(sim)
			}
		},
	}))

	extraHit := warrior.Talents.ImprovedRevenge > 0 && warrior.Env.GetNumTargets() > 1

	warrior.Revenge = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   5,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 5,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(DefensiveStance) && warrior.revengeProcAura.IsActive()
		},

		DamageMultiplier: (1.0 + 0.3*float64(warrior.Talents.ImprovedRevenge)) * core.TernaryFloat64(warrior.HasPrimeGlyph(proto.WarriorPrimeGlyph_GlyphOfRevenge), 1.1, 1.0),
		ThreatMultiplier: 1,
		FlatThreatBonus:  121,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// TODO: check this roll range and ap coefficient, this is from the 4.3.3 simc export
			ap := spell.MeleeAttackPower() * 0.31
			baseDamage := sim.Roll(1618.3, 1977.92) + ap
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			if !result.Landed() {
				spell.IssueRefund(sim)
			}

			if extraHit {
				if sim.RandomFloat("Revenge Target Roll") <= 0.5*float64(warrior.Talents.ImprovedRevenge) {
					otherTarget := sim.Environment.NextTargetUnit(target)
					baseDamage := sim.Roll(1618.3, 1977.92) + ap
					spell.CalcAndDealDamage(sim, otherTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				}
			}

			warrior.revengeProcAura.Deactivate(sim)
		},
	})
}
