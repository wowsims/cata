package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type Mage struct {
	core.Character

	Talents       *proto.MageTalents
	Options       *proto.MageOptions
	ArcaneOptions *proto.ArcaneMage_Options
	FireOptions   *proto.FireMage_Options
	FrostOptions  *proto.FrostMage_Options

	mirrorImage *MirrorImage

	arcaneMissilesTickSpell *core.Spell
	Combustion              *core.Spell
	Ignite                  *core.Spell
	LivingBomb              *core.Spell
	NetherTempest           *core.Spell
	FireBlast               *core.Spell
	FlameOrbExplode         *core.Spell
	Flamestrike             *core.Spell
	FlamestrikeBW           *core.Spell
	FrostfireOrb            *core.Spell
	Pyroblast               *core.Spell
	SummonWaterElemental    *core.Spell
	IcyVeins                *core.Spell
	Icicle                  *core.Spell

	arcanePowerGCDmod *core.SpellMod

	arcaneMissilesProcAura *core.Aura
	arcanePotencyAura      *core.Aura
	arcanePowerAura        *core.Aura
	invocationAura         *core.Aura
	runeOfPowerAura        *core.Aura
	presenceOfMindAura     *core.Aura
	FingersOfFrostAura     *core.Aura
	BrainFreezeAura        *core.Aura
	IcyVeinsAura           *core.Aura
	iceFloesfAura          *core.Aura
	IciclesAura            *core.Aura
	FrostBombAuras         core.AuraArray

	arcaneMissileCritSnapshot float64
	brainFreezeProcChance     float64
	baseHotStreakProcChance   float64

	combustionDotEstimate int32

	ClassSpellScaling float64
	Icicles           []float64

	// Item sets
	T12_4pc *core.Aura
	T13_4pc *core.Aura
	T14_4pc *core.Aura
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

func (mage *Mage) GetFrostMasteryBonus() float64 {
	return (.16 + 0.02*mage.GetMasteryPoints())
}

func (mage *Mage) Initialize() {
	mage.applyArmorSpells()
	mage.applyGlyphs()
	mage.ApplyMastery()
	// mage.registerArcaneBlastSpell()
	// mage.registerArcaneExplosionSpell()
	mage.registerArcaneMissilesSpell()
	mage.registerBlizzardSpell()
	mage.registerConeOfColdSpell()
	mage.registerDeepFreezeSpell()
	mage.registerFireballSpell()
	mage.registerFireBlastSpell()
	mage.registerFlamestrikeSpell()
	mage.registerIceLanceSpell()
	mage.registerScorchSpell()
	mage.registerLivingBombSpell()
	mage.registerNetherTempestSpell()
	mage.registerFrostfireBoltSpell()
	mage.registerEvocation()
	// mage.registerManaGemsCD()
	mage.registerMirrorImageCD()
	// mage.registerCombustionSpell()
	// mage.registerBlastWaveSpell()
	mage.registerDragonsBreathSpell()
	mage.registerFrostBombSpell()
	mage.registerfrostNovaSpell()
	mage.registerIceLanceSpell()
	mage.registerIcyVeinsCD()
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

	core.FillTalentsProto(mage.Talents.ProtoReflect(), options.TalentsString)

	mage.EnableManaBar()

	mage.Icicles = make([]float64, 0)
	mage.mirrorImage = mage.NewMirrorImage()
	// mage.flameOrb = mage.NewFlameOrb()
	// mage.frostfireOrb = mage.NewFrostfireOrb()

	return mage
}

func (mage *Mage) applyArmorSpells() {

	mageArmorEffectCategory := "MageArmors"

	moltenArmor := mage.RegisterAura(core.Aura{
		Label:    "Molten Armor",
		ActionID: core.ActionID{SpellID: 30482},
		Duration: core.NeverExpires,
	}).AttachStatBuff(stats.SpellCritPercent, 5)

	moltenArmor.NewExclusiveEffect(mageArmorEffectCategory, true, core.ExclusiveEffect{})

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 30482},
		SpellSchool:    core.SpellSchoolFire,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellMoltenArmor,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 3,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !moltenArmor.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			moltenArmor.Activate(sim)
		},
	})

	mageArmor := mage.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 6117},
		Label:    "Mage Armor",
		Duration: core.NeverExpires,
	}).AttachStatBuff(stats.MasteryRating, 3000.0)

	mageArmor.NewExclusiveEffect(mageArmorEffectCategory, true, core.ExclusiveEffect{})

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 6117},
		SpellSchool:    core.SpellSchoolArcane,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellMageArmor,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 3,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !mageArmor.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mageArmor.Activate(sim)
		},
	})

	frostArmor := mage.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 7302},
		Label:    "Frost Armor",
		Duration: core.NeverExpires,
	}).AttachMultiplyCastSpeed(1.07)

	frostArmor.NewExclusiveEffect(mageArmorEffectCategory, true, core.ExclusiveEffect{})

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 7302},
		SpellSchool:    core.SpellSchoolFrost,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellFrostArmor,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 3,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !frostArmor.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			frostArmor.Activate(sim)
		},
	})
}

// Agent is a generic way to access underlying mage on any of the agents.
type MageAgent interface {
	GetMage() *Mage
}

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
	MageSpellFrostArmor
	MageSpellFrostbolt
	MageSpellFrostBomb
	MageSpellFrostBombExplosion
	MageSpellFrostfireBolt
	MageSpellFrostfireOrb
	MageSpellFrostNova
	MageSpellFrozenOrb
	MageSpellFrozenOrbTick
	MageSpellIcicle
	MageSpellIceFloes
	MageSpellIceLance
	MageSpellIcyVeins
	MageSpellIgnite
	MageSpellLivingBombExplosion
	MageSpellLivingBombDot
	MageSpellMageArmor
	MageSpellManaGems
	MageSpellMirrorImage
	MageSpellMoltenArmor
	MageSpellNetherTempest
	MageSpellPresenceOfMind
	MageSpellPyroblast
	MageSpellPyroblastDot
	MagespellRuneOfPower
	MageSpellScorch
	MageSpellCombustion
	MageSpellCombustionApplication
	MageWaterElementalSpellWaterBolt
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
		MageSpellManaGems | MageSpellMirrorImage | MageSpellPresenceOfMind | MageSpellFlameOrb
	MageSpellExtraResult = MageSpellLivingBombExplosion | MageSpellArcaneMissilesTick | MageSpellBlizzard
)
