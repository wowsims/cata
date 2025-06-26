package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type Mage struct {
	core.Character

	ClassSpellScaling float64

	Talents       *proto.MageTalents
	Options       *proto.MageOptions
	ArcaneOptions *proto.ArcaneMage_Options
	FireOptions   *proto.FireMage_Options
	FrostOptions  *proto.FrostMage_Options

	mirrorImages []*MirrorImage

	AlterTime            *core.Spell
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

	AlterTimeAura        *core.Aura
	InvocationAura       *core.Aura
	RuneOfPowerAura      *core.Aura
	PresenceOfMindAura   *core.Aura
	FingersOfFrostAura   *core.Aura
	BrainFreezeAura      *core.Aura
	IcyVeinsAura         *core.Aura
	IceFloesAura         *core.Aura
	IciclesAura          *core.Aura
	ArcaneChargesAura    *core.Aura
	HeatingUp            *core.Aura
	InstantPyroblastAura *core.Aura

	T15_4PC_ArcaneChargeEffect  float64
	T15_4PC_FrostboltProcChance float64
	Icicles                     []float64

	// Item sets
	T12_4pc *core.Aura
	T13_4pc *core.Aura
	T14_4pc *core.Aura
	T16_4pc *core.Aura
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
		if sim.Proc(0.15+core.TernaryFloat64(spell.Matches(MageSpellFrostbolt), mage.T15_4PC_FrostboltProcChance, 0), "FingersOfFrostProc") {
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

	// mage.registerArcaneExplosionSpell()
	mage.registerBlizzardSpell()
	mage.registerConeOfColdSpell()
	mage.registerDeepFreezeSpell()
	mage.registerFlamestrikeSpell()
	mage.registerIceLanceSpell()
	mage.registerFrostfireBoltSpell()
	mage.registerEvocation()
	mage.registerFireBlastSpell()
	mage.registerManaGems()
	mage.registerMirrorImageCD()
	mage.registerfrostNovaSpell()
	mage.registerIceLanceSpell()
	mage.registerIcyVeinsCD()
	mage.registerHeatingUp()
	mage.registerAlterTimeCD()
}

func (mage *Mage) registerMastery() {
	mage.registerFrostMastery()
}

func (mage *Mage) Reset(sim *core.Simulation) {
	mage.T15_4PC_ArcaneChargeEffect = 1.0
	mage.T15_4PC_FrostboltProcChance = 0
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
	mage.T15_4PC_ArcaneChargeEffect = 1.0
	mage.T15_4PC_FrostboltProcChance = 0

	return mage
}

// Agent is a generic way to access underlying mage on any of the agents.
type MageAgent interface {
	GetMage() *Mage
}

const (
	FireSpellMaxTimeUntilResult       = 750 * time.Millisecond
	HeatingUpDeactivateBuffer         = 250 * time.Millisecond
	MageSpellFlagNone           int64 = 0
	MageSpellAlterTime          int64 = 1 << iota
	MageSpellArcaneBarrage
	MageSpellArcaneBlast
	MageSpellArcaneExplosion
	MageSpellArcanePower
	MageSpellArcaneMissilesCast
	MageSpellArcaneMissilesTick
	MageSpellBlizzard
	MageSpellConeOfCold
	MageSpellDeepFreeze
	MageSpellDragonsBreath
	MageSpellEvocation
	MageSpellFireBlast
	MageSpellFireball
	MageSpellFlamestrike
	MageSpellFlamestrikeDot
	MageSpellFrostArmor
	MageSpellFrostbolt
	MageSpellFrostBomb
	MageSpellFrostBombExplosion
	MageSpellFrostfireBolt
	MageSpellFrostNova
	MageSpellFrozenOrb
	MageSpellFrozenOrbTick
	MageSpellIcicle
	MageSpellIceFloes
	MageSpellIceLance
	MageSpellIcyVeins
	MageSpellIgnite
	MageSpellInfernoBlast
	MageSpellLivingBombApply
	MageSpellLivingBombExplosion
	MageSpellLivingBombDot
	MageSpellMageArmor
	MageSpellManaGems
	MageSpellMirrorImage
	MageSpellMoltenArmor
	MageSpellNetherTempest
	MageSpellNetherTempestDot
	MageSpellPresenceOfMind
	MageSpellPyroblast
	MageSpellPyroblastDot
	MagespellRuneOfPower
	MageSpellScorch
	MageSpellCombustion
	MageSpellCombustionDot
	MageMirrorImageSpellArcaneBlast
	MageWaterElementalSpellWaterBolt
	MageSpellLast
	MageSpellsAll       = MageSpellLast<<1 - 1
	MageSpellLivingBomb = MageSpellLivingBombDot | MageSpellLivingBombExplosion
	MageSpellFire       = MageSpellDragonsBreath | MageSpellFireball | MageSpellCombustion |
		MageSpellFireBlast | MageSpellFlamestrike | MageSpellFrostfireBolt | MageSpellIgnite |
		MageSpellLivingBomb | MageSpellPyroblast | MageSpellScorch
	MageSpellBrainFreeze  = MageSpellFireball | MageSpellFrostfireBolt
	MageSpellsAllDamaging = MageSpellArcaneBarrage | MageSpellArcaneBlast | MageSpellArcaneExplosion | MageSpellArcaneMissilesTick | MageSpellBlizzard | MageSpellDeepFreeze |
		MageSpellDragonsBreath | MageSpellFireBlast | MageSpellFireball | MageSpellFlamestrike | MageSpellFrostbolt | MageSpellFrostfireBolt | MageSpellFrozenOrbTick |
		MageSpellIceLance | MageSpellLivingBombExplosion | MageSpellLivingBombDot | MageSpellPyroblast | MageSpellPyroblastDot | MageSpellScorch | MageSpellInfernoBlast
	MageSpellInstantCast = MageSpellArcaneBarrage | MageSpellArcaneMissilesCast | MageSpellArcaneMissilesTick |
		MageSpellFireBlast | MageSpellArcaneExplosion | MageSpellInfernoBlast | MageSpellPyroblastDot |
		MageSpellCombustion | MageSpellConeOfCold | MageSpellDeepFreeze |
		MageSpellDragonsBreath | MageSpellIceLance | MageSpellManaGems | MageSpellMirrorImage |
		MageSpellPresenceOfMind | MageSpellLivingBombDot | MageSpellFrostBomb | MageSpellNetherTempest | MageSpellNetherTempestDot
	MageSpellExtraResult = MageSpellLivingBombExplosion | MageSpellArcaneMissilesTick | MageSpellBlizzard
	FireSpellIgnitable   = MageSpellFireball | MageSpellFrostfireBolt | MageSpellInfernoBlast | MageSpellScorch | MageSpellPyroblast
)
