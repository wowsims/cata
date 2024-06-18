package simsignals

type triggerSignal struct {
	channel chan struct{}
}

func (s *triggerSignal) Trigger() {
	select {
	case <-s.channel: // Already closed
	default:
		close(s.channel)
	}
}

func (s *triggerSignal) IsTriggered() bool {
	select {
	case <-s.channel:
		return true
	default:
	}
	return false
}

type Signals struct {
	Abort triggerSignal
}

func CreateSignals() Signals {
	return Signals{
		Abort: triggerSignal{channel: make(chan struct{})},
	}
}
