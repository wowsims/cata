package mage

import (
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

	mirrorImages []*MirrorImage

	Combustion           *core.Spell
	Ignite               *core.Spell
	LivingBomb           *core.Spell
	NetherTempest        *core.Spell
	FireBlast            *core.Spell
	FlameOrbExplode      *core.Spell
	Flamestrike          *core.Spell
	FlamestrikeBW        *core.Spell
	FrostfireOrb         *core.Spell
	Pyroblast            *core.Spell
	SummonWaterElemental *core.Spell
	SummonMirrorImages   *core.Spell
	IcyVeins             *core.Spell
	Icicle               *core.Spell

	InvocationAura     *core.Aura
	RuneOfPowerAura    *core.Aura
	PresenceOfMindAura *core.Aura
	FingersOfFrostAura *core.Aura
	BrainFreezeAura    *core.Aura
	IcyVeinsAura       *core.Aura
	IceFloesAura       *core.Aura
	IciclesAura        *core.Aura
	ArcaneChargesAura  *core.Aura

	arcaneMissileCritSnapshot float64
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

func (mage *Mage) ProcFingersOfFrost(sim *core.Simulation, spell *core.Spell) {
	if mage.FingersOfFrostAura == nil {
		return
	}
	if spell.Matches(MageSpellFrostbolt | MageSpellFrostfireBolt) {
		if sim.Proc(0.15, "FingersOfFrostProc") {
			mage.FingersOfFrostAura.Activate(sim)
			mage.FingersOfFrostAura.AddStack(sim)
		}
	} else if spell.Matches(MageSpellBlizzard) {
		if sim.Proc(0.05, "FingersOfFrostBlizzardProc") {
			mage.FingersOfFrostAura.Activate(sim)
			mage.FingersOfFrostAura.AddStack(sim)
		}
	}
}

func (mage *Mage) Initialize() {
	mage.registerGlyphs()
	mage.registerPassives()
	mage.registerSpells()
}

func (mage *Mage) registerPassives() {
	mage.ApplyArmorSpecializationEffect(stats.Intellect, proto.ArmorType_ArmorTypeCloth, 89744)

	mage.registerMastery()
}

func (mage *Mage) registerSpells() {
	mage.registerArmorSpells()

	// mage.registerArcaneBlastSpell()
	// mage.registerArcaneExplosionSpell()
	mage.registerBlizzardSpell()
	mage.registerConeOfColdSpell()
	mage.registerDeepFreezeSpell()
	mage.registerFireballSpell()
	mage.registerFireBlastSpell()
	mage.registerFlamestrikeSpell()
	mage.registerIceLanceSpell()
	mage.registerScorchSpell()
	mage.registerFrostfireBoltSpell()
	mage.registerEvocation()
	mage.registerManaGems()
	mage.registerMirrorImageCD()
	// mage.registerCombustionSpell()
	// mage.registerBlastWaveSpell()
	mage.registerDragonsBreathSpell()
	mage.registerfrostNovaSpell()
	mage.registerIceLanceSpell()
	mage.registerIcyVeinsCD()
}

func (mage *Mage) registerMastery() {
	mage.registerFrostMastery()
}

func (mage *Mage) Reset(sim *core.Simulation) {
	mage.arcaneMissileCritSnapshot = 0.0
	mage.Icicles = make([]float64, 0)
}

func NewMage(character *core.Character, options *proto.Player, mageOptions *proto.MageOptions) *Mage {
	mage := &Mage{
		Character:         *character,
		Talents:           &proto.MageTalents{},
		Options:           mageOptions,
		ClassSpellScaling: core.GetClassSpellScalingCoefficient(proto.Class_ClassMage),
	}

	core.FillTalentsProto(mage.Talents.ProtoReflect(), options.TalentsString)

	mage.mirrorImages = []*MirrorImage{mage.NewMirrorImage(), mage.NewMirrorImage(), mage.NewMirrorImage()}
	mage.EnableManaBar()
	// Nether Attunement
	// https://www.wowhead.com/mop-classic/spell=117957/nether-attunement
	mage.HasteEffectsManaRegen()

	mage.Icicles = make([]float64, 0)
	// mage.mirrorImage = mage.NewMirrorImage()
	// mage.flameOrb = mage.NewFlameOrb()
	// mage.frostfireOrb = mage.NewFrostfireOrb()

	return mage
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
	MageMirrorImageSpellArcaneBlast
	MageWaterElementalSpellWaterBolt
	MageSpellLast
	MageSpellsAll        = MageSpellLast<<1 - 1
	MageSpellLivingBomb  = MageSpellLivingBombDot | MageSpellLivingBombExplosion
	MageSpellFireMastery = MageSpellLivingBombDot | MageSpellPyroblastDot | MageSpellCombustion // Ignite done manually in spell due to unique mechanic
	MageSpellFire        = MageSpellBlastWave | MageSpellCombustionApplication | MageSpellDragonsBreath | MageSpellFireball |
		MageSpellFireBlast | MageSpellFlamestrike | MageSpellFrostfireBolt | MageSpellIgnite |
		MageSpellLivingBomb | MageSpellPyroblast | MageSpellScorch
	MageSpellChill        = MageSpellFrostbolt | MageSpellFrostfireBolt
	MageSpellBrainFreeze  = MageSpellFireball | MageSpellFrostfireBolt
	MageSpellsAllDamaging = MageSpellArcaneBarrage | MageSpellArcaneBlast | MageSpellArcaneExplosion | MageSpellArcaneMissilesTick | MageSpellBlastWave | MageSpellBlizzard | MageSpellDeepFreeze |
		MageSpellDragonsBreath | MageSpellFireBlast | MageSpellFireball | MageSpellFlamestrike | MageSpellFrostbolt | MageSpellFrostfireBolt |
		MageSpellFrostfireOrb | MageSpellIceLance | MageSpellLivingBombExplosion | MageSpellLivingBombDot | MageSpellPyroblast | MageSpellPyroblastDot | MageSpellScorch
	MageSpellInstantCast = MageSpellArcaneBarrage | MageSpellArcaneMissilesCast | MageSpellArcaneMissilesTick |
		MageSpellFireBlast | MageSpellArcaneExplosion | MageSpellBlastWave |
		MageSpellCombustionApplication | MageSpellConeOfCold | MageSpellDeepFreeze |
		MageSpellDragonsBreath | MageSpellIceLance | MageSpellManaGems | MageSpellMirrorImage |
		MageSpellPresenceOfMind | MageSpellLivingBombDot | MageSpellFrostBomb | MageSpellNetherTempest
	MageSpellExtraResult = MageSpellLivingBombExplosion | MageSpellArcaneMissilesTick | MageSpellBlizzard
)
