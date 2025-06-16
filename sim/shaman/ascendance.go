package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (shaman *Shaman) registerAscendanceSpell() {

	var originalMHSpell *core.Spell
	var originalOHSpell *core.Spell

	var isEnh = shaman.Spec == proto.Spec_SpecEnhancementShaman

	windslashMH := shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 114089, Tag: 1},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskMeleeMHAuto,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete | core.SpellFlagReadinessTrinket,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           shaman.DefaultCritMultiplier(),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		},

		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := spell.Unit.AutoAttacks.MH().CalculateAverageWeaponDamage(spell.MeleeAttackPower())
			return spell.CalcDamage(sim, target, baseDamage, spell.OutcomeExpectedMeleeWeaponSpecialHitAndCrit)
		},
	})

	windslashOH := shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 114089, Tag: 2},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskMeleeOHAuto,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           shaman.DefaultCritMultiplier(),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.OHWeaponDamage(sim, spell.MeleeAttackPower())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
		},
	})

	shaman.AscendanceAura = shaman.GetOrRegisterAura(core.Aura{
		Label:    "Ascendance",
		ActionID: core.ActionID{SpellID: 114049},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if isEnh {
				//TODO weapon swap during ascendance breaks this i think
				originalMHSpell = shaman.AutoAttacks.MHAuto()
				originalOHSpell = shaman.AutoAttacks.OHAuto()
				shaman.AutoAttacks.SetMHSpell(windslashMH)
				shaman.AutoAttacks.SetOHSpell(windslashOH)
			}
			pa := &core.PendingAction{
				NextActionAt: aura.ExpiresAt(),
				Priority:     core.ActionPriorityGCD,
				OnAction:     func(sim *core.Simulation) {},
			}
			sim.AddPendingAction(pa)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			//Lava Beam cast gets cancelled if ascendance fades during it
			if (shaman.Hardcast.ActionID.SpellID == 114074) && shaman.Hardcast.Expires > sim.CurrentTime {
				shaman.CancelHardcast(sim)
			}
			if isEnh {
				shaman.Stormstrike.CD.Set(shaman.Stormblast.CD.ReadyAt())
				shaman.AutoAttacks.SetMHSpell(originalMHSpell)
				shaman.AutoAttacks.SetOHSpell(originalOHSpell)
			}
		},
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask:  SpellMaskLavaBurst,
		Kind:       core.SpellMod_Cooldown_Multiplier,
		FloatValue: -1,
	})

	shaman.Ascendance = shaman.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 114049},
		SpellSchool:    core.SpellSchoolPhysical,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskAscendance,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 5.2,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Minute * 3,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			shaman.AscendanceAura.Activate(sim)
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: shaman.Ascendance,
		Type:  core.CooldownTypeDPS,
	})
}
