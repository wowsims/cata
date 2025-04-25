package demonology

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/wowsims/mop/sim/common"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/simsignals"
	"github.com/wowsims/mop/sim/core/stats"
)

func init() {
	RegisterDemonologyWarlock()
}

func setupFakeSim(lockStats stats.Stats, glyphs *proto.Glyphs, duration float64) *core.Simulation {
	var specOptions = &proto.Player_DemonologyWarlock{
		DemonologyWarlock: &proto.DemonologyWarlock{
			Options: &proto.DemonologyWarlock_Options{
				ClassOptions: &proto.WarlockOptions{
					Summon:       proto.WarlockOptions_NoSummon,
					DetonateSeed: false,
				},
			},
		},
	}

	sim := core.NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{
			RandomSeed: 19137,
		},

		Raid: &proto.Raid{
			Parties: []*proto.Party{
				{Players: []*proto.Player{{
					Name:          "Caster",
					Class:         proto.Class_ClassWarlock,
					Race:          proto.Race_RaceOrc,
					Consumes:      &proto.Consumes{},
					Buffs:         &proto.IndividualBuffs{},
					TalentsString: "-2312222310310212211-33202",
					Glyphs:        glyphs,
					Spec:          specOptions,
					Equipment:     &proto.EquipmentSpec{},
				}}, Buffs: &proto.PartyBuffs{}}},
		},
		Encounter: &proto.Encounter{
			Targets: []*proto.Target{
				{Name: "target", Level: 83, MobType: proto.MobType_MobTypeDemon},
			},
			Duration: duration,
		},
	}, simsignals.CreateSignals())

	sim.Options.Debug = true
	sim.Log = func(message string, vals ...interface{}) {
		fmt.Printf(fmt.Sprintf("[%0.1f] "+message+"\n", append([]interface{}{sim.CurrentTime.Seconds()}, vals...)...))
	}

	sim.Reset()
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)

	lock.AddStatsDynamic(sim, stats.Stats{
		stats.SpellPower: -lock.GetStat(stats.SpellPower),
	})
	lock.AddStatsDynamic(sim, lockStats)

	return sim
}

var defStats = stats.Stats{
	stats.SpellHitPercent:  17,
	stats.SpellCritPercent: -100,
	stats.SpellPower:       6343,
	// this was tested before Blizzard fixed mastery flooring, true mastery in these tests was 1059
	// which was rounded down to 31%. So the "effective" mastery was ((31 / 2.3) - 8) * 179.280040
	stats.MasteryRating: 982.1428278260872,
}

func (lock *DemonologyWarlock) checkSpell(t *testing.T, sim *core.Simulation, spell *core.Spell, expected float64,
	expectedMult float64) {

	damageBefore := spell.SpellMetrics[0].TotalDamage
	if !spell.Cast(sim, lock.CurrentTarget) {
		t.Fatal("Failed to cast spell")
	}

	for {
		if finished := sim.Step(); finished {
			break
		}
	}

	damageAfter := spell.SpellMetrics[0].TotalDamage
	delta := damageAfter - damageBefore

	mult := spell.AttackerDamageMultiplier(lock.AttackTables[lock.CurrentTarget.UnitIndex], false)
	if !core.WithinToleranceFloat64(expectedMult, mult, 0.000001) {
		t.Fatalf("Incorrect multiplier: Expected: %0.6f, Actual: %0.6f", expectedMult, mult)
	}

	if !core.WithinToleranceFloat64(expected, delta, 0.0001) {
		t.Fatalf("Incorrect damage applied: Expected: %0.4f, Actual: %0.4f", expected, delta)
	}
}

func (lock *DemonologyWarlock) checkSpellDamageRange(
	t *testing.T, sim *core.Simulation, spell *core.Spell, min float64, max float64, expectMult float64) {

	mult := spell.AttackerDamageMultiplier(lock.AttackTables[lock.CurrentTarget.UnitIndex], false)
	if !core.WithinToleranceFloat64(expectMult, mult, 0.000001) {
		t.Fatalf("Incorrect multiplier: Expected: %0.6f, Actual: %0.6f", expectMult, mult)
	}

	for i := 1; i <= 300; i++ {
		damageBefore := spell.SpellMetrics[0].TotalDamage

		spell.SkipCastAndApplyEffects(sim, lock.CurrentTarget)
		sim.Step()

		damageAfter := spell.SpellMetrics[0].TotalDamage
		delta := damageAfter - damageBefore

		if delta > max || delta < min {
			t.Fatalf("Incorrect damage applied: %0.4f < %0.4f < %0.4f", min, delta, max)
		}
	}
}

func checkDotTick(t *testing.T, sim *core.Simulation, dot *core.Dot, attackTable *core.AttackTable, expected float64, expectedMult float64) {
	dot.Apply(sim)
	damageBefore := dot.Spell.SpellMetrics[0].TotalDamage
	dot.TickOnce(sim)
	damageAfter := dot.Spell.SpellMetrics[0].TotalDamage

	mult := dot.Spell.AttackerDamageMultiplier(attackTable, true)
	if !core.WithinToleranceFloat64(expectedMult, mult, 0.000001) {
		t.Fatalf("Incorrect multiplier: Expected: %0.6f, Actual: %0.6f", expectedMult, mult)
	}

	delta := damageAfter - damageBefore
	if !core.WithinToleranceFloat64(expected, delta, 0.0001) {
		t.Fatalf("Incorrect damage applied: Expected: %0.4f, Actual: %0.4f", expected, delta)
	}
}

func TestImmolateDoTBase(t *testing.T) {
	sim := setupFakeSim(defStats, &proto.Glyphs{}, 10)
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)

	attackTable := lock.AttackTables[lock.CurrentTarget.UnitIndex]
	checkDotTick(t, sim, lock.Immolate.CurDot(), attackTable, 2166.0618, 1.4076)
}

func TestImmolateDoTGlyphed(t *testing.T) {
	sim := setupFakeSim(defStats, &proto.Glyphs{Prime1: int32(proto.WarlockPrimeGlyph_GlyphOfImmolate)}, 10)
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)

	attackTable := lock.AttackTables[lock.CurrentTarget.UnitIndex]
	checkDotTick(t, sim, lock.Immolate.CurDot(), attackTable, 2346.5670, 1.5249)
}

func TestImmolateNonPeriodic(t *testing.T) {
	sim := setupFakeSim(defStats, &proto.Glyphs{}, 10)
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)

	lock.checkSpell(t, sim, lock.Immolate, 2901.6213, 1.4076)
}

func TestIncinerateBase(t *testing.T) {
	sim := setupFakeSim(defStats, &proto.Glyphs{}, 10)
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)
	incinerate := lock.GetSpell(core.ActionID{SpellID: 29722})

	lock.checkSpellDamageRange(t, sim, incinerate, 5161.68, 5270.35, 1.31376)
}

func TestIncinerateMoltenCore(t *testing.T) {
	sim := setupFakeSim(defStats, &proto.Glyphs{}, 10)
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)
	lock.GetAura("Molten Core Proc Aura").Activate(sim)
	incinerate := lock.GetSpell(core.ActionID{SpellID: 29722})

	lock.checkSpellDamageRange(t, sim, incinerate, 5991.24, 6117.37, 1.5249)
}

func TestIncinerateGlyphed(t *testing.T) {
	sim := setupFakeSim(defStats, &proto.Glyphs{Prime1: int32(proto.WarlockPrimeGlyph_GlyphOfIncinerate)}, 10)
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)
	incinerate := lock.GetSpell(core.ActionID{SpellID: 29722})

	lock.checkSpellDamageRange(t, sim, incinerate, 5392.11, 5505.63, 1.37241)
}

func TestIncinerateGlyphedMoltenCore(t *testing.T) {
	sim := setupFakeSim(defStats, &proto.Glyphs{Prime1: int32(proto.WarlockPrimeGlyph_GlyphOfIncinerate)}, 10)
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)
	lock.GetAura("Molten Core Proc Aura").Activate(sim)
	incinerate := lock.GetSpell(core.ActionID{SpellID: 29722})

	lock.checkSpellDamageRange(t, sim, incinerate, 6221.67, 6352.65, 1.58355)
}

func TestIncinerateGlyphedMoltenCoreMeta(t *testing.T) {
	sim := setupFakeSim(defStats, &proto.Glyphs{Prime1: int32(proto.WarlockPrimeGlyph_GlyphOfIncinerate)}, 10)
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)
	lock.Metamorphosis.SkipCastAndApplyEffects(sim, lock.CurrentTarget)
	lock.GetAura("Molten Core Proc Aura").Activate(sim)
	incinerate := lock.GetSpell(core.ActionID{SpellID: 29722})

	lock.checkSpellDamageRange(t, sim, incinerate, 9394.72, 9592.5, 2.39116)
}

func TestImmolationAura(t *testing.T) {
	sim := setupFakeSim(defStats, &proto.Glyphs{}, 10)
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)
	lock.Metamorphosis.SkipCastAndApplyEffects(sim, lock.CurrentTarget)

	immoAura := lock.GetSpell(core.ActionID{SpellID: 50589})
	attackTable := lock.AttackTables[lock.CurrentTarget.UnitIndex]
	checkDotTick(t, sim, immoAura.AOEDot(), attackTable, 2127.4521, 1.771230)
}

func TestMoltenCoreCastTime(t *testing.T) {
	sim := setupFakeSim(defStats, &proto.Glyphs{}, 10)
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)
	lock.GetAura("Molten Core Proc Aura").Activate(sim)
	incinerate := lock.GetSpell(core.ActionID{SpellID: 29722})

	expectedCasttime := 1575 * time.Millisecond
	if incinerate.CastTime() != expectedCasttime {
		t.Fatalf("Incorrect casttime: Expected: %v, Actual: %v", expectedCasttime, incinerate.CastTime())
	}

	expectedGCD := 1500 * time.Millisecond
	if incinerate.DefaultCast.GCD != expectedGCD {
		t.Fatalf("Incorrect GCD: Expected: %v, Actual: %v", expectedGCD, incinerate.DefaultCast.GCD)
	}
}

func TestMoltenCoreAfterCast(t *testing.T) {
	sim := setupFakeSim(defStats, &proto.Glyphs{}, 10)
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)
	incinerate := lock.GetSpell(core.ActionID{SpellID: 29722})

	if !incinerate.Cast(sim, lock.CurrentTarget) {
		t.Fatal("Failed to cast spell")
	}

	sim.AddPendingAction(&core.PendingAction{
		NextActionAt: sim.CurrentTime + 500*time.Millisecond,
		OnAction: func(sim *core.Simulation) {
			lock.GetAura("Molten Core Proc Aura").Activate(sim)
			lock.GetAura("Molten Core Proc Aura").SetStacks(sim, 3)
		},
	})

	for {
		if finished := sim.Step(); finished {
			break
		}
	}

	if lock.GetAura("Molten Core Proc Aura").GetStacks() != 3 {
		t.Fatalf("Consumed a molten core stack when it shouldn't have")
	}

	expected := 1.5249
	mult := incinerate.AttackerDamageMultiplier(lock.AttackTables[lock.CurrentTarget.UnitIndex], false)
	if !core.WithinToleranceFloat64(expected, mult, 0.0001) {
		t.Fatalf("Incorrect damage applied: Expected: %0.4f, Actual: %0.4f", expected, mult)
	}
}

func TestDecimationCastTime(t *testing.T) {
	sim := setupFakeSim(defStats, &proto.Glyphs{}, 10)
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)
	lock.GetAura("Decimation Proc Aura").Activate(sim)
	soulfire := lock.GetSpell(core.ActionID{SpellID: 6353})

	expectedCasttime := 1800 * time.Millisecond
	if soulfire.CastTime() != expectedCasttime {
		t.Fatalf("Incorrect casttime: Expected: %v, Actual: %v", expectedCasttime, soulfire.CastTime())
	}

	expectedGCD := 1500 * time.Millisecond
	if soulfire.DefaultCast.GCD != expectedGCD {
		t.Fatalf("Incorrect GCD: Expected: %v, Actual: %v", expectedGCD, soulfire.DefaultCast.GCD)
	}
}

func checkTicks(t *testing.T, dot *core.Dot, msg string, expected int32) {
	if dot.RemainingTicks() != expected {
		t.Helper()
		t.Fatalf("%s: Expected: %v, Actual: %v", msg, expected, dot.RemainingTicks())
	}
}

func waitUntilTime(sim *core.Simulation, time time.Duration) {
	for {
		sim.Step()
		if sim.CurrentTime >= time {
			break
		}
	}
}

func TestFelFlameExtension(t *testing.T) {
	sim := setupFakeSim(defStats, &proto.Glyphs{}, 30)
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)
	immoDot := lock.Immolate.CurDot()
	immoDot.Apply(sim)
	felflame := lock.GetSpell(core.ActionID{SpellID: 77799})

	checkTicks(t, immoDot, "Baseline immolate ticks are wrong?", 7)

	felflame.SkipCastAndApplyEffects(sim, lock.CurrentTarget)
	sim.Step()
	checkTicks(t, immoDot, "Incorrect tick count after fel flame extension", 7)

	// applying it again shouldn't change ticks since we're at the cap
	felflame.SkipCastAndApplyEffects(sim, lock.CurrentTarget)
	sim.Step()
	checkTicks(t, immoDot, "Incorrect tick count after fel flame extension", 7)

	waitUntilTime(sim, 12500*time.Millisecond)
	checkTicks(t, immoDot, "Incorrect tick count after waiting 9s", 3)

	felflame.SkipCastAndApplyEffects(sim, lock.CurrentTarget)
	sim.Step()
	checkTicks(t, immoDot, "Incorrect tick count after fel flame extension", 5)
}

func TestShadowflameHasteCap(t *testing.T) {
	sim := setupFakeSim(defStats, &proto.Glyphs{}, 10)
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)
	lock.Unit.MultiplyCastSpeed(1 + 0.05) // 5% haste buff
	lock.Unit.MultiplyCastSpeed(1 + 0.03) // dark intent
	lock.AddStatsDynamic(sim, stats.Stats{
		stats.HasteRating: 1006,
	})
	shadowflame := lock.GetSpell(core.ActionID{SpellID: 47897})
	shadowflameDot := lock.GetSpell(core.ActionID{SpellID: 47897}.WithTag(1)).CurDot()

	shadowflame.SkipCastAndApplyEffects(sim, lock.CurrentTarget)
	checkTicks(t, shadowflameDot, "Incorrect tick count for shadowflame at 1006 haste", 3)
	shadowflameDot.Deactivate(sim)

	lock.AddStatsDynamic(sim, stats.Stats{
		stats.HasteRating: 1,
	})
	shadowflame.SkipCastAndApplyEffects(sim, lock.CurrentTarget)
	checkTicks(t, shadowflameDot, "Incorrect tick count for shadowflame at 1007 haste", 4)
}

func TestImmolateHasteCap(t *testing.T) {
	sim := setupFakeSim(defStats, &proto.Glyphs{}, 10)
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)
	lock.Unit.MultiplyCastSpeed(1 + 0.05) // 5% haste buff
	lock.Unit.MultiplyCastSpeed(1 + 0.03) // dark intent
	lock.AddStatsDynamic(sim, stats.Stats{
		stats.HasteRating: 1572,
	})
	immolate := lock.Immolate
	immolateDot := lock.Immolate.CurDot()

	immolate.SkipCastAndApplyEffects(sim, lock.CurrentTarget)
	checkTicks(t, immolateDot, "Incorrect tick count for immolate at 1572 haste", 8)
	immolateDot.Deactivate(sim)

	lock.AddStatsDynamic(sim, stats.Stats{
		stats.HasteRating: 1,
	})
	immolate.SkipCastAndApplyEffects(sim, lock.CurrentTarget)
	checkTicks(t, immolateDot, "Incorrect tick count for immolate at 1573 haste", 9)
}

func TestCorruptionHasteCap(t *testing.T) {
	sim := setupFakeSim(defStats, &proto.Glyphs{}, 10)
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)
	lock.Unit.MultiplyCastSpeed(1 + 0.05) // 5% haste buff
	lock.Unit.MultiplyCastSpeed(1 + 0.03) // dark intent
	lock.AddStatsDynamic(sim, stats.Stats{
		stats.HasteRating: 1992,
	})
	corruption := lock.Corruption
	corruptionDot := corruption.CurDot()

	corruption.SkipCastAndApplyEffects(sim, lock.CurrentTarget)
	checkTicks(t, corruptionDot, "Incorrect tick count for corruption at 1992 haste", 7)
	corruptionDot.Deactivate(sim)

	lock.AddStatsDynamic(sim, stats.Stats{
		stats.HasteRating: 1,
	})
	corruption.SkipCastAndApplyEffects(sim, lock.CurrentTarget)
	checkTicks(t, corruptionDot, "Incorrect tick count for corruption at 1993 haste", 8)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func TestDemonology(t *testing.T) {
	var defaultDemonologyWarlock = &proto.Player_DemonologyWarlock{
		DemonologyWarlock: &proto.DemonologyWarlock{
			Options: &proto.DemonologyWarlock_Options{
				ClassOptions: &proto.WarlockOptions{
					Summon:       proto.WarlockOptions_Felguard,
					DetonateSeed: false,
				},
			},
		},
	}

	var itemFilter = core.ItemFilter{
		WeaponTypes: []proto.WeaponType{
			proto.WeaponType_WeaponTypeSword,
			proto.WeaponType_WeaponTypeDagger,
			proto.WeaponType_WeaponTypeStaff,
		},
		HandTypes: []proto.HandType{
			proto.HandType_HandTypeOffHand,
		},
		ArmorType: proto.ArmorType_ArmorTypeCloth,
		RangedWeaponTypes: []proto.RangedWeaponType{
			proto.RangedWeaponType_RangedWeaponTypeWand,
		},
	}

	var fullConsumes = &proto.Consumes{
		Flask:             proto.Flask_FlaskOfTheDraconicMind,
		Food:              proto.Food_FoodSeveredSagefish,
		DefaultPotion:     proto.Potions_VolcanicPotion,
		ExplosiveBigDaddy: true,
		TinkerHands:       proto.TinkerHands_TinkerHandsSynapseSprings,
	}

	// Shadow Bolt
	var demonologyTalentsShadowBolt = "-3312222300310212211-33202"
	var demonologyGlyphsShadowBolt = &proto.Glyphs{
		Prime1: int32(proto.WarlockPrimeGlyph_GlyphOfImmolate),
		Prime2: int32(proto.WarlockPrimeGlyph_GlyphOfCorruption),
		Prime3: int32(proto.WarlockPrimeGlyph_GlyphOfMetamorphosis),
		Major1: int32(proto.WarlockMajorGlyph_GlyphOfShadowBolt),
		Major2: int32(proto.WarlockMajorGlyph_GlyphOfLifeTap),
		Major3: int32(proto.WarlockMajorGlyph_GlyphOfFelhunter),
	}

	// Incinerate
	var demonologyTalentsIncenerate = "003-3312222300310212211-03202"
	var demonologyGlyphsIncinerate = &proto.Glyphs{
		Prime1: int32(proto.WarlockPrimeGlyph_GlyphOfImmolate),
		Prime2: int32(proto.WarlockPrimeGlyph_GlyphOfIncinerate),
		Prime3: int32(proto.WarlockPrimeGlyph_GlyphOfMetamorphosis),
		Major1: int32(proto.WarlockMajorGlyph_GlyphOfSoulstone),
		Major2: int32(proto.WarlockMajorGlyph_GlyphOfLifeTap),
		Major3: int32(proto.WarlockMajorGlyph_GlyphOfFelhunter),
	}

	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:       proto.Class_ClassWarlock,
		Race:        proto.Race_RaceOrc,
		OtherRaces:  []proto.Race{proto.Race_RaceTroll, proto.Race_RaceGoblin, proto.Race_RaceHuman},
		GearSet:     core.GetGearSet("../../../ui/warlock/demonology/gear_sets", "p4"),
		ItemSwapSet: core.GetItemSwapGearSet("../../../ui/warlock/demonology/gear_sets", "p4_item_swap"),
		Talents:     demonologyTalentsIncenerate,
		Glyphs:      demonologyGlyphsIncinerate,
		OtherTalentSets: []core.TalentsCombo{
			{
				Label:   "Incinerate",
				Talents: demonologyTalentsShadowBolt,
				Glyphs:  demonologyGlyphsShadowBolt,
			},
		},
		Consumes:         fullConsumes,
		SpecOptions:      core.SpecOptionsCombo{Label: "Demonology Warlock", SpecOptions: defaultDemonologyWarlock},
		OtherSpecOptions: []core.SpecOptionsCombo{},
		Rotation:         core.GetAplRotation("../../../ui/warlock/demonology/apls", "incinerate"),
		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../../ui/warlock/demonology/apls", "shadow-bolt"),
		},
		ItemFilter:       itemFilter,
		StartingDistance: 25,
	}))
}
