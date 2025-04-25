package mage

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

var TalentTreeSizes = [3]int{21, 21, 19}

type Mage struct {
	core.Character

	Talents       *proto.MageTalents
	Options       *proto.MageOptions
	ArcaneOptions *proto.ArcaneMage_Options
	FireOptions   *proto.FireMage_Options
	FrostOptions  *proto.FrostMage_Options

	mirrorImage *MirrorImage
	// flameOrb     *FlameOrb
	// frostfireOrb *FrostfireOrb

	t12MirrorImage *T12MirrorImage
	t13ProcAura    *core.StatBuffAura

	arcaneMissilesTickSpell *core.Spell
	Combustion              *core.Spell
	Ignite                  *core.Spell
	LivingBomb              *core.Spell
	FireBlast               *core.Spell
	FlameOrbExplode         *core.Spell
	Flamestrike             *core.Spell
	FlamestrikeBW           *core.Spell
	FrostfireOrb            *core.Spell
	Pyroblast               *core.Spell
	SummonWaterElemental    *core.Spell
	IcyVeins                *core.Spell

	arcanePowerGCDmod *core.SpellMod

	arcaneMissilesProcAura *core.Aura
	arcanePotencyAura      *core.Aura
	arcanePowerAura        *core.Aura
	presenceOfMindAura     *core.Aura
	FingersOfFrostAura     *core.Aura

	arcaneMissileCritSnapshot float64
	brainFreezeProcChance     float64
	baseHotStreakProcChance   float64

	combustionDotEstimate int32

	ClassSpellScaling float64

	// Item sets
	T12_4pc *core.Aura
	T13_4pc *core.Aura
}

func (mage *Mage) GetCharacter() *core.Character {
	return &mage.Character
}

func (mage *Mage) GetMage() *Mage {
	return mage
}

func (mage *Mage) HasMajorGlyph(glyph proto.MageMajorGlyph) bool {
	return mage.HasGlyph(int32(glyph))
}
func (mage *Mage) HasMinorGlyph(glyph proto.MageMinorGlyph) bool {
	return mage.HasGlyph(int32(glyph))
}

func (mage *Mage) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.ArcaneBrilliance = true
}

func (mage *Mage) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (mage *Mage) ApplyTalents() {
	mage.ApplyArmorSpecializationEffect(stats.Intellect, proto.ArmorType_ArmorTypeCloth, 89744)

	// mage.ApplyArcaneTalents()
	// mage.ApplyFireTalents()
	// mage.ApplyFrostTalents()

	// mage.applyGlyphs()
}

func (mage *Mage) Initialize() {
	// mage.applyArmorSpells()
	// mage.registerArcaneBlastSpell()
	// mage.registerArcaneExplosionSpell()
	mage.registerArcaneMissilesSpell()
	// mage.registerBlizzardSpell()
	mage.registerDeepFreezeSpell()
	mage.registerFireballSpell()
	mage.registerFireBlastSpell()
	// mage.registerFlameOrbSpell()
	// mage.registerFlameOrbExplodeSpell()
	mage.registerFlamestrikeSpell()
	mage.registerFreezeSpell()
	// mage.registerFrostboltSpell()
	// mage.registerFrostfireOrbSpell()
	mage.registerIceLanceSpell()
	mage.registerScorchSpell()
	mage.registerLivingBombSpell()
	// mage.registerFrostfireBoltSpell()
	mage.registerEvocation()
	// mage.registerManaGemsCD()
	mage.registerMirrorImageCD()
	// mage.registerCombustionSpell()
	// mage.registerBlastWaveSpell()
	mage.registerDragonsBreathSpell()
	// mage.registerSummonWaterElementalCD()

	// mage.applyArcaneMissileProc()
}

// TODO: Fix this to work with the new talent system.
// func (mage *Mage) applyArcaneMissileProc() {
// 	if mage.Talents.HotStreak || mage.Talents.BrainFreeze > 0 {
// 		return
// 	}

// 	// Aura for when proc is successful
// 	mage.arcaneMissilesProcAura = mage.RegisterAura(core.Aura{
// 		Label:    "Arcane Missiles Proc",
// 		ActionID: core.ActionID{SpellID: 79683},
// 		Duration: time.Second * 20,
// 		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
// 			if spell.ClassSpellMask == MageSpellArcaneMissilesCast {
// 				aura.Deactivate(sim)
// 			}
// 		},
// 	})

// 	procChance := 0.4

// 	const MageSpellsArcaneMissilesNow = MageSpellArcaneBarrage | MageSpellArcaneBlast |
// 		MageSpellFireball | MageSpellFrostbolt | MageSpellFrostfireBolt | MageSpellFrostfireOrb

// 	// Listener for procs
// 	core.MakePermanent(mage.RegisterAura(core.Aura{
// 		Label: "Arcane Missiles Activation",
// 		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
// 			if spell.ClassSpellMask&MageSpellsArcaneMissilesNow == 0 {
// 				return
// 			}
// 			if sim.Proc(procChance, "Arcane Missiles") {
// 				mage.arcaneMissilesProcAura.Activate(sim)
// 			}
// 		},
// 	}))
// }

func (mage *Mage) Reset(sim *core.Simulation) {
	mage.arcaneMissileCritSnapshot = 0.0
}

func NewMage(character *core.Character, options *proto.Player, mageOptions *proto.MageOptions) *Mage {
	mage := &Mage{
		Character:         *character,
		Talents:           &proto.MageTalents{},
		Options:           mageOptions,
		ClassSpellScaling: core.GetClassSpellScalingCoefficient(proto.Class_ClassMage),
	}

	core.FillTalentsProto(mage.Talents.ProtoReflect(), options.TalentsString, TalentTreeSizes)

	mage.EnableManaBar()
	mage.SetDefaultSpellCritMultiplier(mage.SpellCritMultiplier(1.33, 0.0))

	mage.mirrorImage = mage.NewMirrorImage()
	// mage.flameOrb = mage.NewFlameOrb()
	// mage.frostfireOrb = mage.NewFrostfireOrb()

	if mage.CouldHaveSetBonus(ItemSetFirehawkRobesOfConflagration, 2) {
		mage.t12MirrorImage = mage.NewT12MirrorImage()
	}

	return mage
}

// func (mage *Mage) applyArmorSpells() {
// 	// Molten Armor
// 	// +3% spell crit, +5% with glyph
// 	critPercentToAdd := 3.0
// 	if mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfMoltenArmor) {
// 		critPercentToAdd = 5.0
// 	}

// 	mageArmorEffectCategory := "MageArmors"

// 	moltenArmor := mage.RegisterAura(core.Aura{
// 		Label:    "Molten Armor",
// 		ActionID: core.ActionID{SpellID: 30482},
// 		Duration: core.NeverExpires,
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			mage.AddStatDynamic(sim, stats.SpellCritPercent, critPercentToAdd)
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			mage.AddStatDynamic(sim, stats.SpellCritPercent, -critPercentToAdd)
// 		},
// 	})

// 	moltenArmor.NewExclusiveEffect(mageArmorEffectCategory, true, core.ExclusiveEffect{})

// 	mage.RegisterSpell(core.SpellConfig{
// 		ActionID:       core.ActionID{SpellID: 30482},
// 		SpellSchool:    core.SpellSchoolFire,
// 		Flags:          core.SpellFlagAPL,
// 		ClassSpellMask: MageSpellMoltenArmor,

// 		Cast: core.CastConfig{
// 			DefaultCast: core.Cast{
// 				GCD: core.GCDDefault,
// 			},
// 		},
// 		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
// 			return !moltenArmor.IsActive()
// 		},

// 		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
// 			moltenArmor.Activate(sim)
// 		},
// 	})

// 	// Mage Armor
// 	// Restores 3% of your max mana every 5 seconds (+20% affect with glyph)
// 	mageArmorManaMetric := mage.NewManaMetrics(core.ActionID{SpellID: 6117})
// 	hasGlyph := mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfMageArmor)
// 	manaRegenPer5Second := core.TernaryFloat64(hasGlyph, .036, 0.03)

// 	var pa *core.PendingAction
// 	mageArmor := mage.RegisterAura(core.Aura{
// 		ActionID: core.ActionID{SpellID: 6117},
// 		Label:    "Mage Armor",
// 		Duration: core.NeverExpires,
// 		OnGain: func(aura *core.Aura, sim *core.Simulation) {
// 			pa = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
// 				Period: time.Second * 5,
// 				OnAction: func(sim *core.Simulation) {
// 					mage.AddMana(sim, mage.MaxMana()*manaRegenPer5Second, mageArmorManaMetric)
// 				},
// 			})
// 		},
// 		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
// 			pa.Cancel(sim)
// 		},
// 	})

// 	mageArmor.NewExclusiveEffect(mageArmorEffectCategory, true, core.ExclusiveEffect{})

// 	mage.RegisterSpell(core.SpellConfig{
// 		ActionID:       core.ActionID{SpellID: 6117},
// 		SpellSchool:    core.SpellSchoolArcane,
// 		Flags:          core.SpellFlagAPL,
// 		ClassSpellMask: MageSpellMageArmor,

// 		Cast: core.CastConfig{
// 			DefaultCast: core.Cast{
// 				GCD: core.GCDDefault,
// 			},
// 		},
// 		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
// 			return !mageArmor.IsActive()
// 		},

// 		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
// 			mageArmor.Activate(sim)
// 		},
// 	})

// 	// Frost Armor
// 	// TODO:
// }

// Agent is a generic way to access underlying mage on any of the agents.
type MageAgent interface {
	GetMage() *Mage
}

// func (mage *Mage) hasChillEffect(spell *core.Spell) bool {
// 	return spell.ClassSpellMask&MageSpellChill > 0 || (spell.ClassSpellMask == MageSpellBlizzard && mage.Talents.IceShards > 0)
// }

const (
	MageSpellFlagNone      int64 = 0
	MageSpellArcaneBarrage int64 = 1 << iota
	MageSpellArcaneBlast
	MageSpellArcaneExplosion
	MageSpellArcanePower
	MageSpellArcaneMissilesCast
	MageSpellArcaneMissilesTick
	MageSpellBlastWave
	MageSpellBlizzard
	MageSpellConeOfCold
	MageSpellDeepFreeze
	MageSpellDragonsBreath
	MageSpellEvocation
	MageSpellFireBlast
	MageSpellFireball
	MageSpellFlamestrike
	MageSpellFlameOrb
	MageSpellFocusMagic
	MageSpellFreeze
	MageSpellFrostbolt
	MageSpellFrostfireBolt
	MageSpellFrostfireOrb
	MageSpellIceLance
	MageSpellIcyVeins
	MageSpellIgnite
	MageSpellLivingBombExplosion
	MageSpellLivingBombDot
	MageSpellManaGems
	MageSpellMirrorImage
	MageSpellPresenceOfMind
	MageSpellPyroblast
	MageSpellPyroblastDot
	MageSpellScorch
	MageSpellMoltenArmor
	MageSpellMageArmor
	MageSpellCombustion
	MageSpellCombustionApplication
	MageSpellLast
	MageSpellsAll        = MageSpellLast<<1 - 1
	MageSpellLivingBomb  = MageSpellLivingBombDot | MageSpellLivingBombExplosion
	MageSpellFireMastery = MageSpellLivingBombDot | MageSpellPyroblastDot | MageSpellCombustion // Ignite done manually in spell due to unique mechanic
	MageSpellFire        = MageSpellBlastWave | MageSpellCombustionApplication | MageSpellDragonsBreath | MageSpellFireball |
		MageSpellFireBlast | MageSpellFlameOrb | MageSpellFlamestrike | MageSpellFrostfireBolt | MageSpellIgnite |
		MageSpellLivingBomb | MageSpellPyroblast | MageSpellScorch
	MageSpellChill        = MageSpellFrostbolt | MageSpellFrostfireBolt
	MageSpellBrainFreeze  = MageSpellFireball | MageSpellFrostfireBolt
	MageSpellsAllDamaging = MageSpellArcaneBarrage | MageSpellArcaneBlast | MageSpellArcaneExplosion | MageSpellArcaneMissilesTick | MageSpellBlastWave | MageSpellBlizzard | MageSpellDeepFreeze |
		MageSpellDragonsBreath | MageSpellFireBlast | MageSpellFireball | MageSpellFlamestrike | MageSpellFlameOrb | MageSpellFrostbolt | MageSpellFrostfireBolt |
		MageSpellFrostfireOrb | MageSpellIceLance | MageSpellLivingBombExplosion | MageSpellLivingBombDot | MageSpellPyroblast | MageSpellPyroblastDot | MageSpellScorch
	MageSpellInstantCast = MageSpellArcaneBarrage | MageSpellArcaneMissilesCast | MageSpellArcaneMissilesTick | MageSpellFireBlast | MageSpellArcaneExplosion |
		MageSpellBlastWave | MageSpellCombustionApplication | MageSpellConeOfCold | MageSpellDeepFreeze | MageSpellDragonsBreath | MageSpellIceLance |
		MageSpellManaGems | MageSpellMirrorImage | MageSpellPresenceOfMind | MageSpellMoltenArmor | MageSpellMageArmor | MageSpellFlameOrb
	MageSpellExtraResult = MageSpellLivingBombExplosion | MageSpellArcaneMissilesTick | MageSpellBlizzard
)
