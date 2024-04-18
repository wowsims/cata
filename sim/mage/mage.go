package mage

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/talent_trees"
)

const (
	SpellFlagMage   = core.SpellFlagAgentReserved1
	BarrageSpells   = core.SpellFlagAgentReserved2
	HotStreakSpells = core.SpellFlagAgentReserved3
)

var TalentTreeSizes = [3]int{21, 21, 19}

type Mage struct {
	core.Character

	Talents       *proto.MageTalents
	Options       *proto.MageOptions
	ArcaneOptions *proto.ArcaneMage_Options
	FireOptions   *proto.FireMage_Options
	FrostOptions  *proto.FrostMage_Options

	//waterElemental *WaterElemental
	//mirrorImage    *MirrorImage

	// Cached values for a few mechanics.
	bonusCritDamage float64

	ArcaneBarrage           *core.Spell
	ArcaneBlast             *core.Spell
	ArcaneExplosion         *core.Spell
	ArcaneMissiles          *core.Spell
	ArcaneMissilesTickSpell *core.Spell
	Blizzard                *core.Spell
	DeepFreeze              *core.Spell
	Ignite                  *core.Spell
	LivingBomb              *core.Spell
	Fireball                *core.Spell
	FireBlast               *core.Spell
	Flamestrike             *core.Spell
	FlamestrikeRank8        *core.Spell
	Frostbolt               *core.Spell
	FrostfireBolt           *core.Spell
	IceLance                *core.Spell
	Pyroblast               *core.Spell
	Scorch                  *core.Spell
	MirrorImage             *core.Spell
	BlastWave               *core.Spell
	DragonsBreath           *core.Spell

	IcyVeins             *core.Spell
	SummonWaterElemental *core.Spell

	ArcaneBlastAura    *core.Aura
	ArcanePotencyAura  *core.Aura
	ArcanePowerAura    *core.Aura
	MissileBarrageAura *core.Aura
	ClearcastingAura   *core.Aura
	ScorchAuras        core.AuraArray
	hotStreakCritAura  *core.Aura
	HotStreakAura      *core.Aura
	CombustionAura     *core.Aura
	FingersOfFrostAura *core.Aura
	BrainFreezeAura    *core.Aura

	CritDebuffCategories core.ExclusiveCategoryArray
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

	// if mage.Talents.ArcaneEmpowerment == 3 {
	// 	raidBuffs.ArcaneEmpowerment = true
	// }
}
func (mage *Mage) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (mage *Mage) Initialize() {
	// mage.registerArcaneBarrageSpell()
	// mage.registerArcaneBlastSpell()
	// mage.registerArcaneExplosionSpell()
	// mage.registerArcaneMissilesSpell()
	// mage.registerBlizzardSpell()
	// mage.registerDeepFreezeSpell()
	// mage.registerFireballSpell()
	// mage.registerFireBlastSpell()
	// mage.registerFlamestrikeSpells()
	// mage.registerFrostboltSpell()
	// mage.registerIceLanceSpell()
	// mage.registerPyroblastSpell()
	// mage.registerScorchSpell()
	// mage.registerLivingBombSpell()
	// mage.registerFrostfireBoltSpell()
	// mage.registerEvocation()
	// mage.registerManaGemsCD()
	// mage.registerMirrorImageCD()
	// mage.registerBlastWaveSpell()
	// mage.registerDragonsBreathSpell()
	// mage.registerSummonWaterElementalCD()
}

func (mage *Mage) Reset(sim *core.Simulation) {
}

func NewMage(character *core.Character, options *proto.Player, mageOptions *proto.MageOptions) *Mage {
	mage := &Mage{
		Character: *character,
		Talents:   &proto.MageTalents{},
		Options:   mageOptions,
	}
	// core.FillTalentsProto(mage.Talents.ProtoReflect(), options.TalentsString, TalentTreeSizes)
	mage.FillTalentsData(talent_trees.MageTalentsConfig, options.TalentsString)

	// mage.bonusCritDamage = .25*float64(mage.Talents.SpellPower) + .1*float64(mage.Talents.Burnout)
	// mage.EnableManaBar()

	// if mage.Options.Armor == proto.MageOptions_MageArmor {
	// 	mage.PseudoStats.SpiritRegenRateCasting += .5
	// 	if mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfMageArmor) {
	// 		mage.PseudoStats.SpiritRegenRateCasting += .2
	// 	}
	// 	if mage.HasSetBonus(ItemSetKhadgarsRegalia, 2) {
	// 		mage.PseudoStats.SpiritRegenRateCasting += .1
	// 	}
	// } else if mage.Options.Armor == proto.MageOptions_MoltenArmor {
	// 	//Need to switch to spirit crit calc
	// 	multi := 0.35
	// 	if mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfMoltenArmor) {
	// 		multi += .2
	// 	}
	// 	if mage.HasSetBonus(ItemSetKhadgarsRegalia, 2) {
	// 		multi += .15
	// 	}
	// 	mage.Character.AddStatDependency(stats.Spirit, stats.SpellCrit, multi)
	// }

	// mage.mirrorImage = mage.NewMirrorImage()

	return mage
}

// Agent is a generic way to access underlying mage on any of the agents.
type MageAgent interface {
	GetMage() *Mage
}
