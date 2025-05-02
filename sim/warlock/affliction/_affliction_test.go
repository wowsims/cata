package affliction

import (
	"fmt"
	"testing"
	"time"

	_ "unsafe"

	_ "github.com/wowsims/mop/sim/common"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/simsignals"
	"github.com/wowsims/mop/sim/core/stats"
)

func init() {
	RegisterAfflictionWarlock()
}

const afflictionTalents = "223222003013321321-03-33"

func setupFakeSim(lockStats stats.Stats, talents string, glyphs *proto.Glyphs) *core.Simulation {
	var specOptions = &proto.Player_AfflictionWarlock{
		AfflictionWarlock: &proto.AfflictionWarlock{
			Options: &proto.AfflictionWarlock_Options{
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
					Rotation:      &proto.APLRotation{},
				}}, Buffs: &proto.PartyBuffs{}}},
		},
		Encounter: &proto.Encounter{
			Targets: []*proto.Target{
				{Name: "target", Level: 83, MobType: proto.MobType_MobTypeDemon},
			},
			ExecuteProportion_25: 0.2,
			ExecuteProportion_35: 0.2,
			ExecuteProportion_90: 0.2,
			Duration:             10,
		},
	}, simsignals.CreateSignals())

	sim.Options.Debug = true
	sim.Log = func(message string, vals ...interface{}) {
		fmt.Printf(fmt.Sprintf("[%0.1f] "+message+"\n", append([]interface{}{sim.CurrentTime.Seconds()}, vals...)...))
	}

	sim.Reset()
	lock := sim.Raid.Parties[0].Players[0].(*AfflictionWarlock)

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
	// this was tested before Blizzard fixed mastery flooring, true mastery in these tests was 1059
	// which was rounded down to 22%. So the "effective" mastery was ((22 / 1.625) - 8) * 179.280040
	stats.MasteryRating: 992.9356061538462,
}

func (lock *AfflictionWarlock) checkSpell(t *testing.T, sim *core.Simulation, spell *core.Spell, expected float64,
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

func (lock *AfflictionWarlock) checkSpellDamageRange(
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
	if dot.Spell.Flags.Matches(core.SpellFlagHauntSE) {
		mult *= attackTable.HauntSEDamageTakenMultiplier
	}

	if !core.WithinToleranceFloat64(expectedMult, mult, 0.000001) {
		t.Fatalf("Incorrect multiplier: Expected: %0.6f, Actual: %0.6f", expectedMult, mult)
	}

	delta := damageAfter - damageBefore
	if !core.WithinToleranceFloat64(expected, delta, 0.0001) {
		t.Fatalf("Incorrect damage applied: Expected: %0.4f, Actual: %0.4f", expected, delta)
	}
}

func TestDrainLifeBase(t *testing.T) {
	sim := setupFakeSim(defStats, afflictionTalents, &proto.Glyphs{})
	lock := sim.Raid.Parties[0].Players[0].(*AfflictionWarlock)
	drainLife := lock.GetSpell(core.ActionID{SpellID: 689})

	attackTable := lock.AttackTables[lock.CurrentTarget.UnitIndex]
	checkDotTick(t, sim, drainLife.CurDot(), attackTable, 1310.1845, 1.586000)
}

func TestDrainLife1SE(t *testing.T) {
	sim := setupFakeSim(defStats, afflictionTalents, &proto.Glyphs{})
	lock := sim.Raid.Parties[0].Players[0].(*AfflictionWarlock)
	drainLife := lock.GetSpell(core.ActionID{SpellID: 689})
	lock.ShadowEmbraceDebuffAura(lock.CurrentTarget).Activate(sim)
	lock.ShadowEmbraceDebuffAura(lock.CurrentTarget).SetStacks(sim, 1)

	attackTable := lock.AttackTables[lock.CurrentTarget.UnitIndex]
	checkDotTick(t, sim, drainLife.CurDot(), attackTable, 1375.6938, 1.665300)
}

func TestDrainLife3SE(t *testing.T) {
	sim := setupFakeSim(defStats, afflictionTalents, &proto.Glyphs{})
	lock := sim.Raid.Parties[0].Players[0].(*AfflictionWarlock)
	drainLife := lock.GetSpell(core.ActionID{SpellID: 689})
	lock.ShadowEmbraceDebuffAura(lock.CurrentTarget).Activate(sim)
	lock.ShadowEmbraceDebuffAura(lock.CurrentTarget).SetStacks(sim, 3)

	attackTable := lock.AttackTables[lock.CurrentTarget.UnitIndex]
	checkDotTick(t, sim, drainLife.CurDot(), attackTable, 1506.7122, 1.823900)
}

func TestDrainLife1SEHaunt(t *testing.T) {
	sim := setupFakeSim(defStats, afflictionTalents, &proto.Glyphs{Prime1: int32(proto.WarlockPrimeGlyph_GlyphOfHaunt)})
	lock := sim.Raid.Parties[0].Players[0].(*AfflictionWarlock)
	drainLife := lock.GetSpell(core.ActionID{SpellID: 689})
	lock.ShadowEmbraceDebuffAura(lock.CurrentTarget).Activate(sim)
	lock.ShadowEmbraceDebuffAura(lock.CurrentTarget).SetStacks(sim, 1)
	lock.HauntDebuffAuras[lock.CurrentTarget.UnitIndex].Activate(sim)

	attackTable := lock.AttackTables[lock.CurrentTarget.UnitIndex]
	// true multiplier is 2.171218, but siphon soul is not accounted for
	checkDotTick(t, sim, drainLife.CurDot(), attackTable, 1793.6295, 2.048319)
}

func TestDrainLife3SEHauntSiphonSoul(t *testing.T) {
	sim := setupFakeSim(defStats, afflictionTalents, &proto.Glyphs{Prime1: int32(proto.WarlockPrimeGlyph_GlyphOfHaunt)})
	lock := sim.Raid.Parties[0].Players[0].(*AfflictionWarlock)
	drainLife := lock.GetSpell(core.ActionID{SpellID: 689})
	lock.ShadowEmbraceDebuffAura(lock.CurrentTarget).Activate(sim)
	lock.ShadowEmbraceDebuffAura(lock.CurrentTarget).SetStacks(sim, 3)
	lock.HauntDebuffAuras[lock.CurrentTarget.UnitIndex].Activate(sim)
	lock.Corruption.CurDot().Apply(sim)
	lock.BaneOfDoom.CurDot().Apply(sim)

	attackTable := lock.AttackTables[lock.CurrentTarget.UnitIndex]
	// true multiplier is 2.647208, but siphon soul is not accounted for
	checkDotTick(t, sim, drainLife.CurDot(), attackTable, 2186.8421, 2.243397)
}

func TestDrainLife3SEHauntSiphonSoulDemonSoul(t *testing.T) {
	sim := setupFakeSim(defStats, afflictionTalents, &proto.Glyphs{Prime1: int32(proto.WarlockPrimeGlyph_GlyphOfHaunt)})
	lock := sim.Raid.Parties[0].Players[0].(*AfflictionWarlock)
	drainLife := lock.GetSpell(core.ActionID{SpellID: 689})
	lock.ShadowEmbraceDebuffAura(lock.CurrentTarget).Activate(sim)
	lock.ShadowEmbraceDebuffAura(lock.CurrentTarget).SetStacks(sim, 3)
	lock.HauntDebuffAuras[lock.CurrentTarget.UnitIndex].Activate(sim)
	lock.Corruption.CurDot().Apply(sim)
	lock.BaneOfDoom.CurDot().Apply(sim)
	lock.GetAura("Demon Soul: Felhunter").Activate(sim)

	attackTable := lock.AttackTables[lock.CurrentTarget.UnitIndex]
	// true multiplier is 3.081177, but siphon soul is not accounted for
	checkDotTick(t, sim, drainLife.CurDot(), attackTable, 2545.3408, 2.611167)
}

func TestDrainSoulBase(t *testing.T) {
	sim := setupFakeSim(defStats, afflictionTalents, &proto.Glyphs{})
	lock := sim.Raid.Parties[0].Players[0].(*AfflictionWarlock)
	drainSoul := lock.GetSpell(core.ActionID{SpellID: 1120})

	attackTable := lock.AttackTables[lock.CurrentTarget.UnitIndex]
	checkDotTick(t, sim, drainSoul.CurDot(), attackTable, 3578.8642, 1.586000)
}

func TestDrainSoul1SE(t *testing.T) {
	sim := setupFakeSim(defStats, afflictionTalents, &proto.Glyphs{})
	lock := sim.Raid.Parties[0].Players[0].(*AfflictionWarlock)
	drainSoul := lock.GetSpell(core.ActionID{SpellID: 1120})
	lock.ShadowEmbraceDebuffAura(lock.CurrentTarget).Activate(sim)
	lock.ShadowEmbraceDebuffAura(lock.CurrentTarget).SetStacks(sim, 1)

	attackTable := lock.AttackTables[lock.CurrentTarget.UnitIndex]
	checkDotTick(t, sim, drainSoul.CurDot(), attackTable, 3757.8074, 1.665300)
}

func TestDrainSoul3SE(t *testing.T) {
	sim := setupFakeSim(defStats, afflictionTalents, &proto.Glyphs{})
	lock := sim.Raid.Parties[0].Players[0].(*AfflictionWarlock)
	drainSoul := lock.GetSpell(core.ActionID{SpellID: 1120})
	lock.ShadowEmbraceDebuffAura(lock.CurrentTarget).Activate(sim)
	lock.ShadowEmbraceDebuffAura(lock.CurrentTarget).SetStacks(sim, 3)

	attackTable := lock.AttackTables[lock.CurrentTarget.UnitIndex]
	checkDotTick(t, sim, drainSoul.CurDot(), attackTable, 4115.6938, 1.823900)
}

//go:linkname sim_advance github.com/wowsims/mop/sim/core.(*Simulation).advance
func sim_advance(*core.Simulation, time.Duration)

func TestDrainSoulExecute3SE(t *testing.T) {
	sim := setupFakeSim(defStats, afflictionTalents, &proto.Glyphs{})
	sim_advance(sim, 9*time.Second)
	lock := sim.Raid.Parties[0].Players[0].(*AfflictionWarlock)
	drainSoul := lock.GetSpell(core.ActionID{SpellID: 1120})
	lock.ShadowEmbraceDebuffAura(lock.CurrentTarget).Activate(sim)
	lock.ShadowEmbraceDebuffAura(lock.CurrentTarget).SetStacks(sim, 3)

	attackTable := lock.AttackTables[lock.CurrentTarget.UnitIndex]
	checkDotTick(t, sim, drainSoul.CurDot(), attackTable, 8725.2709, 3.866668)
}

func TestDrainSoul3SEHauntSiphonSoul(t *testing.T) {
	sim := setupFakeSim(defStats, afflictionTalents, &proto.Glyphs{Prime1: int32(proto.WarlockPrimeGlyph_GlyphOfHaunt)})
	lock := sim.Raid.Parties[0].Players[0].(*AfflictionWarlock)
	drainSoul := lock.GetSpell(core.ActionID{SpellID: 1120})
	lock.ShadowEmbraceDebuffAura(lock.CurrentTarget).Activate(sim)
	lock.ShadowEmbraceDebuffAura(lock.CurrentTarget).SetStacks(sim, 3)
	lock.HauntDebuffAuras[lock.CurrentTarget.UnitIndex].Activate(sim)
	lock.Corruption.CurDot().Apply(sim)
	lock.BaneOfDoom.CurDot().Apply(sim)

	attackTable := lock.AttackTables[lock.CurrentTarget.UnitIndex]
	// true multiplier is 2.647208, but siphon soul is not accounted for
	checkDotTick(t, sim, drainSoul.CurDot(), attackTable, 5973.5180, 2.243397)
}

func TestDrainSoulExecute3SEHauntSiphonSoul(t *testing.T) {
	sim := setupFakeSim(defStats, afflictionTalents, &proto.Glyphs{Prime1: int32(proto.WarlockPrimeGlyph_GlyphOfHaunt)})
	sim_advance(sim, 9*time.Second)
	lock := sim.Raid.Parties[0].Players[0].(*AfflictionWarlock)
	drainSoul := lock.GetSpell(core.ActionID{SpellID: 1120})
	lock.ShadowEmbraceDebuffAura(lock.CurrentTarget).Activate(sim)
	lock.ShadowEmbraceDebuffAura(lock.CurrentTarget).SetStacks(sim, 3)
	lock.HauntDebuffAuras[lock.CurrentTarget.UnitIndex].Activate(sim)
	lock.Corruption.CurDot().Apply(sim)
	lock.BaneOfDoom.CurDot().Apply(sim)

	attackTable := lock.AttackTables[lock.CurrentTarget.UnitIndex]
	// true multiplier is 5.612082, but siphon soul is not accounted for
	checkDotTick(t, sim, drainSoul.CurDot(), attackTable, 12663.8582, 4.756002)
}

func TestDrainSoulExecute3SEHauntSiphonSoulDemonSoul(t *testing.T) {
	sim := setupFakeSim(defStats, afflictionTalents, &proto.Glyphs{Prime1: int32(proto.WarlockPrimeGlyph_GlyphOfHaunt)})
	sim_advance(sim, 9*time.Second)
	lock := sim.Raid.Parties[0].Players[0].(*AfflictionWarlock)
	drainSoul := lock.GetSpell(core.ActionID{SpellID: 1120})
	lock.ShadowEmbraceDebuffAura(lock.CurrentTarget).Activate(sim)
	lock.ShadowEmbraceDebuffAura(lock.CurrentTarget).SetStacks(sim, 3)
	lock.HauntDebuffAuras[lock.CurrentTarget.UnitIndex].Activate(sim)
	lock.Corruption.CurDot().Apply(sim)
	lock.BaneOfDoom.CurDot().Apply(sim)
	lock.GetAura("Demon Soul: Felhunter").Activate(sim)

	attackTable := lock.AttackTables[lock.CurrentTarget.UnitIndex]
	// true multiplier is 6.532095, but siphon soul is not accounted for
	checkDotTick(t, sim, drainSoul.CurDot(), attackTable, 14739.9005, 5.535674)
}

func checkTicks(t *testing.T, dot *core.Dot, msg string, expected int32) {
	if dot.RemainingTicks() != expected {
		t.Helper()
		t.Fatalf("%s: Expected: %v, Actual: %v", msg, expected, dot.RemainingTicks())
	}
}

func TestCorruptionHasteCap(t *testing.T) {
	sim := setupFakeSim(defStats, afflictionTalents, &proto.Glyphs{})
	lock := sim.Raid.Parties[0].Players[0].(*AfflictionWarlock)
	lock.Unit.MultiplyCastSpeed(1 + 0.05) // 5% haste buff
	lock.Unit.MultiplyCastSpeed(1 + 0.03) // dark intent
	lock.AddStatsDynamic(sim, stats.Stats{
		stats.HasteRating: 2588,
	})
	unstableAff := lock.UnstableAffliction
	unstableAffDot := unstableAff.CurDot()

	unstableAff.SkipCastAndApplyEffects(sim, lock.CurrentTarget)
	checkTicks(t, unstableAffDot, "Incorrect tick count for UA at 2588 haste", 6)
	unstableAffDot.Deactivate(sim)

	lock.AddStatsDynamic(sim, stats.Stats{
		stats.HasteRating: 1,
	})
	unstableAff.SkipCastAndApplyEffects(sim, lock.CurrentTarget)
	checkTicks(t, unstableAffDot, "Incorrect tick count for UA at 2589 haste", 7)
}

func TestAffliction(t *testing.T) {
	var defaultAfflictionWarlock = &proto.Player_AfflictionWarlock{
		AfflictionWarlock: &proto.AfflictionWarlock{
			Options: &proto.AfflictionWarlock_Options{
				ClassOptions: &proto.WarlockOptions{
					Summon:       proto.WarlockOptions_Felhunter,
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

	var fullConsumesSpec = &proto.ConsumesSpec{
		FlaskId:     58086, // Flask of the Draconic Mind
		FoodId:      62671, // Severed Sagefish Head
		PotId:       58091, // Volcanic Potion
		ExplosiveId: 89637, // Big Daddy Explosive
		TinkerId:    82174, // Synapse Springs
	}

	var afflictionGlyphs = &proto.Glyphs{
		Prime1: int32(proto.WarlockPrimeGlyph_GlyphOfHaunt),
		Prime2: int32(proto.WarlockPrimeGlyph_GlyphOfUnstableAffliction),
		Prime3: int32(proto.WarlockPrimeGlyph_GlyphOfCorruption),
		Major1: int32(proto.WarlockMajorGlyph_GlyphOfFelhunter),
		Major2: int32(proto.WarlockMajorGlyph_GlyphOfShadowBolt),
		Major3: int32(proto.WarlockMajorGlyph_GlyphOfSoulSwap),
	}

	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:            proto.Class_ClassWarlock,
		Race:             proto.Race_RaceOrc,
		OtherRaces:       []proto.Race{proto.Race_RaceTroll, proto.Race_RaceGoblin, proto.Race_RaceHuman},
		GearSet:          core.GetGearSet("../../../ui/warlock/affliction/gear_sets", "p4"),
		Talents:          afflictionTalents,
		Glyphs:           afflictionGlyphs,
		Consumables:      fullConsumesSpec,
		SpecOptions:      core.SpecOptionsCombo{Label: "Affliction Warlock", SpecOptions: defaultAfflictionWarlock},
		OtherSpecOptions: []core.SpecOptionsCombo{},
		Rotation:         core.GetAplRotation("../../../ui/warlock/affliction/apls", "default"),
		OtherRotations:   []core.RotationCombo{},
		ItemFilter:       itemFilter,
		StartingDistance: 25,
	}))
}
