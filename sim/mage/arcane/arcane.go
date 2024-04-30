package arcane

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/mage"
)

func RegisterArcaneMage() {
	core.RegisterAgentFactory(
		proto.Player_ArcaneMage{},
		proto.Spec_SpecArcaneMage,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewArcaneMage(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_ArcaneMage)
			if !ok {
				panic("Invalid spec value for Arcane Mage!")
			}
			player.Spec = playerSpec
		},
	)
}

type ArcaneMage struct {
	*mage.Mage

	Options *proto.ArcaneMage_Options
}

func NewArcaneMage(character *core.Character, options *proto.Player) *ArcaneMage {
	arcaneOptions := options.GetArcaneMage().Options

	arcaneMage := &ArcaneMage{
		Mage: mage.NewMage(character, options, arcaneOptions.ClassOptions),
	}
	arcaneMage.ArcaneOptions = arcaneOptions

	return arcaneMage
}

func (arcaneMage *ArcaneMage) GetMage() *mage.Mage {
	return arcaneMage.Mage
}

func (arcaneMage *ArcaneMage) Reset(sim *core.Simulation) {
	arcaneMage.Mage.Reset(sim)
}

func (arcaneMage *ArcaneMage) Initialize() {
	arcaneMage.Mage.Initialize()

	arcaneMage.registerArcaneBarrageSpell()
	arcaneMage.applyArcaneMissileProc()
}

func (arcane *ArcaneMage) applyArcaneMissileProc() {
	// Aura for when proc is successful
	arcane.ArcaneMissilesProcAura = arcane.RegisterAura(core.Aura{
		Label:    "Arcane Missiles Proc",
		ActionID: core.ActionID{SpellID: 79683},
		Duration: time.Second * 20,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ClassSpellMask == mage.MageSpellArcaneMissilesCast {
				aura.Deactivate(sim)
			}
		},
	})

	procChance := 0.4

	const MageSpellsArcaneMissilesNow = mage.MageSpellArcaneBarrage | mage.MageSpellArcaneBlast |
		mage.MageSpellFireball | mage.MageSpellFrostbolt | mage.MageSpellFrostfireBolt | mage.MageSpellFrostfireOrb

	// Listener for procs
	core.MakePermanent(arcane.RegisterAura(core.Aura{
		Label: "Arcane Missiles Activation",
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ClassSpellMask&MageSpellsArcaneMissilesNow == 0 {
				return
			}
			if sim.Proc(procChance, "Arcane Missiles") {
				arcane.ArcaneMissilesProcAura.Activate(sim)
			}
		},
	}))
}

func (arcaneMage *ArcaneMage) ApplyTalents() {

	arcaneMage.Mage.ApplyTalents()
	// Arcane Specialization Bonus
	arcaneMage.AddStaticMod(core.SpellModConfig{
		School:     core.SpellSchoolArcane,
		FloatValue: 0.25,
		Kind:       core.SpellMod_DamageDone_Pct,
	})

	// Arcane Mastery

	arcaneMastery := arcaneMage.AddDynamicMod(core.SpellModConfig{
		ClassMask: mage.MageSpellsAll,
		Kind:      core.SpellMod_DamageDone_Pct,
	})

	arcaneMage.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery, newMastery float64) {
		arcaneMastery.UpdateFloatValue(arcaneMage.ArcaneMasteryValue())
	})

	core.MakePermanent(arcaneMage.GetOrRegisterAura(core.Aura{
		Label: "Mana Adept",
		//ActionID: core.ActionID{SpellID: 76547},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			arcaneMastery.UpdateFloatValue(arcaneMage.ArcaneMasteryValue())
			arcaneMastery.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			arcaneMastery.Deactivate()
		},
	}))

	core.MakeProcTriggerAura(&arcaneMage.Unit, core.ProcTrigger{
		Name:     "Arcane Mastery Mana Updater",
		Callback: core.CallbackOnCastComplete,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			arcaneMastery.UpdateFloatValue(arcaneMage.ArcaneMasteryValue())
		},
	})

}

func (arcaneMage *ArcaneMage) GetArcaneMasteryBonus() float64 {
	return (0.12 + 0.015*arcaneMage.GetMasteryPoints())
}

func (arcaneMage *ArcaneMage) ArcaneMasteryValue() float64 {
	return arcaneMage.GetArcaneMasteryBonus() * (arcaneMage.CurrentMana() / arcaneMage.MaxMana())
}
