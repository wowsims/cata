package marksmanship

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/hunter"
)

func RegisterMarksmanshipHunter() {
	core.RegisterAgentFactory(
		proto.Player_MarksmanshipHunter{},
		proto.Spec_SpecMarksmanshipHunter,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewMarksmanshipHunter(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_MarksmanshipHunter)
			if !ok {
				panic("Invalid spec value for Marksmanship Hunter!")
			}
			player.Spec = playerSpec
		},
	)
}
func (mm *MarksmanshipHunter) applyMastery() {
	actionID := core.ActionID{SpellID: 76659}

	wqSpell := mm.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagRanged,

		DamageMultiplier: 0.8, // Wowwiki says it remains 80%
		CritMultiplier:   mm.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.RangedWeaponDamage(sim, spell.RangedAttackPower())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
		},
	})

	mm.RegisterAura(core.Aura{
		Label:    "Wild Quiver Mastery",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask != core.ProcMaskRangedSpecial && spell != mm.AutoAttacks.RangedAuto() {
				return
			}
			procChance := (mm.CalculateMasteryPoints() + 8) * 0.02
			if spell.ClassSpellMask == hunter.HunterSpellBarrage {
				procChance = procChance / 6
			}
			if sim.RandomFloat("Wild Quiver") < procChance {
				wqSpell.Cast(sim, result.Target)
			}
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ClassSpellMask != hunter.HunterSpellBarrage {
				return
			}

			procChance := (mm.CalculateMasteryPoints() + 8) * 0.02
			if sim.RandomFloat("Wild Quiver") < procChance/6 {
				wqSpell.Cast(sim, result.Target)
			}
		},
	})
}
func NewMarksmanshipHunter(character *core.Character, options *proto.Player) *MarksmanshipHunter {
	mmOptions := options.GetMarksmanshipHunter().Options

	mmHunter := &MarksmanshipHunter{
		Hunter: hunter.NewHunter(character, options, mmOptions.ClassOptions),
	}
	mmHunter.MarksmanshipOptions = mmOptions
	return mmHunter
}
func (mmHunter *MarksmanshipHunter) Initialize() {
	mmHunter.Hunter.Initialize()
	// MM Hunter Spec Bonus
	// mmHunter.AddStaticMod(core.SpellModConfig{
	// 	Kind:       core.SpellMod_DamageDone_Flat,
	// 	ProcMask:   core.ProcMaskRangedAuto,
	// 	FloatValue: 0.15,
	// })

	mmHunter.registerAimedShotSpell()
	mmHunter.registerChimeraShotSpell()
	mmHunter.registerSteadyShotSpell()
	mmHunter.applyMastery()
	mmHunter.ApplySpecialization()
}

type MarksmanshipHunter struct {
	*hunter.Hunter
}

func (mmHunter *MarksmanshipHunter) GetHunter() *hunter.Hunter {
	return mmHunter.Hunter
}

func (mmHunter *MarksmanshipHunter) Reset(sim *core.Simulation) {
	mmHunter.Hunter.Reset(sim)
}
