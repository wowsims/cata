package simsignals

import (
	"errors"
	"fmt"
	"sync"
)

var trackedIds = map[string]Signals{}
var lock = sync.RWMutex{}

func RegisterWithId(id string) (Signals, error) {
	if id == "" {
		return Signals{}, errors.New("id is empty")
	}

	lock.Lock()
	defer lock.Unlock()

	_, exists := trackedIds[id]
	if exists {
		return Signals{}, fmt.Errorf("id %s is not unique", id)
	}

	signals := CreateSignals()

	trackedIds[id] = signals
	return signals, nil
}

func UnregisterId(id string) {
	lock.Lock()
	defer lock.Unlock()
	delete(trackedIds, id)
}

func getById(id string) (Signals, bool) {
	lock.RLock()
	defer lock.RUnlock()
	signals, ok := trackedIds[id]
	return signals, ok
}

func AbortById(id string) bool {
	sig, ok := getById(id)
	if !ok {
		return false
	}
	sig.Abort.Trigger()
	return true
}
