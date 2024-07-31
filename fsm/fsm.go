package fsm

type FSMState uint8

const (
	FSMStateHandshake FSMState = iota
	FSMStateStatus
	FSMStateLogin
	FSMStateConfiguration
	FSMStatePlay
)

type FSM struct {
	currentState FSMState
}

func NewFSM() *FSM {
	return &FSM{
		currentState: FSMStateHandshake,
	}
}

func (fsm *FSM) SetState(state FSMState) {
	fsm.currentState = state
}

func (fsm *FSM) State() FSMState {
	return fsm.currentState
}
