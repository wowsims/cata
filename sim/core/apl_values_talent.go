package core

import (
	"fmt"

	"github.com/wowsims/cata/sim/core/proto"
)

type APLTalent struct {
	id    int32
	known bool
}

type APLValueTalentKnown struct {
	DefaultAPLValueImpl
	talent APLTalent
}

func (rot *APLRotation) newValueTalentIsKnown(config *proto.APLValueTalentIsKnown) APLValue {
	character := rot.unit.Env.Raid.GetPlayerFromUnit(rot.unit).GetCharacter()
	talents := character.talents
	if talents == nil {
		rot.ValidationWarning("%s does not have talent data created", rot.unit.Label)
		return nil
	}
	spellId := config.GetTalentId().GetSpellId()
	return &APLValueTalentKnown{
		talent: APLTalent{
			id:    spellId,
			known: talents.IsKnown(spellId),
		},
	}
}
func (value *APLValueTalentKnown) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueTalentKnown) GetBool(sim *Simulation) bool {
	return value.talent.known
}
func (value *APLValueTalentKnown) String() string {
	return fmt.Sprintf("Talent Known(%d)", value.talent.id)
}
