package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

const (
	SpellFlagMage       = core.SpellFlagAgentReserved1
	ArcaneMissileSpells = core.SpellFlagAgentReserved2
	HotStreakSpells     = core.SpellFlagAgentReserved3
	BrainFreezeSpells   = core.SpellFlagAgentReserved4
)

var TalentTreeSizes = [3]int{21, 21, 19}

type Mage struct {
	core.Character

	moltenArmorMod    *core.SpellMod
	arcanePowerGCDmod *core.SpellMod

	Talents       *proto.MageTalents
	Options       *proto.MageOptions
	ArcaneOptions *proto.ArcaneMage_Options
	FireOptions   *proto.FireMage_Options
	FrostOptions  *proto.FrostMage_Options

	//waterElemental *WaterElemental
	mirrorImage *MirrorImage
	flameOrb    *FlameOrb

	// Cached values for a few mechanics.
	bonusCritDamage float64

	ArcaneBarrage           *core.Spell
	ArcaneBlast             *core.Spell
	ArcaneExplosion         *core.Spell
	ArcaneMissiles          *core.Spell
	ArcaneMissilesTickSpell *core.Spell
	ArcanePower             *core.Spell
	Blizzard                *core.Spell
	Combustion              *core.Spell
	DeepFreeze              *core.Spell
	Ignite                  *core.Spell
	LivingBomb              *core.Spell
	Fireball                *core.Spell
	FireBlast               *core.Spell
	FlameOrb                *core.Spell
	FlameOrbExplode         *core.Spell
	Flamestrike             *core.Spell
	Frostbolt               *core.Spell
	FrostfireBolt           *core.Spell
	FrostfireOrb            *core.Spell
	FrostfireOrbTickSpell   *core.Spell
	IceLance                *core.Spell
	Pyroblast               *core.Spell
	PyroblastDot            *core.Spell
	Scorch                  *core.Spell
	MirrorImage             *core.Spell
	BlastWave               *core.Spell
	DragonsBreath           *core.Spell
	IcyVeins                *core.Spell
	SummonWaterElemental    *core.Spell

	ArcaneBlastAura        *core.Aura
	ArcaneMissilesProcAura *core.Aura
	ArcanePotencyAura      *core.Aura
	ArcanePowerAura        *core.Aura
	BrainFreezeAura        *core.Aura
	ClearcastingAura       *core.Aura
	CriticalMassAuras      core.AuraArray
	FingersOfFrostAura     *core.Aura
	FlameOrbTimer          *core.Aura
	FrostArmorAura         *core.Aura
	GlyphedFrostArmorPA    *core.PendingAction
	hotStreakCritAura      *core.Aura
	HotStreakAura          *core.Aura
	MageArmorAura          *core.Aura
	MageArmorPA            *core.PendingAction
	PyromaniacAura         *core.Aura

	ScalingBaseDamage float64

	CritDebuffCategories core.ExclusiveCategoryArray
}

func (mage *Mage) GetCharacter() *core.Character {
	return &mage.Character
}

func (mage *Mage) GetMage() *Mage {
	return mage
}

func (mage *Mage) HasPrimeGlyph(glyph proto.MagePrimeGlyph) bool {
	return mage.HasGlyph(int32(glyph))
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
	mage.applyGlyphs()
	mage.applyArmor()
	mage.registerArcaneBarrageSpell()
	mage.registerArcaneBlastSpell()
	mage.registerArcaneExplosionSpell()
	mage.registerArcaneMissilesSpell()
	mage.registerBlizzardSpell()
	mage.registerDeepFreezeSpell()
	mage.registerFireballSpell()
	mage.registerFireBlastSpell()
	mage.registerFlameOrbSpell()
	mage.registerFlameOrbExplodeSpell()
	mage.registerFlamestrikeSpell()
	mage.registerFrostboltSpell()
	mage.registerFrostfireOrbSpell()
	mage.registerIceLanceSpell()
	mage.registerPyroblastSpell()
	mage.registerScorchSpell()
	mage.registerLivingBombSpell()
	mage.registerFrostfireBoltSpell()
	mage.registerEvocation()
	mage.registerManaGemsCD()
	mage.registerMirrorImageCD()
	mage.registerCombustionSpell()
	mage.registerBlastWaveSpell()
	mage.registerDragonsBreathSpell()
	// mage.registerSummonWaterElementalCD()

	mage.applyArcaneMastery()
	mage.applyFireMastery()
	mage.applyArcaneMissileProc()
	mage.ScalingBaseDamage = 937.330078125
}

func (mage *Mage) Reset(sim *core.Simulation) {
}

func (mage *Mage) applyArmor() {
	if mage.Options.Armor == proto.MageOptions_MoltenArmor {
		if mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfMoltenArmor) {
			mage.AddStaticMod(core.SpellModConfig{
				ClassMask:  MageSpellsAll,
				FloatValue: 7 * core.CritRatingPerCritChance,
				Kind:       core.SpellMod_BonusCrit_Rating,
			})
		} else {
			mage.AddStaticMod(core.SpellModConfig{
				ClassMask:  MageSpellsAll,
				FloatValue: 5 * core.CritRatingPerCritChance,
				Kind:       core.SpellMod_BonusCrit_Rating,
			})
		}
	} else if mage.Options.Armor == proto.MageOptions_MageArmor {
		hasGlyph := mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfMageArmor)
		manaRegenPerSecond := mage.MaxMana() * core.TernaryFloat64(hasGlyph, .072, 0.06)
		// TODO regen 3% max mana as mp5 aka 0.6% max mana per second
		mage.MageArmorAura = core.MakePermanent(mage.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 6117},
			Label:    "Mage Armor",
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				mage.MageArmorPA = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
					Period: time.Second * 1,
					OnAction: func(sim *core.Simulation) {
						mage.AddMana(sim, manaRegenPerSecond, mage.NewManaMetrics(core.ActionID{SpellID: 6117}))
					},
				})
			},
		}))
	}

}
func NewMage(character *core.Character, options *proto.Player, mageOptions *proto.MageOptions) *Mage {
	mage := &Mage{
		Character: *character,
		Talents:   &proto.MageTalents{},
		Options:   mageOptions,
	}

	core.FillTalentsProto(mage.Talents.ProtoReflect(), options.TalentsString, TalentTreeSizes)

	// mage.EnableManaBar()

	mage.mirrorImage = mage.NewMirrorImage()
	mage.flameOrb = mage.NewFlameOrb()
	mage.EnableManaBar()
	return mage
}

/*
--------------------------------------
	Arcane Mastery
---------------------------------------
*/
//Increases all spell damage done by up to 12%, based on the amount of mana the Mage has unspent.
//Each point of Mastery increases damage done by up to an additional 1.5%.

func (mage *Mage) GetArcaneMasteryBonus() float64 {
	return (1.12 + 0.015*mage.GetMasteryPoints())
}

func (mage *Mage) applyArcaneMastery() {
	// Arcane Mastery Mod
	arcaneMastery := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellsAll,
		FloatValue: mage.CurrentMana() / mage.MaxMana() * mage.GetArcaneMasteryBonus(), //take current % of mana, get that portion of damage bonus
		Kind:       core.SpellMod_DamageDone_Pct,
	})

	mage.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery, newMastery float64) {
		arcaneMastery.UpdateFloatValue(1.22 + 0.28*core.MasteryRatingToMasteryPoints(newMastery))
	})

	core.MakePermanent(mage.GetOrRegisterAura(core.Aura{
		Label:    "Mana Adept",
		ActionID: core.ActionID{SpellID: 76547},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			arcaneMastery.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			arcaneMastery.Deactivate()
		},
	}))
}

/*
--------------------------------------

	Fire Mastery

---------------------------------------
*/
func (mage *Mage) GetFireMasteryBonus() float64 {
	return (1.22 + 0.28*mage.GetMasteryPoints())
}

func (mage *Mage) applyFireMastery() {
	// Fire Mastery Mod
	fireMastery := mage.AddDynamicMod(core.SpellModConfig{
		ClassMask:  MageSpellFireDoT,
		FloatValue: mage.GetFireMasteryBonus(),
		Kind:       core.SpellMod_DamageDone_Pct,
	})

	mage.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery, newMastery float64) {
		fireMastery.UpdateFloatValue(1.22 + 0.28*core.MasteryRatingToMasteryPoints(newMastery))
	})
}

/* --------------------------------------
				Frost Mastery
---------------------------------------*/

// Agent is a generic way to access underlying mage on any of the agents.
type MageAgent interface {
	GetMage() *Mage
}

func (mage *Mage) applyArcaneMissileProc() {
	if mage.Talents.HotStreak || mage.Talents.BrainFreeze > 0 {
		return
	}

	t10ProcAura := mage.BloodmagesRegalia2pcAura()

	// Aura for when proc is successful
	mage.ArcaneMissilesProcAura = mage.RegisterAura(core.Aura{
		Label:    "Arcane Missiles Proc",
		ActionID: core.ActionID{SpellID: 79683},
		Duration: time.Second * 20,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == mage.ArcaneMissiles {
				aura.Deactivate(sim)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if t10ProcAura != nil {
				t10ProcAura.Activate(sim)
			}
		},
	})

	procChance := 0.4

	// Listener for procs
	mage.RegisterAura(core.Aura{
		Label:    "Arcane Missiles Activation",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Flags.Matches(ArcaneMissileSpells) {
				return
			}
			if sim.Proc(procChance, "Arcane Missiles") {
				mage.ArcaneMissilesProcAura.Activate(sim)
			}
		},
	})
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
	MageSpellCombustion
	MageSpellConeOfCold
	MageSpellDeepFreeze
	MageSpellDragonsBreath
	MageSpellEvocation
	MageSpellFireBlast
	MageSpellFireball
	MageSpellFlamestrike
	MageSpellFlameOrb
	MageSpellFocusMagic
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
	MageSpellPyroblast
	MageSpellPyroblastDot
	MageSpellScorch
	MageSpellLast

	MageSpellsAll         = MageSpellLast<<1 - 1
	MageSpellLivingBomb   = MageSpellLivingBombDot | MageSpellLivingBombExplosion
	MageSpellFireDoT      = MageSpellLivingBombDot | MageSpellPyroblastDot | MageSpellIgnite | MageSpellCombustion
	MageSpellChill        = MageSpellFrostbolt | MageSpellFrostfireBolt
	MageSpellBrainFreeze  = MageSpellFireball | MageSpellFrostfireBolt
	MageSpellsAllDamaging = MageSpellArcaneBarrage | MageSpellArcaneBlast | MageSpellArcaneExplosion | /*MageSpellArcaneMissiles | */ MageSpellBlastWave | MageSpellBlizzard | MageSpellDeepFreeze |
		MageSpellDragonsBreath | MageSpellFireBlast | MageSpellFireball | MageSpellFlamestrike | MageSpellFlameOrb | MageSpellFrostbolt | MageSpellFrostfireBolt |
		MageSpellFrostfireOrb | MageSpellIceLance | MageSpellLivingBombExplosion | MageSpellLivingBombDot | MageSpellPyroblast | MageSpellPyroblastDot | MageSpellScorch
)
