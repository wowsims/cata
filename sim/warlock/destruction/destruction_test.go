package destruction

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
	RegisterDestructionWarlock()
}

const minimalDestroTalents = "--3"
const destroTalents = "003-23002-3320202312230310211"

var immoGlyph = proto.Glyphs{Prime1: int32(proto.WarlockPrimeGlyph_GlyphOfImmolate)}

func setupFakeSim(lockStats stats.Stats, talents string, glyphs *proto.Glyphs) *core.Simulation {
	var specOptions = &proto.Player_DestructionWarlock{
		DestructionWarlock: &proto.DestructionWarlock{
			Options: &proto.DestructionWarlock_Options{
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
					TalentsString: talents,
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
	lock := sim.Raid.Parties[0].Players[0].(*DestructionWarlock)

	lock.AddStatsDynamic(sim, stats.Stats{
		stats.SpellPower: -lock.GetStat(stats.SpellPower),
	})
	lock.AddStatsDynamic(sim, lockStats)

	return sim
}

var defStats = stats.Stats{
	stats.SpellHitPercent:  17,
	stats.SpellCritPercent: -100,
	stats.SpellPower:       5766,
	stats.MasteryRating:    1059,
}

var defHastedStats = stats.Stats{
	stats.SpellHitPercent:  17,
	stats.SpellCritPercent: -100,
	stats.SpellPower:       5766,
	stats.MasteryRating:    1059,
	stats.HasteRating:      325,
}

func (lock *DestructionWarlock) checkSpell(t *testing.T, sim *core.Simulation, spell *core.Spell, expected float64,
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

func (lock *DestructionWarlock) checkSpellDamageRange(
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

func TestImmolateNonPeriodic(t *testing.T) {
	sim := setupFakeSim(defStats, destroTalents, &proto.Glyphs{})
	lock := sim.Raid.Parties[0].Players[0].(*DestructionWarlock)

	lock.checkSpell(t, sim, lock.Immolate, 3446.4580, 1.781616)
}

func TestImmolateDoTBase(t *testing.T) {
	sim := setupFakeSim(defStats, destroTalents, &proto.Glyphs{})
	lock := sim.Raid.Parties[0].Players[0].(*DestructionWarlock)

	attackTable := lock.AttackTables[lock.CurrentTarget.UnitIndex]
	checkDotTick(t, sim, lock.Immolate.CurDot(), attackTable, 2560.6833, 1.781616)
}

func TestImmolateDoTGlyphed(t *testing.T) {
	sim := setupFakeSim(defStats, destroTalents, &proto.Glyphs{Prime1: int32(proto.WarlockPrimeGlyph_GlyphOfImmolate)})
	lock := sim.Raid.Parties[0].Players[0].(*DestructionWarlock)

	attackTable := lock.AttackTables[lock.CurrentTarget.UnitIndex]
	checkDotTick(t, sim, lock.Immolate.CurDot(), attackTable, 2774.0736, 1.930084)
}

func TestIncinerateBase(t *testing.T) {
	sim := setupFakeSim(defStats, destroTalents, &proto.Glyphs{})
	lock := sim.Raid.Parties[0].Players[0].(*DestructionWarlock)
	incinerate := lock.GetSpell(core.ActionID{SpellID: 29722})

	lock.checkSpellDamageRange(t, sim, incinerate, 6016.0539, 6153.5921, 1.662842)
}

func TestIncinerateGlyphed(t *testing.T) {
	sim := setupFakeSim(defStats, destroTalents, &proto.Glyphs{Prime1: int32(proto.WarlockPrimeGlyph_GlyphOfIncinerate)})
	lock := sim.Raid.Parties[0].Players[0].(*DestructionWarlock)
	incinerate := lock.GetSpell(core.ActionID{SpellID: 29722})

	lock.checkSpellDamageRange(t, sim, incinerate, 6284.6278, 6428.3061, 1.737076)
}

func TestIncinerateImmoCoE(t *testing.T) {
	var testStats = stats.Stats{
		stats.SpellHitPercent:  17,
		stats.SpellCritPercent: -100,
		stats.SpellPower:       5073,
		stats.MasteryRating:    832,
	}
	sim := setupFakeSim(testStats, destroTalents, &proto.Glyphs{})
	lock := sim.Raid.Parties[0].Players[0].(*DestructionWarlock)
	incinerate := lock.GetSpell(core.ActionID{SpellID: 29722})
	lock.Immolate.CurDot().Apply(sim)
	lock.CurseOfElementsAuras.Get(lock.CurrentTarget).Activate(sim)

	lock.checkSpellDamageRange(t, sim, incinerate, 7704.7629, 7901.1877, 1.884747)
}

// TODO: incinerate test with immolation

func TestConflagrateBase(t *testing.T) {
	sim := setupFakeSim(defStats, minimalDestroTalents, &proto.Glyphs{})
	lock := sim.Raid.Parties[0].Players[0].(*DestructionWarlock)
	lock.Immolate.Dot(lock.CurrentTarget).Apply(sim)

	// these values are slightly off to ingame, probably different rounding
	lock.checkSpellDamageRange(t, sim, lock.Conflagrate, 6401.7083, 6401.7084, 1.484680)
}

func TestConflagrateBaseGlyphed(t *testing.T) {
	sim := setupFakeSim(defStats, minimalDestroTalents, &immoGlyph)
	lock := sim.Raid.Parties[0].Players[0].(*DestructionWarlock)
	lock.Immolate.Dot(lock.CurrentTarget).Apply(sim)

	// these values are slightly off to ingame, probably different rounding
	lock.checkSpellDamageRange(t, sim, lock.Conflagrate, 7041.8792, 7041.8793, 1.633148)
}

func TestConflagrateTalented(t *testing.T) {
	sim := setupFakeSim(defStats, destroTalents, &proto.Glyphs{})
	lock := sim.Raid.Parties[0].Players[0].(*DestructionWarlock)
	lock.Immolate.Dot(lock.CurrentTarget).Apply(sim)

	// these values are slightly off to ingame, probably different rounding
	lock.checkSpellDamageRange(t, sim, lock.Conflagrate, 7682.0500, 7682.0501, 1.781616)
}

func TestConflagrateTalentedGlyphed(t *testing.T) {
	sim := setupFakeSim(defStats, destroTalents, &immoGlyph)
	lock := sim.Raid.Parties[0].Players[0].(*DestructionWarlock)
	lock.Immolate.Dot(lock.CurrentTarget).Apply(sim)

	// these values are slightly off to ingame, probably different rounding
	lock.checkSpellDamageRange(t, sim, lock.Conflagrate, 8322.2208, 8322.2209, 1.930084)
}

func TestConflagrateTalentedGlyphedHasted(t *testing.T) {
	sim := setupFakeSim(defHastedStats, destroTalents, &immoGlyph)
	lock := sim.Raid.Parties[0].Players[0].(*DestructionWarlock)
	lock.Immolate.Dot(lock.CurrentTarget).Apply(sim)

	// these values are slightly off to ingame, probably different rounding
	lock.checkSpellDamageRange(t, sim, lock.Conflagrate, 8322.2208, 8322.2209, 1.930084)
}

func TestBackdraftCastTime(t *testing.T) {
	sim := setupFakeSim(defStats, destroTalents, &proto.Glyphs{})
	lock := sim.Raid.Parties[0].Players[0].(*DestructionWarlock)
	lock.GetAura("Backdraft").Activate(sim)
	lock.GetAura("Backdraft").SetStacks(sim, 3)
	chaosBolt := lock.GetSpell(core.ActionID{SpellID: 50796})

	expectedCasttime := 1400 * time.Millisecond
	if chaosBolt.CastTime() != expectedCasttime {
		t.Fatalf("Incorrect casttime: Expected: %v, Actual: %v", expectedCasttime, chaosBolt.CastTime())
	}

	expectedGCD := 1500 * time.Millisecond
	if chaosBolt.DefaultCast.GCD != expectedGCD {
		t.Fatalf("Incorrect GCD: Expected: %v, Actual: %v", expectedGCD, chaosBolt.DefaultCast.GCD)
	}
}

func checkTicks(t *testing.T, dot *core.Dot, msg string, expected int32) {
	if dot.RemainingTicks() != expected {
		t.Helper()
		t.Fatalf("%s: Expected: %v, Actual: %v", msg, expected, dot.RemainingTicks())
	}
}

func TestImmolateHasteCap(t *testing.T) {
	sim := setupFakeSim(defStats, destroTalents, &immoGlyph)
	lock := sim.Raid.Parties[0].Players[0].(*DestructionWarlock)
	lock.Unit.MultiplyCastSpeed(1 + 0.05) // 5% haste buff
	lock.Unit.MultiplyCastSpeed(1 + 0.03) // dark intent
	lock.AddStatsDynamic(sim, stats.Stats{
		stats.HasteRating: 2588,
	})
	immolate := lock.Immolate
	immolateDot := lock.Immolate.CurDot()

	immolate.SkipCastAndApplyEffects(sim, lock.CurrentTarget)
	checkTicks(t, immolateDot, "Incorrect tick count for immolate at 2588 haste", 6)
	immolateDot.Deactivate(sim)

	lock.AddStatsDynamic(sim, stats.Stats{
		stats.HasteRating: 1,
	})
	immolate.SkipCastAndApplyEffects(sim, lock.CurrentTarget)
	checkTicks(t, immolateDot, "Incorrect tick count for immolate at 2589 haste", 7)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func TestDestruction(t *testing.T) {
	var defaultDestructionWarlock = &proto.Player_DestructionWarlock{
		DestructionWarlock: &proto.DestructionWarlock{
			Options: &proto.DestructionWarlock_Options{
				ClassOptions: &proto.WarlockOptions{
					Summon:       proto.WarlockOptions_Imp,
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

	var destructionTalents = "003-03202-3320202312201312211"
	var destructionGlyphs = &proto.Glyphs{
		Prime1: int32(proto.WarlockPrimeGlyph_GlyphOfConflagrate),
		Prime2: int32(proto.WarlockPrimeGlyph_GlyphOfImmolate),
		Prime3: int32(proto.WarlockPrimeGlyph_GlyphOfImp),
		Major1: int32(proto.WarlockMajorGlyph_GlyphOfLifeTap),
		Major3: int32(proto.WarlockMajorGlyph_GlyphOfSoulLink),
		Major2: int32(proto.WarlockMajorGlyph_GlyphOfHealthstone),
	}

	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:            proto.Class_ClassWarlock,
		Race:             proto.Race_RaceOrc,
		OtherRaces:       []proto.Race{proto.Race_RaceTroll, proto.Race_RaceGoblin, proto.Race_RaceHuman},
		GearSet:          core.GetGearSet("../../../ui/warlock/destruction/gear_sets", "p4"),
		Talents:          destructionTalents,
		Glyphs:           destructionGlyphs,
		Consumes:         fullConsumes,
		SpecOptions:      core.SpecOptionsCombo{Label: "Destruction Warlock", SpecOptions: defaultDestructionWarlock},
		OtherSpecOptions: []core.SpecOptionsCombo{},
		Rotation:         core.GetAplRotation("../../../ui/warlock/destruction/apls", "default"),
		ItemFilter:       itemFilter,
		StartingDistance: 25,
	}))
}
