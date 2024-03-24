package arms

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (war *ArmsWarrior) RegisterBladestorm() {
	if !war.Talents.Bladestorm {
		return
	}
	actionID := core.ActionID{SpellID: 46924}
	numHits := war.Env.GetNumTargets() // 1 hit per target
	results := make([]*core.SpellResult, numHits)

	war.Bladestorm = war.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagChanneled | core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,
		RageCost: core.RageCostOptions{
			Cost: 25,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: core.TernaryDuration(war.HasPrimeGlyph(proto.WarriorPrimeGlyph_GlyphOfBladestorm), time.Second*75, time.Second*90),
			},
			IgnoreHaste: true,
		},
		CritMultiplier:   1,
		DamageMultiplier: 1.5, // deals 150% weapon damage
		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Bladestorm",
			},
			NumberOfTicks: 6,
			TickLength:    time.Second * 1,
			OnTick: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot) {
				target := war.CurrentTarget
				spell := dot.Spell
				curTarget := target
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					baseDamage := 0 +
						spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
						spell.BonusWeaponDamage()
					results[hitIndex] = spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}
				curTarget = target
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					spell.DealDamage(sim, results[hitIndex])
					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			dot := spell.AOEDot()
			dot.Apply(sim)
			dot.TickOnce(sim)
		},
	})
	war.AddMajorCooldown(core.MajorCooldown{
		Spell: war.Bladestorm,
		Type:  core.CooldownTypeDPS,
	})
}
