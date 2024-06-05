package demonology

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/wowsims/cata/sim/common"
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/simsignals"
	"github.com/wowsims/cata/sim/core/stats"
)

func init() {
	RegisterDemonologyWarlock()
}

func setupFakeSim(lockStats stats.Stats, glyphs *proto.Glyphs) *core.Simulation {
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
			Duration: 10,
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
	stats.SpellHit:   17 * core.SpellHitRatingPerHitChance,
	stats.SpellCrit:  -100 * core.CritRatingPerCritChance,
	stats.SpellPower: 6343,
	stats.Mastery:    1059,
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
	sim := setupFakeSim(defStats, &proto.Glyphs{})
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)

	attackTable := lock.AttackTables[lock.CurrentTarget.UnitIndex]
	checkDotTick(t, sim, lock.ImmolateDot.CurDot(), attackTable, 2166.0618, 1.4076)
}

func TestImmolateDoTGlyphed(t *testing.T) {
	sim := setupFakeSim(defStats, &proto.Glyphs{Prime1: int32(proto.WarlockPrimeGlyph_GlyphOfImmolate)})
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)

	attackTable := lock.AttackTables[lock.CurrentTarget.UnitIndex]
	checkDotTick(t, sim, lock.ImmolateDot.CurDot(), attackTable, 2346.5670, 1.5249)
}

func TestImmolateNonPeriodic(t *testing.T) {
	sim := setupFakeSim(defStats, &proto.Glyphs{})
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)

	lock.checkSpell(t, sim, lock.Immolate, 2901.6213, 1.4076)
}

func TestIncinerateBase(t *testing.T) {
	sim := setupFakeSim(defStats, &proto.Glyphs{})
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)
	incinerate := lock.GetSpell(core.ActionID{SpellID: 29722})

	lock.checkSpellDamageRange(t, sim, incinerate, 5161.68, 5270.35, 1.31376)
}

func TestIncinerateMoltenCore(t *testing.T) {
	sim := setupFakeSim(defStats, &proto.Glyphs{})
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)
	lock.GetAura("Molten Core Proc Aura").Activate(sim)
	incinerate := lock.GetSpell(core.ActionID{SpellID: 29722})

	lock.checkSpellDamageRange(t, sim, incinerate, 5991.24, 6117.37, 1.5249)
}

func TestIncinerateGlyphed(t *testing.T) {
	sim := setupFakeSim(defStats, &proto.Glyphs{Prime1: int32(proto.WarlockPrimeGlyph_GlyphOfIncinerate)})
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)
	incinerate := lock.GetSpell(core.ActionID{SpellID: 29722})

	lock.checkSpellDamageRange(t, sim, incinerate, 5392.11, 5505.63, 1.37241)
}

func TestIncinerateGlyphedMoltenCore(t *testing.T) {
	sim := setupFakeSim(defStats, &proto.Glyphs{Prime1: int32(proto.WarlockPrimeGlyph_GlyphOfIncinerate)})
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)
	lock.GetAura("Molten Core Proc Aura").Activate(sim)
	incinerate := lock.GetSpell(core.ActionID{SpellID: 29722})

	lock.checkSpellDamageRange(t, sim, incinerate, 6221.67, 6352.65, 1.58355)
}

func TestIncinerateGlyphedMoltenCoreMeta(t *testing.T) {
	sim := setupFakeSim(defStats, &proto.Glyphs{Prime1: int32(proto.WarlockPrimeGlyph_GlyphOfIncinerate)})
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)
	lock.Metamorphosis.SkipCastAndApplyEffects(sim, lock.CurrentTarget)
	lock.GetAura("Molten Core Proc Aura").Activate(sim)
	incinerate := lock.GetSpell(core.ActionID{SpellID: 29722})

	lock.checkSpellDamageRange(t, sim, incinerate, 9394.72, 9592.5, 2.39116)
}

func TestImmolationAura(t *testing.T) {
	sim := setupFakeSim(defStats, &proto.Glyphs{})
	lock := sim.Raid.Parties[0].Players[0].(*DemonologyWarlock)
	lock.Metamorphosis.SkipCastAndApplyEffects(sim, lock.CurrentTarget)

	immoAura := lock.GetSpell(core.ActionID{SpellID: 50589})
	attackTable := lock.AttackTables[lock.CurrentTarget.UnitIndex]
	checkDotTick(t, sim, immoAura.AOEDot(), attackTable, 2127.4521, 1.771230)
}

func TestMoltenCoreCastTime(t *testing.T) {
	sim := setupFakeSim(defStats, &proto.Glyphs{})
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
	sim := setupFakeSim(defStats, &proto.Glyphs{})
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
	sim := setupFakeSim(defStats, &proto.Glyphs{})
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

	var demonologyTalents = "-3312222300310212211-33202"
	var demonologyGlyphs = &proto.Glyphs{
		Prime1: int32(proto.WarlockPrimeGlyph_GlyphOfImmolate),
		Prime2: int32(proto.WarlockPrimeGlyph_GlyphOfCorruption),
		Prime3: int32(proto.WarlockPrimeGlyph_GlyphOfMetamorphosis),
		Major1: int32(proto.WarlockMajorGlyph_GlyphOfLifeTap),
		Major2: int32(proto.WarlockMajorGlyph_GlyphOfShadowBolt),
		Major3: int32(proto.WarlockMajorGlyph_GlyphOfSoulLink),
	}

	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:            proto.Class_ClassWarlock,
		Race:             proto.Race_RaceOrc,
		OtherRaces:       []proto.Race{proto.Race_RaceTroll, proto.Race_RaceGoblin, proto.Race_RaceHuman},
		GearSet:          core.GetGearSet("../../../ui/warlock/demonology/gear_sets", "p1"),
		Talents:          demonologyTalents,
		Glyphs:           demonologyGlyphs,
		Consumes:         fullConsumes,
		SpecOptions:      core.SpecOptionsCombo{Label: "Demonology Warlock", SpecOptions: defaultDemonologyWarlock},
		OtherSpecOptions: []core.SpecOptionsCombo{},
		Rotation:         core.GetAplRotation("../../../ui/warlock/demonology/apls", "default"),
		OtherRotations:   []core.RotationCombo{},
		ItemFilter:       itemFilter,
		StartingDistance: 25,
	}))
}
