package core

import (
	"time"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

// This is just the static bonuses. Most professions are handled elsewhere.
func (character *Character) applyProfessionEffects() {
	if character.HasProfession(proto.Profession_Mining) {
		character.AddStat(stats.Stamina, 120)
	}

	if character.HasProfession(proto.Profession_Skinning) {
		character.AddStats(stats.Stats{stats.CritRating: 80})
	}

	if character.HasProfession(proto.Profession_Herbalism) {
		actionID := ActionID{SpellID: 74497}

		aura := character.NewTemporaryStatsAura(
			"Lifeblood",
			actionID,
			stats.Stats{stats.HasteRating: 2880},
			time.Second*20,
		)

		spell := character.RegisterSpell(SpellConfig{
			ActionID:    actionID,
			SpellSchool: SpellSchoolNature,
			ProcMask:    ProcMaskSpellHealing,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			CritMultiplier:   character.DefaultCritMultiplier(),
			ApplyEffects: func(sim *Simulation, _ *Unit, spell *Spell) {
				amount := sim.RollWithLabel(720, 2160, "Healing Roll")
				spell.CalcAndDealHealing(sim, spell.Unit, amount, spell.OutcomeHealingCrit)
				aura.Activate(sim)
			},
		})
		character.AddMajorCooldown(MajorCooldown{
			Type:  CooldownTypeDPS,
			Spell: spell,
		})
	}
}
