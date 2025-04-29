package demonology

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/warlock"
)

func RegisterDemonologyWarlock() {
	core.RegisterAgentFactory(
		proto.Player_DemonologyWarlock{},
		proto.Spec_SpecDemonologyWarlock,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewDemonologyWarlock(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_DemonologyWarlock)
			if !ok {
				panic("Invalid spec value for Demonology Warlock!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewDemonologyWarlock(character *core.Character, options *proto.Player) *DemonologyWarlock {
	demoOptions := options.GetDemonologyWarlock().Options

	demonology := &DemonologyWarlock{
		Warlock: warlock.NewWarlock(character, options, demoOptions.ClassOptions),
	}

	return demonology
}

type DemonologyWarlock struct {
	*warlock.Warlock
}

func (demonology *DemonologyWarlock) GetWarlock() *warlock.Warlock {
	return demonology.Warlock
}

func (demonology *DemonologyWarlock) Initialize() {
	demonology.Warlock.Initialize()

	// demonology.registerHandOfGuldan()
	// demonology.registerMetamorphosis()
	// demonology.registerSummonFelguard()
}

func (demonology *DemonologyWarlock) ApplyTalents() {
	demonology.Warlock.ApplyTalents()

	// Demonic Knowledge
	demonology.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		School:     core.SpellSchoolShadow | core.SpellSchoolFire,
		FloatValue: 0.15,
	})
}

// func (demonology *DemonologyWarlock) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
// 	raidBuffs.DemonicPact = demonology.Talents.DemonicPact && demonology.Options.Summon != proto.WarlockOptions_NoSummon
// }

func (demonology *DemonologyWarlock) Reset(sim *core.Simulation) {
	demonology.Warlock.Reset(sim)
}

// func (demonology *DemonologyWarlock) registerSummonFelguard() {
// 	stunActionID := core.ActionID{SpellID: 32752}

// 	demonology.Felguard.RegisterAura(demonology.GetSummonStunAura())
// 	demonology.RegisterSpell(core.SpellConfig{
// 		ActionID:       core.ActionID{SpellID: 30146},
// 		SpellSchool:    core.SpellSchoolShadow,
// 		ProcMask:       core.ProcMaskEmpty,
// 		Flags:          core.SpellFlagAPL,
// 		ClassSpellMask: warlock.WarlockSpellSummonFelguard,

// 		ManaCost: core.ManaCostOptions{BaseCostPercent: 80},
// 		Cast: core.CastConfig{
// 			DefaultCast: core.Cast{
// 				GCD:      core.GCDDefault,
// 				CastTime: 6 * time.Second,
// 			},
// 			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
// 				demonology.ActivatePetSummonStun(sim, stunActionID)
// 			},
// 		},

// 		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 			demonology.SoulBurnAura.Deactivate(sim)
// 			demonology.ChangeActivePet(sim, demonology.Warlock.Felguard)
// 		},
// 	})
// }
