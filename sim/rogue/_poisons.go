package rogue

import (
	"fmt"
	"strconv"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (rogue *Rogue) applyPoisons() {
	rogue.applyDeadlyPoison()
	rogue.applyInstantPoison()
	rogue.applyWoundPoison()
}

func (rogue *Rogue) registerPoisonAuras() {
	if rogue.Talents.SavageCombat > 0 {
		rogue.SavageCombatDebuffAuras = rogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
			return core.SavageCombatAura(target, rogue.Talents.SavageCombat)
		})
	}
	if rogue.Talents.MasterPoisoner {
		rogue.MasterPoisonerDebuffAuras = rogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
			aura := core.MasterPoisonerDebuff(target)
			aura.Duration = core.NeverExpires
			return aura
		})
	}
}

func (rogue *Rogue) registerDeadlyPoisonSpell() {
	baseDamage := 0.11999999732 * rogue.ClassSpellScaling
	apScaling := 0.03500000015

	rogue.DeadlyPoison = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 2818},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamageProc,
		ClassSpellMask: RogueSpellDeadlyPoison,
		Flags:          core.SpellFlagPassiveSpell,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1 + 0.12*float64(rogue.Talents.VilePoisons),
		CritMultiplier:           1,
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "Deadly Poison",
				MaxStacks: 5,
				Duration:  time.Second * 12,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					if rogue.Talents.SavageCombat > 0 {
						rogue.SavageCombatDebuffAuras.Get(aura.Unit).Activate(sim)
					}
					if rogue.Talents.MasterPoisoner {
						rogue.MasterPoisonerDebuffAuras.Get(aura.Unit).Activate(sim)
					}
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					if rogue.Talents.SavageCombat > 0 {
						rogue.SavageCombatDebuffAuras.Get(aura.Unit).Deactivate(sim)
					}
					if rogue.Talents.MasterPoisoner {
						rogue.MasterPoisonerDebuffAuras.Get(aura.Unit).Deactivate(sim)
					}
				},
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 3,

			OnSnapshot: func(_ *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				if stacks := dot.GetStacks(); stacks > 0 {
					dot.SnapshotBaseDamage = (baseDamage + apScaling*dot.Spell.MeleeAttackPower()) * float64(stacks)
					attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
					dot.SnapshotCritChance = dot.Spell.SpellCritChance(attackTable.Defender)
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable, true)
				}
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
			if !dot.IsActive() {
				dot.Apply(sim)
				dot.SetStacks(sim, 1)
				dot.TakeSnapshot(sim, false)
				return
			}

			if dot.GetStacks() < 5 {
				dot.Refresh(sim)
				dot.AddStack(sim)
				dot.TakeSnapshot(sim, false)
				return
			}

			if rogue.lastDeadlyPoisonProcMask.Matches(core.ProcMaskMeleeMH) {
				switch rogue.Options.OhImbue {
				case proto.RogueOptions_InstantPoison:
					rogue.InstantPoison[DeadlyProc].Cast(sim, target)
				case proto.RogueOptions_WoundPoison:
					rogue.WoundPoison[DeadlyProc].Cast(sim, target)
				}
			}

			// Confirmed: Thrown Deadly Poison proc only the MH poison, and is not proc'd from MH/OH Deadly Poison
			if rogue.lastDeadlyPoisonProcMask.Matches(core.ProcMaskMeleeOH | core.ProcMaskRanged | core.ProcMaskMeleeProc) {
				switch rogue.Options.MhImbue {
				case proto.RogueOptions_InstantPoison:
					rogue.InstantPoison[DeadlyProc].Cast(sim, target)
				case proto.RogueOptions_WoundPoison:
					rogue.WoundPoison[DeadlyProc].Cast(sim, target)
				}
			}

			dot.Refresh(sim)
			dot.TakeSnapshot(sim, false)
		},
	})
}

func (rogue *Rogue) procDeadlyPoison(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	rogue.lastDeadlyPoisonProcMask = spell.ProcMask
	rogue.DeadlyPoison.Cast(sim, result.Target)
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
					rogue.procDeadlyPoison(sim, spell, result)
				}
			},
		})
	}

}

func (rogue *Rogue) applyWoundPoison() {
	const basePPM = 0.5 / (1.4 / 60) // ~21.43, the former 50% normalized to a 1.4 speed weapon

	for _, itemSlot := range core.AllWeaponSlots() {
		procMask := rogue.getPoisonProcMaskForSlot(proto.RogueOptions_WoundPoison, itemSlot)
		rogue.woundPoisonPPMM[itemSlot] = rogue.AutoAttacks.NewPPMManager(basePPM, procMask)

		core.MakeProcTriggerAura(&rogue.Unit, core.ProcTrigger{
			Name:     fmt.Sprintf("Wound Poison %s", itemSlot),
			Outcome:  core.OutcomeLanded,
			Callback: core.CallbackOnSpellHitDealt,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procMask := core.Ternary(spell.SpellID == 86392, core.ProcMaskMeleeMH, spell.ProcMask)
				if rogue.woundPoisonPPMM[itemSlot].Proc(sim, procMask, "Wound Poison") {
					rogue.WoundPoison[NormalProc].Cast(sim, result.Target)
				}
			},
		})
	}
}

type PoisonProcSource int

const (
	NormalProc PoisonProcSource = iota
	DeadlyProc
	ShivProc
	VilePoisonsProc
)

func (rogue *Rogue) makeInstantPoison(procSource PoisonProcSource) *core.Spell {
	isShivProc := procSource == ShivProc
	ipBaseDamage := 0.31299999356 * rogue.ClassSpellScaling
	ipVariance := 0.28000000119 * ipBaseDamage
	minHit := ipBaseDamage - ipVariance/2
	apScaling := 0.09000000358

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 8680, Tag: int32(procSource)},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamageProc,
		ClassSpellMask: RogueSpellInstantPoison,
		Flags:          core.SpellFlagPassiveSpell,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1 + 0.12*float64(rogue.Talents.VilePoisons),
		CritMultiplier:           rogue.DefaultCritMultiplier()(),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := minHit +
				sim.RandomFloat("Instant Poison")*ipVariance +
				apScaling*spell.MeleeAttackPower()
			if isShivProc {
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHit)
			} else {
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}

func (rogue *Rogue) makeWoundPoison(procSource PoisonProcSource) *core.Spell {
	isShivProc := procSource == ShivProc
	wpBaseDamage := 0.24500000477 * rogue.ClassSpellScaling
	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 13218, Tag: int32(procSource)},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamageProc,
		ClassSpellMask: RogueSpellWoundPoison,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1 + 0.12*float64(rogue.Talents.VilePoisons),
		CritMultiplier:           rogue.DefaultCritMultiplier(),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := wpBaseDamage + 0.04*spell.MeleeAttackPower()

			var result *core.SpellResult
			if isShivProc {
				result = spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHit)
			} else {
				result = spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			}

			if result.Landed() {
				rogue.WoundPoisonDebuffAuras.Get(target).Activate(sim)
			}
		},
	})
}

var WoundPoisonActionID = core.ActionID{SpellID: 13219}

func (rogue *Rogue) registerWoundPoisonSpell() {
	woundPoisonDebuffAura := core.Aura{
		Label:    "WoundPoison-" + strconv.Itoa(int(rogue.Index)),
		ActionID: WoundPoisonActionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if rogue.Talents.SavageCombat > 0 {
				rogue.SavageCombatDebuffAuras.Get(aura.Unit).Activate(sim)
			}
			if rogue.Talents.MasterPoisoner {
				rogue.MasterPoisonerDebuffAuras.Get(aura.Unit).Activate(sim)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if rogue.Talents.SavageCombat > 0 {
				rogue.SavageCombatDebuffAuras.Get(aura.Unit).Deactivate(sim)
			}
			if rogue.Talents.MasterPoisoner {
				rogue.MasterPoisonerDebuffAuras.Get(aura.Unit).Deactivate(sim)
			}
		},
	}

	rogue.WoundPoisonDebuffAuras = rogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.RegisterAura(woundPoisonDebuffAura)
	})
	rogue.WoundPoison = [4]*core.Spell{
		rogue.makeWoundPoison(NormalProc),
		rogue.makeWoundPoison(DeadlyProc),
		rogue.makeWoundPoison(ShivProc),
		rogue.makeWoundPoison(VilePoisonsProc),
	}
}

func (rogue *Rogue) registerInstantPoisonSpell() {
	rogue.InstantPoison = [4]*core.Spell{
		rogue.makeInstantPoison(NormalProc),
		rogue.makeInstantPoison(DeadlyProc),
		rogue.makeInstantPoison(ShivProc),
		rogue.makeInstantPoison(VilePoisonsProc),
	}
}

func (rogue *Rogue) GetDeadlyPoisonProcChance() float64 {
	return 0.3 + core.TernaryFloat64(rogue.Spec == proto.Spec_SpecAssassinationRogue, 0.2, 0)
}

func (rogue *Rogue) UpdateDeadlyPoisonPPH(bonusChance float64) {
	pph := rogue.GetDeadlyPoisonProcChance() + bonusChance
	for _, itemSlot := range core.AllWeaponSlots() {
		procMask := rogue.getPoisonProcMaskForSlot(proto.RogueOptions_DeadlyPoison, itemSlot)
		rogue.deadlyPoisonPPHM[itemSlot] = rogue.AutoAttacks.NewStaticDynamicProcManager(pph, procMask)
	}
}

func (rogue *Rogue) UpdateInstantPoisonPPM(bonusChance float64) {
	const basePPM = 0.2 / (1.4 / 60) // ~8.57, the former 20% normalized to a 1.4 speed weapon
	ppm := basePPM * (1 + core.TernaryFloat64(rogue.Spec == proto.Spec_SpecAssassinationRogue, 0.5, 0) + bonusChance)

	for _, itemSlot := range core.AllWeaponSlots() {
		procMask := rogue.getPoisonProcMaskForSlot(proto.RogueOptions_InstantPoison, itemSlot)
		rogue.instantPoisonPPMM[itemSlot] = rogue.AutoAttacks.NewStaticPPMManager(ppm, procMask)
	}
}

func (rogue *Rogue) applyInstantPoison() {
	rogue.UpdateInstantPoisonPPM(0)

	for _, itemSlot := range core.AllWeaponSlots() {
		core.MakeProcTriggerAura(&rogue.Unit, core.ProcTrigger{
			Name:     fmt.Sprintf("Instant Poison %s", itemSlot),
			Outcome:  core.OutcomeLanded,
			Callback: core.CallbackOnSpellHitDealt,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procMask := core.Ternary(spell.SpellID == 86392, core.ProcMaskMeleeMH, spell.ProcMask)
				if rogue.instantPoisonPPMM[itemSlot].Proc(sim, procMask, "Instant Poison") {
					rogue.InstantPoison[NormalProc].Cast(sim, result.Target)
				}
			},
		})
	}
}
