package hunter

func (hunter *Hunter) RegisterLynxRushSpell() {
	if !hunter.Talents.LynxRush {
		return
	}
}
