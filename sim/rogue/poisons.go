package rogue

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (rogue *Rogue) applyPoisons() {
	rogue.applyDeadlyPoison()
	rogue.applyWoundPoison()
}

func (rogue *Rogue) registerPoisonAuras() {
	rogue.MasterPoisonerDebuffAuras = rogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		aura := core.MasterPoisonerDebuff(target)
		aura.Duration = core.NeverExpires
		return aura
	})
}

func (rogue *Rogue) getPoisonProcMaskForSlot(imbue proto.RogueOptions_PoisonImbue, itemSlot proto.ItemSlot) core.ProcMask {
	mask := core.ProcMaskUnknown
	switch {
	case itemSlot == proto.ItemSlot_ItemSlotMainHand && rogue.Options.MhImbue == imbue:
		mask |= core.ProcMaskMeleeMH | core.ProcMaskMeleeProc
	case itemSlot == proto.ItemSlot_ItemSlotOffHand && rogue.Options.OhImbue == imbue:
		mask |= core.ProcMaskMeleeOH
	case itemSlot == proto.ItemSlot_ItemSlotRanged && rogue.Options.ThImbue == imbue:
		mask |= core.ProcMaskRanged
	}

	return mask
}

func (rogue *Rogue) GetLethalPoisonProcChance() float64 {
	return 0.3 + core.TernaryFloat64(rogue.Spec == proto.Spec_SpecAssassinationRogue, 0.2, 0)
}

func (rogue *Rogue) UpdateDeadlyPoisonPPH(bonusChance float64) {
	pph := rogue.GetLethalPoisonProcChance() + bonusChance
	for _, itemSlot := range core.AllWeaponSlots() {
		procMask := rogue.getPoisonProcMaskForSlot(proto.RogueOptions_DeadlyPoison, itemSlot)
		rogue.deadlyPoisonPPHM[itemSlot] = rogue.AutoAttacks.NewStaticDynamicProcManager(pph, procMask)
	}
}

func (rogue *Rogue) UpdateWoundPoisonPPH(bonusChance float64) {
	pph := rogue.GetLethalPoisonProcChance() + bonusChance
	for _, itemSlot := range core.AllWeaponSlots() {
		procMask := rogue.getPoisonProcMaskForSlot(proto.RogueOptions_WoundPoison, itemSlot)
		rogue.deadlyPoisonPPHM[itemSlot] = rogue.AutoAttacks.NewStaticDynamicProcManager(pph, procMask)
	}
}

func (rogue *Rogue) registerDeadlyPoisonSpell() {
	dot_baseDamage := 0.11999999732 * rogue.ClassSpellScaling
	dot_apScaling := 0.03500000015

	hit_baseDamage := 263.0
	hit_apScaling := 0.109

	rogue.DeadlyPoison = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 2818},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamageProc,
		ClassSpellMask: RogueSpellDeadlyPoison,
		Flags:          core.SpellFlagPassiveSpell,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           1,
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Deadly Poison DoT",
				MaxStacks: 1,
				Duration:  time.Second * 12,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					rogue.MasterPoisonerDebuffAuras.Get(aura.Unit).Activate(sim)

				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					rogue.MasterPoisonerDebuffAuras.Get(aura.Unit).Deactivate(sim)

				},
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 3,

			OnSnapshot: func(_ *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = (dot_baseDamage + dot_apScaling*dot.Spell.MeleeAttackPower())
				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				dot.SnapshotCritChance = dot.Spell.SpellCritChance(attackTable.Defender)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable, true)
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHitAndCrit)
			if !result.Landed() {
				return
			}

			dot := spell.Dot(target)
			if dot.IsActive() {
				baseDamage := hit_baseDamage + hit_apScaling*spell.MeleeAttackPower()
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHit)
			}

			dot.Apply(sim)
			dot.Refresh(sim)
			dot.TakeSnapshot(sim, false)
		},
	})
}

func (rogue *Rogue) registerWoundPoisonSpell() {
	woundPoisonDebuffAura := core.Aura{
		Label:    "Wound Poison",
		ActionID: core.ActionID{SpellID: 13219},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			rogue.MasterPoisonerDebuffAuras.Get(aura.Unit).Activate(sim)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			rogue.MasterPoisonerDebuffAuras.Get(aura.Unit).Deactivate(sim)
		},
	}

	rogue.WoundPoisonDebuffAuras = rogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.RegisterAura(woundPoisonDebuffAura)
	})

	wpBaseDamage := 0.24500000477 * rogue.ClassSpellScaling
	rogue.WoundPoison = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 13218},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamageProc,
		ClassSpellMask: RogueSpellWoundPoison,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           rogue.DefaultCritMultiplier(),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := wpBaseDamage + 0.04*spell.MeleeAttackPower()

			var result *core.SpellResult
			result = spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				rogue.WoundPoisonDebuffAuras.Get(target).Activate(sim)
			}
		},
	})
}

func (rogue *Rogue) applyDeadlyPoison() {
	rogue.UpdateDeadlyPoisonPPH(0)

	for _, itemSlot := range core.AllWeaponSlots() {
		core.MakeProcTriggerAura(&rogue.Unit, core.ProcTrigger{
			Name:     fmt.Sprintf("Deadly Poison %s", itemSlot),
			Outcome:  core.OutcomeLanded,
			Callback: core.CallbackOnSpellHitDealt,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procMask := core.Ternary(spell.SpellID == 86392, core.ProcMaskMeleeMH, spell.ProcMask)
				if rogue.deadlyPoisonPPHM[itemSlot].Proc(sim, procMask, "Deadly Poison") {
					rogue.DeadlyPoison.Cast(sim, result.Target)
				}
			},
		})
	}

}

func (rogue *Rogue) applyWoundPoison() {
	rogue.UpdateWoundPoisonPPH(0)

	for _, itemSlot := range core.AllWeaponSlots() {
		core.MakeProcTriggerAura(&rogue.Unit, core.ProcTrigger{
			Name:     fmt.Sprintf("Wound Poison %s", itemSlot),
			Outcome:  core.OutcomeLanded,
			Callback: core.CallbackOnSpellHitDealt,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procMask := core.Ternary(spell.SpellID == 86392, core.ProcMaskMeleeMH, spell.ProcMask)
				if rogue.woundPoisonPPHM[itemSlot].Proc(sim, procMask, "Wound Poison") {
					rogue.WoundPoison.Cast(sim, result.Target)
				}
			},
		})
	}
}
