package frost

// func (dk *DeathKnight) ThreatOfThassarianProc(sim *core.Simulation, result *core.SpellResult, ohSpell *core.Spell) {
// 	if dk.Talents.ThreatOfThassarian == 0 || dk.GetOHWeapon() == nil {
// 		return
// 	}
// 	if sim.Proc([]float64{0.0, 0.3, 0.6, 1.0}[dk.Talents.ThreatOfThassarian], "Threat of Thassarian") {
// 		ohSpell.Cast(sim, result.Target)
// 	}
// }

// ohSpell := dk.GetOrRegisterSpell(core.SpellConfig{
// 		ActionID:       PlagueStrikeActionID.WithTag(2),
// 		SpellSchool:    core.SpellSchoolPhysical,
// 		ProcMask:       core.ProcMaskMeleeOHSpecial,
// 		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell,
// 		ClassSpellMask: DeathKnightSpellPlagueStrike,

// 		DamageMultiplier: 1,
// 		CritMultiplier:   dk.DefaultCritMultiplier(),
// 		ThreatMultiplier: 1,

// 		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 			baseDamage := dk.CalcScalingSpellDmg(0.18700000644) +
// 				spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

// 			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
// 		},
// 	})

// ohSpell := dk.GetOrRegisterSpell(core.SpellConfig{
// 		ActionID:       BloodStrikeActionID.WithTag(2),
// 		SpellSchool:    core.SpellSchoolPhysical,
// 		ProcMask:       core.ProcMaskMeleeOHSpecial,
// 		Flags:          core.SpellFlagMeleeMetrics,
// 		ClassSpellMask: DeathKnightSpellBloodStrike,

// 		DamageMultiplier:         0.8,
// 		DamageMultiplierAdditive: 1,
// 		CritMultiplier:           dk.DefaultCritMultiplier(),
// 		ThreatMultiplier:         1,

// 		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 			baseDamage := dk.ClassSpellScaling*0.37799999118 +
// 				spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

// 			baseDamage *= dk.GetDiseaseMulti(target, 1.0, 0.025)

// 			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
// 		},
// 	})

// ohSpell := dk.GetOrRegisterSpell(core.SpellConfig{
// 	ActionID:       DeathStrikeActionID.WithTag(2),
// 	SpellSchool:    core.SpellSchoolPhysical,
// 	ProcMask:       core.ProcMaskMeleeOHSpecial,
// 	Flags:          core.SpellFlagMeleeMetrics,
// 	ClassSpellMask: DeathKnightSpellDeathStrike,

// 	DamageMultiplier: 1.5,
// 	CritMultiplier:   dk.DefaultCritMultiplier(),
// 	ThreatMultiplier: 1,

// 	ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 		baseDamage := dk.ClassSpellScaling*0.14699999988 +
// 			spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

// 		spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
// 		doHealing(sim, 0.05)
// 	},
// })

// ohSpell := dk.GetOrRegisterSpell(core.SpellConfig{
// 		ActionID:       obliterateActionID.WithTag(2),
// 		SpellSchool:    core.SpellSchoolPhysical,
// 		ProcMask:       core.ProcMaskMeleeOHSpecial,
// 		Flags:          core.SpellFlagMeleeMetrics,
// 		ClassSpellMask: DeathKnightSpellObliterate,

// 		DamageMultiplier:         1.5,
// 		DamageMultiplierAdditive: 1,
// 		CritMultiplier:           dk.DefaultCritMultiplier(),
// 		ThreatMultiplier:         1,

// 		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 			baseDamage := dk.ClassSpellScaling*0.28900000453 +
// 				spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

// 			baseDamage *= dk.GetDiseaseMulti(target, 1.0, 0.125)

// 			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
// 		},
// 	})

// ohSpell := dk.GetOrRegisterSpell(core.SpellConfig{
// 		ActionID:       frostStrikeActionID.WithTag(2),
// 		SpellSchool:    core.SpellSchoolFrost,
// 		ProcMask:       core.ProcMaskMeleeOHSpecial,
// 		Flags:          core.SpellFlagMeleeMetrics,
// 		ClassSpellMask: death_knight.DeathKnightSpellFrostStrike,

// 		DamageMultiplier:         1.3,
// 		DamageMultiplierAdditive: 1,
// 		CritMultiplier:           dk.DefaultCritMultiplier(),
// 		ThreatMultiplier:         1,

// 		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 			baseDamage := dk.ClassSpellScaling*0.12399999797 +
// 				spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

// 			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
// 		},
// 	})
